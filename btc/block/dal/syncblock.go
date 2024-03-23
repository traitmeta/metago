package dal

import (
	"gorm.io/gorm"

	"github.com/traitmeta/metago/btc/block/model"
)

type Dal struct {
	DB *gorm.DB
}

func NewDal(db *gorm.DB) *Dal {
	return &Dal{DB: db}
}

func (d *Dal) GetSyncBlockByName(name string) (*model.SyncBlock, error) {
	var resp = &model.SyncBlock{}
	err := d.DB.Model(model.SyncBlock{}).Where("name = ?", name).Take(&resp).Error
	return resp, err
}

func (d *Dal) UpdateBlockByName(name, blockHash string, blockHeight int64) error {
	var update = &model.SyncBlock{
		BlockHeight: blockHeight,
		BlockHash:   blockHash,
	}
	return d.DB.Model(model.SyncBlock{}).Where("name = ?", name).Updates(&update).Error
}
