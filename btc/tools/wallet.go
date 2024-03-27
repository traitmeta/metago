package tools

import (
	"encoding/hex"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
)

// AddressType represents the various address types waddrmgr is currently able
// to generate, and maintain.
//
// NOTE: These MUST be stable as they're used for scope address schema
// recognition within the database.
type AddressType uint8

const (
	// PubKeyHash is a regular p2pkh address.
	PubKeyHash AddressType = iota

	// Script reprints a raw script address.
	Script

	// RawPubKey is just raw public key to be used within scripts, This
	// type indicates that a scoped manager with this address type
	// shouldn't be consulted during historical rescans.
	RawPubKey

	// NestedWitnessPubKey represents a p2wkh output nested within a p2sh
	// output. Using this address type, the wallet can receive funds from
	// other wallet's which don't yet recognize the new segwit standard
	// output types. Receiving funds to this address maintains the
	// scalability, and malleability fixes due to segwit in a backwards
	// compatible manner.
	NestedWitnessPubKey

	// WitnessPubKey represents a p2wkh (pay-to-witness-key-hash) address
	// type.
	WitnessPubKey

	// WitnessScript represents a p2wsh (pay-to-witness-script-hash) address
	// type.
	WitnessScript

	// TaprootPubKey represents a p2tr (pay-to-taproot) address type that
	// uses BIP-0086 (for the derivation path and for calculating the tap
	// root hash/tweak).
	TaprootPubKey

	// TaprootScript represents a p2tr (pay-to-taproot) address type that
	// commits to a script and not just a single key.
	TaprootScript
)

const (
	// witnessVersionV0 is the SegWit v0 witness version used for p2wpkh and
	// p2wsh outputs and addresses.
	witnessVersionV0 byte = 0x00

	// witnessVersionV1 is the SegWit v1 witness version used for p2tr
	// outputs and addresses.
	witnessVersionV1 byte = 0x01
)

var ErrNotSupportBtcAddressType = errors.New("not support btc address type")

func VerifyPublicKeyBVMAddress(pubKey, addr string, net *chaincfg.Params) bool {
	addressTypeList := []AddressType{PubKeyHash, WitnessPubKey, TaprootPubKey}
	for _, addressType := range addressTypeList {
		addrFromPBK, err := BtcPubKeyToBtcAddress(pubKey, addressType, net)
		if err != nil {
			return false
		}

		bvmAddr := BtcAddressToBvmAddress(addrFromPBK)
		if strings.EqualFold(bvmAddr, addr) {
			return true
		}
	}
	return false
}

// BtcPubKeyToBtcAddressWithType for veda public key to address
func BtcPubKeyToBtcAddressWithType(pubkeyHex string, addressType string, net *chaincfg.Params) (string, error) {
	addrType, err := StrToAddressType(addressType)
	if err != nil {
		return "", err
	}

	return BtcPubKeyToBtcAddress(pubkeyHex, addrType, net)
}

func StrToAddressType(str string) (AddressType, error) {
	switch str {
	case "p2pkh":
		return PubKeyHash, nil
	case "p2wpkh":
		return WitnessPubKey, nil
	case "p2sh":
		return WitnessPubKey, nil
	case "p2tr":
		return TaprootPubKey, nil
	default:
		return PubKeyHash, ErrNotSupportBtcAddressType
	}
}

// BtcPubKeyToBtcAddress for veda public key to address
func BtcPubKeyToBtcAddress(pubkeyHex string, addressType AddressType, net *chaincfg.Params) (string, error) {
	scriptPubKey, err := hex.DecodeString(pubkeyHex)
	if err != nil {
		return "", errors.Wrap(err, "failed on decoding pkscript")
	}

	pubkey, err := btcec.ParsePubKey(scriptPubKey)
	if err != nil {
		return "", errors.Wrap(err, "failed on decoding pkscript")
	}

	addr, err := newAddressFromPubKey(pubkey, btcec.IsCompressedPubKey(scriptPubKey), addressType, net)
	if err != nil {
		return "", errors.Wrap(err, "failed on decoding pkscript")
	}

	return addr.EncodeAddress(), nil
}

// newAddressFromPubKey returns a new address based on the
// passed account, public key, and whether or not the public key should be
// compressed.
func newAddressFromPubKey(pubKey *btcec.PublicKey, compressed bool, addrType AddressType, net *chaincfg.Params) (btcutil.Address, error) {

	// Create a pay-to-pubkey-hash address from the public key.
	var pubKeyHash []byte
	if compressed {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	} else {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	}

	var address btcutil.Address
	var err error

	switch addrType {

	case NestedWitnessPubKey:
		// For this address type we'll generate an address which is
		// backwards compatible to Bitcoin nodes running 0.6.0 onwards, but
		// allows us to take advantage of segwit's scripting improvements,
		// and malleability fixes.

		// First, we'll generate a normal p2wkh address from the pubkey hash.
		witAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, net)
		if err != nil {
			return nil, err
		}

		// Next we'll generate the witness program which can be used as a
		// pkScript to pay to this generated address.
		witnessProgram, err := txscript.PayToAddrScript(witAddr)
		if err != nil {
			return nil, err
		}

		// Finally, we'll use the witness program itself as the pre-image
		// to a p2sh address. In order to spend, we first use the
		// witnessProgram as the sigScript, then present the proper
		// <sig, pubkey> pair as the witness.
		address, err = btcutil.NewAddressScriptHash(witnessProgram, net)
		if err != nil {
			return nil, err
		}

	case PubKeyHash:
		address, err = btcutil.NewAddressPubKeyHash(pubKeyHash, net)
		if err != nil {
			return nil, err
		}

	case WitnessPubKey:
		address, err = btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, net)
		if err != nil {
			return nil, err
		}

	case TaprootPubKey:
		tapKey := txscript.ComputeTaprootKeyNoScript(pubKey)
		address, err = btcutil.NewAddressTaproot(schnorr.SerializePubKey(tapKey), net)
		if err != nil {
			return nil, err
		}
	}

	return address, nil
}

type rlpEntry struct {
	WalletAddress common.Address
	AccountNonce  uint64
}

func BtcAddressToContractAddress(btcAddress string, nonce uint64) string {
	addr := BtcAddressToBvmAddress(btcAddress)
	entry := rlpEntry{
		WalletAddress: common.HexToAddress(addr),
		AccountNonce:  nonce,
	}

	enc, err := rlp.EncodeToBytes(entry)
	if err != nil {
		return err.Error()
	}

	hash := crypto.Keccak256Hash(enc)
	return common.BytesToAddress(hash.Bytes()).Hex()
}

func BtcAddressToBvmAddress(address string) string {
	hash := crypto.Keccak256Hash([]byte(address))
	addr := common.BytesToAddress(hash.Bytes())
	return addr.Hex()
}
