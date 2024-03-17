package ord

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	ProtocolId = "ord"
)

type ParsedEnvelope struct {
	Input   uint32
	Offset  uint32
	Payload Inscription
	Pushnum bool
	Stutter bool
}

type Inscription struct {
	Body                  []byte
	ContentEncoding       []byte
	ContentType           []byte
	Delegate              []byte
	DuplicateField        bool
	IncompleteField       bool
	Metadata              []byte
	MetaProtocol          []byte
	Parent                []byte
	Pointer               []byte
	UnRecognizedEvenField bool
}

func ParsedEnvelopeFromRaw(data Envelope) {

}

// Envelope
// content_type, with a tag of 1, whose value is the MIME type of the body.
// pointer, with a tag of 2, see pointer docs.
// parent, with a tag of 3, see provenance.
// TODO metadata, with a tag of 5, see metadata.
// TODO metaprotocol, with a tag of 7, whose value is the metaprotocol identifier.
// TODO content_encoding, with a tag of 9, whose value is the encoding of the body.
// TODO delegate, with a tag of 11, see delegate.
type Envelope struct {
	Input       uint32
	Offset      uint32
	TypeDataMap map[int][]byte
	Payload     [][]byte
	Pushnum     bool
	Stutter     bool
}

func (e *Envelope) GetContent() []byte {
	if v, ok := e.TypeDataMap[0]; ok {
		return v
	}

	return nil
}

func (e *Envelope) GetContentType() string {
	if v, ok := e.TypeDataMap[1]; ok {
		return string(v)
	}

	return ""
}

// GetPointer Pointer
func (e *Envelope) GetPointer() uint64 {
	if v, ok := e.TypeDataMap[2]; ok {
		return binary.LittleEndian.Uint64(v)
	}

	return 0
}

// GetProvenance is parent little-endian OP_PUSH 3 TXID INDEX
// TXID = 32-byte INDEX = 4-byte
func (e *Envelope) GetProvenance() string {
	v, ok := e.TypeDataMap[3]
	if !ok {
		return ""
	}

	return covLittleEndianToOrdIdStr(v)
}

func (e *Envelope) GetContentEncoding() string {
	if v, ok := e.TypeDataMap[9]; ok {
		return string(v)
	}

	return ""
}

// TODO
//func (e *Envelope) GetMetadata() string {
//	if v, ok := e.TypeDataMap[5]; ok {
//		err := cbor.Unmarshal(v, &atomicalToken)
//		return string(v)
//	}
//
//	return ""
//}

func (e *Envelope) GetDelegate() string {
	v, ok := e.TypeDataMap[11]
	if !ok {
		return ""
	}

	return covLittleEndianToOrdIdStr(v)
}

func covLittleEndianToOrdIdStr(v []byte) string {
	bigEndian := make([]byte, 32)
	for i := 0; i < 32; i++ {
		bigEndian[i] = v[32-i-1]
	}

	txId := hex.EncodeToString(bigEndian)
	var index uint64 = 0
	if len(v) > 32 {
		index = binary.LittleEndian.Uint64(v[32:])
	}

	return fmt.Sprintf("%si%d", txId, index)
}

func FromTransaction(transaction *wire.MsgTx) []Envelope {
	envelopes := make([]Envelope, 0)
	for i, input := range transaction.TxIn {
		if len(input.Witness) != 3 {
			continue
		}

		if len(input.Witness[1]) == 0 {
			continue
		}

		if inputEnvelopes, err := FromTapScript(input.Witness[1], i); err == nil {
			envelopes = append(envelopes, inputEnvelopes...)
		}
	}

	return envelopes
}

func FromTapScript(tapScript []byte, input int) ([]Envelope, error) {
	envelopes := make([]Envelope, 0)

	instructions := txscript.MakeScriptTokenizer(0, tapScript)
	stuttered := false
	for {
		instruction := instructions.Next()
		if !instruction {
			break
		}
		opcode := instructions.Opcode()
		if bytes.Equal([]byte{opcode}, []byte{txscript.OP_0}) {
			if !instructions.Next() {
				break
			}
			stutter, envelope := FromInstructions(&instructions, input, len(envelopes), stuttered)
			if envelope != nil {
				envelopes = append(envelopes, *envelope)
			} else {
				stuttered = stutter
			}
		}
	}

	return envelopes, nil
}

