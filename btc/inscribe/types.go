package ord

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"

	"github.com/traitmeta/metago/pkg/btcapi"
)

type InscriptionData struct {
	ContentType string
	Body        []byte
	Destination string

	// extra data
	MetaProtocol string
}

type InscriptionRequest struct {
	// a local signature is required for committing the commit tx.
	// Currently, CommitTxPrivateKeyList[i] sign CommitTxOutPointList[i]
	CommitFeeRate  int64 // note: 给矿工的手续费率，在构建commit tx时使用
	FeeRate        int64 // note: 交易费率，相当于gas price
	DataList       []InscriptionData
	RevealOutValue int64
	PrivateKey     string
}

type inscriptionTxCtxData struct {
	privateKey              *btcec.PrivateKey
	inscriptionScript       []byte
	commitTxAddress         btcutil.Address
	commitTxAddressPkScript []byte
	controlBlockWitness     []byte
	recoveryPrivateKeyWIF   string
	revealTxPrevOutput      *wire.TxOut
	middleTxPrevOutput      *wire.TxOut
}

type InscribeTxPrevOutput struct {
	RevealTxPrevOutput *wire.TxOut
	MiddleTxPrevOutput *wire.TxOut
}

type InscribeTxBase struct {
	FeeRate        int64
	RevealOutValue int64
}

type InscribeTxsPreview struct {
	InscribeTxBase
	Destination string
	Previews    []InscribeTxPreview
}

type InscribeTxPreview struct {
	InscribeTxPrevOutput

	Destination              string
	InscriptionWitnessScript []byte
	CommitTxPkScript         []byte
	CommitTxAddress          btcutil.Address
	PrivateKey               *btcec.PrivateKey
	ControlBlockWitness      []byte
	RecoveryPrivateKeyWIF    string
}

type InscribeTxDetail struct {
	InscribeTxPrevOutput
	Tx                       *wire.MsgTx
	Destination              string
	InscriptionWitnessScript []byte
	CommitTxPkScript         []byte
	CommitTxAddress          btcutil.Address
	PrivateKey               *btcec.PrivateKey
	ControlBlockWitness      []byte
	RecoveryPrivateKeyWIF    string
}

type BlockchainClient struct {
	RpcClient    *rpcclient.Client
	BtcApiClient btcapi.Client
}
