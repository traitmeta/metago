package ord

import (
	"encoding/base64"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

type StandardInscriber struct {
	net                   *chaincfg.Params
	client                *BlockchainClient
	serviceFeeReceiveAddr string // note: 用于存储平台手续费的地址

	txCtxDataList             *InscribeTxsPreview
	revealTxPrevOutputFetcher *txscript.MultiPrevOutFetcher // note: 用于获取reveal tx的输入
	revealTx                  []*wire.MsgTx                 // note: reveal tx
	commitTx                  *wire.MsgTx                   // note: commit tx
	middleTxPrevOutputFetcher *txscript.MultiPrevOutFetcher // note: 用于获取middle tx的输入
	middleTx                  *wire.MsgTx                   // note: 连接commit tx和（除了第一笔）reveal tx的中间tx： 包含reveal tx1 + commit tx2... + 手续费 + 找零
}

func NewInscribeTool(net *chaincfg.Params, rpcclient *rpcclient.Client, serviceFeeReceiveAddr string) (*StandardInscriber, error) {
	if serviceFeeReceiveAddr == "" {
		return nil, errors.New("service fee receive address is empty")
	}

	tool := &StandardInscriber{
		net: net,
		client: &BlockchainClient{
			RpcClient: rpcclient,
		},
		serviceFeeReceiveAddr:     serviceFeeReceiveAddr,
		revealTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		middleTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
	}
	return tool, nil
}

func NewInscribeToolWithBtcApiClient(net *chaincfg.Params, rpcClient *rpcclient.Client, serviceFeeReceiveAddr string) (*StandardInscriber, error) {
	if serviceFeeReceiveAddr == "" {
		return nil, errors.New("service fee receive address is empty")
	}

	tool := &StandardInscriber{
		net: net,
		client: &BlockchainClient{
			RpcClient: rpcClient,
		},
		serviceFeeReceiveAddr:     serviceFeeReceiveAddr,
		revealTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		middleTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
	}
	return tool, nil
}

type Info struct {
	PayAddr       string
	PayPrivateKey string
	InscribeFee   int64
	ServiceFee    int64
	MinerFee      int64
}

func (ins *StandardInscriber) AddInOutToTx(tx *wire.MsgTx, index int, destination string, revealOutValue int64) error {
	in := wire.NewTxIn(&wire.OutPoint{Index: uint32(index)}, nil, nil) // note: 先构造reveal tx的空交易输入
	in.Sequence = DefaultSequenceNum
	tx.AddTxIn(in)
	receiver, err := btcutil.DecodeAddress(destination, ins.net) // note: 生成铭文的接收地址, 这里是destination[index]的P2TR地址
	if err != nil {
		return errors.Wrap(err, "decode address error")
	}

	scriptPubKey, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return errors.Wrap(err, "pay to address script error")
	}

	// note: 再构造reveal tx的交易输出，也就是铭文的接收地址
	tx.AddTxOut(wire.NewTxOut(revealOutValue, scriptPubKey))

	return nil
}

func (ins *StandardInscriber) Init(req InscriptionRequest) error {
	ins.InitPrevOutputFetcher()
	// 1. 构建空的交易
	txs, err := ins.InitAllRevealTx(req.DataList, req.RevealOutValue)
	if err != nil {
		return err
	}

	// 2. 构建所有的 铭文witness
	// TODO add first private key
	privKeys, err := ins.InitPrivateKeys(len(req.DataList), req.PrivateKey)
	if err != nil {
		return err
	}

	for i := 0; i < len(req.DataList); i++ {
		txs[i].PrivateKey = privKeys[i]
	}

	// 构建所有reveal 交易的铭文witness
	allWitness, err := ins.InitAllWitness(req.DataList, privKeys)
	if err != nil {
		return err
	}

	for i := 0; i < len(req.DataList); i++ {
		txs[i].WitnessScript = allWitness[i]
	}

	// 3. 计算所有的
	//scripts, err := ins.InitAllScript(txs)
	//if err != nil {
	//	return err
	//}

	var totalPrevOutValue int64
	for i, tx := range txs {
		tx.SetSize(int64(tx.Raw.SerializeSize()))
		script, err := ins.InitScript(tx)
		if err != nil {
			return err
		}

		prevOutValue := tx.CalcPrevOutput(req.RevealOutValue, req.FeeRate)
		if i != 0 {
			totalPrevOutValue += prevOutValue
		}

		tx.SetTxPrevOutput(script.CommitTxPkScript, prevOutValue)
	}

	return nil
}

func (ins *StandardInscriber) InitPrevOutputFetcher() {
	ins.revealTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil) // note: 初始化，reveal tx的输入
	ins.middleTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil) // note: 初始化，middle tx的输入
}