func Accept(instructions *txscript.ScriptTokenizer, instruction []byte) bool {
	opCode := instructions.Opcode()
	if bytes.Equal([]byte{opCode}, instruction) && instructions.Next() {
		return true
	}

	return false
}

func FromInstructions(instructions *txscript.ScriptTokenizer, input int, offset int, stutter bool) (bool, *Envelope) {
	if !Accept(instructions, []byte{txscript.OP_IF}) {
		stutter := Accept(instructions, []byte{})
		if stutter {
			return stutter, nil
		}
		return stutter, nil
	}
	opcode := instructions.Opcode()
	if !bytes.Equal([]byte{opcode}, []byte{txscript.OP_DATA_3}) || !bytes.Equal(instructions.Data(), []byte(ProtocolId)) {
		stutter := Accept(instructions, []byte{})
		if stutter {
			return stutter, nil
		}
		return stutter, nil
	}

	pushnum := false
	payload := [][]byte{}
	typeDataMap := make(map[int][]byte)
	currentType := 0
	for {
		instruction := instructions.Next()
		if !instruction {
			return false, nil
		}

		opcode := instructions.Opcode()
		switch opcode {
		case txscript.OP_ENDIF:
			return false, &Envelope{
				Input:       uint32(input),
				Offset:      uint32(offset),
				TypeDataMap: typeDataMap,
				Payload:     payload,
				Pushnum:     pushnum,
				Stutter:     stutter,
			}
		case txscript.OP_1NEGATE:
			pushnum = true
			payload = append(payload, []byte{0x81})
		case txscript.OP_1:
			pushnum = true
			payload = append(payload, []byte{0x01})
		case txscript.OP_2:
			pushnum = true
			payload = append(payload, []byte{0x02})
		case txscript.OP_3:
			pushnum = true
			payload = append(payload, []byte{0x03})
		case txscript.OP_4:
			pushnum = true
			payload = append(payload, []byte{0x04})
		case txscript.OP_5:
			pushnum = true
			payload = append(payload, []byte{0x05})
		case txscript.OP_6:
			pushnum = true
			payload = append(payload, []byte{0x06})
		case txscript.OP_7:
			pushnum = true
			payload = append(payload, []byte{0x07})
		case txscript.OP_8:
			pushnum = true
			payload = append(payload, []byte{0x08})
		case txscript.OP_9:
			pushnum = true
			payload = append(payload, []byte{0x09})
		case txscript.OP_10:
			pushnum = true
			payload = append(payload, []byte{0x0a})
		case txscript.OP_11:
			pushnum = true
			payload = append(payload, []byte{0x0b})
		case txscript.OP_12:
			pushnum = true
			payload = append(payload, []byte{0x0c})
		case txscript.OP_13:
			pushnum = true
			payload = append(payload, []byte{0x0d})
		case txscript.OP_14:
			pushnum = true
			payload = append(payload, []byte{0x0e})
		case txscript.OP_15:
			pushnum = true
			payload = append(payload, []byte{0x0f})
		case txscript.OP_16:
			pushnum = true
			payload = append(payload, []byte{0x10})
		case txscript.OP_PUSHDATA1, txscript.OP_PUSHDATA2, txscript.OP_PUSHDATA4:
			if _, ok := typeDataMap[currentType]; !ok {
				typeDataMap[currentType] = []byte{}
			} else {
				typeDataMap[currentType] = append(typeDataMap[currentType], instructions.Data()...)
			}
			payload = append(payload, instructions.Data())
		case txscript.OP_DATA_1:
			data := instructions.Data()
			currentType = int(data[0])
			if _, ok := typeDataMap[currentType]; !ok {
				typeDataMap[currentType] = []byte{}
			} else {
				typeDataMap[currentType] = append(typeDataMap[currentType], data...)
			}
			payload = append(payload, instructions.Data())
		//	The next opcode bytes is data to be pushed onto the stack
		case txscript.OP_0:
			currentType = 0
			if _, ok := typeDataMap[currentType]; !ok {
				typeDataMap[currentType] = []byte{}
			}
			payload = append(payload, instructions.Data())
		default:
			if opcode > txscript.OP_DATA_1 && opcode <= txscript.OP_DATA_75 {
				data := instructions.Data()
				if _, ok := typeDataMap[currentType]; !ok {
					typeDataMap[currentType] = []byte{}
				} else {
					typeDataMap[currentType] = append(typeDataMap[currentType], data...)
				}
				payload = append(payload, instructions.Data())
			} else {
				return false, nil
			}
		}
	}
}
