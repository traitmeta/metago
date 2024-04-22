package walletmgr

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	log "github.com/sirupsen/logrus"
)

type WalletMgr struct {
	net     *chaincfg.Params
	wallets []*Wallet
	cache   *Cache
}

func InitWalletMgr(cacheDir string, net *chaincfg.Params) (*WalletMgr, error) {
	cache, err := InitCache(cacheDir)
	if err != nil {
		return nil, err
	}

	var wallets []*Wallet
	wifs, err := cache.ReadAllWalletWiF()
	if err != nil {
		return nil, err
	}

	for _, wif := range wifs {
		wallet, err := InitWallet(wif, cache, net)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	return &WalletMgr{
		net:     net,
		wallets: wallets,
		cache:   cache,
	}, nil
}

func (wm *WalletMgr) MintsRunes(cli *rpcclient.Client, runeId string, destination string) {
	for i, wallet := range wm.wallets {
		err := wallet.MintRunes(cli, runeId, destination)
		if err != nil {
			log.Info(fmt.Sprintf("wallet %d mint runes %s to %s failed", i, runeId, destination))
			continue
		}
	}
}

func (wm *WalletMgr) GetMintProcessing(runeId string, destination string) {
}
