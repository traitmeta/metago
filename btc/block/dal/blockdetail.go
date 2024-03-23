package dal

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/traitmeta/metago/btc/block/model"
)

func (d *Dal) UpsertBlockInfo(data *model.BlockInfo) error {
	return d.DB.Model(model.BlockInfo{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "block_number"}},
		DoUpdates: clause.AssignmentColumns([]string{"update_time", "block_hash", "bits"}),
	}).Create(&data).Error
}

func (d *Dal) GetElementBlocks(pattern string) ([]model.BlockInfo, error) {
	var results []model.BlockInfo
	var result []model.BlockInfo
	if err := d.DB.Select("id,block_number,bits").
		Where("bits like ?", fmt.Sprintf("%%%s%%", pattern)).
		FindInBatches(&result, 1000, func(tx *gorm.DB, batch int) error {
			results = append(results, result...)
			return nil
		}).Error; err != nil {
		return nil, err
	}

	return results, nil
}
