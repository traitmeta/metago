package ord

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"sync"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
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
}

type inscriptionTxCtxData struct {
	privateKey              *btcec.PrivateKey
	inscriptionScript       []byte
	commitTxAddress         btcutil.Address
	commitTxAddressPkScript []byte
	controlBlockWitness     []byte
	recoveryPrivateKeyWIF   string
	revealTxPrevOutput      *wire.TxOut
	middleTxPrevOutput      *wire.TxOut
}

type blockchainClient struct {
	rpcClient    *rpcclient.Client
	btcApiClient btcapi.BTCAPIClient
}

type InscribeTool interface {
	GetPayAddrAndFee(request *InscriptionRequest) (payAddr, payAddrPK string, inscFee, serviceFee, minerFee int64, err error)
	Inscribe(commitTxHash string, actualMiddlePrevOutputFee int64, payAddrPK string, request *InscriptionRequest) (ctxTxData *CtxTxData, err error)
}

type InscriptionTool struct {
	net                   *chaincfg.Params
	client                *blockchainClient
	serviceFeeReceiveAddr string // note: 用于存储平台手续费的地址

	txCtxDataList             []*inscriptionTxCtxData
	revealTxPrevOutputFetcher *txscript.MultiPrevOutFetcher // note: 用于获取reveal tx的输入
	revealTx                  []*wire.MsgTx                 // note: reveal tx
	commitTx                  *wire.MsgTx                   // note: commit tx
	middleTxPrevOutputFetcher *txscript.MultiPrevOutFetcher // note: 用于获取middle tx的输入
	middleTx                  *wire.MsgTx                   // note: 连接commit tx和（除了第一笔）reveal tx的中间tx： 包含reveal tx1 + commit tx2... + 手续费 + 找零
}

const (
	defaultSequenceNum = wire.MaxTxInSequenceNum - 10

	MaxStandardTxWeight = blockchain.MaxBlockWeight / 10
)

func NewInscribeTool(net *chaincfg.Params, rpcclient *rpcclient.Client, serviceFeeReceiveAddr string) (InscribeTool, error) {
	if serviceFeeReceiveAddr == "" {
		return nil, errors.New("service fee receive address is empty")
	}

	tool := &InscriptionTool{
		net: net,
		client: &blockchainClient{
			rpcClient: rpcclient,
		},
		serviceFeeReceiveAddr:     serviceFeeReceiveAddr,
		revealTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		middleTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
	}
	return tool, nil
}

func NewInscribeToolWithBtcApiClient(net *chaincfg.Params, rpcClient *rpcclient.Client, serviceFeeReceiveAddr string) (InscribeTool, error) {
	if serviceFeeReceiveAddr == "" {
		return nil, errors.New("service fee receive address is empty")
	}

	tool := &InscriptionTool{
		net: net,
		client: &blockchainClient{
			rpcClient: rpcClient,
		},
		serviceFeeReceiveAddr:     serviceFeeReceiveAddr,
		revealTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		middleTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
	}
	return tool, nil
}

