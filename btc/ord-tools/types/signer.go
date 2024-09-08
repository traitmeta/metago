package types

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
)

type InscribeTxPreview struct {
	Destination              string
	InscriptionWitnessScript []byte
	PrivateKey               *btcec.PrivateKey
	ControlBlockWitness      []byte

	CommitTxPkScript      []byte
	CommitTxAddress       btcutil.Address
	RecoveryPrivateKeyWIF string
}
