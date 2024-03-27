package ord

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/wire"
)

const (
	defaultSequenceNum = wire.MaxTxInSequenceNum - 10

	MaxStandardTxWeight = blockchain.MaxBlockWeight / 10
)

const (
	defaultRevealOutValue = int64(546) // 500 sat, ord default 10000
	minRevealOutValue     = int64(294)
)
