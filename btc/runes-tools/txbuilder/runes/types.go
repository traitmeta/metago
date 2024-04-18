package runes

import "math/big"

type RuneStone struct {
	Edicts  []Edicts    `json:"edicts"`
	Etching *Etching    `json:"etching"`
	Mint    interface{} `json:"mint"` // blockHeight:txIndex
	Pointer int         `json:"pointer"`
}

// NewMintRuneStone runeID = blockHeight:txIndex
func NewMintRuneStone(runeID string) *RuneStone {
	return &RuneStone{
		Edicts:  []Edicts{},
		Etching: nil,
		Mint:    runeID,
		Pointer: 1,
	}
}

type Edicts struct {
	Id     string   `json:"id"` // blockHeight:txIndex
	Amount *big.Int `json:"amount"`
	Output int      `json:"output"`
}

type Etching struct {
	Divisibility interface{} `json:"divisibility"` // 默认是nil, 有数据类型为int
	Premine      *big.Int    `json:"premine"`
	Rune         interface{} `json:"rune"`    // runes name，默认是nil, 有数据类型为string
	Spacers      interface{} `json:"spacers"` // 默认是nil, 有数据类型为int
	Symbol       interface{} `json:"symbol"`  // 默认是nil, 有数据需要时string，单个字符
	Terms        *Terms      `json:"terms"`
	Turbo        bool        `json:"turbo"`
}

type Terms struct {
	Amount *big.Int   `json:"amount"`
	Cap    int        `json:"cap"`
	Height []*big.Int `json:"height"`
	Offset []*big.Int `json:"offset"`
}

type EtchRequest struct {
	FeeRate     int64
	RuneID      string
	Destination string
}