func (tool *InscriptionTool) GetPayAddrAndFee(request *InscriptionRequest) (payAddr, payAddrPK string, inscFee, serviceFee, minerFee int64, err error) {
	tool.txCtxDataList = make([]*inscriptionTxCtxData, len(request.DataList)) // not1e: 初始化，铭文列表
	tool.revealTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil)     // note: 初始化，reveal tx的输入
	tool.middleTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil)     // note: 初始化，middle tx的输入

	revealOutValue := GetRevealOutValue(request)
	tool.txCtxDataList = make([]*inscriptionTxCtxData, len(request.DataList)) // note: 初始化，铭文列表
	destinations := make([]string, len(request.DataList))                     // note: 初始化，铭文的接收地址

	for i := 0; i < len(request.DataList); i++ {
		privateKey, err := btcec.NewPrivateKey() // note: 创建一个密钥对，用来构建reveal tx
		if err != nil {
			return "", "", 0, 0, 0, errors.Wrap(err, "create private key error")
		}
		if i == 0 { // warn: 保存构建第一个铭文的私钥
			privKeyBytes := privateKey.Serialize()                      // 将私钥转换为字节数组
			payAddrPK = base64.StdEncoding.EncodeToString(privKeyBytes) // 使用Base64编码将字节数组转换为字符串
		}
		txCtxData, err := createInscriptionTxCtxData(tool.net, request.DataList[i], privateKey) // note: 创建commit交易及包含铭文信息的Taproot脚本信息

		if err != nil {
			return "", "", 0, 0, 0, errors.Wrap(err, "create inscription tx ctx data error")
		}
		tool.txCtxDataList[i] = txCtxData
		destinations[i] = request.DataList[i].Destination
	}
	totalRevealPrevOutput, err := tool.buildEmptyRevealTx(destinations, revealOutValue, request.FeeRate)
	if err != nil {
		return "", "", 0, 0, 0, errors.Wrap(err, "build empty reveal tx error")
	}

	totalMiddlePrevOutput, serviceFee, minerFee, err := tool.buildEmptyMiddleTx(totalRevealPrevOutput, destinations[0], revealOutValue, request.FeeRate, int64(len(request.DataList)))
	if err != nil {
		return "", "", 0, 0, 0, errors.Wrap(err, "build empty middle tx error")
	}

	return tool.txCtxDataList[0].commitTxAddress.String(), payAddrPK, totalMiddlePrevOutput, serviceFee, minerFee, nil
}

func (tool *InscriptionTool) Inscribe(commitTxHash string, actualMiddlePrevOutputFee int64, payAddrPK string, request *InscriptionRequest) (ctxTxData *CtxTxData, err error) {
	tool.txCtxDataList = make([]*inscriptionTxCtxData, len(request.DataList)) // note: 初始化，铭文列表
	tool.revealTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil)     // note: 初始化，reveal tx的输入
	tool.middleTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil)     // note: 初始化，middle tx的输入

	revealOutValue := GetRevealOutValue(request)
	tool.txCtxDataList = make([]*inscriptionTxCtxData, len(request.DataList)) // note: 初始化，铭文列表
	destinations := make([]string, len(request.DataList))                     // note: 初始化，铭文的接收地址
	for i := 0; i < len(request.DataList); i++ {
		var privateKey *btcec.PrivateKey
		var err error
		if i == 0 { // warn: 构建第一个铭文的私钥，用之前生成的
			decodedPrivKey, err := base64.StdEncoding.DecodeString(payAddrPK)
			if err != nil {
				return nil, errors.Wrap(err, "decode private key error")
			}

			// 从字节数组还原私钥
			privateKey, _ = btcec.PrivKeyFromBytes(decodedPrivKey)
		} else {
			privateKey, err = btcec.NewPrivateKey() // note: 创建一个密钥对，用来构建reveal tx
			if err != nil {
				return nil, errors.Wrap(err, "create private key error")
			}
		}
		txCtxData, err := createInscriptionTxCtxData(tool.net, request.DataList[i], privateKey) // note: 创建commit交易及包含铭文信息的Taproot脚本信息
		if err != nil {
			return nil, errors.Wrap(err, "create inscription tx ctx data error")
		}
		tool.txCtxDataList[i] = txCtxData
		destinations[i] = request.DataList[i].Destination
	}
	totalRevealPrevOutput, err := tool.buildEmptyRevealTx(destinations, revealOutValue, request.FeeRate)
	if err != nil {
		return nil, errors.Wrap(err, "build empty reveal tx error")
	}
	fmt.Println("totalRevealPrevOutput, ", totalRevealPrevOutput)
	totalMiddlePrevOutput, _, _, err := tool.buildEmptyMiddleTx(totalRevealPrevOutput, destinations[0], revealOutValue, request.FeeRate, int64(len(request.DataList)))
	if err != nil {
		return nil, errors.Wrap(err, "build empty middle tx error")
	}
	log.Println("totalMiddlePrevOutput, ", totalMiddlePrevOutput)

	if totalMiddlePrevOutput > actualMiddlePrevOutputFee {
		return nil, errors.New("actualMiddlePrevOutputFee is not enough")
	}

	err = tool.completeMiddleTx(commitTxHash, actualMiddlePrevOutputFee)
	if err != nil {
		return nil, errors.Wrap(err, "complete middle tx error")
	}

	err = tool.completeRevealTx()
	if err != nil {
		return nil, errors.Wrap(err, "complete reveal tx error")
	}

	// ====================log====================
	{
		recoveryKeyWIFList := tool.GetRecoveryKeyWIFList()
		for i, recoveryKeyWIF := range recoveryKeyWIFList {
			log.Printf("recoveryKeyWIF %d %s \n", i, recoveryKeyWIF)
		}

		middleTxHex, err := tool.GetMiddleTxHex()
		if err != nil {
			log.Fatalf("get middle tx hex err, %v", err)
		}
		log.Printf("middleTxHex %s \n", middleTxHex)

		revealTxHexList, err := tool.GetRevealTxHexList()
		if err != nil {
			log.Fatalf("get reveal tx hex err, %v", err)
		}

		for i, revealTxHex := range revealTxHexList {
			log.Printf("revealTxHex %d %s \n", i, revealTxHex)
		}
	}
	// ====================log====================
	ctxTxData = &CtxTxData{}
	ctxTxData.PrevCommitTxHash = commitTxHash

	ctxTxData.MiddleTx = tool.middleTx
	middleTxHash, err := tool.sendRawTransaction(tool.middleTx)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("send middle tx error: tx_hash %s", tool.middleTx.TxHash().String()))
	}

	ctxTxData.MiddleTxHash = middleTxHash.String()
	ctxTxData.MiddleInscriptionId = fmt.Sprintf("%si0", middleTxHash)

	// 返回所有reveal tx的结构数据
	revealTxData, err := tool.saveRevealTx()
	ctxTxData.RevealTxData = revealTxData

	// warn: 最多只发送前23笔reveal交易， 因为之后的交易不会成功发送， 会报错：{"code":-26,"message":"too-long-mempool-chain, too many descendants for tx... [limit: 25]"}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	minTx := min(len(tool.revealTx), 23)
	for i := 0; i < minTx; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_revealTxHash, err := tool.sendRawTransaction(tool.revealTx[i])
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				return
			}
			ctxTxData.RevealTxData[i].IsSend = true
			ctxTxData.RevealTxData[i].RevealTxHash = _revealTxHash.String()
			ctxTxData.RevealTxData[i].InscriptionId = fmt.Sprintf("%si0", _revealTxHash)
		}(i)
	}
	wg.Wait()

	if firstErr != nil {
		return ctxTxData, firstErr
	}

	return ctxTxData, nil
}

