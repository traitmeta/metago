package txbuilder

import "encoding/hex"

type PrevInfo struct {
	PrevTxId   string
	PrevScript []byte
	PreAmount  int64
}

func BuilderPrevInfo(txId, script string, amount int64) PrevInfo {
	scriptBytes, err := hex.DecodeString(script)
	if err != nil {
		return PrevInfo{}
	}

	return PrevInfo{
		PrevTxId:   txId,
		PrevScript: scriptBytes,
		PreAmount:  amount,
	}
}
