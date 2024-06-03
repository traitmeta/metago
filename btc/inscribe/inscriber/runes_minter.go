package inscriber

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
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/traitmeta/metago/btc/runes-tools/txbuilder/runes"
)

type CtxTxData struct {
	//commit tx imfo
	PrevCommitTxHash string // 相应的commit tx的哈希

	//middle tx info
	MiddleTx            *wire.MsgTx // middle tx的具体数据
	MiddleTxHash        string      // middle tx的哈希
	MiddleInscriptionId string      // 在 middle tx 铭刻的铭文Id

	//reveal tx info
	RevealTxData []*RevealTxData // reveal tx data
}

type RevealTxData struct {
	CtxPrivateKey         string      // private key, Base64编码
	RecoveryPrivateKeyWIF string      // 用于恢复私钥的WIF格式的私钥
	RevealTx              *wire.MsgTx // 揭示交易的具体数据

	// inscribed inscriptions info
	IsSend        bool   // 是否已成功发送交易
	RevealTxHash  string // 已成功发送后，揭示交易的哈希
	InscriptionId string // 已成功发送后，已铭刻铭文的铭文Id
}

type inscriptionTxCtxData struct {
	privateKey              *btcec.PrivateKey
	commitTxAddress         btcutil.Address
	commitTxAddressPkScript []byte
	recoveryPrivateKeyWIF   string
	revealTxPrevOutput      *wire.TxOut
	middleTxPrevOutput      *wire.TxOut
}

type RunesMinter struct {
	net                       *chaincfg.Params
	client                    BTCBaseClient
	txCtxDataList             []*inscriptionTxCtxData
	revealTxPrevOutputFetcher *txscript.MultiPrevOutFetcher // note: 用于获取reveal tx的输入
	revealTx                  []*wire.MsgTx                 // note: reveal tx
	commitTx                  *wire.MsgTx                   // note: commit tx
	middleTxPrevOutputFetcher *txscript.MultiPrevOutFetcher // note: 用于获取middle tx的输入
	middleTx                  *wire.MsgTx                   // note: 连接commit tx和（除了第一笔）reveal tx的中间tx： 包含reveal tx1 + commit tx2... + 手续费 + 找零
}

func NewRunesMintInscribeTool(net *chaincfg.Params, btcClient BTCBaseClient, runesCli *runes.Client) (*RunesMinter, error) {
	tool := &RunesMinter{
		net:                       net,
		client:                    btcClient,
		revealTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		middleTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
	}
	return tool, nil
}

type PayInfo struct {
	Addr           string
	PkScript       string
	InscriptionFee int64
	MinerFee       int64
}

func (tool *RunesMinter) GetPayAddrAndFee(req MintReq) (*PayInfo, error) {
	var payInfo = &PayInfo{}
	tool.txCtxDataList = make([]*inscriptionTxCtxData, req.Count)         // not1e: 初始化，铭文列表
	tool.revealTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil) // note: 初始化，reveal tx的输入
	tool.middleTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil) // note: 初始化，middle tx的输入
	destinations := make([]string, req.Count)                             // note: 初始化，铭文的接收地址

	for i := 0; i < req.Count; i++ {
		privateKey, err := btcec.NewPrivateKey() // note: 创建一个密钥对，用来构建reveal tx
		if err != nil {
			return nil, errors.Wrap(err, "create private key error")
		}
		if i == 0 { // warn: 保存构建第一个铭文的私钥
			privKeyBytes := privateKey.Serialize()                             // 将私钥转换为字节数组
			payInfo.PkScript = base64.StdEncoding.EncodeToString(privKeyBytes) // 使用Base64编码将字节数组转换为字符串
		}
		txCtxData, err := createRuneMintTxCtxData(tool.net, privateKey) // note: 创建commit交易及包含铭文信息的Taproot脚本信息

		if err != nil {
			return nil, errors.Wrap(err, "create inscription tx ctx data error")
		}
		tool.txCtxDataList[i] = txCtxData
		destinations[i] = req.Receiver
	}

	err := tool.buildEmptyRevealTx(destinations, req.FeeRate, req.RuneId)
	if err != nil {
		return nil, errors.Wrap(err, "build empty reveal tx error")
	}

	payInfo.InscriptionFee, payInfo.MinerFee, err = tool.buildEmptyMiddleTx(destinations[0], req)
	if err != nil {
		return nil, errors.Wrap(err, "build empty middle tx error")
	}

	payInfo.Addr = tool.txCtxDataList[0].commitTxAddress.String()
	return payInfo, nil
}

