package inscriber

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
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

type MintTxs struct {
	PayTxHash string        `json:"pay_tx_hash"` // user who pay for mint address's transaction hash
	MiddleTx  *wire.MsgTx   `json:"middle_tx"`   // which include one mint runes and other UTXO for reveals txs input
	RevealTxs []*wire.MsgTx `json:"reveal_txs"`  // all txs which have mint runes
}

type WrapTx struct {
	TxPrevOutputFetcher *txscript.MultiPrevOutFetcher `json:"tx_prev_output_fetcher"`
	WireTx              *wire.MsgTx                   `json:"wire_tx"`
}
