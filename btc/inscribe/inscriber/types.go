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

type MintReq struct {
	RuneId   string `json:"rune_id"`  // blockNum:TxIdx
	Receiver string `json:"receiver"` // receiver address
	FeeRate  int64  `json:"fee_rate"`
	Count    int    `json:"count"` // number of mint times
}