func createInscriptionTxCtxData(net *chaincfg.Params, data InscriptionData, privateKey *btcec.PrivateKey) (*inscriptionTxCtxData, error) {
	inscriptionBuilder := txscript.NewScriptBuilder(). // 构建reveal脚本（铭文内容包含在这个脚本中），这里先创建一个脚本构建者
								AddData(schnorr.SerializePubKey(privateKey.PubKey())). // 将公钥添加到脚本中
								AddOp(txscript.OP_CHECKSIG).                           // 添加一个OP_CHECKSIG操作码（用于验证交易的签名）：在这里，它会验证提供的签名是否和脚本中嵌入的公钥匹配，从而确定交易是否又公钥的持有者授权
								AddOp(txscript.OP_FALSE).                              //现在将铭文内容以ordinal的方式添加到脚本中，
								AddOp(txscript.OP_IF).                                 // 1. 添加一个OP_FALSE操作码和一个OP_IF操作码，（这里是为了构建一个条件语句块，OP_IF操作码会将栈顶的元素弹出并检查它是否为0，如果是0，则执行条件语句块中的代码，否则跳过条件语句块中的代码）
								AddData([]byte("ord"))                                 // 2. 这里是铭文的协议ID：ord

	if data.MetaProtocol != "" {
		inscriptionBuilder.AddOp(txscript.OP_DATA_1).
			AddOp(txscript.OP_DATA_7).
			AddData([]byte(data.MetaProtocol))
	}

	// Two OP_DATA_1 should be OP_1. However, in the following link, it's not set as OP_1:
	// https://github.com/casey/ord/blob/0.5.1/src/inscription.rs#L17
	// Therefore, we use two OP_DATA_1 to maintain consistency with ord.
	inscriptionBuilder.AddOp(txscript.OP_DATA_1). // 3. 添加一个OP_DATA_1操作码: 表示OP_PUSHBYTES_1操作码，它会将一个字节的数据推送到栈中
							AddOp(txscript.OP_DATA_1).         // 4. 添加一个OP_DATA_1操作码: 表示插入的一个字节数据是: 0x01
							AddData([]byte(data.ContentType)). // 5. 添加铭文的内容类型
							AddOp(txscript.OP_0)               // 6. 添加一个OP_0操作码: 推送一个零长度的字节数组，在ordinal中，这个操作码会被用来标记铭文的开始
	maxChunkSize := 520
	bodySize := len(data.Body)
	for i := 0; i < bodySize; i += maxChunkSize {
		end := i + maxChunkSize
		if end > bodySize {
			end = bodySize
		}
		// to skip txscript.MaxScriptSize 10000
		inscriptionBuilder.AddFullData(data.Body[i:end]) // 7. 添加铭文的内容,最大长度为520
	}
	inscriptionScript, err := inscriptionBuilder.Script() // note: 铭文内容基本构建完成， 生成铭文脚本
	if err != nil {
		return nil, errors.Wrap(err, "create inscription script error")
	}
	// to skip txscript.MaxScriptSize 10000
	inscriptionScript = append(inscriptionScript, txscript.OP_ENDIF) // note: 脚本最后添加一个OP_ENDIF操作码，表示条件语句块的结束，也是铭文脚本的结束

	leafNode := txscript.NewBaseTapLeaf(inscriptionScript) // note: 创建一个新的tapscript叶子节点，将刚才构造的铭文脚本添加到叶子节点中
	proof := &txscript.TapscriptProof{
		TapLeaf:  leafNode,
		RootNode: leafNode,
	}

	// note: 利用前面生成的证明对象和公钥生成Control block
	controlBlock := proof.ToControlBlock(privateKey.PubKey())
	controlBlockWitness, err := controlBlock.ToBytes()
	//fmt.Println("controlBlockWitness: ", hex.EncodeToString(controlBlockWitness))
	if err != nil {
		return nil, errors.Wrap(err, "control block to bytes error")
	}

	// note: 生成最终的 Taproot 地址（commit tx 的输出地址）和 Pay-to-Taproot(P2TR) 地址的脚本。
	tapHash := proof.RootNode.TapHash()
	commitTxAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(privateKey.PubKey(), tapHash[:])), net)
	log.Println("commitTxAddress: ", commitTxAddress)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address error")
	}
	commitTxAddressPkScript, err := txscript.PayToAddrScript(commitTxAddress)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address pk script error")
	}

	recoveryPrivateKeyWIF, err := btcutil.NewWIF(txscript.TweakTaprootPrivKey(*privateKey, tapHash[:]), net, true)
	if err != nil {
		return nil, errors.Wrap(err, "create recovery private key wif error")
	}

	return &inscriptionTxCtxData{
		privateKey:              privateKey,
		inscriptionScript:       inscriptionScript,
		commitTxAddress:         commitTxAddress,
		commitTxAddressPkScript: commitTxAddressPkScript,
		controlBlockWitness:     controlBlockWitness,
		recoveryPrivateKeyWIF:   recoveryPrivateKeyWIF.String(),
	}, nil
}

