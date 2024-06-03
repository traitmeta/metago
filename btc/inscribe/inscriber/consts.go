package inscriber

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/wire"
)

const (
	defaultRevealOutValue = int64(546) // 500 sat, ord default 10000
	minRevealOutValue     = int64(330)
)

const (
	defaultSequenceNum  = wire.MaxTxInSequenceNum - 10
	EnableRbfNoLockTime = wire.MaxTxInSequenceNum - 2
	MaxStandardTxWeight = blockchain.MaxBlockWeight / 10
)
