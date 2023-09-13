package dal

import (
	"context"

	"github.com/traitmeta/gotos/lib/db"
	"github.com/traitmeta/metago/core/common"
	"github.com/traitmeta/metago/core/models"
)

var Event *eventDal

type eventDal struct{}

func InitEventDal() {
	Event = &eventDal{}
}

func (b *eventDal) Insert(ctx context.Context, event models.Event) error {
	if err := db.DBEngine.WithContext(ctx).Create(&event).Error; err != nil {
		return err
	}
	return nil
}

func (t *eventDal) Inserts(ctx context.Context, events []models.Event) error {
	if err := db.DBEngine.WithContext(ctx).CreateInBatches(events, common.BatchSize).Error; err != nil {
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