func (tool *InscriptionTool) buildEmptyRevealTx(destination []string, revealOutValue, feeRate int64) (totalPrevOutput int64, err error) {
	var revealTx []*wire.MsgTx // note: 初始化， 存储reveal tx
	totalPrevOutput = 0        // note: 初始化，存储所有reveal tx的总费用
	total := len(tool.txCtxDataList)
	addTxInTxOutIntoRevealTx := func(tx *wire.MsgTx, index int) error {
		in := wire.NewTxIn(&wire.OutPoint{Index: uint32(index)}, nil, nil) // note: 先构造reveal tx的空交易输入
		in.Sequence = defaultSequenceNum
		tx.AddTxIn(in)
		receiver, err := btcutil.DecodeAddress(destination[index], tool.net) // note: 生成铭文的接收地址, 这里是destination[index]的P2TR地址
		if err != nil {
			return errors.Wrap(err, "decode address error")
		}
		scriptPubKey, err := txscript.PayToAddrScript(receiver)
		if err != nil {
			return errors.Wrap(err, "pay to address script error")
		}
		out := wire.NewTxOut(revealOutValue, scriptPubKey) // note: 再构造reveal tx的交易输出，也就是铭文的接收地址
		tx.AddTxOut(out)
		return nil
	}

	revealTx = make([]*wire.MsgTx, total-1) // note: 初始化，创建多个reveal tx

	for i := 1; i < total; i++ { // modified: 循环除了第一个铭文数据，构造除了第一个reveal tx 的信息
		tx := wire.NewMsgTx(wire.TxVersion)    // note: 创建一个新的reveal tx
		err := addTxInTxOutIntoRevealTx(tx, i) // note: 往每个tx添加“空的交易输入”和输出
		if err != nil {
			return 0, errors.Wrap(err, "add tx in tx out into reveal tx error")
		}
		prevOutput := revealOutValue + int64(tx.SerializeSize())*feeRate // note: 计算单个reveal tx的文本费用： 铭文的value + 铭文的交易大小 * feeRate
		{
			emptySignature := make([]byte, 64)
			emptyControlBlockWitness := make([]byte, 33)
			// 初始化一个空的签名和控制块，计算单个铭文交易，witness部分的额外手续费，并更新totalPrevOutput
			fee := (int64(wire.TxWitness{emptySignature, tool.txCtxDataList[i].inscriptionScript, emptyControlBlockWitness}.SerializeSize()+2+3) / 4) * feeRate
			prevOutput += fee // note: 计算单个铭文的总费用： 铭文的value + 铭文的交易大小 * feeRate + 额外手续费

			tool.txCtxDataList[i].revealTxPrevOutput = &wire.TxOut{
				PkScript: tool.txCtxDataList[i].commitTxAddressPkScript,
				Value:    prevOutput,
			}

			totalPrevOutput += prevOutput
		}
		revealTx[i-1] = tx
	}

	tool.revealTx = revealTx
	return totalPrevOutput, nil
}

