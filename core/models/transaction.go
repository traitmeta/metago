package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	Id          uint64 `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BlockNumber uint64 `json:"block_number"`
	TxHash      string `json:"tx_hash" gorm:"type:char(66)" `
	From        string `json:"from" gorm:"type:char(42)" `
	To          string `json:"to" gorm:"type:char(42)" `
	Value       string `json:"value" gorm:"type:varchar(256)" `
	Contract    string `json:"contract" gorm:"type:char(42)" `
	Status      uint64 `json:"status"`
	InputData   string `json:"input_data" gorm:"type:varchar(4096)"`
	*gorm.Model
}

func (tx *Transaction) TableName() string {
	return "transactions"
}
