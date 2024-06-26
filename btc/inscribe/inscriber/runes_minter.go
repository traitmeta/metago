package inscriber

import (
	"bytes"
	"encoding/hex"
	"log"
	"sync"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

type RunesMinter struct {
	net    *chaincfg.Params
	client BTCBaseClient
}

func NewRunesMinter(net *chaincfg.Params, btcClient BTCBaseClient) *RunesMinter {
	tool := &RunesMinter{
		net:    net,
		client: btcClient,
	}

	return tool
}

type PayInfo struct {
	Addr           string
	PkScript       string
	InscriptionFee int64
	MinerFee       int64
}

func (rm *RunesMinter) Inscribe(commitTxHash string, actualMiddlePrevOutputFee int64, payAddrPK string, req MintReq) (*SendResult, error) {
	middleTx, revealTxs, err := rm.BuildRunesTxs(commitTxHash, actualMiddlePrevOutputFee, payAddrPK, req)
	if err != nil {
		return nil, err
	}

	sendMap := rm.SendRunesTxs(middleTx, revealTxs)

	var sendResult = &SendResult{
		MiddleTx:  middleTx,
		RevealTxs: revealTxs,
		TxsStatus: sendMap,
	}

	return sendResult, nil
}

func (rm *RunesMinter) sendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	return rm.client.SendRawTransaction(tx)
}

func getTxHex(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func (rm *RunesMinter) CalcRunesTxsFee(payAddrPK string, req MintReq) (int64, error) {
	var builder = NewBuilder(rm.net)
	allWallet, err := builder.BuildAllUsedWallet(req, payAddrPK)
	if err != nil {
		return 0, errors.Wrap(err, "build all used wallet error")
	}

	revealWrapTxs, err := builder.BuildRevealTxsWithEmptyInput(allWallet, req)
	if err != nil {
		return 0, errors.Wrap(err, "build all empty reveal tx error")
	}

	middleWrapTx, err := builder.BuildMiddleTxWithEmptyInput(req, revealWrapTxs, allWallet[0].PkScript)
	if err != nil {
		return 0, errors.Wrap(err, "build empty middle tx error")
	}

	return middleWrapTx.PrevOutput.Value, nil
}

func (rm *RunesMinter) BuildRunesTxs(commitTxHash string, actualMiddlePrevOutputFee int64, payAddrPK string, req MintReq) (*WrapTx, []*WrapTx, error) {
	var builder = NewBuilder(rm.net)
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

func (rm *RunesMinter) SendRunesTxs(middleTx *WrapTx, revealTxs []*WrapTx) map[string]bool {
	var sendTxMap = make(map[string]bool)

	middleTxHash, err := rm.sendRawTransaction(middleTx.WireTx)
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
			revealTxHash, err := rm.sendRawTransaction(revealTxs[i].WireTx)
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
