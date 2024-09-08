package witness

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/traitmeta/metago/btc/ord-tools/ord"
)

type InscriptionRawTx struct {
	TxPrevOutput   *wire.TxOut
	WitnessScript  *RevealWitness
	Size           int64
	Raw            *wire.MsgTx
	RevealOutValue int64
	FeeRate        int64
	PrivateKey     *btcec.PrivateKey
}

type SignInfo struct {
	PrivateKey    *btcec.PrivateKey
	RevealWitness *RevealWitness
	RevealAccount *RevealAccount
}

type RevealWitness struct {
	SignatureWitness    []byte
	InsWitnessScript    []byte
	ControlBlockWitness []byte
}

type RevealAccount struct {
	CommitTxAddress       btcutil.Address
	CommitTxPkScript      []byte
	RecoveryPrivateKeyWIF string
}

func NewInscriptionWitness() *RevealWitness {
	return &RevealWitness{
		SignatureWitness:    make([]byte, 64),
		InsWitnessScript:    nil,
		ControlBlockWitness: make([]byte, 33),
	}
}

func NewInscriptionRawTx() *InscriptionRawTx {
	return &InscriptionRawTx{
		WitnessScript: NewInscriptionWitness(),
	}
}

func (irt *InscriptionRawTx) SetTxPrevOutput(pkScript []byte, prevOutput int64) {
	irt.TxPrevOutput = &wire.TxOut{
		PkScript: pkScript,
		Value:    prevOutput,
	}
}

func (irt *InscriptionRawTx) SetWitnessScript(inscriptionWitnessScript []byte) {
	irt.WitnessScript.InsWitnessScript = inscriptionWitnessScript
}

func (irt *InscriptionRawTx) SetSize(txSize int64) {
	irt.Size = txSize
}

func (irt *InscriptionRawTx) CalcPrevOutput(revealOutValue, feeRate int64) int64 {
	txFee := irt.Size * feeRate
	prevOutput := revealOutValue + txFee
	emptySignature := make([]byte, 64)
	emptyControlBlockWitness := make([]byte, 33)
	witnessSize := (wire.TxWitness{emptySignature, irt.WitnessScript.InsWitnessScript, emptyControlBlockWitness}.SerializeSize() + 2 + 3) / 4
	// 初始化一个空的签名和控制块，计算单个铭文交易，witness部分的额外手续费，并更新totalPrevOutput
	witnessFee := int64(witnessSize) * feeRate
	prevOutput += witnessFee

	return prevOutput
}

func BuildInscriptionWitness(datas []ord.InscriptionData, privateKey *btcec.PrivateKey, revealOutValue int64) ([]byte, error) {
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
			satBytes, err := ord.GetSatBytes(revealOutValue * int64(i))
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