func (tool *InscriptionTool) getServiceFeePkScript() (*[]byte, error) {
	// 解析接收地址
	addr, err := btcutil.DecodeAddress(tool.serviceFeeReceiveAddr, tool.net)
	if err != nil {
		log.Fatalf("Invalid address: %v", err)
	}

	// 创建一个支付到该地址的脚本
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		log.Fatalf("Failed to create pkScript: %v", err)
	}
	return &pkScript, nil
}

// note: 构造middleTx: 连接commit tx和（除了第一笔）reveal tx的中间tx： 包含reveal tx1 + commit tx2... + 手续费 + 找零
func (tool *InscriptionTool) buildEmptyMiddleTx(totalRevealPrevOutput int64, firstDestination string, revealOutValue, feeRate, inscAmount int64) (totalPrevOutput, serviceFee, minerFee int64, err error) {
	tx := wire.NewMsgTx(wire.TxVersion) // note: 初始化，创建一个新的middle tx

	addTxInTxOutOfRevealTx1IntoMiddleTx := func(tx *wire.MsgTx, index int) error {
		in := wire.NewTxIn(&wire.OutPoint{Index: uint32(index)}, nil, nil) // note: middle tx的空交易输入
		in.Sequence = defaultSequenceNum
		tx.AddTxIn(in)
		receiver, err := btcutil.DecodeAddress(firstDestination, tool.net) // note: 生成铭文的接收地址, 这里是destinations[0]的P2TR地址
		if err != nil {
			return errors.Wrap(err, "decode address error")
		}
		scriptPubKey, err := txscript.PayToAddrScript(receiver)
		if err != nil {
			return errors.Wrap(err, "pay to address script error")
		}
		out := wire.NewTxOut(revealOutValue, scriptPubKey) // note: 再构造reveal tx的交易输出，也就是铭文的接收地址
		tx.AddTxOut(out)
		return nil
	}

	// 输出
	err = addTxInTxOutOfRevealTx1IntoMiddleTx(tx, 0) // note: 添加第一个reveal tx1的输出
	if err != nil {
		return 0, 0, 0, errors.Wrap(err, "add tx in tx out of reveal tx1 into middle tx error")
	}
	// 计算第一个reveal tx1的费用
	prevOutput := revealOutValue // note: 计算单个reveal tx的文本费用： 铭文的value
	{
		emptySignature := make([]byte, 64)
		emptyControlBlockWitness := make([]byte, 33)
		// 初始化一个空的签名和控制块，计算单个铭文交易，witness部分的额外手续费，并更新totalPrevOutput
		fee := (int64(wire.TxWitness{emptySignature, tool.txCtxDataList[0].inscriptionScript, emptyControlBlockWitness}.SerializeSize()+2+3) / 4) * feeRate
		prevOutput += fee // note: 计算单个铭文的总费用： 铭文的value + 见证文本费用
		minerFee += fee   // note: minerFee 加上 见证文本费用
	}
	totalPrevOutput += prevOutput

	// modified: middle tx 的输出再加上reveal tx的preOutput
	for i := range tool.txCtxDataList { // reveal tx的输入 是 commit tx的输出
		if i == 0 {
			continue
		}
		tx.AddTxOut(tool.txCtxDataList[i].revealTxPrevOutput)
		totalPrevOutput += tool.txCtxDataList[i].revealTxPrevOutput.Value
		minerFee += int64(tool.txCtxDataList[i].revealTxPrevOutput.Value) - revealOutValue
	}

	// 在middle tx的末尾，添加一个给平台手续费操作的uxto输出
	serviceFee = getServiceFee(inscAmount)
	servicePkScript, err := tool.getServiceFeePkScript()
	tx.AddTxOut(wire.NewTxOut(serviceFee, *servicePkScript))
	totalPrevOutput += serviceFee

	// 总费用加上交易费用
	txSizeFee := int64(tx.SerializeSize()) * feeRate
	minerFee += txSizeFee
	totalPrevOutput += txSizeFee

	// 设置middleTxPrevOutput:
	tool.txCtxDataList[0].middleTxPrevOutput = &wire.TxOut{
		PkScript: tool.txCtxDataList[0].commitTxAddressPkScript,
		Value:    totalPrevOutput,
	}

	tool.middleTx = tx
	return totalPrevOutput, serviceFee, minerFee, nil
}

