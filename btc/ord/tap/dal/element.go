package dal

import (
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"

	"github.com/traitmeta/metago/btc/ord/tap/model"
)

func (d *Dal) GetElements() ([]model.TapElement, error) {
	var tapElems []model.TapElement
	if err := d.DB.Model(model.TapElement{}).
		Select("name,element,element_inscription_id,field,pattern").
		Scan(&tapElems).Error; err != nil {
		return nil, errors.Wrap(err, "failed on query tap element")
	}

	return tapElems, nil
}

func (d *Dal) GetElementsBits() ([]model.TapElement, error) {
	var tapElems []model.TapElement
	if err := d.DB.Model(model.TapElement{}).
		Select("name,element,element_inscription_id,field,pattern").
		Where("field = 11").
		Scan(&tapElems).Error; err != nil {
		return nil, errors.Wrap(err, "failed on query tap element")
	}

	return tapElems, nil
}

func (d *Dal) UpsertElements(elems []model.TapElement) error {
	if len(elems) <= 0 {
		return nil
	}

	for i := 0; i < len(elems); i += 100 {
		end := i + 100
		if end > len(elems) {
			end = len(elems)
		}
		insertData := elems[i:end]
		if err := d.DB.Model(model.TapElement{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tick"}},
			DoUpdates: clause.AssignmentColumns([]string{"inscription_height"}),
		}).Create(&insertData).Error; err != nil {
			return err
		}
	}

	return nil
}