func (tool *RunesMinter) Inscribe(commitTxHash string, actualMiddlePrevOutputFee int64, payAddrPK string, req MintReq) (ctxTxData *CtxTxData, err error) {
	tool.revealTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil) // note: 初始化，reveal tx的输入
	tool.middleTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil) // note: 初始化，middle tx的输入
	tool.txCtxDataList = make([]*inscriptionTxCtxData, req.Count)         // note: 初始化，铭文列表
	destinations := make([]string, req.Count)                             // note: 初始化，铭文的接收地址
	for i := 0; i < req.Count; i++ {
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
		txCtxData, err := createRuneMintTxCtxData(tool.net, privateKey) // note: 创建commit交易及包含铭文信息的Taproot脚本信息
		if err != nil {
			return nil, errors.Wrap(err, "create inscription tx ctx data error")
		}
		tool.txCtxDataList[i] = txCtxData
		destinations[i] = req.Receiver
	}

	err = tool.buildEmptyRevealTx(destinations, req.FeeRate, req.RuneId)
	if err != nil {
		return nil, errors.Wrap(err, "build empty reveal tx error")
	}

	totalMiddlePrevOutput, _, err := tool.buildEmptyMiddleTx(destinations[0], req)
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

	ctxTxData = &CtxTxData{}
	ctxTxData.PrevCommitTxHash = commitTxHash

	ctxTxData.MiddleTx = tool.middleTx
	middleTxHash, err := tool.sendRawTransaction(tool.middleTx)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("send middle tx error: tx_hash %s", tool.middleTx.TxHash().String()))
	}
	log.Printf("middleTxHash %s \n", middleTxHash.String())

	ctxTxData.MiddleTxHash = middleTxHash.String()
	ctxTxData.MiddleInscriptionId = fmt.Sprintf("%si0", middleTxHash)

	// return all reveal tx data
	ctxTxData.RevealTxData, err = tool.saveRevealTx()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("saveRevealTx error: tx_hash %s", tool.middleTx.TxHash().String()))
	}

	// warn: err info：{"code":-26,"message":"too-long-mempool-chain, too many descendants for tx... [limit: 25]"}
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
			log.Printf("revealTxHash %d %s \n", i, _revealTxHash.String())
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

func createRuneMintTxCtxData(net *chaincfg.Params, privateKey *btcec.PrivateKey) (*inscriptionTxCtxData, error) {
	// note: 生成最终的 Taproot 地址（commit tx 的输出地址）和 Pay-to-Taproot(P2TR) 地址的脚本。
	commitTxAddress, err := btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(privateKey.PubKey()).SerializeCompressed()[1:], net)
	log.Println("commitTxAddress: ", commitTxAddress)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address error")
	}

	commitTxAddressPkScript, err := txscript.PayToAddrScript(commitTxAddress)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address pk script error")
	}

	recoveryPrivateKeyWIF, err := btcutil.NewWIF(txscript.TweakTaprootPrivKey(*privateKey, []byte{}), net, true)
	if err != nil {
		return nil, errors.Wrap(err, "create recovery private key wif error")
	}

	return &inscriptionTxCtxData{
		privateKey:              privateKey,
		commitTxAddress:         commitTxAddress,
		commitTxAddressPkScript: commitTxAddressPkScript,
		recoveryPrivateKeyWIF:   recoveryPrivateKeyWIF.String(),
	}, nil
}

