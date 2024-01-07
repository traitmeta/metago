package eip712

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func Verifying(dataHash []byte, userAddress common.Address, signature string) error {
	signatureBytes, _ := hex.DecodeString(signature)

	if len(signature) != 65 {
		return fmt.Errorf("invalid signature length: %d", len(signature))
	}

	if signature[64] != 27 && signature[64] != 28 {
		return fmt.Errorf("invalid recovery id: %d", signature[64])
	}
	signatureBytes[64] -= 27

	pubKeyRaw, err := crypto.Ecrecover(dataHash, signatureBytes)
	if err != nil {
		return fmt.Errorf("invalid signature: %s", err.Error())
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		return err
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	if !bytes.Equal(userAddress.Bytes(), recoveredAddr.Bytes()) {
		return errors.New("addresses do not match")
	}

	return nil
}
