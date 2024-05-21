package runes

import (
	"math/bits"

	"github.com/holiman/uint256"
)

type Tag uint8

const (
	Body        Tag = 0
	Flags       Tag = 2
	Rune        Tag = 4
	Premine     Tag = 6
	Cap         Tag = 8
	Amount      Tag = 10
	HeightStart Tag = 12
	HeightEnd   Tag = 14
	OffsetStart Tag = 16
	OffsetEnd   Tag = 18
	Mint        Tag = 20
	Pointer     Tag = 22
	Cenotaph    Tag = 126

	Divisibility Tag = 1
	Spacers      Tag = 3
	Symbol       Tag = 5
	Nop          Tag = 127
)

// func (t Tag) takeN(N int, fields map[*uint256.Int][]*uint256.Int, with func([]*uint256.Int) T) T {
// 	field, ok := fields[uint256.NewInt(uint64(t))]
// 	if !ok {
// 		return nil
// 	}

// 	values := make([]*uint256.Int, N)
// 	for i := 0; i < N; i++ {
// 		values[i] = field[i]
// 	}

// 	value := with(values)
// 	fields[uint256.NewInt(uint64(t))] = fields[uint256.NewInt(uint64(t))][N:]

// 	if len(fields[uint256.NewInt(uint64(t))]) == 0 {
// 		delete(fields, uint256.NewInt(uint64(t)))
// 	}

// 	return value
// }

// func (t Tag) encodeN(N int, values []uint128, payload *[]byte) {
// 	for _, value := range values {
// 		encodeVarint(uint256.NewInt(uint64(t)), payload)
// 		encodeVarint(value, payload)
// 	}
// }

// func (t Tag) encodeOption(value interface{}, payload *[]byte) {
// 	if value != nil {
// 		t.encodeN(1, []uint128{tagToUint128(value)}, payload)
// 	}
// }

func (t Tag) intoUint128() *uint256.Int {
	return uint256.NewInt(uint64(t))
}

// func (t Tag) equalsInt(value uint128) bool {
// 	return t.intoUint128() == value
// }

func encodeVarint(n uint128, v *[]byte) {
	for n.lo>>7 > 0 {
		*v = append(*v, byte(n.lo&0x7F|0x80))
		n.lo >>= 7
	}
	for n.hi>>7 > 0 {
		*v = append(*v, byte(n.hi&0x7F|0x80))
		n.hi >>= 7
	}

	*v = append(*v, byte(n.hi))
}

type uint128 struct {
	lo uint64
	hi uint64
}

func (u uint128) bits() int {
	if u.hi == 0 {
		return bits.Len64(u.lo)
	}
	return 64 + bits.Len64(u.hi)
}

func makeUint128(lo, hi uint64) uint128 {
	return uint128{lo, hi}
}

func tagToUint128(value Tag) uint128 {
	return uint128{uint64(value), 0}
}
