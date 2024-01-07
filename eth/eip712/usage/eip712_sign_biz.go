package usage

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethmath "github.com/ethereum/go-ethereum/common/math"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/pkg/errors"

	"github.com/traitmeta/metago/eth/eip712"
)

func SignForOrder(priv string, order string) (signature string, err error) {
	var td apitypes.TypedData
	err = json.Unmarshal([]byte(order), &td)
	prv, err := ethcrypto.HexToECDSA(priv)
	if err != nil {
		return signature, errors.Wrap(err, "private key error")
	}

	signature, err = eip712.EIP712Sign(prv, &td)
	if err != nil {
		return signature, errors.Wrap(err, "SignForOrder")
	}

	return signature, err
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

func VerifySigForBind(ctx context.Context, signer, tronAddress, sig string) (isValid bool, err error) {
	var myTypedData = `
	{
		"types": {
		  "OrderComponents": [
			{ "name": "offerer", "type": "address" },
			{ "name": "zone", "type": "address" },
			{ "name": "offer", "type": "OfferItem[]" },
			{ "name": "consideration", "type": "ConsiderationItem[]" },
			{ "name": "orderType", "type": "uint8" },
			{ "name": "startTime", "type": "uint256" },
			{ "name": "endTime", "type": "uint256" },
			{ "name": "zoneHash", "type": "bytes32" },
			{ "name": "salt", "type": "uint256" },
			{ "name": "conduitKey", "type": "bytes32" },
			{ "name": "counter", "type": "uint256" }
		  ],
		  "OfferItem": [
			{ "name": "itemType", "type": "uint8" },
			{ "name": "token", "type": "address" },
			{ "name": "identifierOrCriteria", "type": "uint256" },
			{ "name": "startAmount", "type": "uint256" },
			{ "name": "endAmount", "type": "uint256" }
		  ],
		  "ConsiderationItem": [
			{ "name": "itemType", "type": "uint8" },
			{ "name": "token", "type": "address" },
			{ "name": "identifierOrCriteria", "type": "uint256" },
			{ "name": "startAmount", "type": "uint256" },
			{ "name": "endAmount", "type": "uint256" },
			{ "name": "recipient", "type": "address" }
		  ],
		  "EIP712Domain": [
			{ "name": "name", "type": "string" },
			{ "name": "version", "type": "string" },
			{ "name": "chainId", "type": "uint256" },
			{ "name": "verifyingContract", "type": "address" }
		  ]
		},
		"domain": {
		  "name": "Seaport",
		  "version": "1.5",
		  "chainId": "2494104990",
		  "verifyingContract": "0x9afa6139c383e7a3796131a43dcb86caa9178170"
		},
		"primaryType": "OrderComponents",
		"message": {
			"offerer": "TC5vH9Sdp6Ybo6y9QqSmLAoDsz7jxwrfaQ",
			"zone": "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb",
			"offer": [
				{
					"itemType": 2,
					"token": "TMNcHdDMEUAZ4bcXTHEbRwYKkwnDGjaMB3",
					"identifierOrCriteria": "385",
					"startAmount": "1",
					"endAmount": "1"
				}
			],
			"consideration": [
				{
					"itemType": 0,
					"token": "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb",
					"identifierOrCriteria": "0",
					"startAmount": "32340000",
					"endAmount": "32340000",
					"recipient": "TC5vH9Sdp6Ybo6y9QqSmLAoDsz7jxwrfaQ"
				},
				{
					"itemType": 0,
					"token": "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb",
					"identifierOrCriteria": "0",
					"startAmount": "495000",
					"endAmount": "495000",
					"recipient": "TAXH4o5jj1mhB4aDNfZNpUKNTkGrEeEySV"
				},
				{
					"itemType": 0,
					"token": "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb",
					"identifierOrCriteria": "0",
					"startAmount": "165000",
					"endAmount": "165000",
					"recipient": "TRxMn4k6ZpPhHLEqVojaXRzr3PP3p3B9KE"
				}
			],
			"orderType": 0,
			"startTime": "1692347711",
			"endTime": "1694939694",
			"zoneHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
			"salt": "43296667924576948581825440725076177030533144458658729516795853069932738938401",
			"conduitKey": "0x0000000000000000000000000000000000000000000000000000000000000000",
			"totalOriginalConsiderationItems": "3",
			"counter": "0"
			}
		}
	`
	var signPayload apitypes.TypedData
	err = json.Unmarshal([]byte(myTypedData), &signPayload)
	_, realSigner, err := eip712.RecoverEIP712Signature(sig, &signPayload)
	if err != nil {
		errors.Wrap(err, "recover failed")
		return
	}

	if !strings.EqualFold(realSigner.Hex(), signer) {
		err = errors.New(fmt.Sprintf("EIP712 recovered signer=%s not equal with %s", realSigner.Hex(), signer))
		return
	}

	return true, nil
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

// SignFoTest will print every step for 712
func SignFoTest(priv string) (signature string, err error) {
	var myTypedData = `
	{
		"types": {
		  "OrderComponents": [
			{ "name": "offerer", "type": "address" },
			{ "name": "zone", "type": "address" },
			{ "name": "offer", "type": "OfferItem[]" },
			{ "name": "consideration", "type": "ConsiderationItem[]" },
			{ "name": "orderType", "type": "uint8" },
			{ "name": "startTime", "type": "uint256" },
			{ "name": "endTime", "type": "uint256" },
			{ "name": "zoneHash", "type": "bytes32" },
			{ "name": "salt", "type": "uint256" },
			{ "name": "conduitKey", "type": "bytes32" },
			{ "name": "counter", "type": "uint256" }
		  ],
		  "OfferItem": [
			{ "name": "itemType", "type": "uint8" },
			{ "name": "token", "type": "address" },
			{ "name": "identifierOrCriteria", "type": "uint256" },
			{ "name": "startAmount", "type": "uint256" },
			{ "name": "endAmount", "type": "uint256" }
		  ],
		  "ConsiderationItem": [
			{ "name": "itemType", "type": "uint8" },
			{ "name": "token", "type": "address" },
			{ "name": "identifierOrCriteria", "type": "uint256" },
			{ "name": "startAmount", "type": "uint256" },
			{ "name": "endAmount", "type": "uint256" },
			{ "name": "recipient", "type": "address" }
		  ],
		  "EIP712Domain": [
			{ "name": "name", "type": "string" },
			{ "name": "version", "type": "string" },
			{ "name": "chainId", "type": "uint256" },
			{ "name": "verifyingContract", "type": "address" }
		  ]
		},
		"domain": {
		  "name": "Seaport",
		  "version": "1.5",
		  "chainId": "2494104990",
		  "verifyingContract": "0x9afa6139c383e7a3796131a43dcb86caa9178170"
		},
		"primaryType": "OrderComponents",
		"message": {
		  "offerer": "0x33a4f229bd34ea7783302c99ffd6e26324bd2789",
		  "zone": "0x0000000000000000000000000000000000000000",
		  "offer": [
			{
			  "itemType": "2",
			  "token": "0xfb8d0f7d033268b76a5077a5462ae711d1e48a0b",
			  "identifierOrCriteria": "4",
			  "startAmount": "1",
			  "endAmount": "1"
			}
		  ],
		  "consideration": [
			{
			  "itemType": "0",
			  "token": "0x0000000000000000000000000000000000000000",
			  "identifierOrCriteria": "0",
			  "startAmount": "20000000",
			  "endAmount": "20000000",
			  "recipient": "0x33a4f229bd34ea7783302c99ffd6e26324bd2789"
			},
			{
			  "itemType": "0",
			  "token": "0x0000000000000000000000000000000000000000",
			  "identifierOrCriteria": "0",
			  "startAmount": "1000000",
			  "endAmount": "1000000",
			  "recipient": "0x95b467b0d33c34d5bc2ab3fb005cf9aca4033f00"
			}
		  ],
		  "orderType": "0",
		  "startTime": "1591123844",
		  "endTime": "1791123844",
		  "zoneHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		  "salt": "24446860302761739304752683030156737591518664810215442929816108075358245610000",
		  "conduitKey": "0x0000000000000000000000000000000000000000000000000000000000000000",
		  "counter": "0"
		}
	  }
	`
	var td apitypes.TypedData
	err = json.Unmarshal([]byte(myTypedData), &td)
	prv, err := ethcrypto.HexToECDSA(priv)
	if err != nil {
		panic("init signer private key failed")
	}

	var EIP712DomainField = "EIP712Domain"
	domainSep, err := td.HashStruct(EIP712DomainField, td.Domain.Map())
	fmt.Println(domainSep.String())

	dataHash, err := td.HashStruct(td.PrimaryType, td.Message)
	fmt.Println(dataHash.String())

	signData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSep), string(dataHash)))
	signDataDigest := ethcrypto.Keccak256Hash(signData)
	fmt.Println(signDataDigest.String())

	return eip712.EIP712Sign(prv, &td)
}

