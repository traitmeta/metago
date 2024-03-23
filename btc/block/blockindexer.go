package block

import (
	"context"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/traitmeta/metago/btc/block/dal"
	"github.com/traitmeta/metago/btc/block/model"
	"github.com/traitmeta/metago/btc/sync"
)

const TapBlock = "block"
const RollBackBlockNumber = 3

type Indexer struct {
	ctx    context.Context
	dal    *dal.Dal
	client *rpcclient.Client
	*sync.BaseSync
}

func New(ctx context.Context, client *rpcclient.Client, db *gorm.DB) *Indexer {
	return &Indexer{
		ctx:      ctx,
		dal:      dal.NewDal(db),
		client:   client,
		BaseSync: sync.NewBaseSync(),
	}
}

func (s *Indexer) Start() {
	s.Send()
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

func (s *Indexer) SyncBlock() error {
	indexBlock, err := s.dal.GetSyncBlockByName(TapBlock)
	if err != nil {
		return err
	}

	blockHeight := indexBlock.BlockHeight + 1
	start := time.Now()
	blockHash, err := s.client.GetBlockHash(blockHeight)
	if err != nil {
		return err
	}

	start = time.Now()
	blockHeader, err := s.client.GetBlockHeader(blockHash)
	if err != nil {
		return err
	}

	log.WithContext(s.ctx).WithField("time_spend", time.Since(start)).Info("SyncBlock get block hash end")
	// check fork
	lastBlock, err := s.dal.GetSyncBlockByName(TapBlock)
	if err != nil {
		return errors.New("cannot found block height info in db")
	}

	if lastBlock.BlockHeight > 0 && lastBlock.BlockHash != blockHeader.PrevBlock.String() {
		// rollback
		if err := s.RollBack(blockHeight, RollBackBlockNumber); err != nil {
			return err
		}
		blockHeight -= RollBackBlockNumber // 回滚后重新检查这个区块
		return nil
	}

	blockBits := model.BlockInfo{
		BlockNumber: blockHeight,
		BlockHash:   blockHash.String(),
		Bits:        strconv.FormatUint(uint64(blockHeader.Bits), 16),
	}

	if err := s.dal.UpsertBlockInfo(&blockBits); err != nil {
		return err
	}

	if err := s.dal.UpdateBlockByName(TapBlock, blockHash.String(), blockHeight); err != nil {
		return err
	}

	return nil
}

func (s *Indexer) RollBack(blockHeight int64, backNumber int64) error {
	rollbackNumber := blockHeight - backNumber
	blockHash, err := s.client.GetBlockHash(rollbackNumber)
	if err != nil {
		return err
	}

	if err := s.dal.UpdateBlockByName(TapBlock, blockHash.String(), rollbackNumber); err != nil {
		return err
	}
	return nil
}
