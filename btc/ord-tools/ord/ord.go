package ord

import (
	"encoding/hex"
	"fmt"
)

type InscriptionData struct {
	ContentType string
	Body        []byte
	Destination string

	// extra data
	MetaProtocol string
}

type InscriptionRequest struct {
	// a local signature is required for committing the commit tx.
	// Currently, CommitTxPrivateKeyList[i] sign CommitTxOutPointList[i]
	CommitFeeRate  int64 // note: 给矿工的手续费率，在构建commit tx时使用
	FeeRate        int64 // note: 交易费率，相当于gas price
	DataList       []InscriptionData
	RevealOutValue int64
	PrivateKey     string
}

func GetSatBytes(decimalNum int64) ([]byte, error) {
	// Convert the decimal number to hexadecimal
	hexStr := fmt.Sprintf("%x", decimalNum)
	if len(hexStr)%2 != 0 {
		// Pad with '0' if the length is odd
		hexStr = "0" + hexStr
	}

	// Convert the hexadecimal string to byte array
	hexBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	// Reverse the byte array by swapping elements in-place
	for i, j := 0, len(hexBytes)-1; i < j; i, j = i+1, j-1 {
		hexBytes[i], hexBytes[j] = hexBytes[j], hexBytes[i]
	}

	return hexBytes, nil
}
