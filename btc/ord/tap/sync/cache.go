package sync

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/traitmeta/metago/btc/ord/tap/model"
)

type Cache struct {
	db *leveldb.DB
}

func NewCache(dir string) (*Cache, error) {
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed on connect leveldb")
	}

	return &Cache{
		db: db,
	}, nil
}

func (c *Cache) SetTickDeployToDetail(deployedTick string, detail TickDetail) error {
	deployedTick = strings.ToLower(deployedTick)
	rawDetail, err := json.Marshal(detail)
	if err != nil {
		return errors.Wrap(err, "failed on marshal deploy tick detail")
	}

	return c.db.Put([]byte(deployedTick), rawDetail, nil)
}

func (c *Cache) GetTickDeployDetail(deployedTick string) (*TickDetail, error) {
	deployedTick = strings.ToLower(deployedTick)
	data, err := c.db.Get([]byte(deployedTick), nil)
	if err != nil {
		return nil, err
	}

	var detail TickDetail
	if err := json.Unmarshal(data, &detail); err != nil {
		return nil, err
	}

	return &detail, nil
}

func (c *Cache) SetTickMintedBlock(ticker, blockHeight string, inscriptionId string) error {
	ticker = strings.ToLower(ticker)
	return c.db.Put([]byte(fmt.Sprintf("%s:%s", ticker, blockHeight)), []byte(inscriptionId), nil)
}

func (c *Cache) GetTickMintedBlock(ticker, blockHeight string) (string, error) {
	ticker = strings.ToLower(ticker)
	data, err := c.db.Get([]byte(fmt.Sprintf("%s:%s", ticker, blockHeight)), nil)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Cache) SetNameToElement(name, element string) error {
	name = strings.ToLower(name)
	return c.db.Put([]byte(name+ElementId), []byte(element), nil)
}

func (c *Cache) GetNameToElement(name string) (string, error) {
	name = strings.ToLower(name)
	data, err := c.db.Get([]byte(name+ElementId), nil)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Cache) SetNoNameToElement(noName, element string) error {
	noName = strings.ToLower(noName)
	return c.db.Put([]byte(noName+ElementId), []byte(element), nil)
}

func (c *Cache) GetNoNameToElement(noName string) (string, error) {
	noName = strings.ToLower(noName)
	data, err := c.db.Get([]byte(noName+ElementId), nil)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Cache) GetElementById(elementId string) (string, error) {
	elementId = strings.ToLower(elementId)
	data, err := c.db.Get([]byte(elementId+ElementId), nil)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// BatchSet
// 1. process Cache in last step
// 2. record ths block's hash and all cache keys
// 3. batch write Cache
// 4. if fork is happened, remove all keys in this block height
func (c *Cache) BatchSet(blockHeight int64, elements []model.TapElement, mintActivities []model.TapActivity, deployTicks []model.TapElementTick) error {
	batch := new(leveldb.Batch)
	var cachedKeys CachedKeys
	for _, element := range elements {
		nameKey := strings.ToLower(element.Name + ElementId)
		cachedKeys.NameKeys = append(cachedKeys.NameKeys, nameKey)
		noNameKey := strings.ToLower(ElementNoName(element.Pattern, element.Field) + ElementId)
		cachedKeys.NoNameKeys = append(cachedKeys.NoNameKeys, noNameKey)
		batch.Put([]byte(nameKey), []byte(element.Element))
		batch.Put([]byte(noNameKey), []byte(element.Element))
	}

	for _, deploy := range deployTicks {
		deployKey := strings.ToLower(deploy.Tick)
		cachedKeys.DeployKeys = append(cachedKeys.DeployKeys, deployKey)
		rawDetail, err := json.Marshal(TickDetail{
			ElementInscriptionId: deploy.ElementInscriptionId,
			TickInscriptionId:    deploy.TickInscriptionId,
			InscriptionHeight:    deploy.InscriptionHeight,
		})
		if err != nil {
			return errors.Wrap(err, "failed on marshal deploy tick detail")
		}

		batch.Put([]byte(deployKey), rawDetail)
	}

	for _, mint := range mintActivities {
		mintKey := fmt.Sprintf("%s:%s", strings.ToLower(mint.Tick), mint.BlockNumber)
		cachedKeys.MintKeys = append(cachedKeys.MintKeys, mintKey)

		batch.Put([]byte(mintKey), []byte(mint.InscriptionId))
	}

	tapCacheKey := fmt.Sprintf("%s:%d", TapCacheId, blockHeight)
	allCachedKeys, err := json.Marshal(cachedKeys)
	if err != nil {
		return errors.Wrap(err, "failed on marshal allCachedKeys")
	}

	batch.Put([]byte(tapCacheKey), allCachedKeys)
	if err := c.db.Write(batch, nil); err != nil {
		return errors.Wrap(err, "failed on cache tap")
	}

	return nil
}

func (c *Cache) BatchRollBack(blockHeight int64) error {
	tapCacheKey := fmt.Sprintf("%s:%d", TapCacheId, blockHeight)

	get, err := c.db.Get([]byte(tapCacheKey), nil)
	if err != nil {
		return err
	}

	var cachedKeys CachedKeys
	err = json.Unmarshal(get, &cachedKeys)
	if err != nil {
		return errors.Wrap(err, "failed on marshal allCachedKeys")
	}

	batch := new(leveldb.Batch)
	for _, key := range cachedKeys.NameKeys {
		batch.Delete([]byte(key))
	}

	for _, key := range cachedKeys.NoNameKeys {
		batch.Delete([]byte(key))
	}

	for _, key := range cachedKeys.DeployKeys {
		batch.Delete([]byte(key))
	}
	for _, key := range cachedKeys.MintKeys {
		batch.Delete([]byte(key))
	}

	if err := c.db.Write(batch, nil); err != nil {
		return errors.Wrap(err, "failed on rollback tap")
	}

	return nil
}