// InitAllRevealTx middle tx is a special reveal tx
func (ins *StandardInscriber) InitAllRevealTx(dataList []InscriptionData, revealOutValue int64) ([]*InscriptionRawTx, error) {
	revealOutValue = GetRevealOutValue(revealOutValue)
	revealTxs, err := ins.batchNewRawTxs(dataList, revealOutValue)
	if err != nil {
		return nil, errors.Wrap(err, "build empty reveal tx error")
	}

	return revealTxs, nil
}

func (ins *StandardInscriber) InitPrivateKeys(size int, payPrivateKey string) ([]*btcec.PrivateKey, error) {
	privateKeys := make([]*btcec.PrivateKey, size) // note: 初始化，铭文列表
	for i := 0; i < size; i++ {
		if i == 0 {
			//  构建第一个铭文的私钥，用之前生成的
			decodedPrivKey, err := base64.StdEncoding.DecodeString(payPrivateKey)
			if err != nil {
				return nil, errors.Wrap(err, "decode private key error")
			}

			privateKey, _ := btcec.PrivKeyFromBytes(decodedPrivKey)
			privateKeys[i] = privateKey
			continue
		}

		privateKey, err := btcec.NewPrivateKey() // note: 创建一个密钥对，用来构建reveal tx
		if err != nil {
			return nil, errors.Wrap(err, "create private key error")
		}

		privateKeys[i] = privateKey
	}

	return privateKeys, nil
}

func (ins *StandardInscriber) InitAllWitness(dataList []InscriptionData, privateKeyList []*btcec.PrivateKey) ([]*InscriptionWitness, error) {
	var allWitness = make([]*InscriptionWitness, len(dataList))
	for i := 0; i < len(dataList); i++ {
		witness := NewInscriptionWitness()
		inscriptionScript, err := BuildInscriptionWitness([]InscriptionData{dataList[i]}, privateKeyList[i], 0) // note: 铭文内容基本构建完成， 生成铭文脚本
		if err != nil {
			return nil, errors.Wrap(err, "create inscription script error")
		}

		//  创建一个新的taproot script叶子节点，将刚才构造的铭文脚本添加到叶子节点中
		leafNode := txscript.NewBaseTapLeaf(inscriptionScript)
		proof := &txscript.TapscriptProof{
			TapLeaf:  leafNode,
			RootNode: leafNode,
		}

		// 利用前面生成的证明对象和公钥生成Control block
		controlBlock := proof.ToControlBlock(privateKeyList[i].PubKey())
		controlBlockWitness, err := controlBlock.ToBytes()
		if err != nil {
			return nil, errors.Wrap(err, "control block to bytes error")
		}
		witness.InsWitnessScript = inscriptionScript
		witness.ControlBlockWitness = controlBlockWitness
		allWitness = append(allWitness, witness)
	}

	return allWitness, nil
}

func (ins *StandardInscriber) InitAllScript(rawTxs []*InscriptionRawTx) ([]*InscribeTxPreview, error) {
	var scripts = make([]*InscribeTxPreview, len(rawTxs))
	for i, tx := range rawTxs {
		//  创建一个新的taproot script叶子节点，将刚才构造的铭文脚本添加到叶子节点中
		leafNode := txscript.NewBaseTapLeaf(tx.WitnessScript.InsWitnessScript)
		proof := &txscript.TapscriptProof{
			TapLeaf:  leafNode,
			RootNode: leafNode,
		}

		// 生成最终的 Taproot 地址（commit tx 的输出地址）和 Pay-to-Taproot(P2TR) 地址的脚本。
		tapHash := proof.RootNode.TapHash()
		commitTxAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(tx.PrivateKey.PubKey(), tapHash[:])), ins.net)
		if err != nil {
			return nil, errors.Wrap(err, "create commit tx address error")
		}

		commitTxAddressPkScript, err := txscript.PayToAddrScript(commitTxAddress)
		if err != nil {
			return nil, errors.Wrap(err, "create commit tx address pk script error")
		}

		recoveryPrivateKeyWIF, err := btcutil.NewWIF(txscript.TweakTaprootPrivKey(*tx.PrivateKey, tapHash[:]), ins.net, true)
		if err != nil {
			return nil, errors.Wrap(err, "create recovery private key wif error")
		}
		scripts[i] = &InscribeTxPreview{
			CommitTxPkScript:      commitTxAddressPkScript,
			CommitTxAddress:       commitTxAddress,
			RecoveryPrivateKeyWIF: recoveryPrivateKeyWIF.String(),
		}
	}

	return scripts, nil
}

