package usage

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethmath "github.com/ethereum/go-ethereum/common/math"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/pkg/errors"

	"github.com/traitmeta/metago/eth/eip712"
)

func SignData(priv string, order string) (signature string, err error) {
	var td apitypes.TypedData
	err = json.Unmarshal([]byte(order), &td)
	if err != nil {
		return signature, errors.Wrap(err, "unmarshal order error")
	}

	prv, err := ethcrypto.HexToECDSA(priv)
	if err != nil {
		return signature, errors.Wrap(err, "private key error")
	}

	return eip712.EIP712Sign(prv, &td)
}

func VerifySignature(sig string, order string) (signature string, err error) {
	var td apitypes.TypedData
	err = json.Unmarshal([]byte(order), &td)
	if err != nil {
		return signature, errors.Wrap(err, "unmarshal order error")
	}

	_, address, err := eip712.RecoverEIP712Signature(sig, &td)
	if err != nil {
		return "", err
	}

	return address.Hex(), nil
}

var typedDomain = apitypes.TypedDataDomain{
	Name:              "Seaport",
	Version:           "1.5",
	ChainId:           ethmath.NewHexOrDecimal256(int64(2494104990)),
	VerifyingContract: "0xfd74465415f10afb8c941195821a9d5eec63df2c",
}

var typedDomainFieldTypes = []apitypes.Type{
	{Name: "name", Type: "string"},
	{Name: "version", Type: "string"},
	{Name: "chainId", Type: "uint256"},
	{Name: "verifyingContract", Type: "address"},
}

// SignForBatchClaim array in types
func SignForBatchClaim(ctx context.Context, signPrv *ecdsa.PrivateKey, sender string, tokenIds []interface{}) (signature string, err error) {
	primaryType := "BatchClaim"
	fieldTypes := apitypes.Types{
		"EIP712Domain": typedDomainFieldTypes,
		primaryType: []apitypes.Type{
			{Name: "sender", Type: "address"},
			{Name: "tokenIds", Type: "uint256[]"},
		},
	}
	signPayload := &apitypes.TypedData{
		Domain:      typedDomain,
		Types:       fieldTypes,
		PrimaryType: primaryType,
		Message: apitypes.TypedDataMessage{
			"sender":   sender,
			"tokenIds": tokenIds,
		},
	}

	return eip712.EIP712Sign(signPrv, signPayload)
}

func recover(signature, dataHash string) (recoveredAddr string, err error) {
	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		return recoveredAddr, errors.Wrap(err, "should be a hex string with 0x prefix")
	}
	sigBytes[64] -= 27

	signDataDigest, err := hexutil.Decode(dataHash)
	if err != nil {
		return recoveredAddr, errors.Wrap(err, "decode dataHash failed")
	}

	pubKeyRaw, err := ethcrypto.Ecrecover(signDataDigest, sigBytes)
	if err != nil {
		return recoveredAddr, errors.Wrap(err, "recover signer failed")
	}

	pubKey, err := ethcrypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		return recoveredAddr, errors.Wrap(err, "invalid public key bytes")
	}

	recoveredAddr = ethcrypto.PubkeyToAddress(*pubKey).String()
	return
}