func (tool *RunesMinter) buildEmptyRevealTx(destination []string, feeRate int64, runeId string) (err error) {
	var revealTx []*wire.MsgTx
	total := len(tool.txCtxDataList)
	addTxInTxOutIntoRevealTx := func(tx *wire.MsgTx, index int, runeId string) error {
		in := wire.NewTxIn(&wire.OutPoint{Index: uint32(index)}, nil, nil)
		in.Sequence = defaultSequenceNum
		tx.AddTxIn(in)

		output, err := runes.CreateMintRuneStoneOutput(runeId)
		if err != nil {
			return errors.Wrap(err, "build runestone script fail")
		}
		tx.AddTxOut(output)

		receiver, err := btcutil.DecodeAddress(destination[index], tool.net) // note: 生成铭文的接收地址, 这里是destination[index]的P2TR地址
		if err != nil {
			return errors.Wrap(err, "decode address error")
		}

		scriptPubKey, err := txscript.PayToAddrScript(receiver)
		if err != nil {
			return errors.Wrap(err, "pay to address script error")
		}
		out := wire.NewTxOut(defaultRevealOutValue, scriptPubKey) // note: 再构造reveal tx的交易输出，也就是铭文的接收地址
		tx.AddTxOut(out)
		return nil
	}

	revealTx = make([]*wire.MsgTx, total-1) // note: 初始化，创建多个reveal tx

	for i := 1; i < total; i++ { // modified: 循环除了第一个铭文数据，构造除了第一个reveal tx 的信息
		tx := wire.NewMsgTx(int32(2))
		err := addTxInTxOutIntoRevealTx(tx, i, runeId)
		if err != nil {
			return errors.Wrap(err, "add tx in tx out into reveal tx error")
		}

		prevOutput := defaultRevealOutValue + int64(tx.SerializeSize())*feeRate // note: 铭文的value + 铭文的交易大小 * feeRate
		{
			emptySignature := make([]byte, 64)
			witness_fee := (int64(wire.TxWitness{emptySignature}.SerializeSize()+2+3) / 4) * feeRate
			prevOutput += witness_fee

			tool.txCtxDataList[i].revealTxPrevOutput = &wire.TxOut{
				PkScript: tool.txCtxDataList[i].commitTxAddressPkScript,
				Value:    prevOutput,
			}

		}
		revealTx[i-1] = tx
	}

	tool.revealTx = revealTx
	return nil
}

// note: include reveal tx1 + commit tx2... + 手续费 + 找零
func (tool *RunesMinter) buildEmptyMiddleTx(firstDestination string, req MintReq) (totalPrevOutput, minerFee int64, err error) {
	tx := wire.NewMsgTx(int32(2))
	addTxInTxOutOfRevealTx1IntoMiddleTx := func(tx *wire.MsgTx, index int) error {
		in := wire.NewTxIn(&wire.OutPoint{Index: uint32(index)}, nil, nil)
		in.Sequence = defaultSequenceNum
		tx.AddTxIn(in)

		receiver, err := btcutil.DecodeAddress(firstDestination, tool.net)
		if err != nil {
			return errors.Wrap(err, "decode address error")
		}

		scriptPubKey, err := txscript.PayToAddrScript(receiver)
		if err != nil {
			return errors.Wrap(err, "pay to address script error")
		}

		out := wire.NewTxOut(defaultRevealOutValue, scriptPubKey)
		tx.AddTxOut(out)
		return nil
	}

	// TODO 使用OUTPOINT，而不是默认指定0
	err = addTxInTxOutOfRevealTx1IntoMiddleTx(tx, 0)
	if err != nil {
		return 0, 0, errors.Wrap(err, "add tx in tx out of reveal tx1 into middle tx error")
	}

	prevOutput := defaultRevealOutValue
	{
		emptySignature := make([]byte, 64)
		witnessFee := (int64(wire.TxWitness{emptySignature}.SerializeSize()+2+3) / 4) * req.FeeRate
		prevOutput += witnessFee
		minerFee += witnessFee
	}
	totalPrevOutput += prevOutput

	for i := range tool.txCtxDataList {
		if i == 0 {
			continue
		}
		tx.AddTxOut(tool.txCtxDataList[i].revealTxPrevOutput)
		totalPrevOutput += tool.txCtxDataList[i].revealTxPrevOutput.Value
		minerFee += int64(tool.txCtxDataList[i].revealTxPrevOutput.Value) - defaultRevealOutValue
	}

	output, err := runes.CreateMintRuneStoneOutput(req.RuneId)
	if err != nil {
		return 0, 0, errors.Wrap(err, "build runestone script fail")
	}

	tx.AddTxOut(output)
	txSizeFee := int64(tx.SerializeSize()) * req.FeeRate
	minerFee += txSizeFee
	totalPrevOutput += txSizeFee

	tool.txCtxDataList[0].middleTxPrevOutput = &wire.TxOut{
		PkScript: tool.txCtxDataList[0].commitTxAddressPkScript,
		Value:    totalPrevOutput,
	}

	tool.middleTx = tx
	return totalPrevOutput, minerFee, nil
}

