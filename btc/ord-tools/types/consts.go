package types

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/wire"
)

const (
	DefaultSequenceNum  = wire.MaxTxInSequenceNum - 10
	MaxStandardTxWeight = blockchain.MaxBlockWeight / 10
)

const (
	DefaultRevealOutValue = int64(546)
	MinRevealOutValue     = int64(330)
)
