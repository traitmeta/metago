package ord

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func GetRevealOutValue(revealOutValue int64) int64 {
	result := defaultRevealOutValue // note: 铭文所在 UTXO 的 sats 数量
	if revealOutValue >= minRevealOutValue {
		result = revealOutValue
	}

	return result
}

func EncodeTxToHex(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func DecodeTxFromHex(txHex string) (*wire.MsgTx, error) {
	decodeString, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	if err := tx.Deserialize(bytes.NewReader(decodeString)); err != nil {
		return nil, err
	}

	return tx, nil
}

func VerifySign(signature []byte, hash []byte, pubKey *btcec.PublicKey) (bool, error) {
	signatureStruct, err := schnorr.ParseSignature(signature)
	if err != nil {
		return false, err
	}

	return signatureStruct.Verify(hash, pubKey), nil
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

func BuildInscriptionWitness(datas []InscriptionData, privateKey *btcec.PrivateKey, revealOutValue int64) ([]byte, error) {
	totalInscriptionScript := make([]byte, 0)

	for i, data := range datas {
		// 构建reveal脚本（铭文内容包含在这个脚本中），这里先创建一个脚本构建者
		builder := txscript.NewScriptBuilder()

		if i == 0 {
			builder.
				AddData(schnorr.SerializePubKey(privateKey.PubKey())). // 将公钥添加到脚本中
				AddOp(txscript.OP_CHECKSIG)                            // 添加一个OP_CHECKSIG操作码（用于验证交易的签名）：在这里，它会验证提供的签名是否和脚本中嵌入的公钥匹配，从而确定交易是否又公钥的持有者授权
		}

		builder.
			AddOp(txscript.OP_FALSE). //现在将铭文内容以ordinal的方式添加到脚本中，
			AddOp(txscript.OP_IF).    // 1. 添加一个OP_FALSE操作码和一个OP_IF操作码，（这里是为了构建一个条件语句块，OP_IF操作码会将栈顶的元素弹出并检查它是否为0，如果是0，则执行条件语句块中的代码，否则跳过条件语句块中的代码）
			AddData([]byte("ord"))    // 2. 这里是铭文的协议ID：ord

		if data.MetaProtocol != "" {
			builder.
				AddOp(txscript.OP_DATA_1).
				AddOp(txscript.OP_DATA_7).
				AddData([]byte(data.MetaProtocol))
		}

		// Two OP_DATA_1 should be OP_1. However, in the following link, it's not set as OP_1:
		// https://github.com/casey/ord/blob/0.5.1/src/inscription.rs#L17
		// Therefore, we use two OP_DATA_1 to maintain consistency with ord.
		builder.
			AddOp(txscript.OP_DATA_1).        // 3. 添加一个OP_DATA_1操作码: 表示OP_PUSHBYTES_1操作码，它会将一个字节的数据推送到栈中
			AddOp(txscript.OP_DATA_1).        // 4. 添加一个OP_DATA_1操作码: 表示插入的一个字节数据是: 0x01
			AddData([]byte(data.ContentType)) // 5. 添加铭文的内容类型

		if i != 0 { // note: 除了第一个铭文数据，其他的铭文数据都会添加一个Pointer信息
			// modified: revealOutValue 转 十六进制 ；转 []byte 再倒叙
			satBytes, err := GetSatBytes(revealOutValue * int64(i))
			if err != nil {
				return nil, err
			}

			builder.
				AddOp(txscript.OP_DATA_1).
				AddOp(txscript.OP_DATA_2).
				AddData(satBytes) // warn revealOutValue 的十六进制 倒叙 字节数组
		}

		builder.AddOp(txscript.OP_0) // 6. 添加一个OP_0操作码: 推送一个零长度的字节数组，在ordinal中，这个操作码会被用来标记铭文的开始
		maxChunkSize := 520
		bodySize := len(data.Body)
		for i := 0; i < bodySize; i += maxChunkSize {
			end := i + maxChunkSize
			if end > bodySize {
				end = bodySize
			}
			// to skip txscript.MaxScriptSize 10000
			builder.AddFullData(data.Body[i:end]) // 7. 添加铭文的内容,最大长度为520
		}

		inscriptionScript, err := builder.Script() // note: 铭文内容基本构建完成， 生成铭文脚本
		if err != nil {
			return nil, err
		}

		// to skip txscript.MaxScriptSize 10000
		inscriptionScript = append(inscriptionScript, txscript.OP_ENDIF) // note: 脚本最后添加一个OP_ENDIF操作码，表示条件语句块的结束，也是铭文脚本的结束
		totalInscriptionScript = append(totalInscriptionScript, inscriptionScript...)
	}

	return totalInscriptionScript, nil
}

func SendRawTransaction(client *BlockchainClient, tx *wire.MsgTx) (*chainhash.Hash, error) {
	if client.RpcClient != nil {
		return client.RpcClient.SendRawTransaction(tx, false)
	} else {
		return client.BtcApiClient.BroadcastTx(tx)
	}
}

func GetServiceFee(inscAmount int64) int64 {
	if inscAmount <= 11 {
		return 1000
	} else if inscAmount <= 50 {
		return inscAmount * 95
	} else if inscAmount <= 100 {
		return inscAmount * 90
	} else if inscAmount <= 200 {
		return inscAmount * 80
	} else if inscAmount <= 400 {
		return inscAmount * 70
	} else if inscAmount <= 600 {
		return inscAmount * 60
	} else if inscAmount <= 1000 {
		return inscAmount * 50
	}

	return inscAmount * 50
}

func getServiceFeePkScript(address string, net *chaincfg.Params) (*[]byte, error) {
	// 解析接收地址
	addr, err := btcutil.DecodeAddress(address, net)
	if err != nil {
		return nil, err
	}

	// 创建一个支付到该地址的脚本
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	return &pkScript, nil
}
