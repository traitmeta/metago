package walletmgr

import (
	"github.com/btcsuite/btcd/chaincfg"
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

func (wm *WalletMgr) MintsRunes() {

}
