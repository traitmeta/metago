package dal

import (
	"github.com/traitmeta/metago/core/common"
	"github.com/traitmeta/metago/core/models"
	"github.com/traitmeta/metago/pkg/db"
)

var Transaction *transactionDal

type transactionDal struct{}

func InitTransactionDal() {
	Transaction = &transactionDal{}
}

func (b *transactionDal) Insert(tranaction models.Transaction) error {
	if err := db.DBEngine.Create(&tranaction).Error; err != nil {
		return err
	}
	return nil
}

func (t *transactionDal) Inserts(tranactions []models.Transaction) error {
	if err := db.DBEngine.CreateInBatches(tranactions, common.BatchSize).Error; err != nil {
		return err
	}

	return nil
}