func (tool *InscriptionTool) completeMiddleTx(commitTxHash string, actualMiddlePrevOutputFee int64) error {
	newCommitTxHash, err := chainhash.NewHashFromStr(commitTxHash)
	if err != nil {
		// 处理错误
	}
	fmt.Println("newCommitTxHash", newCommitTxHash)
	// note: 完成中间交易的构建过程，主要通过计算签名哈希(hash)、生成签名，并为每个中间交易设置见证数据(witness data)。

	//modified: 更新commit tx的输出金额为用户实际的转账金额
	tool.txCtxDataList[0].middleTxPrevOutput.Value = actualMiddlePrevOutputFee

	tool.middleTxPrevOutputFetcher.AddPrevOut(wire.OutPoint{ // note: 将commit tx添加到middleTxPrevOutputFetcher的输入列表中
		Hash:  *newCommitTxHash,
		Index: uint32(0),
	}, tool.txCtxDataList[0].middleTxPrevOutput)

	tool.middleTx.TxIn[0].PreviousOutPoint.Hash = *newCommitTxHash // note: 更新交易的输入指向相应的commit Tx哈希

	middleTx := tool.middleTx

	// note: 计算签名哈希： 基于当前的reveal tx, 关联的前一个交易输出，以及铭文的脚本，计算签名哈希； 这是准备签名验证数据的关键，保证交易能在网络中被正确验证
	witnessArray, err := txscript.CalcTapscriptSignaturehash(txscript.NewTxSigHashes(middleTx, tool.middleTxPrevOutputFetcher),
		txscript.SigHashDefault, middleTx, 0, tool.middleTxPrevOutputFetcher, txscript.NewBaseTapLeaf(tool.txCtxDataList[0].inscriptionScript))
	if err != nil {
		return errors.Wrap(err, "calc tapscript signaturehash error")
	}

	// note: 使用私钥对签名哈希进行签名
	signature, err := schnorr.Sign(tool.txCtxDataList[0].privateKey, witnessArray)
	if err != nil {
		return errors.Wrap(err, "schnorr sign error")
	}
	witness := wire.TxWitness{signature.Serialize(), tool.txCtxDataList[0].inscriptionScript, tool.txCtxDataList[0].controlBlockWitness}
	tool.middleTx.TxIn[0].Witness = witness

	revealWeight := blockchain.GetTransactionWeight(btcutil.NewTx(middleTx))
	if revealWeight > MaxStandardTxWeight {
		return errors.New(fmt.Sprintf("middle transaction weight greater than %d (MAX_STANDARD_TX_WEIGHT): %d", MaxStandardTxWeight, revealWeight))
	}

	return nil
}

