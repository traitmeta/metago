package eip712

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/pkg/errors"

	"github.com/traitmeta/metago/eth/common"
)

// EIP712Sign eth EIP712 签名
func EIP712Sign(prv *ecdsa.PrivateKey, data *apitypes.TypedData) (string, error) {
	var signature string
	dataHash, err := data.HashStruct(data.PrimaryType, data.Message)
	if err != nil {
		fmt.Println(dataHash.String())
		return signature, errors.Wrap(err, "EIP712Sign calculate data hash falied")
	}

	domainSep, err := data.HashStruct(common.EIP712DomainField, data.Domain.Map())
	if err != nil {
		fmt.Println(domainSep.String())
		return signature, errors.Wrap(err, "EIP712Sign calculate domain hash falied")
	}

	signData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSep), string(dataHash)))
	signDataDigest := crypto.Keccak256Hash(signData)
	fmt.Println(signDataDigest.String())

	sig, err := crypto.Sign(signDataDigest.Bytes(), prv)
	if err != nil {
		return signature, errors.Wrap(err, "EIP712Sign sign falied")
	}

	sig[64] += 27
	signature = hexutil.Encode(sig)
	return signature, err
}

// EIP712SignatureRecover eth EIP712 恢复 signer 地址
func RecoverEIP712Signature(sig string, data *apitypes.TypedData) (dataHash hexutil.Bytes, recoveredAddr ethcommon.Address, err error) {
	dataHash, err = data.HashStruct(data.PrimaryType, data.Message)
	if err != nil {
		err = errors.Wrap(err, "RecoverEIP712Signature calculate order hash falied")
		return
	}

	domainSep, err := data.HashStruct(common.EIP712DomainField, data.Domain.Map())
	if err != nil {
		panic(err)
	}

	signData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSep), string(dataHash)))
	signDataDigest := crypto.Keccak256Hash(signData)
	sigBytes, err := hexutil.Decode(sig)
	if err != nil {
		err = errors.Wrap(err, "RecoverEIP712Signature without 0x prefix")
		return
	}

	sigBytes[64] -= 27
	pubKeyRaw, err := crypto.Ecrecover(signDataDigest.Bytes(), sigBytes)
	if err != nil {
		err = errors.Wrap(err, "RecoverEIP712Signature recover signer failed")
		return
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		err = errors.Wrap(err, "RecoverEIP712Signature invalid public key bytes")
		return
	}

	recoveredAddr = crypto.PubkeyToAddress(*pubKey)
	return
}

func RecoverRawMsgSignature(sig, msg string) (signerAddr string, err error) {
	msgHash := crypto.Keccak256Hash([]byte(msg))
	bytes1 := msgHash.Bytes()
	bytes2 := []byte(common.ETH_MESSAGE_HEADER)
	toDigest := bytes.Join([][]byte{bytes2, bytes1}, []byte(""))
	digest := crypto.Keccak256Hash(toDigest)
	digestByte := digest.Bytes()

	sigBytes, err := hexutil.Decode(sig)
	if err != nil {
		err = errors.Wrap(err, "RecoverRawMsgSignature bad signature")
		return
	}

	sigBytes[64] -= 27
	pubKeyRaw, err := crypto.Ecrecover(digestByte, sigBytes)
	if err != nil {
		err = errors.Wrap(err, "RecoverRawMsgSignature Ecrecover public failed")
		return
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		err = errors.Wrap(err, "RecoverRawMsgSignature UnmarshalPubkey failed")
		return
	}

	signerAddr = crypto.PubkeyToAddress(*pubKey).Hex()
	return signerAddr, nil
}
