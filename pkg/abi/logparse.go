package abi

import (
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/traitmeta/metago/config"
	"github.com/traitmeta/metago/pkg/abi/erc1155"
	"github.com/traitmeta/metago/pkg/abi/erc20"
)

func GetErc20Metadata(contractAddress string) (totalSupply *big.Int, name, symbol string, decimals uint8, err error) {
	tokenAddress := common.HexToAddress(contractAddress)
	instance, err := erc20.NewErc20(tokenAddress, config.EthRpcClient)
	if err != nil {
		log.Fatal(err)
	}

	totalSupply, err = instance.TotalSupply(nil)
	if err != nil {
		return
	}

	name, err = instance.Name(nil)
	if err != nil {
		return
	}

	decimals, err = instance.Decimals(nil)
	if err != nil {
		return
	}

	symbol, err = instance.Symbol(nil)
	return
}

func GetErc20TotalSupply(contractAddress string) (*big.Int, error) {
	tokenAddress := common.HexToAddress(contractAddress)
	instance, err := erc20.NewErc20(tokenAddress, config.EthRpcClient)
	if err != nil {
		log.Fatal(err)
	}

	return instance.TotalSupply(nil)
}

func ParseErc20TransferLog(data []byte) (*big.Int, error) {
	contractAbi, err := abi.JSON(strings.NewReader(string(erc20.Erc20ABI)))
	if err != nil {
		return nil, err
	}

	transferEvent := struct {
		From  common.Address `json:"from"`
		To    common.Address `json:"to"`
		Value *big.Int       `json:"value"`
	}{}

	err = contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", data)
	if err != nil {
		return nil, err
	}

	return transferEvent.Value, nil
}

func ParseErc1155BatchTransferLog(data []byte) ([]*big.Int, []*big.Int, error) {
	contractAbi, err := abi.JSON(strings.NewReader(string(erc1155.Erc1155ABI)))
	if err != nil {
		return nil, nil, err
	}

	transferEvent := struct {
		Operator common.Address
		From     common.Address
		To       common.Address
		Ids      []*big.Int
		Values   []*big.Int
	}{}

	err = contractAbi.UnpackIntoInterface(&transferEvent, "TransferBatch", data)
	if err != nil {
		return nil, nil, err
	}

	return transferEvent.Ids, transferEvent.Values, nil
}

func ParseErc1155SignleTransferLog(data []byte) (*big.Int, *big.Int, error) {
	contractAbi, err := abi.JSON(strings.NewReader(string(erc1155.Erc1155ABI)))
	if err != nil {
		return nil, nil, err
	}

	transferEvent := struct {
		Operator common.Address
		From     common.Address
		To       common.Address
		Id       *big.Int
		Value    *big.Int
	}{}

	err = contractAbi.UnpackIntoInterface(&transferEvent, "TransferSingle", data)
	if err != nil {
		return nil, nil, err
	}

	return transferEvent.Id, transferEvent.Value, nil
}