// note: 完成铭文揭示交易（Reveal Transaction）的构建过程，主要通过计算签名哈希(hash)、生成签名，并为每个揭示交易设置见证数据（witness data）。
// note: 此外，它还验证每个揭示交易的大小是否符合比特币网络的标准交易重量限制。
func (tool *InscriptionTool) completeRevealTx() error {
	for i := range tool.txCtxDataList {
		if i == 0 {
			continue
		}
		tool.revealTxPrevOutputFetcher.AddPrevOut(wire.OutPoint{ // note: 将middle tx添加到revealTxPrevOutputFetcher的输入列表中
			Hash:  tool.middleTx.TxHash(),
			Index: uint32(i),
		}, tool.txCtxDataList[i].revealTxPrevOutput)

		//note: 更新交易的输入指向相应的commit Tx哈希
		tool.revealTx[i-1].TxIn[0].PreviousOutPoint.Hash = tool.middleTx.TxHash()
	}

	witnessList := make([]wire.TxWitness, len(tool.txCtxDataList)-1) // note: 初始化，存储铭文的见证数据
	for i := range tool.txCtxDataList {
		if i == 0 {
			continue
		}
		revealTx := tool.revealTx[i-1]
		idx := 0

		// warn: idx的含义
		// note: 计算签名哈希： 基于当前的reveal tx, 关联的前一个交易输出，以及铭文的脚本，计算签名哈希； 这是准备签名验证数据的关键，保证交易能在网络中被正确验证
		witnessArray, err := txscript.CalcTapscriptSignaturehash(txscript.NewTxSigHashes(revealTx, tool.revealTxPrevOutputFetcher),
			txscript.SigHashDefault, revealTx, idx, tool.revealTxPrevOutputFetcher, txscript.NewBaseTapLeaf(tool.txCtxDataList[i].inscriptionScript))
		if err != nil {
			return errors.Wrap(err, "calc tapscript signaturehash error")
		}

		// note: 使用私钥对签名哈希进行签名
		signature, err := schnorr.Sign(tool.txCtxDataList[i].privateKey, witnessArray)
		if err != nil {
			return errors.Wrap(err, "schnorr sign error")
		}
		witnessList[i-1] = wire.TxWitness{signature.Serialize(), tool.txCtxDataList[i].inscriptionScript, tool.txCtxDataList[i].controlBlockWitness}
	}
	for i := range witnessList { // note: 为每个揭示交易设置见证数据
		tool.revealTx[i].TxIn[0].Witness = witnessList[i]
	}
	// check tx max tx wight
	for i, tx := range tool.revealTx { // note: 验证每个揭示交易的大小是否符合比特币网络的标准交易重量限制
		revealWeight := blockchain.GetTransactionWeight(btcutil.NewTx(tx))
		if revealWeight > MaxStandardTxWeight {
			return errors.New(fmt.Sprintf("reveal(index %d) transaction weight greater than %d (MAX_STANDARD_TX_WEIGHT): %d", i, MaxStandardTxWeight, revealWeight))
		}
	}
	return nil
}

