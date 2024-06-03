package inscriber

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

type BTCBaseClient interface {
	GetRawTransaction(txHash *chainhash.Hash) (*btcutil.Tx, error)
	SendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error)
}
