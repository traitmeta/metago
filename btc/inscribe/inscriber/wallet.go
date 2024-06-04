package inscriber

import (
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/pkg/errors"
)

func CreateWallet(net *chaincfg.Params, privateKey *btcec.PrivateKey) (*WalletInfo, error) {
	commitTxAddress, err := btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(privateKey.PubKey()).SerializeCompressed()[1:], net)
	log.Println("commitTxAddress: ", commitTxAddress)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address error")
	}

	commitTxAddressPkScript, err := txscript.PayToAddrScript(commitTxAddress)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address pk script error")
	}

	recoveryPrivateKeyWIF, err := btcutil.NewWIF(txscript.TweakTaprootPrivKey(*privateKey, []byte{}), net, true)
	if err != nil {
		return nil, errors.Wrap(err, "create recovery private key wif error")
	}

	return &WalletInfo{
		PrivateKey:      privateKey,
		Address:         commitTxAddress,
		PkScript:        commitTxAddressPkScript,
		RecoveryPKofWIF: recoveryPrivateKeyWIF.String(),
	}, nil
}
