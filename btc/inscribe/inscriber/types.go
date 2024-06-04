package inscriber

import (
	"github.com/btcsuite/btcd/btcec/v2"
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

// NOTE : len(middleTx) + len(RevealTxs) = MintReq.Count
type MintTxs struct {
	PayTxHash string    `json:"pay_tx_hash"` // user who pay for mint address's transaction hash
	MiddleTx  *WrapTx   `json:"middle_tx"`   // which include one mint runes and other UTXO for reveals txs input
	RevealTxs []*WrapTx `json:"reveal_txs"`  // all txs which have mint runes
}

type WrapTx struct {
	PrevOutput          *wire.TxOut                   `json:"prev_output"`
	TxPrevOutputFetcher *txscript.MultiPrevOutFetcher `json:"tx_prev_output_fetcher"`
	WireTx              *wire.MsgTx                   `json:"wire_tx"`
	MinerFee            int64                         `json:"miner_fee"`
}

type WalletInfo struct {
	PrivateKey      *btcec.PrivateKey `json:"private_key"`
	Address         btcutil.Address   `json:"address"`
	PkScript        []byte            `json:"pk_script"`
	RecoveryPKofWIF string            `json:"recovery_pk_of_wif"`
}
