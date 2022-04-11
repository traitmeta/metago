package eip712

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func Verifying(storedChallengeHash []byte, userAddress common.Address, incomingMetamaskSignature string) error {
	// Fetch the previously stored challenge hash from your database
	// var storedChallengeHash []byte = ...
	// Fetch the ETH address whose signature you will be verifying
	// var userAddress common.Address{} = ...
	// Decode the hex-encoded signature from metamask.
	signature, _ := hex.DecodeString(incomingMetamaskSignature)

	if len(signature) != 65 {
		return fmt.Errorf("invalid signature length: %d", len(signature))
	}

	if signature[64] != 27 && signature[64] != 28 {
		return fmt.Errorf("invalid recovery id: %d", signature[64])
	}
	signature[64] -= 27

	pubKeyRaw, err := crypto.Ecrecover(storedChallengeHash, signature)
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
