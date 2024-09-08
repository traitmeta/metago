package walletmgr

import "github.com/pkg/errors"

var (
	ErrNotFindWalletWifInCache = errors.New("get wallet wif cache failed")
	ErrNotFindPrevInfoInCache  = errors.New("get wallet previous tx info in cache failed")
	ErrWritePrevInfoInCache    = errors.New("write wallet previous tx info in cache failed")
)
