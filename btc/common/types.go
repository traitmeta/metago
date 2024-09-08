package common

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/traitmeta/metago/pkg/btcapi"
)

type BlockchainClient struct {
	RpcClient    *rpcclient.Client
	BtcApiClient btcapi.Client
}