type CtxTxData struct {
	//commit tx imfo
	PrevCommitTxHash string // 相应的commit tx的哈希

	//middle tx info
	MiddleTx            *wire.MsgTx // middle tx的具体数据
	MiddleTxHash        string      // middle tx的哈希
	MiddleInscriptionId string      // 在 middle tx 铭刻的铭文Id

	//reveal tx info
	RevealTxData []*RevealTxData // 揭示交易的具体数据

}

type RevealTxData struct {
	CtxPrivateKey         string      // 签名见证信息的私钥, Base64编码
	RecoveryPrivateKeyWIF string      // 用于恢复私钥的WIF格式的私钥
	RevealTx              *wire.MsgTx // 揭示交易的具体数据

	// inscribed inscriptions info
	IsSend        bool   // 是否已成功发送交易
	RevealTxHash  string // 已成功发送后，揭示交易的哈希
	InscriptionId string // 已成功发送后，已铭刻铭文的铭文Id
}

func (tool *InscriptionTool) saveRevealTx() ([]*RevealTxData, error) {
	//newCommitTxHash, err := chainhash.NewHashFromStr(prevCommitTxHash)
	//if err != nil {
	//	return nil, errors.Wrap(err, "new commit tx hash from str error")
	//}

	revealTxData := make([]*RevealTxData, len(tool.txCtxDataList)-1)
	for i := range tool.txCtxDataList {
		if i == 0 {
			continue
		}

		ctxPrivateKey := base64.StdEncoding.EncodeToString(tool.txCtxDataList[i].privateKey.Serialize())
		revealTxData[i-1] = &RevealTxData{
			CtxPrivateKey:         ctxPrivateKey,
			RecoveryPrivateKeyWIF: tool.txCtxDataList[i].recoveryPrivateKeyWIF,
			RevealTx:              tool.revealTx[i-1],
		}
	}
	return revealTxData, nil
}

func (tool *InscriptionTool) GetRecoveryKeyWIFList() []string {
	wifList := make([]string, len(tool.txCtxDataList))
	for i := range tool.txCtxDataList {
		wifList[i] = tool.txCtxDataList[i].recoveryPrivateKeyWIF
	}
	return wifList
}

func getTxHex(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func (tool *InscriptionTool) GetCommitTxHex() (string, error) {
	return getTxHex(tool.commitTx)
}

func (tool *InscriptionTool) GetMiddleTxHex() (string, error) {
	return getTxHex(tool.middleTx)
}

func (tool *InscriptionTool) GetRevealTxHexList() ([]string, error) {
	txHexList := make([]string, len(tool.revealTx))
	for i := range tool.revealTx {
		txHex, err := getTxHex(tool.revealTx[i])
		if err != nil {
			return nil, err
		}
		txHexList[i] = txHex
	}
	return txHexList, nil
}

func (tool *InscriptionTool) sendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	if tool.client.rpcClient != nil {
		return tool.client.rpcClient.SendRawTransaction(tx, false)
	} else {
		return tool.client.btcApiClient.BroadcastTx(tx)
	}
}

func (tool *InscriptionTool) calculateFee() int64 {
	fees := int64(0)
	for _, in := range tool.middleTx.TxIn {
		fees += tool.middleTxPrevOutputFetcher.FetchPrevOutput(in.PreviousOutPoint).Value
	}
	for _, out := range tool.middleTx.TxOut {
		fees -= out.Value
	}

	for _, tx := range tool.revealTx {
		for _, in := range tx.TxIn {
			fees += tool.revealTxPrevOutputFetcher.FetchPrevOutput(in.PreviousOutPoint).Value
		}
		for _, out := range tx.TxOut {
			fees -= out.Value
		}
	}
	return fees
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getServiceFee(inscAmount int64) int64 {
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
