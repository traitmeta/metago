package rlp

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type Entity struct {
	WalletAddress common.Address
	AccountNonce  uint64
}

func EncodeToBytes(item interface{}) ([]byte, error) {
	b, err := rlp.EncodeToBytes(item)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("RLP encode output: \n", common.Bytes2Hex(b))
	return b, nil
}

func ContractAddress(input []byte) string {
	hash := crypto.Keccak256Hash(input)
	return common.BytesToAddress(hash.Bytes()).Hex()
}
