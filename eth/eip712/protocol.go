package eip712

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// Replace this with the address of the user's wallet
const (
	walletAddress = "0x61e0499cF10d341A5E45FA9c211aF3Ba9A2b50ef"
	salt          = "some-random-string-or-hash-here"
)

var timestamp = strconv.FormatInt(time.Now().Unix(), 10)

func Eip712Hash() (*common.Hash, error) {
	// Generate a random nonce to include in our challenge
	nonceBytes := make([]byte, 32)
	n, err := rand.Read(nonceBytes)
	if n != 32 {
		return nil, errors.New("nonce: n != 64 (bytes)")
	} else if err != nil {
		return nil, err
	}
	nonce := hex.EncodeToString(nonceBytes)

	signerData := apitypes.TypedData{
		Types: apitypes.Types{
			"Challenge": []apitypes.Type{
				{Name: "address", Type: "address"},
				{Name: "nonce", Type: "string"},
				{Name: "timestamp", Type: "string"},
			},
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "version", Type: "string"},
				{Name: "salt", Type: "string"},
			},
		},
		PrimaryType: "Challenge",
		Domain: apitypes.TypedDataDomain{
			Name:              "ETHChallenger",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(1),
			VerifyingContract: "",
			Salt:              salt,
		},
		Message: apitypes.TypedDataMessage{
			"timestamp": timestamp,
			"address":   walletAddress,
			"nonce":     nonce,
		},
	}

	typedDataHash, _ := signerData.HashStruct(signerData.PrimaryType, signerData.Message)
	domainSeparator, _ := signerData.HashStruct("EIP712Domain", signerData.Domain.Map())

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	challengeHash := crypto.Keccak256Hash(rawData)
	return &challengeHash, nil
}
