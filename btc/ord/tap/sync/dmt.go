package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/traitmeta/metago/btc/ord/common"
	"github.com/traitmeta/metago/btc/ord/tap/dal"
	"github.com/traitmeta/metago/btc/ord/tap/model"
)

type DMTIndexer struct {
	ctx       context.Context
	dao       *dal.Dal
	rds       *redis.Client
	client    *rpcclient.Client
	processor *DMTProcessor
	cache     *Cache
	BaseSync
}

func NewDMTIndexer(ctx context.Context, client *rpcclient.Client, db *gorm.DB, cache *Cache, rds *redis.Client) *DMTIndexer {
	return &DMTIndexer{
		ctx:       ctx,
		dao:       dal.NewDal(db),
		rds:       rds,
		client:    client,
		processor: NewDMTProcessor(ctx, db, cache),
		cache:     cache,
		BaseSync: BaseSync{
			ch: make(chan bool, 1),
		},
	}
}

func (s *DMTIndexer) Start() {
	s.ch <- true
	for {
		select {
		case <-s.Receive():
			err := s.SyncBlock()
			if err != nil {
				log.WithContext(s.ctx).WithField("err", err).Error("SyncBlock failed on query block table")
				time.Sleep(10 * time.Second)
				s.Send()
			}
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *DMTIndexer) SyncBlock() error {
	indexBlock, err := s.dao.GetSyncBlockByName(common.Tap)
	if err != nil {
		return err
	}

	blockHeight := indexBlock.BlockHeight + 1
	blockHash, err := s.client.GetBlockHash(blockHeight)
	if err != nil {
		return err
	}

	block, err := s.client.GetBlock(blockHash)
	if err != nil {
		return err
	}

	// 检查分叉
	lastBlock, err := s.dao.GetSyncBlockByName(common.Tap)
	if err != nil {
		return errors.New("not found block height in db")
	}

	if lastBlock.BlockHeight > 0 && lastBlock.BlockHash != block.Header.PrevBlock.String() {
		// 发生了分叉，执行回滚操作
		if err := s.RollBack(blockHeight, common.RollBackBlockNumber); err != nil {
			return err
		}
		blockHeight -= common.RollBackBlockNumber // 回滚后重新检查这个区块
		return nil
	}

	if err := s.processBlock(blockHeight, block); err != nil {
		return err
	}

	return nil
}

func (s *DMTIndexer) processBlock(blockHeight int64, block *wire.MsgBlock) (err error) {
	var elements []model.TapElement
	var deployActivities []model.TapActivity
	var mintActivities []model.TapActivity
	var deployTicks []model.TapElementTick
	for txIdx := 0; txIdx < len(block.Transactions); txIdx++ {
		tx := block.Transactions[txIdx]
		elems, deployActs, mintActs, ticks, err := s.HandleTx(blockHeight, block.Header.Timestamp.Unix(), tx)
		if err != nil {
			return err
		}

		elements = append(elements, elems...)
		deployActivities = append(deployActivities, deployActs...)
		mintActivities = append(mintActivities, mintActs...)
		deployTicks = append(deployTicks, ticks...)
	}

	if len(elements) > 0 {
		elements = s.processor.FilterValidElement(elements)
	}

	if len(deployTicks) > 0 {
		deployTicks, deployActivities, err = s.processor.FilterValidDmtDeploy(deployTicks, deployActivities)
		if err != nil {
			return errors.Wrap(err, "FilterValidDmtDeploy")
		}
	}

	if len(mintActivities) > 0 {
		mintActivities, err = s.processor.FilterValidDmtMint(mintActivities, deployTicks)
		if err != nil {
			return errors.Wrap(err, "FilterValidDmtMint")
		}
	}
	var validActivities []model.TapActivity
	if len(deployActivities) > 0 {
		validActivities = append(validActivities, deployActivities...)
	}

	if len(mintActivities) > 0 {
		validActivities = append(validActivities, mintActivities...)
	}

	if len(validActivities) == 0 && len(elements) == 0 {
		return s.dao.UpdateBlockByName(common.Tap, block.BlockHash().String(), blockHeight)
	}

	err = s.cache.BatchSet(blockHeight, elements, mintActivities, deployTicks)
	if err != nil {
		return errors.Wrap(err, "DMTIndexer cache BatchSet")
	}

	err = s.DBHandle(blockHeight, block.BlockHash().String(), elements, validActivities, deployTicks)
	if err != nil {
		return s.cache.BatchRollBack(blockHeight)
	}

	s.HandleElementsCache(blockHeight, strconv.FormatUint(uint64(block.Header.Bits), 16), elements)
	s.HandleTickElementMintCache(mintActivities)
	return nil
}

func (s *DMTIndexer) DBHandle(blockHeight int64, blockHash string, elements []model.TapElement, validActivities []model.TapActivity, deployTicks []model.TapElementTick) error {
	dbTx := s.dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			dbTx.Rollback()
		}
	}()

	dao := dal.NewDal(dbTx)
	if err := dao.UpsertElements(elements); err != nil {
		dbTx.Rollback()
		return errors.Wrap(err, "DMTIndexer UpsertElements")
	}

	if err := dao.UpsertActivities(validActivities); err != nil {
		dbTx.Rollback()
		return errors.Wrap(err, "DMTIndexer UpsertActivities")
	}

	if err := dao.UpsertElementTick(deployTicks); err != nil {
		dbTx.Rollback()
		return errors.Wrap(err, "DMTIndexer UpsertElementTick")
	}

	if err := dao.UpdateBlockByName(common.Tap, blockHash, blockHeight); err != nil {
		dbTx.Rollback()
		return errors.Wrap(err, "DMTIndexer UpdateBlockIndex")
	}

	if err := dbTx.Commit().Error; err != nil {
		dbTx.Rollback()
		return errors.Wrap(err, "DMTIndexer Commit")
	}

	return nil
}

