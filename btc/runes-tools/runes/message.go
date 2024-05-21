package runes

import (
	"errors"

	"github.com/shopspring/decimal"
)

type Edict struct {
	id     RuneId
	amount decimal.Decimal
	output uint64
}

func fromIntegersSimple1(outputLen uint64, id RuneId, amount decimal.Decimal, output uint64) (*Edict, error) {
	if output > outputLen {
		return nil, errors.New("edict output greater than transaction output count")
	}

	return &Edict{id, amount, output}, nil
}

type Flaw int

const (
	EdictOutput Flaw = iota
	EdictRuneId
	InvalidScript
	Opcode
	SupplyOverflow
	TrailingIntegers
	TruncatedField
	UnrecognizedEvenTag
	UnrecognizedFlag
	Varint
)

func (f Flaw) String() string {
	switch f {
	case EdictOutput:
		return "edict output greater than transaction output count"
	case EdictRuneId:
		return "invalid rune ID in edict"
	case InvalidScript:
		return "invalid script in OP_RETURN"
	case Opcode:
		return "non-pushdata opcode in OP_RETURN"
	case SupplyOverflow:
		return "supply overflows u128"
	case TrailingIntegers:
		return "trailing integers in body"
	case TruncatedField:
		return "field with missing value"
	case UnrecognizedEvenTag:
		return "unrecognized even tag"
	case UnrecognizedFlag:
		return "unrecognized field"
	case Varint:
		return "invalid varint"
	default:
		return ""
	}
}

type Message struct {
	flaw   *Flaw
	edicts []Edict
	fields map[decimal.Decimal][]decimal.Decimal
}

// func fromIntegersSimple(outputLen uint64, payload []decimal.Decimal) Message {
// 	edicts := make([]Edict, 0)
// 	fields := make(map[decimal.Decimal][]decimal.Decimal)
// 	var flaw *Flaw

// 	for i := 0; i < len(payload); i += 2 {
// 		tag := payload[i]

// 		if tag == Body {
// 			var id RuneId
// 			for j := i + 1; j < len(payload); j += 4 {
// 				if j+4 > len(payload) {
// 					flaw = &TrailingIntegers
// 					break
// 				}

// 				next, err := id.next(chunk[0], chunk[1])
// 				if err != nil {
// 					flaw = &EdictRuneId
// 					break
// 				}

// 				edict, err := fromIntegersSimple(outputLen, next, chunk[2], chunk[3])
// 				if err != nil {
// 					flaw = &EdictOutput
// 					break
// 				}

// 				id = next
// 				edicts = append(edicts, edict)
// 			}
// 			break
// 		}

// 		if i+1 >= len(payload) {
// 			flaw = &TruncatedField
// 			break
// 		}

// 		value := payload[i+1]
// 		fields[tag] = append(fields[tag], value)
// 	}

// 	return Message{
// 		flaw:   flaw,
// 		edicts: edicts,
// 		fields: fields,
// 	}
// }
