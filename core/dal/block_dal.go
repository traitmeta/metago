package dal

import (
	"context"

	"github.com/traitmeta/gotos/lib/db"
	"github.com/traitmeta/metago/core/common"
	"github.com/traitmeta/metago/core/models"
)

var Block *blockDal

type blockDal struct{}

func InitBlockDal() {
	Block = &blockDal{}
}

func (b *blockDal) Insert(ctx context.Context, block models.Block) error {
	if err := db.DBEngine.WithContext(ctx).Create(&block).Error; err != nil {
		return err
	}
	return nil
}

func (t *blockDal) Inserts(ctx context.Context, blocks []models.Block) error {
	if err := db.DBEngine.WithContext(ctx).CreateInBatches(blocks, common.BatchSize).Error; err != nil {
		return err
	}

	return nil
}

func (b *blockDal) Counts() (int64, error) {
	var count int64
	if err := db.DBEngine.Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func (b *blockDal) GetLatest() (*models.Block, error) {
	var block *models.Block
	if err := db.DBEngine.Last(block).Error; err != nil {
		return block, err
	}
	return block, nil
}
