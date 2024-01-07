package common

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

type RSVStr struct {
	R string
	S string
	V int
}

type RSVBytes struct {
	R []byte
	S []byte
	V int
}

func SigatureToRSVStr(sig string) (rs *RSVStr, err error) {
	rsv, err := SigatureToRSVBytes(sig)
	if err != nil {
		return
	}
	rs = &RSVStr{
		R: hexutil.Encode(rsv.R),
		S: hexutil.Encode(rsv.S),
		V: rsv.V,
	}
	return
}

func SigatureToRSVBytes(sig string) (rs *RSVBytes, err error) {
	if Has0xPrefix(sig) {
		sig = sig[2:]
	}

	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return nil, errors.Wrap(err, "SigatureToRSVBytes without 0x prefix")
	}

	if len(sigBytes) != 65 {
		return nil, errors.Wrap(err, "SigatureToRSVBytes sigBytes length should be 132")
	}

	v := int(sigBytes[64])
	if v < 27 || v > 28 {
		return nil, errors.Wrap(err, fmt.Sprintf("SigatureToRSVBytes version should be %d or %d", 27, 28))

	}

	rs = &RSVBytes{
		R: sigBytes[:32],
		S: sigBytes[32:64],
		V: v,
	}
	return
}
