package tap

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/wire"
)

// Element format : <name>.<pattern>.<field>.element
type Element struct {
	name    string
	pattern string
	field   string
}

func (e *Element) String() string {
	if e.pattern == "" {
		return fmt.Sprintf("%s.%s.element", e.name, e.field)
	}

	return fmt.Sprintf("%s.%s.%s.element", e.name, e.pattern, e.field)
}

func (e *Element) IsValid(elements []Element) bool {
	for _, char := range InvalidChar {
		if strings.ContainsRune(e.name, char) {
			return false
		}
	}

	for _, src := range elements {
		if strings.EqualFold(e.name, src.name) || (strings.EqualFold(e.pattern, e.pattern) && strings.EqualFold(e.field, e.field)) {
			return false
		}
	}

	return true
}

func (e *Element) IsAvailable(block *wire.MsgBlock) bool {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, block.Header.Bits)
	hexBits := hex.EncodeToString(a)
	return strings.Contains(hexBits, e.pattern)
}