func recover() (recoveredAddr string, err error) {
	sigBytes, err := hexutil.Decode("0xc5f7a27fb56690c5ca607b2ddc5efd58ca8b0290dab78a72a478e5d32e7facb562284309264b9cc0a09473e20ba5271587e3aa0bf9bc00b2296c560eb7b6035b1b")
	if err != nil {
		return recoveredAddr, errors.Wrap(err, "should be a hex string with 0x prefix")
	}
	sigBytes[64] -= 27

	signDataDigest, err := hexutil.Decode("0x846c2aa6277c50980556cccc77d2c9bcde1258b00228ce062da733268802fa01")
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

func SignForZone(priv string) (signature string, err error) {
	var myTypedData = `
	{
		"types": {
		  "SignedOrder": [
			{ "name": "fulfiller", "type": "address" },
			{ "name": "expiration", "type": "uint64" },
			{ "name": "orderHash", "type": "bytes32" },
			{ "name": "context", "type": "bytes" }
		  ],
		  "EIP712Domain": [
			{ "name": "name", "type": "string" },
			{ "name": "version", "type": "string" },
			{ "name": "chainId", "type": "uint256" },
			{ "name": "verifyingContract", "type": "address" }
		  ]
		},
		"domain": {
		  "name": "SignedZone",
		  "version": "1.0",
		  "chainId": "1",
		  "verifyingContract": "0xd182d0a388f4923c478395dfb4ea889e55013967"
		},
		"primaryType": "SignedOrder",
		"message": {
		  "fulfiller": "0x4b20993bc481177ec7e8f571cecae8a9e22c02db",
		  "expiration": "1751641513",
		  "orderHash": "0x846c2aa6277c50980556cccc77d2c9bcde1258b00228ce062da733268802fa01",
		  "context": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}
	  }
	`
	var td apitypes.TypedData
	err = json.Unmarshal([]byte(myTypedData), &td)
	prv, err := ethcrypto.HexToECDSA(priv)
	if err != nil {
		panic("init signer private key failed")
	}
	var EIP712DomainField = "EIP712Domain"
	domainSep, err := td.HashStruct(EIP712DomainField, td.Domain.Map())
	fmt.Println(domainSep.String())

	dataHash, err := td.HashStruct(td.PrimaryType, td.Message)
	fmt.Println(dataHash.String())

	signData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSep), string(dataHash)))
	signDataDigest := ethcrypto.Keccak256Hash(signData)
	fmt.Println(signDataDigest.String())

	return eip712.EIP712Sign(prv, &td)
}
