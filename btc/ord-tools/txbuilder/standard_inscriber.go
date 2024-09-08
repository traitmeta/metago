package txbuilder

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/traitmeta/metago/btc/common"
	"github.com/traitmeta/metago/btc/ord-tools/ord"
	"github.com/traitmeta/metago/btc/ord-tools/types"
	"github.com/traitmeta/metago/btc/ord-tools/witness"
)

func GetRevealOutValue(revealOutValue int64) int64 {
	result := types.DefaultRevealOutValue // note: 铭文所在 UTXO 的 sats 数量
	if revealOutValue >= types.MinRevealOutValue {
		result = revealOutValue
	}

	return result
}

type StandardInscriber struct {
	net                       *chaincfg.Params
	client                    *common.BlockchainClient
	serviceFeeReceiveAddr     string                        // note: 用于存储平台手续费的地址
	revealTxPrevOutputFetcher *txscript.MultiPrevOutFetcher // note: 用于获取reveal tx的输入
	middleTxPrevOutputFetcher *txscript.MultiPrevOutFetcher // note: 用于获取middle tx的输入
}

func NewInscribeTool(net *chaincfg.Params, rpcclient *rpcclient.Client, serviceFeeReceiveAddr string) (*StandardInscriber, error) {
	if serviceFeeReceiveAddr == "" {
		return nil, errors.New("service fee receive address is empty")
	}

	tool := &StandardInscriber{
		net: net,
		client: &common.BlockchainClient{
			RpcClient: rpcclient,
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

func (ins *StandardInscriber) Init(req ord.InscriptionRequest) error {
	signer := witness.NewSignerBuilder(ins.net)
	signInfos, err := signer.InitSigner(req.DataList)
	if err != nil {
		return err
	}

	ins.InitPrevOutputFetcher()
	// 1. 构建空的交易
	txs, err := ins.InitAllRevealTx(req.DataList, req.RevealOutValue)
	if err != nil {
		return err
	}

	// 2. 构建所有的 铭文witness
	for i, signInfo := range signInfos {
		txs[i].PrivateKey = signInfo.PrivateKey
		txs[i].WitnessScript = signInfo.RevealWitness
	}

	var totalPrevOutValue int64
	for i, tx := range txs {
		tx.SetSize(int64(tx.Raw.SerializeSize()))
		prevOutValue := tx.CalcPrevOutput(req.RevealOutValue, req.FeeRate)
		if i != 0 {
			totalPrevOutValue += prevOutValue
		}

		tx.SetTxPrevOutput(signInfos[i].RevealAccount.CommitTxPkScript, prevOutValue)
	}

	return nil
}

func (ins *StandardInscriber) InitPrevOutputFetcher() {
	ins.revealTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil) // note: 初始化，reveal tx的输入
	ins.middleTxPrevOutputFetcher = txscript.NewMultiPrevOutFetcher(nil) // note: 初始化，middle tx的输入
}

// InitAllRevealTx middle tx is a special reveal tx
func (ins *StandardInscriber) InitAllRevealTx(ords []ord.InscriptionData, revealOutValue int64) ([]*InscriptionRawTx, error) {
	revealOutValue = GetRevealOutValue(revealOutValue)
	revealTxs, err := ins.batchNewEmptyRawTxs(ords, revealOutValue)
	if err != nil {
		return nil, errors.Wrap(err, "build empty reveal tx error")
	}

	return revealTxs, nil
}

func (ins *StandardInscriber) batchNewEmptyRawTxs(ords []ord.InscriptionData, revealOutValue int64) ([]*InscriptionRawTx, error) {
	size := len(ords)
	revealTx := make([]*InscriptionRawTx, size) // 初始化，创建多个reveal tx
	for idx := 0; idx < size; idx++ {
		rawTx := NewInscriptionRawTx()
		// modified: 循环除了第一个铭文数据，构造除了第一个reveal tx 的信息
		tx, err := ins.newEmptyRevealTx(idx, ords[idx].Destination, revealOutValue) // note: 往每个tx添加“空的交易输入”和输出
		if err != nil {
			return nil, errors.Wrap(err, "add tx in tx out into reveal tx error")
		}

		rawTx.Raw = tx
		revealTx[idx] = rawTx
	}

	return revealTx, nil
}

func (ins *StandardInscriber) newEmptyRevealTx(index int, destination string, revealOutValue int64) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)                                      // note: 创建一个新的reveal tx
	err := ins.AddInOutToTx(tx, index, destination, revealOutValue, ins.net) // note: 往每个tx添加“空的交易输入”和输出
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
	if err != nil {
		return 0, 0, 0, errors.Wrap(err, "get service fee pk script error")
	}

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
	addr, err := btcutil.DecodeAddress(address, net)
	if err != nil {
		return nil, err
	}

	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	return &pkScript, nil
}

func (ins *StandardInscriber) AddInOutToTx(tx *wire.MsgTx, index int, destination string, revealOutValue int64, net *chaincfg.Params) error {
	// note: 先构造reveal tx的空交易输入
	in := wire.NewTxIn(&wire.OutPoint{Index: uint32(index)}, nil, nil)
	in.Sequence = types.DefaultSequenceNum
	tx.AddTxIn(in)

	// note: 生成铭文的接收地址, 这里是destination[index]的P2TR地址
	receiver, err := btcutil.DecodeAddress(destination, net)
	if err != nil {
		return errors.Wrap(err, "decode address error")
	}

	scriptPubKey, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return errors.Wrap(err, "pay to address script error")
	}

	tx.AddTxOut(wire.NewTxOut(revealOutValue, scriptPubKey))

	return nil
}