func (s *DMTIndexer) HandleElementsCache(blockHeight int64, blockBits string, elements []model.TapElement) {
	for _, ele := range elements {
		if strings.Contains(blockBits, ele.Pattern) {
			err := s.rds.SAdd(s.ctx, common.CacheKey(fmt.Sprintf(common.TapElement, ele.ElementInscriptionId)), blockHeight).Err()
			if err != nil {
				log.WithContext(s.ctx).WithFields(
					log.Fields{
						"block_height": blockHeight,
						"block_bits":   blockBits,
						"element":      ele.ElementInscriptionId,
					},
				).Warn("DMTIndexer HandleCache")
			}
		}
	}
}

func (s *DMTIndexer) HandleTickElementMintCache(mintActivities []model.TapActivity) {
	// activity 归类
	var mintMap = make(map[string][]interface{})
	for _, mint := range mintActivities {
		number, err := strconv.ParseInt(mint.BlockNumber, 10, 64)
		if err != nil {
			continue
		}

		key := fmt.Sprintf("%s:%s", mint.ElementInscriptionId, strings.ToLower(mint.Tick))
		if v, ok := mintMap[key]; ok {
			v = append(v, number)
			mintMap[key] = v
		} else {
			mintMap[key] = []interface{}{number}
		}
	}

	for k, v := range mintMap {
		err := s.rds.SAdd(s.ctx, common.CacheKey(fmt.Sprintf(common.TapElement, k)), v...)
		if err != nil {
			log.WithContext(s.ctx).WithFields(
				log.Fields{
					"key": k,
				},
			).Warn("DMTIndexer HandleCache")
		}
	}
}

func (s *DMTIndexer) HandleTx(blockHeight, blockTime int64, tx *wire.MsgTx) ([]model.TapElement, []model.TapActivity, []model.TapActivity, []model.TapElementTick, error) {
	if len(tx.TxOut) == 0 {
		log.WithContext(s.ctx).WithField("tx_hash", tx.TxHash().String()).Warn("tx out number 0")
		return nil, nil, nil, nil, nil
	}

	envelopes := FromTransaction(tx)
	if len(envelopes) == 0 {
		return nil, nil, nil, nil, nil
	}

	elements := s.processor.ProcessElement(blockHeight, tx.TxHash().String(), envelopes)
	deployTicks, deployActivities, err := s.processor.ProcessDeploy(blockHeight, blockTime, tx.TxHash().String(), envelopes)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	mintActivities, err := s.processor.ProcessMint(blockHeight, tx.TxHash().String(), envelopes)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return elements, deployActivities, mintActivities, deployTicks, nil
}

func MatchElementPattern(content string) bool {
	re := regexp.MustCompile(common.ElementPattern)
	return re.MatchString(content)
}

