package models

import (
	"gorm.io/gorm"
)

type Event struct {
	*gorm.Model

	Address     string   `json:"address" gorm:"type:char(42)" `
	Topics      []string `json:"topics" gorm:"type:longtext" `
	Data        string   `json:"data" gorm:"type:longtext" `
	BlockNumber uint64   `json:"block_number"`
	TxHash      string   `json:"tx_hash" gorm:"type:char(66)" `
	TxIndex     uint     `json:"tx_index" `
	BlockHash   string   `json:"block_hash" gorm:"type:varchar(256)" `
	LogIndex    uint     `json:"log_index"`
	Removed     bool     `json:"removed"`
}

func (e *Event) TableName() string {
	return "events"
}
