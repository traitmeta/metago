package witness

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/pkg/errors"

	"github.com/traitmeta/metago/btc/ord-tools/ord"
)

type SignerBuilder struct {
	net *chaincfg.Params
}

func NewSignerBuilder(net *chaincfg.Params) *SignerBuilder {
	tool := &SignerBuilder{
		net: net,
	}
	return tool
}

func (ins *SignerBuilder) InitSigner(dataList []ord.InscriptionData) ([]*SignInfo, error) {
	size := len(dataList)
	signInfos := make([]*SignInfo, size)
	for i := 0; i < size; i++ {
		privateKey, err := btcec.NewPrivateKey()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("create private key error, idx: %d", i))
		}

		signInfo, err := ins.BuildSignInfo(dataList[i], privateKey)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("build sign info error, idx: %d", i))
		}

		signInfo.PrivateKey = privateKey
		signInfos[i] = signInfo
	}

	return signInfos, nil
}

func (ins *SignerBuilder) BuildSignInfo(insData ord.InscriptionData, privateKey *btcec.PrivateKey) (*SignInfo, error) {
	witness := NewInscriptionWitness()
	inscriptionScript, err := BuildInscriptionWitness([]ord.InscriptionData{insData}, privateKey, 0) // note: 铭文内容基本构建完成， 生成铭文脚本
	if err != nil {
		return nil, errors.Wrap(err, "create inscription script error")
	}

	// 创建一个新的taproot script叶子节点，将刚才构造的铭文脚本添加到叶子节点中
	leafNode := txscript.NewBaseTapLeaf(inscriptionScript)
	proof := &txscript.TapscriptProof{
		TapLeaf:  leafNode,
		RootNode: leafNode,
	}

	// 利用前面生成的证明对象和公钥生成Control block
	controlBlock := proof.ToControlBlock(privateKey.PubKey())
	controlBlockWitness, err := controlBlock.ToBytes()
	if err != nil {
		return nil, errors.Wrap(err, "control block to bytes error")
	}
	witness.InsWitnessScript = inscriptionScript
	witness.ControlBlockWitness = controlBlockWitness

	revealAccount, err := ins.BuildRevealAccount(proof, privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "build reveal account error")
	}

	signInfo := &SignInfo{
		RevealWitness: witness,
		RevealAccount: revealAccount,
	}

	return signInfo, nil
}

func (ins *SignerBuilder) BuildRevealAccount(proof *txscript.TapscriptProof, privatekey *btcec.PrivateKey) (*RevealAccount, error) {
	// 生成最终的 Taproot 地址（commit tx 的输出地址）和 Pay-to-Taproot(P2TR) 地址的脚本。
	tapHash := proof.RootNode.TapHash()
	commitTxAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(privatekey.PubKey(), tapHash[:])), ins.net)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address error")
	}

	commitTxAddressPkScript, err := txscript.PayToAddrScript(commitTxAddress)
	if err != nil {
		return nil, errors.Wrap(err, "create commit tx address pk script error")
	}

	recoveryPrivateKeyWIF, err := btcutil.NewWIF(txscript.TweakTaprootPrivKey(*privatekey, tapHash[:]), ins.net, true)
	if err != nil {
		return nil, errors.Wrap(err, "create recovery private key wif error")
	}

	return &RevealAccount{
		CommitTxPkScript:      commitTxAddressPkScript,
		CommitTxAddress:       commitTxAddress,
		RecoveryPrivateKeyWIF: recoveryPrivateKeyWIF.String(),
	}, nil
}