func (s *DMTIndexer) handlerDmtInscription(blockHeight int64, txId string, envelopes []Envelope) ([]model.TapActivity, []model.TapElementTick, error) {
	var activities []model.TapActivity
	deployTickToElementTick := make(map[string]model.TapElementTick)
	mintTickBlockToBool := make(map[string]bool)
	for _, envelope := range envelopes {
		insData := envelope.ConvertToInscriptionData()
		dmtOpr := &DmtOpr{}
		err := json.Unmarshal(insData.Body, dmtOpr)
		if err != nil {
			continue
		}

		if !strings.EqualFold(dmtOpr.Protocol, common.TapProtocol) {
			continue
		}
		body := ""
		content := envelope.GetContent()
		if content != nil {
			body = string(content)
			log.WithContext(s.ctx).WithField("body", body).Info("DMTIndexer handlerDmtInscription")
		}

		switch dmtOpr.Operation {
		case common.DmtDeploy:
			// 817708 00000000000000000001efa12ef8cf17d8521b8ba3de54433f2c14f3f929224d
			// Check deploy是否合法
			// 1. 缓存中存在表示本次deploy 无效
			if _, ok := deployTickToElementTick[dmtOpr.Ticker]; ok {
				continue
			}

			dbElemTick, err := s.dao.GetElementTickByTick(dmtOpr.Ticker, blockHeight)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil, err
			}
			// 2. DB中存在比当前高度低的数据表示本次deploy 无效
			if dbElemTick != nil {
				deployTickToElementTick[dbElemTick.Tick] = *dbElemTick
				continue
			}

			// 新的并且合法的Deploy 记录到缓存中
			elemTick := model.TapElementTick{
				ElementInscriptionId: dmtOpr.Element,
				Tick:                 dmtOpr.Ticker,
				TickInscriptionId:    fmt.Sprintf("%si%d", txId, envelope.Offset),
				InscriptionHeight:    blockHeight,
			}
			deployTickToElementTick[dmtOpr.Ticker] = elemTick
			activities = append(activities, model.TapActivity{
				ElementInscriptionId: dmtOpr.Element,
				Type:                 model.DeployType,
				Tick:                 dmtOpr.Ticker,
				Body:                 body,
				InscriptionHeight:    blockHeight,
				InscriptionId:        fmt.Sprintf("%si%d", txId, envelope.Offset),
			})
		case common.DmtMint:
			// Check mint是否合法, 补全ElementInscriptionId
			// 缓存中存在表示本次mint 无效
			if _, ok := mintTickBlockToBool[fmt.Sprintf("%s:%s", dmtOpr.Ticker, dmtOpr.Block)]; ok {
				continue
			} else {
				count, err := s.dao.CountMintActivityWithBlock(dmtOpr.Ticker, dmtOpr.Block, blockHeight)
				if err != nil {
					return nil, nil, err
				}
				// 数据库存在 mint无效
				if count > 0 {
					continue
				}
				// 数据库不存在本次有效，记录到缓存中
				mintTickBlockToBool[fmt.Sprintf("%s:%s", dmtOpr.Ticker, dmtOpr.Block)] = true
			}

			elementInscriptionId := ""
			// 这笔交易缓存中存在deploy的信息，则更新mint
			if v, ok := deployTickToElementTick[dmtOpr.Ticker]; ok {
				elementInscriptionId = v.ElementInscriptionId
				v.Minted += 1
			} else {
				// 缓存没有，db没有，那么mint无效
				dbElemTick, err := s.dao.GetElementTick(dmtOpr.Deploy, blockHeight)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					log.WithContext(s.ctx).
						WithFields(log.Fields{"deploy": dmtOpr.Deploy, "block_height": blockHeight}).
						Warn("DMTIndexer handlerDmtInscription")
					return nil, nil, err
				}
				if dbElemTick == nil {
					continue
				}
				// 缓存中没有，DB中有，则加入缓存，更新mint
				if dbElemTick != nil {
					elementInscriptionId = dbElemTick.ElementInscriptionId
					dbElemTick.Minted += 1
					deployTickToElementTick[dbElemTick.Tick] = *dbElemTick
				}
			}

			activities = append(activities, model.TapActivity{
				ElementInscriptionId: elementInscriptionId,
				Type:                 model.MintType,
				Tick:                 dmtOpr.Ticker,
				Body:                 body,
				BlockNumber:          dmtOpr.Block,
				InscriptionHeight:    blockHeight,
				InscriptionId:        fmt.Sprintf("%si%d", txId, envelope.Offset),
			})
		case common.DmtTransfer:
		default:
			continue
		}
	}

	var deployTicks []model.TapElementTick
	for _, v := range deployTickToElementTick {
		deployTicks = append(deployTicks, v)
	}

	return activities, deployTicks, nil
}

func (s *DMTIndexer) RollBack(blockHeight int64, rollbackNumber int) error {
	for i := 0; i < rollbackNumber; i++ {
		err := s.cache.BatchRollBack(blockHeight - int64(i))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("DMTIndexer BatchRollBack: %d", blockHeight-int64(i)))
		}
	}

	rollbackBlockNumber := blockHeight - int64(rollbackNumber)
	log.WithContext(s.ctx).Info("DMTIndexer Rolling Back Block", zap.Int64("block_height", rollbackBlockNumber))
	blockHash, err := s.client.GetBlockHash(rollbackBlockNumber)
	if err != nil {
		log.WithContext(s.ctx).WithField("block_height", rollbackBlockNumber).Warn("DMTIndexer Rolling Back Block failed on get block hash:%v", err)
		return err
	}

	if err := s.dao.UpdateBlockByName(common.Tap, blockHash.String(), rollbackBlockNumber); err != nil {
		log.WithContext(s.ctx).Errorf("DMTIndexer Rolling Back Block failed on update block tap init index block number:%v", err)
		return err
	}

	return nil
}
