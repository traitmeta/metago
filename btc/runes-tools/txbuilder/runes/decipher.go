package runes

import (
	"github.com/btcsuite/btcd/txscript"
)

func Decipher(script []byte) ([]byte, error) {
	instructions := txscript.MakeScriptTokenizer(0, script)
	instruction := instructions.Next()
	if !instruction {
		return nil, nil
	}

	opcode := instructions.Opcode()
	if opcode != txscript.OP_RETURN {
		return nil, nil
	}

	if !instructions.Next() {
		return nil, nil
	}

	if instructions.Opcode() != txscript.OP_13 {
		return nil, nil
	}

	return PushBytesInstructions(&instructions), nil
}

func PushBytesInstructions(instructions *txscript.ScriptTokenizer) []byte {
	payload := []byte{}
	for {
		instruction := instructions.Next()
		if !instruction {
			return payload
		}

		opcode := instructions.Opcode()
		if opcode > txscript.OP_DATA_1 && opcode <= txscript.OP_DATA_75 {
			payload = append(payload, instructions.Data()...)
		} else {
			return payload
		}
	}
}
