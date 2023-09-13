package dal

import (
	"github.com/traitmeta/metago/core/common"
	"github.com/traitmeta/metago/core/models"
	"github.com/traitmeta/metago/pkg/db"
)

var Event *eventDal

type eventDal struct{}

func InitEventDal() {
	Event = &eventDal{}
}

func (b *eventDal) Insert(event models.Event) error {
	if err := db.DBEngine.Create(&event).Error; err != nil {
		return err
	}
	return nil
}

func (t *eventDal) Inserts(events []models.Event) error {
	if err := db.DBEngine.CreateInBatches(events, common.BatchSize).Error; err != nil {
		return err
	}

	return nil
}

func (e *eventDal) GetEventByTxHash(txHash string) (*models.Event, error) {
	var event models.Event
	if err := db.DBEngine.Where("tx_hash = ?", txHash).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}
