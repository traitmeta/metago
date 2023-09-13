package config

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

var EthRpcClient *ethclient.Client

func SetupEthClient() {
	var err error
	EthRpcClient, err = NewEthRpcClient()
	if err != nil {
		log.Panic("config.NewEthRpcClient error : ", err)
	}
}

func NewEthRpcClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial(BlockChain.RpcUrl)
	if err != nil {
		return nil, err
	}
	return client, nil
}
