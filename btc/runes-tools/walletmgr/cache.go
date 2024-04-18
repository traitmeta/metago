package walletmgr

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/traitmeta/metago/btc/runes-tools/txbuilder"
)

const walletWIFKey = "key:of:wallet:wif"
const MempoolGasFee = "mempool:gas:fee:cache"

type Cache struct {
	cache *leveldb.DB
}

func InitCache(cacheDir string) (*Cache, error) {
	db, err := leveldb.OpenFile(cacheDir, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed on connect leveldb")
	}

	return &Cache{
		cache: db,
	}, nil
}

func (c *Cache) ReadAllWalletWiF() ([]string, error) {
	walletWifs, err := c.cache.Get([]byte(walletWIFKey), nil)
	if err != nil {
		return nil, ErrNotFindWalletWifInCache
	}

	var wifList []string
	if err := json.Unmarshal(walletWifs, &wifList); err != nil {
		return nil, errors.Wrap(err, "failed on unmarshal protocal data")
	}

	return wifList, nil
}

func (c *Cache) ReadWalletPrevInfo(address string) (*txbuilder.PrevInfo, error) {
	key := fmt.Sprintf("%s:previous:info", address)
	cachedData, err := c.cache.Get([]byte(key), nil)
	if err != nil {
		return nil, ErrNotFindWalletWifInCache
	}

	var info txbuilder.PrevInfo
	if err := json.Unmarshal(cachedData, &info); err != nil {
		return nil, errors.Wrap(err, "failed on unmarshal protocal data")
	}

	return &info, nil
}

func (c *Cache) ReadMemPoolGas() (int64, error) {
	cachedData, err := c.cache.Get([]byte(MempoolGasFee), nil)
	if err != nil {
		return 0, err
	}

	gas, err := strconv.ParseInt(string(cachedData), 10, 32)
	if err != nil {
		return 0, err
	}

	return gas, nil
}

func (c *Cache) CacheWalletPrevInfo(address string, data txbuilder.PrevInfo) error {
	key := fmt.Sprintf("%s:previous:info", address)
	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed on unmarshal protocal data")
	}

	err = c.cache.Put([]byte(key), bytes, nil)
	if err != nil {
		return ErrWritePrevInfoInCache
	}

	return nil
}

func (c *Cache) CacheWalletMintTxInfo(address string, txHash string, txHex string) error {
	key := fmt.Sprintf("%s:%s", address, txHash)
	bytes, err := json.Marshal(txHex)
	if err != nil {
		return errors.Wrap(err, "failed on unmarshal protocal data")
	}

	err = c.cache.Put([]byte(key), bytes, nil)
	if err != nil {
		return ErrWritePrevInfoInCache
	}

	return nil
}
