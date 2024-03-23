package dal

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/traitmeta/metago/btc/ord/tap/model"
)

func (d *Dal) CountMintActivityWithBlock(tick, opBlockNumber string, inscriptionHeight int64) (int64, error) {
	var count int64
	if err := d.DB.Model(model.TapActivity{}).
		Where("tick = ? and block_number = ? and inscription_height <= ?", tick, opBlockNumber, inscriptionHeight).Count(&count).Error; err != nil {
		return count, errors.Wrap(err, "failed on query tap tick element")
	}

	return count, nil
}

func (d *Dal) UpsertActivities(elems []model.TapActivity) error {
	if len(elems) <= 0 {
		return nil
	}

	lens := len(elems)
	for i := 0; i < lens; i += 100 {
		end := i + 100
		if end > lens {
			end = lens
		}
		insertData := elems[i:end]
		if err := d.DB.Model(model.TapActivity{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "inscription_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"update_time", "inscription_height"}),
		}).Create(&insertData).Error; err != nil {
			return err
		}
	}

	return nil
}

func (d *Dal) UpsertTapActivity(activity model.TapActivity) error {
	if err := d.DB.Model(model.TapActivity{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "inscription_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"update_time"}),
	}).Create(&activity).Error; err != nil {
		return errors.Wrap(err, "failed on upsert tap activity")
	}

	return nil
}

func (d *Dal) GetTickElementActivityBlocks(deploy model.TapElementTick) ([]model.TapActivity, error) {
	var results []model.TapActivity
	var result []model.TapActivity
	if err := d.DB.Select("id,block_number").
		// TODO add deploy_inscription_id deploy.TickInscriptionId
		Where("type = 2 and tick = ?", deploy.Tick).
		FindInBatches(&result, 1000, func(tx *gorm.DB, batch int) error {
			results = append(results, result...)
			return nil
		}).Error; err != nil {
		return nil, err
	}

	return results, nil
}
