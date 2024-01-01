package mempool

import (
	"github.com/btcsuite/btcd/rpcclient"
)

type Client struct {
	client *rpcclient.Client
}

func NewClient(connCfg rpcclient.ConnConfig) *Client {
	client, err := rpcclient.New(&connCfg, nil)
	if err != nil {
		panic("connect client failed")
	}
	return &Client{
		client: client,
	}
}

// func (c *Client) GetMempoolInfo() (map[string]*btcjson.GetRawMempoolVerboseResult, error) {
// 	RawMempoolVerbose
// 	(func(mp *TxPool) RawMempoolVerbose)()
// }
