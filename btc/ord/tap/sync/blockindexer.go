package sync

import (
	"context"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/traitmeta/metago/btc/ord/common"
	"github.com/traitmeta/metago/btc/ord/tap/dal"
	"github.com/traitmeta/metago/btc/ord/tap/model"
)

type BlockBitsIndexer struct {
	ctx    context.Context
	dal    *dal.Dal
	client *rpcclient.Client
	BaseSync
}

func New(ctx context.Context, client *rpcclient.Client, db *gorm.DB) *BlockBitsIndexer {
	return &BlockBitsIndexer{
		ctx:    ctx,
		dal:    dal.NewDal(db),
		client: client,
		BaseSync: BaseSync{
			ch: make(chan bool, 1),
		},
	}
}

func (s *BlockBitsIndexer) Start() {
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

func (s *BlockBitsIndexer) SyncBlock() error {
	indexBlock, err := s.dal.GetSyncBlockByName(common.TapBlock)
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
	lastBlock, err := s.dal.GetSyncBlockByName(common.TapBlock)
	if err != nil {
		return errors.New("cannot found block height info in db")
	}

	if lastBlock.BlockHeight > 0 && lastBlock.BlockHash != blockHeader.PrevBlock.String() {
		// rollback
		if err := s.RollBack(blockHeight, common.RollBackBlockNumber); err != nil {
			return err
		}
		blockHeight -= common.RollBackBlockNumber - 1 // 回滚后重新检查这个区块
		return nil
	}

	blockBits := model.TapBlockBits{
		BlockNumber: blockHeight,
		BlockHash:   blockHash.String(),
		Bits:        strconv.FormatUint(uint64(blockHeader.Bits), 16),
	}

	if err := s.dal.UpsertBlockInfo(&blockBits); err != nil {
		return err
	}

	if err := s.dal.UpdateBlockByName(common.TapBlock, blockHash.String(), blockHeight); err != nil {
		return err
	}

	return nil
}

func (s *BlockBitsIndexer) RollBack(blockHeight int64, backNumber int64) error {
	rollbackNumber := blockHeight - backNumber
	blockHash, err := s.client.GetBlockHash(rollbackNumber)
	if err != nil {
		return err
	}

	if err := s.dal.UpdateBlockByName(common.TapBlock, blockHash.String(), rollbackNumber); err != nil {
		return err
	}
	return nil
}