func (ins *StandardInscriber) InitScript(tx *InscriptionRawTx) (*InscribeTxPreview, error) {
	//  创建一个新的taproot script叶子节点，将刚才构造的铭文脚本添加到叶子节点中
	leafNode := txscript.NewBaseTapLeaf(tx.WitnessScript.InsWitnessScript)
	proof := &txscript.TapscriptProof{
		TapLeaf:  leafNode,
		RootNode: leafNode,
	}

	// 生成最终的 Taproot 地址（commit tx 的输出地址）和 Pay-to-Taproot(P2TR) 地址的脚本。
	tapHash := proof.RootNode.TapHash()
	commitTxAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(tx.PrivateKey.PubKey(), tapHash[:])), ins.net)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address error")
	}

	commitTxAddressPkScript, err := txscript.PayToAddrScript(commitTxAddress)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address pk script error")
	}

	recoveryPrivateKeyWIF, err := btcutil.NewWIF(txscript.TweakTaprootPrivKey(*tx.PrivateKey, tapHash[:]), ins.net, true)
	if err != nil {
		return nil, errors.Wrap(err, "create recovery private key wif error")
	}

	return &InscribeTxPreview{
		CommitTxPkScript:      commitTxAddressPkScript,
		CommitTxAddress:       commitTxAddress,
		RecoveryPrivateKeyWIF: recoveryPrivateKeyWIF.String(),
	}, nil
}

func (ins *StandardInscriber) batchNewRawTxs(destinations []InscriptionData, revealOutValue int64) ([]*InscriptionRawTx, error) {
	total := len(destinations)
	revealTx := make([]*InscriptionRawTx, total) // 初始化，创建多个reveal tx
	for destIdx := 0; destIdx < total; destIdx++ {
		rawTx := NewInscriptionRawTx()
		// modified: 循环除了第一个铭文数据，构造除了第一个reveal tx 的信息
		tx, err := ins.newEmptyRevealTx(destIdx, destinations[destIdx].Destination, revealOutValue) // note: 往每个tx添加“空的交易输入”和输出
		if err != nil {
			return nil, errors.Wrap(err, "add tx in tx out into reveal tx error")
		}

		rawTx.Raw = tx
		revealTx[destIdx] = rawTx
	}

	return revealTx, nil
}

func (ins *StandardInscriber) newEmptyRevealTx(index int, destination string, revealOutValue int64) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)                             // note: 创建一个新的reveal tx
	err := ins.AddInOutToTx(tx, index, destination, revealOutValue) // note: 往每个tx添加“空的交易输入”和输出
	if err != nil {
		return nil, errors.Wrap(err, "add tx in tx out into reveal tx error")
	}

	return tx, nil
}

// InitEmptyMiddleTx : 构造middleTx: 连接commit tx和（除了第一笔）reveal tx的中间tx： 包含reveal tx1 + commit tx2... + 手续费 + 找零
func (ins *StandardInscriber) InitEmptyMiddleTx(txs []*InscriptionRawTx, totalRevealPrevOutput int64, revealOutValue, feeRate, inscAmount int64) (totalPrevOutput, serviceFee, minerFee int64, err error) {
	emptySignature := make([]byte, 64)
	emptyControlBlockWitness := make([]byte, 33)
	fee := (int64(wire.TxWitness{emptySignature, txs[0].WitnessScript.InsWitnessScript, emptyControlBlockWitness}.SerializeSize()+2+3) / 4) * feeRate
	minerFee += fee // note: minerFee 加上 见证文本费用
	totalPrevOutput += txs[0].TxPrevOutput.Value

	// modified: middle tx 的输出再加上reveal tx的preOutput
	for _, tx := range txs[1:] { // reveal tx的输入 是 commit tx的输出
		txs[0].Raw.AddTxOut(tx.TxPrevOutput)
		totalPrevOutput += tx.TxPrevOutput.Value
		minerFee += tx.TxPrevOutput.Value - revealOutValue
	}

	// 在middle tx的末尾，添加一个给平台手续费操作的uxto输出
	serviceFee = GetServiceFee(inscAmount)
	servicePkScript, err := getServiceFeePkScript(ins.serviceFeeReceiveAddr, ins.net)
	txs[0].Raw.AddTxOut(wire.NewTxOut(serviceFee, *servicePkScript))
	totalPrevOutput += serviceFee

	// 总费用加上交易费用
	txSizeFee := int64(txs[0].Raw.SerializeSize()) * feeRate
	minerFee += txSizeFee
	totalPrevOutput += txSizeFee

	// 设置middleTxPrevOutput:
	txs[0].TxPrevOutput = &wire.TxOut{
		PkScript: txs[0].TxPrevOutput.PkScript,
		Value:    totalPrevOutput,
	}

	return totalPrevOutput, serviceFee, minerFee, nil
}
