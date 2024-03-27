package ord

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/wire"
)

type InscriptionRawTx struct {
	TxPrevOutput   *wire.TxOut
	WitnessScript  *InscriptionWitness
	Size           int64
	Raw            *wire.MsgTx
	RevealOutValue int64
	FeeRate        int64
	PrivateKey     *btcec.PrivateKey
}

type InscriptionWitness struct {
	SignatureWitness    []byte
	InsWitnessScript    []byte
	ControlBlockWitness []byte
}

func NewInscriptionWitness() *InscriptionWitness {
	return &InscriptionWitness{
		SignatureWitness:    make([]byte, 64),
		InsWitnessScript:    nil,
		ControlBlockWitness: make([]byte, 33),
	}
}

func NewInscriptionRawTx() *InscriptionRawTx {
	return &InscriptionRawTx{
		WitnessScript: NewInscriptionWitness(),
	}
}

func (irt *InscriptionRawTx) SetTxPrevOutput(pkScript []byte, prevOutput int64) {
	irt.TxPrevOutput = &wire.TxOut{
		PkScript: pkScript,
		Value:    prevOutput,
	}
}

func (irt *InscriptionRawTx) SetWitnessScript(inscriptionWitnessScript []byte) {
	irt.WitnessScript.InsWitnessScript = inscriptionWitnessScript
}

func (irt *InscriptionRawTx) SetSize(txSize int64) {
	irt.Size = txSize
}

func (irt *InscriptionRawTx) CalcPrevOutput(revealOutValue, feeRate int64) int64 {
	txFee := irt.Size * feeRate
	prevOutput := revealOutValue + txFee
	emptySignature := make([]byte, 64)
	emptyControlBlockWitness := make([]byte, 33)
	witnessSize := (wire.TxWitness{emptySignature, irt.WitnessScript.InsWitnessScript, emptyControlBlockWitness}.SerializeSize() + 2 + 3) / 4
	// 初始化一个空的签名和控制块，计算单个铭文交易，witness部分的额外手续费，并更新totalPrevOutput
	witnessFee := int64(witnessSize) * feeRate
	prevOutput += witnessFee

	return prevOutput
}
