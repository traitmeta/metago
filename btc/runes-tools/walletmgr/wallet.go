package walletmgr

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/pkg/errors"

	"github.com/traitmeta/metago/btc/runes-tools/txbuilder"
	"github.com/traitmeta/metago/btc/runes-tools/txbuilder/runes"
)

type Wallet struct {
	net          *chaincfg.Params
	privateKey   *btcec.PrivateKey
	scriptPubKey []byte
	cache        *Cache
	cachePrefix  string
}

func InitWallet(wif string, cache *Cache, net *chaincfg.Params) (*Wallet, error) {
	var privateKey *btcec.PrivateKey
	decodedPrivKey, err := base64.StdEncoding.DecodeString(wif)
	if err != nil {
		return nil, errors.Wrap(err, "decode private key error")
	}

	privateKey, _ = btcec.PrivKeyFromBytes(decodedPrivKey)
	address, err := btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(privateKey.PubKey()).SerializeCompressed()[1:], net)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address error")
	}

	scriptPubKey, err := txscript.PayToAddrScript(address)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address pk script error")
	}

	return &Wallet{
		net:          net,
		privateKey:   privateKey,
		scriptPubKey: scriptPubKey,
		cache:        cache,
		cachePrefix:  address.EncodeAddress(),
	}, nil
}

func (w *Wallet) MintRunes(cli *rpcclient.Client, runeId string, destination string) {
	builder := txbuilder.NewMintTxBuilder(w.privateKey, w.net)
	prevInfo, err := w.cache.ReadWalletPrevInfo(w.cachePrefix)
	if err != nil {
		return
	}

	gasFee, err := w.cache.ReadMemPoolGas()
	if err != nil {
		return
	}

	req := runes.EtchRequest{
		FeeRate:     gasFee,
		RuneID:      runeId,
		Destination: destination,
	}

	tx, err := builder.BuildMintTx(*prevInfo, req)
	if err != nil {
		return
	}

	txHash, err := cli.SendRawTransaction(tx, false)
	if err != nil {
		return
	}

	var txBuf bytes.Buffer
	err = tx.Serialize(&txBuf)
	if err != nil {
		return
	}

	err = w.cache.CacheWalletMintTxInfo(w.cachePrefix, txHash.String(), hex.EncodeToString(txBuf.Bytes()))
	if err != nil {
		return
	}
}
