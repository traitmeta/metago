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

	return &WalletMgr{
		net:     net,
		wallets: nil,
		cache:   cache,
	}, nil
}
