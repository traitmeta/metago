package abi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

const RawABI = `[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "owner",
				"type": "address"
			}
		],
		"name": "List",
		"outputs": [
			{
				"internalType": "address[]",
				"name": "receiver",
				"type": "address[]"
			},
			{
				"internalType": "uint256[]",
				"name": "values",
				"type": "uint256[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "owner",
				"type": "address"
			}
		],
		"name": "Value",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "values",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

func ABI() {
	parsed, err := abi.JSON(strings.NewReader(RawABI))
	if err != nil {
		panic(err)
	}

	{
		address := common.HexToAddress("0x80819B3F30e9D77DE6BE3Df9d6EfaA88261DfF9c")
		uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
		bytes32Ty, _ := abi.NewType("bytes32", "", nil)
		addressTy, _ := abi.NewType("address", "", nil)

		arguments := abi.Arguments{
			{
				Type: addressTy,
			},
			{
				Type: bytes32Ty,
			},
			{
				Type: uint256Ty,
			},
		}

		inputs := crypto.Keccak256(
			common.LeftPadBytes(common.HexToAddress("0000000000000000000000000000000000000000").Bytes(), 32),
			common.LeftPadBytes([]byte{'I', 'D', '1'}, 32),
			common.LeftPadBytes(big.NewInt(42).Bytes(), 32),
		)
		log.Println("inputs ", common.LeftPadBytes(common.HexToAddress("0000000000000000000000000000000000000000").Bytes(), 32), common.LeftPadBytes([]byte{'I', 'D', '1'}, 32), common.LeftPadBytes(big.NewInt(42).Bytes(), 32), hexutil.Encode(inputs))

		bytes1, _ := arguments.Pack(
			common.HexToAddress("0x0000000000000000000000000000000000000000"),
			[32]byte{'I', 'D', '1'},
			big.NewInt(42),
		)
		fmt.Println("pack bytes", bytes1)

		hash2 := sha3.NewLegacyKeccak256()
		hash2.Write(bytes1)
		buf2 := hash2.Sum(nil)

		fmt.Println("buf2:", hexutil.Encode(buf2))

		// Value 参数编码
		valueInput, err := parsed.Pack("Value", address)
		if err != nil {
			panic(err)
		}
		fmt.Println("should value", valueInput)
		// Value 参数解码
		// var addrwant common.Address
		// if err := parsed.Methods["Value"].Inputs.Unpack(&addrwant, valueInput[4:]); err != nil {
		// 	panic(err)
		// }
		// fmt.Println("should equals", addrwant == address)

		// Value 返回值解码
		var balance *big.Int
		var returns = common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000005f5e100")
		if err := parsed.UnpackIntoInterface(&balance, "Value", returns); err != nil {
			panic(err)
		}
		fmt.Println("Value 返回值", balance)
	}

	// List 返回值编码
	{
		// 注意：字段名称需要与 ABI 编码的定义的一致
		// 比如，这里 ABI 编码返回值第一个为 receiver 那么转化为 Go 就是首字母大写的 Receiver
		var res struct {
			Receiver []common.Address // 返回值名称
			Values   []*big.Int       // 返回值名称
		}

		// {"Receiver":["0x80819b3f30e9d77de6be3df9d6efaa88261dff9c"],"Values":[10]}
		raw := common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000100000000000000000000000080819b3f30e9d77de6be3df9d6efaa88261dff9c0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000a")
		if err := parsed.UnpackIntoInterface(&res, "List", raw); err != nil {
			panic(err)
		}
		_ = json.NewEncoder(os.Stdout).Encode(&res)
	}
}

func encodePacked(input ...[]byte) []byte {
	return bytes.Join(input, nil)
}

func encodeBytesString(v string) []byte {
	decoded, err := hex.DecodeString(v)
	if err != nil {
		panic(err)
	}
	return decoded
}

func encodeUint256(v string) []byte {
	bn := new(big.Int)
	bn.SetString(v, 10)
	return math.U256Bytes(bn)
}

func encodeUint256Array(arr []string) []byte {
	var res [][]byte
	for _, v := range arr {
		b := encodeUint256(v)
		res = append(res, b)
	}

	return bytes.Join(res, nil)
}

func encodeHash() {
	// bytes32 stateRoots
	stateRoots := "3a53dc4890241dbe03e486e785761577d1c369548f6b09aa38017828dcdf5c27"
	// uint256[2] calldata signatures
	signatures := []string{
		"3402053321874964899321528271743396700217057178612185975187363512030360053932",
		"1235124644010117237054094970590473241953434069965207718920579820322861537001",
	}
	// uint256 feeReceivers,
	feeReceivers := "0"
	// bytes calldata txss
	txss := "000000000000000100010000"

	result := encodePacked(
		encodeBytesString(stateRoots),
		encodeUint256Array(signatures),
		encodeUint256(feeReceivers),
		encodeBytesString(txss),
	)

	got := hex.EncodeToString(result)
	want := "3a53dc4890241dbe03e486e785761577d1c369548f6b09aa38017828dcdf5c2707857e73108d077c5b7ef89540d6493f70d940f1763a9d34c9d98418a39d28ac02bb0e4743a7d0586711ee3dd6311256579ab7abcd53c9c76f040bfde4d6d6e90000000000000000000000000000000000000000000000000000000000000000000000000000000100010000"
	fmt.Println(got == want) // true
}
