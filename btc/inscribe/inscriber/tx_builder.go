package inscriber

import (
	"encoding/base64"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/traitmeta/metago/btc/runes-tools/txbuilder/runes"
)

type Builder struct {
	net *chaincfg.Params
}

var addTxInTxOutIntoTx = func(tx *wire.MsgTx, index int, destination, runeId string, net *chaincfg.Params) error {
	in := wire.NewTxIn(&wire.OutPoint{Index: uint32(index)}, nil, nil)
	in.Sequence = defaultSequenceNum
	tx.AddTxIn(in)

	output, err := runes.CreateMintRuneStoneOutput(runeId)
	if err != nil {
		return errors.Wrap(err, "build runestone script fail")
	}
	tx.AddTxOut(output)

	receiver, err := btcutil.DecodeAddress(destination, net) // note: 生成铭文的接收地址, 这里是destination[index]的P2TR地址
	if err != nil {
		return errors.Wrap(err, "decode address error")
	}

	scriptPubKey, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return errors.Wrap(err, "pay to address script error")
	}
	out := wire.NewTxOut(defaultRevealOutValue, scriptPubKey) // note: 再构造reveal tx的交易输出，也就是铭文的接收地址
	tx.AddTxOut(out)
	return nil
}

func (b *Builder) BuildPrivateAndScriptInfo(req MintReq, payAddrPK string) ([]*PrivateKeyAndScriptInfo, error) {
	var pkAndScripts []*PrivateKeyAndScriptInfo
	for i := 0; i < req.Count; i++ {
		var privateKey *btcec.PrivateKey
		var err error
		if i == 0 {
			decodedPrivKey, err := base64.StdEncoding.DecodeString(payAddrPK)
			if err != nil {
				return nil, errors.Wrap(err, "decode private key error")
			}

			privateKey, _ = btcec.PrivKeyFromBytes(decodedPrivKey)
		} else {
			privateKey, err = btcec.NewPrivateKey()
			if err != nil {
				return nil, errors.Wrap(err, "create private key error")
			}
		}

		info, err := createPkAndScriptInfo(b.net, privateKey)
		if err != nil {
			return nil, errors.Wrap(err, "create inscription tx ctx data error")
		}

		pkAndScripts = append(pkAndScripts, info)
	}

	return pkAndScripts, nil
}

// note: include reveal tx1 + commit tx2... + 手续费 + 找零
func (b *Builder) BuildMiddleTxWithEmptyInput(req MintReq, revealTxs []*WrapTx, prevScript []byte) (*WrapTx, error) {
	tx := wire.NewMsgTx(int32(2))
	// TODO 使用OUTPOINT，而不是默认指定0
	if err := addTxInTxOutIntoTx(tx, 0, req.Receiver, req.RuneId, b.net); err != nil {
		return nil, errors.Wrap(err, "add tx in tx out of reveal tx1 into middle tx error")
	}

	prevOutput := defaultRevealOutValue
	minerFee := int64(0)
	emptySignature := make([]byte, 64)
	witnessFee := (int64(wire.TxWitness{emptySignature}.SerializeSize()+2+3) / 4) * req.FeeRate
	prevOutput += witnessFee
	minerFee += witnessFee

	for _, v := range revealTxs {
		tx.AddTxOut(v.PrevOutput)
		prevOutput += v.PrevOutput.Value
		minerFee += v.MinerFee
	}

	output, err := runes.CreateMintRuneStoneOutput(req.RuneId)
	if err != nil {
		return nil, errors.Wrap(err, "build runestone script fail")
	}

	tx.AddTxOut(output)
	txSizeFee := int64(tx.SerializeSize()) * req.FeeRate
	minerFee += txSizeFee
	prevOutput += txSizeFee

	var wrapTx = &WrapTx{
		PrevOutput:          &wire.TxOut{PkScript: prevScript, Value: prevOutput},
		TxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		WireTx:              tx,
		MinerFee:            minerFee,
	}

	return wrapTx, nil
}

func (b *Builder) BuildRevealTxsWithEmptyInput(pkAndScripts []*PrivateKeyAndScriptInfo, req MintReq) ([]*WrapTx, error) {
	// Note: first rune in middle tx，others in reveal txs
	revealTx := make([]*WrapTx, req.Count-1)

	for i := 1; i < req.Count; i++ {
		wrapTx, err := b.buildRevealTxWithEmptyInput(i, pkAndScripts, req)
		if err != nil {
			return nil, errors.Wrap(err, "build empty reveal tx")
		}

		revealTx[i-1] = wrapTx
	}

	return revealTx, nil
}

func (b *Builder) buildRevealTxWithEmptyInput(idx int, pkAndScripts []*PrivateKeyAndScriptInfo, req MintReq) (*WrapTx, error) {
	var revealWrapTx *WrapTx

	tx := wire.NewMsgTx(int32(2))
	err := addTxInTxOutIntoTx(tx, idx, req.Receiver, req.RuneId, b.net)
	if err != nil {
		return nil, errors.Wrap(err, "add tx in tx out into reveal tx error")
	}

	prevOutput := defaultRevealOutValue + int64(tx.SerializeSize())*req.FeeRate // note: 铭文的value + 铭文的交易大小 * feeRate
	emptySignature := make([]byte, 64)
	witness_fee := (int64(wire.TxWitness{emptySignature}.SerializeSize()+2+3) / 4) * req.FeeRate
	prevOutput += witness_fee

	revealWrapTx = &WrapTx{
		PrevOutput:          &wire.TxOut{PkScript: pkAndScripts[idx].PkScript, Value: prevOutput},
		TxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		WireTx:              tx,
		MinerFee:            prevOutput - defaultRevealOutValue,
	}

	return revealWrapTx, nil
}

func createPkAndScriptInfo(net *chaincfg.Params, privateKey *btcec.PrivateKey) (*PrivateKeyAndScriptInfo, error) {
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

	return &PrivateKeyAndScriptInfo{
		PrivateKey:      privateKey,
		Address:         commitTxAddress,
		PkScript:        commitTxAddressPkScript,
		RecoveryPKofWIF: recoveryPrivateKeyWIF.String(),
	}, nil
}
