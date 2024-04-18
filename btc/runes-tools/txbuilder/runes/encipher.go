package runes

import (
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func CreateMintRuneStoneOutput(runeId string) (*wire.TxOut, error) {
	encipherScript, err := Encipher(runeId)
	if err != nil {
		return nil, err
	}

	return wire.NewTxOut(0, encipherScript), nil
}

func Encipher(runesId string) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(txscript.OP_DATA_13)
	mint, err := encipherMint(runesId)
	if err != nil {
		return nil, err
	}

	return builder.AddData(mint).Script()
}

func encipherMint(runesId string) ([]byte, error) {
	var v []byte
	parts := strings.Split(runesId, ":")
	blockIdxBytes, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	v = append(v, 20)
	v = append(v, encodeToSlice(blockIdxBytes)...)
	txIdxBytes, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	v = append(v, 20)
	v = append(v, encodeToSlice(txIdxBytes)...)

	v = append(v, 22)
	v = append(v, 1)

	return v, nil
}

func encodeToSlice(n uint64) []byte {
	v := []byte{}
	for n >= 128 {
		v = append(v, byte(n)|0b10000000)
		n >>= 7
	}

	v = append(v, byte(n))
	return v
}
