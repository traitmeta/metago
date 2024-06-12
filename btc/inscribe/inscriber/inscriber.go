package inscriber

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Inscriber struct {
	minter *RunesMinter
	task   *SendTask
	net    *chaincfg.Params
}

func NewInscriber(btcClient BTCBaseClient) *Inscriber {
	minter := NewRunesMinter(&chaincfg.MainNetParams, btcClient)
	return &Inscriber{
		minter: minter,
		task:   nil,
	}
}

func (i *Inscriber) SetTask(task *SendTask) {
	i.task = task
}

type Order struct {
	Id            string
	Who           string
	Runes         string
	Count         int64
	FeeRate       int64
	WifPrivateKey string
	PayAddress    string
	MintFee       int64
}

func (i *Inscriber) CreateOrder(who, runes string, count int64) (*Order, error) {
	uuidV4 := uuid.New()

	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "create private key error")
	}

	payAddress, err := btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(privateKey.PubKey()).SerializeCompressed()[1:], i.net)
	if err != nil {
		return nil, errors.Wrap(err, "create pay tx address error")
	}

	wif, err := btcutil.NewWIF(txscript.TweakTaprootPrivKey(*privateKey, []byte{}), i.net, true)
	if err != nil {
		return nil, errors.Wrap(err, "create recovery private key wif error")
	}

	return &Order{
		Id:            uuidV4.String(),
		Who:           who,
		Runes:         runes,
		Count:         count,
		WifPrivateKey: wif.String(),
		PayAddress:    payAddress.EncodeAddress(),
	}, nil
}

func (i *Inscriber) CalcFee(order *Order) (int64, string, error) {
	mintFee, err := i.minter.CalcRunesTxsFee(order.WifPrivateKey, MintReq{
		RuneId:   order.Runes,
		Receiver: order.Who,
		FeeRate:  order.FeeRate,
		Count:    int(order.Count),
	})
	if err != nil {
		return 0, "", errors.Wrap(err, "CalcRunesTxsFee")
	}

	order.MintFee = mintFee

	return mintFee, order.PayAddress, nil
}

func (i *Inscriber) FirstStep(who, runes string, count int64) error {
	order, err := i.CreateOrder(who, runes, count)
	if err != nil {
		return errors.Wrap(err, "CalcRunesTxsFee")
	}

	mintFee, payAddress, err := i.CalcFee(order)
	if err != nil {
		return errors.Wrap(err, "CalcRunesTxsFee")
	}

	log.Info("you need pay", "to_address", payAddress, "value", mintFee)

	// TODO waitting for pay tx is in mempool or in block
	payTxHash := ""
	task := NewSendTask(i.minter.client, order.Who, order.Runes, order.Id, order.Count)
	sr, err := i.minter.Inscribe(payTxHash, mintFee, order.WifPrivateKey, MintReq{
		RuneId:   order.Runes,
		Receiver: order.Who,
		FeeRate:  order.FeeRate,
		Count:    int(order.Count),
	})
	if err != nil {
		return errors.Wrap(err, "Inscribe")
	}

	task.LoopSendTxs(sr)
	return nil
}
