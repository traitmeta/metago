package transaction

import (
	"bytes"
	"encoding/hex"
	"regexp"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
)

const (
	DefaultTxVersion   = 2
	DefaultSequenceNum = 0xfffffffd

	MaxStandardTxWeight = 4000000 / 10
	WitnessScaleFactor  = 4
)

// GetTxVirtualSize computes the virtual size of a given transaction. A
// transaction's virtual size is based off its weight, creating a discount for
// any witness data it contains, proportional to the current
// blockchain.WitnessScaleFactor value.
func GetTxVirtualSize(tx *btcutil.Tx) int64 {
	// vSize := (weight(tx) + 3) / 4
	//       := (((baseSize * 3) + totalSize) + 3) / 4
	// We add 3 here as a way to compute the ceiling of the prior arithmetic
	// to 4. The division by 4 creates a discount for wit witness data.
	return (GetTransactionWeight(tx) + (WitnessScaleFactor - 1)) / WitnessScaleFactor
}

// GetTransactionWeight computes the value of the weight metric for a given
// transaction. Currently the weight metric is simply the sum of the
// transactions's serialized size without any witness data scaled
// proportionally by the WitnessScaleFactor, and the transaction's serialized
// size including any witness data.
func GetTransactionWeight(tx *btcutil.Tx) int64 {
	msgTx := tx.MsgTx()

	baseSize := msgTx.SerializeSizeStripped()
	totalSize := msgTx.SerializeSize()

	// (baseSize * 3) + totalSize
	return int64((baseSize * (WitnessScaleFactor - 1)) + totalSize)
}

func GetTxHex(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func GetPsbtFromString(psbtStr string) (*psbt.Packet, error) {
	isHex := IsHexString(psbtStr)
	var bs []byte
	var err error
	if isHex {
		bs, err = hex.DecodeString(psbtStr)
	} else {
		bs = []byte(psbtStr)
	}
	p, err := psbt.NewFromRawBytes(bytes.NewReader(bs), !isHex)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func IsHexString(s string) bool {
	if len(s) <= 1 {
		return false
	}
	if s[:2] != "0x" {
		s = "0x" + s
	}
	res, err := regexp.MatchString("^0x[0-9a-fA-F]+$", s)
	if err != nil {
		return false
	}
	return res
}
