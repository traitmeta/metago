package inscriber

import (
	"bytes"
	"encoding/hex"
	"log"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
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

func (tool *RunesMinter) Inscribe(commitTxHash string, actualMiddlePrevOutputFee int64, payAddrPK string, req MintReq) (*SendResult, error) {
	middleTx, revealTxs, err := tool.BuildRunesTxs(commitTxHash, actualMiddlePrevOutputFee, payAddrPK, req)
	if err != nil {
		return nil, err
	}

	sendMap := tool.SendRunesTxs(middleTx, revealTxs)

	var sendResult = &SendResult{
		MiddleTx:  middleTx,
		RevealTxs: revealTxs,
		TxsStatus: sendMap,
	}

	return sendResult, nil
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

func (tool *RunesMinter) BuildRunesTxs(commitTxHash string, actualMiddlePrevOutputFee int64, payAddrPK string, req MintReq) (*WrapTx, []*WrapTx, error) {
	var builder = NewBuilder(tool.net)
	allWallet, err := builder.BuildAllUsedWallet(req, payAddrPK)
	if err != nil {
		return nil, nil, errors.Wrap(err, "build all used wallet error")
	}

	revealWrapTxs, err := builder.BuildRevealTxsWithEmptyInput(allWallet, req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "build all empty reveal tx error")
	}

	middleWrapTx, err := builder.BuildMiddleTxWithEmptyInput(req, revealWrapTxs, allWallet[0].PkScript)
	if err != nil {
		return nil, nil, errors.Wrap(err, "build empty middle tx error")
	}

	if middleWrapTx.PrevOutput.Value > actualMiddlePrevOutputFee {
		return nil, nil, errors.New("actualMiddlePrevOutputFee is not enough")
	}

	if err = builder.CompleteMiddleTx(*allWallet[0].PrivateKey, middleWrapTx, commitTxHash, actualMiddlePrevOutputFee); err != nil {
		return nil, nil, errors.Wrap(err, "complete middle tx error")
	}

	if err = builder.CompleteRevealTxs(revealWrapTxs, allWallet, middleWrapTx.WireTx.TxHash()); err != nil {
		return nil, nil, errors.Wrap(err, "complete reveal tx error")
	}

	return middleWrapTx, revealWrapTxs, nil
}

func (tool *RunesMinter) SendRunesTxs(middleTx *WrapTx, revealTxs []*WrapTx) map[string]bool {
	var sendTxMap = make(map[string]bool)

	middleTxHash, err := tool.sendRawTransaction(middleTx.WireTx)
	if err != nil {
		log.Printf("send middle tx error %s \n", middleTxHash.String())
		return nil
	}
	log.Printf("middleTxHash %s \n", middleTxHash.String())

	sendTxMap[middleTxHash.String()] = true
	var wg sync.WaitGroup
	minTx := min(len(revealTxs), 23)
	for i := 0; i < minTx; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			revealTxHash, err := tool.sendRawTransaction(revealTxs[i].WireTx)
			if err != nil {
				log.Printf("revealTxHash %d %s , err: %v \n", i, revealTxHash.String(), err)
				return
			}

			sendTxMap[revealTxHash.String()] = true
		}(i)
	}
	wg.Wait()

	return sendTxMap
}
