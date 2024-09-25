package txbuilder

import (
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/traitmeta/metago/btc/runes-tools/txbuilder/runes"
)

const EnableRbfNoLockTime = wire.MaxTxInSequenceNum - 2
const MaxStandardTxWeight = blockchain.MaxBlockWeight / 10
const defaultOutputValue int64 = 546

type MintTxBuilder struct {
	tx                *wire.MsgTx
	net               *chaincfg.Params
	privateKey        *btcec.PrivateKey
	prevOutputFetcher *txscript.MultiPrevOutFetcher
}

func CreateEmptyTx() *wire.MsgTx {
	tx := wire.NewMsgTx(wire.TxVersion)
	in := wire.NewTxIn(&wire.OutPoint{Index: uint32(0)}, nil, nil) // note: 先构造reveal tx的空交易输入
	in.Sequence = EnableRbfNoLockTime
	tx.AddTxIn(in)
	return tx
}

func NewMintTxBuilder(private *btcec.PrivateKey, net *chaincfg.Params) *MintTxBuilder {
	return &MintTxBuilder{tx: CreateEmptyTx(), privateKey: private, net: net, prevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil)}
}

func (mtb *MintTxBuilder) BuildMintTx(prev PrevInfo, req runes.EtchRequest) (*wire.MsgTx, error) {
	prevTxHash, err := chainhash.NewHashFromStr(prev.PrevTxId)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert to hash from giving string")
	}

	mtb.tx.TxIn[0].PreviousOutPoint.Hash = *prevTxHash
	prevOutput := wire.NewTxOut(prev.PreAmount, prev.PrevScript)
	mtb.prevOutputFetcher.AddPrevOut(wire.OutPoint{
		Hash:  *prevTxHash,
		Index: uint32(0),
	}, prevOutput)

	err = mtb.buildMintTxOutput(req.RuneID, prev.PrevScript, req.Destination, defaultOutputValue)
	if err != nil {
		return nil, err
	}

	minerFee := mtb.calcMinerFee(req.FeeRate)
	mtb.tx.TxOut[0].Value = prev.PreAmount - minerFee
	err = mtb.sign()
	if err != nil {
		return nil, err
	}

	return mtb.tx, nil
}

func (mtb *MintTxBuilder) buildMintTxOutput(runeId string, payScript []byte, destination string, changeAmount int64) error {
	// 0 = change address
	out := wire.NewTxOut(changeAmount, payScript)
	mtb.tx.AddTxOut(out)

	// 1 = receiver address
	receiver, err := btcutil.DecodeAddress(destination, mtb.net) // note: 生成铭文的接收地址, 这里是destination[index]的P2TR地址
	if err != nil {
		return errors.Wrap(err, "decode address error")
	}

	scriptPubKey, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return errors.Wrap(err, "pay to address script error")
	}

	mtb.tx.AddTxOut(wire.NewTxOut(defaultOutputValue, scriptPubKey))

	// 2 op_return
	output, err := runes.CreateMintRuneStoneOutput(runeId)
	if err != nil {
		return errors.Wrap(err, "build runestone script fail")
	}
	mtb.tx.AddTxOut(output)

	return nil
}

func (mtb *MintTxBuilder) calcMinerFee(feeRate int64) int64 {
	minerFee := int64(mtb.tx.SerializeSize()) * feeRate
	emptySignature := make([]byte, 64)
	// 初始化一个空的签名和控制块，计算单个铭文交易，witness部分的额外手续费，并更新totalPrevOutput
	fee := (int64(wire.TxWitness{emptySignature}.SerializeSize()+2+3) / 4) * feeRate
	minerFee += fee

	return minerFee
}

func (mtb *MintTxBuilder) sign() error {
	// 更新交易的输入指向相应的commit Tx哈希
	witnessArray, err := txscript.CalcTaprootSignatureHash(txscript.NewTxSigHashes(mtb.tx, mtb.prevOutputFetcher),
		txscript.SigHashDefault, mtb.tx, 0, mtb.prevOutputFetcher)
	if err != nil {
		return errors.Wrap(err, "calc tapscript signaturehash error")
	}

	// 使用私钥对签名哈希进行签名
	signature, err := schnorr.Sign(mtb.privateKey, witnessArray)
	if err != nil {
		return errors.Wrap(err, "schnorr sign error")
	}

	witness := wire.TxWitness{signature.Serialize()}
	mtb.tx.TxIn[0].Witness = witness

	// check tx max tx wight
	revealWeight := blockchain.GetTransactionWeight(btcutil.NewTx(mtb.tx))
	if revealWeight > MaxStandardTxWeight {
		return errors.New(fmt.Sprintf("reveal(index %d) transaction weight greater than %d (MAX_STANDARD_TX_WEIGHT): %d", 0, MaxStandardTxWeight, revealWeight))
	}

	return nil
}