// 1. fill commit tx value
// 2. sign transaction
func (tool *RunesMinter) completeMiddleTx(commitTxHash string, actualMiddlePrevOutputFee int64) error {
	newCommitTxHash, err := chainhash.NewHashFromStr(commitTxHash)
	if err != nil {
		return errors.Wrap(err, "failed converting transaction hash")
	}

	fmt.Println("newCommitTxHash", newCommitTxHash)
	tool.txCtxDataList[0].middleTxPrevOutput.Value = actualMiddlePrevOutputFee

	tool.middleTxPrevOutputFetcher.AddPrevOut(wire.OutPoint{
		Hash:  *newCommitTxHash,
		Index: uint32(0),
	}, tool.txCtxDataList[0].middleTxPrevOutput)

	tool.middleTx.TxIn[0].PreviousOutPoint.Hash = *newCommitTxHash

	middleTx := tool.middleTx

	witnessArray, err := txscript.CalcTaprootSignatureHash(txscript.NewTxSigHashes(middleTx, tool.middleTxPrevOutputFetcher),
		txscript.SigHashDefault, middleTx, 0, tool.middleTxPrevOutputFetcher)
	if err != nil {
		return errors.Wrap(err, "calc tapscript signaturehash error")
	}

	priv := txscript.TweakTaprootPrivKey(*tool.txCtxDataList[0].privateKey, []byte{})
	signature, err := schnorr.Sign(priv, witnessArray)
	if err != nil {
		return errors.Wrap(err, "schnorr sign error")
	}
	witness := wire.TxWitness{signature.Serialize()}
	tool.middleTx.TxIn[0].Witness = witness

	revealWeight := blockchain.GetTransactionWeight(btcutil.NewTx(middleTx))
	if revealWeight > MaxStandardTxWeight {
		return errors.New(fmt.Sprintf("middle transaction weight greater than %d (MAX_STANDARD_TX_WEIGHT): %d", MaxStandardTxWeight, revealWeight))
	}

	return nil
}

// note: 完成铭文揭示交易（Reveal Transaction）的构建过程，主要通过计算签名哈希(hash)、生成签名，并为每个揭示交易设置见证数据（witness data）。
// note: 此外，它还验证每个揭示交易的大小是否符合比特币网络的标准交易重量限制。
func (tool *RunesMinter) completeRevealTx() error {
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

		witnessArray, err := txscript.CalcTaprootSignatureHash(txscript.NewTxSigHashes(revealTx, tool.revealTxPrevOutputFetcher),
			txscript.SigHashDefault, revealTx, idx, tool.revealTxPrevOutputFetcher)
		if err != nil {
			return errors.Wrap(err, "calc tapscript signaturehash error")
		}

		// note: 使用私钥对签名哈希进行签名
		priv := txscript.TweakTaprootPrivKey(*tool.txCtxDataList[i].privateKey, []byte{})
		signature, err := schnorr.Sign(priv, witnessArray)
		if err != nil {
			return errors.Wrap(err, "schnorr sign error")
		}
		witnessList[i-1] = wire.TxWitness{signature.Serialize()}
	}
	for i := range witnessList {
		tool.revealTx[i].TxIn[0].Witness = witnessList[i]
	}

	// check tx max tx wight
	for i, tx := range tool.revealTx {
		revealWeight := blockchain.GetTransactionWeight(btcutil.NewTx(tx))
		if revealWeight > MaxStandardTxWeight {
			return errors.New(fmt.Sprintf("reveal(index %d) transaction weight greater than %d (MAX_STANDARD_TX_WEIGHT): %d", i, MaxStandardTxWeight, revealWeight))
		}
	}
	return nil
}

func (tool *RunesMinter) saveRevealTx() ([]*RevealTxData, error) {
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

func (tool *RunesMinter) GetRecoveryKeyWIFList() []string {
	wifList := make([]string, len(tool.txCtxDataList))
	for i := range tool.txCtxDataList {
		wifList[i] = tool.txCtxDataList[i].recoveryPrivateKeyWIF
	}
	return wifList
}

func (tool *RunesMinter) GetCommitTxHex() (string, error) {
	return getTxHex(tool.commitTx)
}

func (tool *RunesMinter) GetMiddleTxHex() (string, error) {
	return getTxHex(tool.middleTx)
}

func (tool *RunesMinter) GetRevealTxHexList() ([]string, error) {
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

func (tool *RunesMinter) sendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	return tool.client.SendRawTransaction(tx)
}

func getTxHex(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}
