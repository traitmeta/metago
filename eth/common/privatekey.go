package common

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func PrivateFromHex(hexPriv string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hexPriv)
}

func AddressFromPrivHex(hexPriv string) (string, error) {
	priv, err := PrivateFromHex(hexPriv)
	if err != nil {
		return "", err
	}

	return GetPublicAddressFromPrivateKey(priv)
}

func GetPublicAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) (string, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("publicKey is not of type *ecdsa.PublicKey")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex(), nil
}

func SignMsgByPrivHex(hexPriv string, msg string) (string, error) {
	priv, err := PrivateFromHex(hexPriv)
	if err != nil {
		return "", err
	}

	data := accounts.TextHash([]byte(msg))
	sig, err := crypto.Sign(data, priv)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(sig), nil
}

func VerifySignature(publicAddress, sigHex string, msg []byte) bool {
	sig, err := hexutil.Decode(sigHex)
	if err != nil {
		return false
	}
	if len(sig) <= crypto.RecoveryIDOffset {
		return false
	}

	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	msg = accounts.TextHash(msg)

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	// use `strings.EqualFold(publicAddress, recoveredAddr.Hex())` if case-insensitive
	return bool(publicAddress == recoveredAddr.Hex())
}
