package runes

import (
	"fmt"
	"strconv"
	"strings"
)

type RuneId struct {
	block uint64
	tx    uint64
}

func NewRuneId(block uint64, tx uint64) *RuneId {
	if block == 0 && tx > 0 {
		return nil
	}
	return &RuneId{block, tx}
}

func (r *RuneId) Delta(next *RuneId) (uint64, uint64) {
	var blockDiff uint64
	var txDiff uint64
	if next.block >= r.block {
		blockDiff = next.block - r.block
	}
	if next.tx >= r.tx {
		txDiff = next.tx - r.tx
	}

	return blockDiff, txDiff
}

func (r *RuneId) Next(block uint64, tx uint64) *RuneId {
	var newBlock uint64
	var newTx uint64

	if block == 0 && r.block+block >= r.block {
		newBlock = r.block + block
	} else {
		return nil
	}

	if block == 0 && r.tx+tx >= r.tx {
		newTx = r.tx + tx
	} else {
		newTx = tx
	}

	return &RuneId{newBlock, newTx}
}

func (r *RuneId) String() string {
	return fmt.Sprintf("%d:%d", r.block, r.tx)
}

func ParseRuneId(input string) (*RuneId, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format")
	}

	block, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid block: %v", err)
	}

	tx, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid tx: %v", err)
	}

	return &RuneId{block, tx}, nil
}
