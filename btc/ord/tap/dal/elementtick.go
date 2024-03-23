package dal

import (
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"

	"github.com/traitmeta/metago/btc/ord/tap/model"
)

func (d *Dal) GetsElementTick() ([]model.TapElementTick, error) {
	var tapElems []model.TapElementTick
	if err := d.DB.Model(model.TapElementTick{}).
		Select("tick,element_inscription_id,tick_inscription_id").
		Scan(&tapElems).Error; err != nil {
		return nil, errors.Wrap(err, "failed on query tap element")
	}

	return tapElems, nil
}

func (d *Dal) GetElementTick(tickInsId string, inscriptionHeight int64) (*model.TapElementTick, error) {
	var res *model.TapElementTick
	if err := d.DB.Model(model.TapElementTick{}).
		Where("tick_inscription_id = ? and inscription_height <= ?", tickInsId, inscriptionHeight).
		Take(&res).Error; err != nil {
		return nil, errors.Wrap(err, "failed on query tap tick element")
	}

	return res, nil
}

func (d *Dal) GetElementTickByTick(tick string, inscriptionHeight int64) (*model.TapElementTick, error) {
	var res *model.TapElementTick
	if err := d.DB.Model(model.TapElementTick{}).
		Where("tick = ? and inscription_height <= ?", tick, inscriptionHeight).
		Take(&res).Error; err != nil {
		return nil, errors.Wrap(err, "failed on query tap tick element")
	}

	return res, nil
}

func (d *Dal) UpsertElementTick(elemTicks []model.TapElementTick) error {
	if len(elemTicks) <= 0 {
		return nil
	}

	lens := len(elemTicks)
	for i := 0; i < lens; i += 100 {
		end := i + 100
		if end > lens {
			end = lens
		}
		insertData := elemTicks[i:end]
		if err := d.DB.Model(model.TapElementTick{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tick"}},
			DoUpdates: clause.AssignmentColumns([]string{"minted", "inscription_height"}),
		}).Create(&insertData).Error; err != nil {
			return err
		}
	}

	return nil
}
