package dal

import (
	"context"

	"github.com/traitmeta/gotos/lib/db"
	"github.com/traitmeta/metago/core/common"
	"github.com/traitmeta/metago/core/models"
)

var Transaction *transactionDal

type transactionDal struct{}

func InitTransactionDal() {
	Transaction = &transactionDal{}
}

func (b *transactionDal) Insert(ctx context.Context, tranaction models.Transaction) error {
	if err := db.DBEngine.WithContext(ctx).Create(&tranaction).Error; err != nil {
		return err
	}
	return nil
}

func (t *transactionDal) Inserts(ctx context.Context, tranactions []models.Transaction) error {
	if err := db.DBEngine.WithContext(ctx).CreateInBatches(tranactions, common.BatchSize).Error; err != nil {
		return err
	}

	return nil
}
