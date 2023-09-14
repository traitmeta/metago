package models

import (
	"gorm.io/gorm"
)

type Event struct {
	*gorm.Model

	Address     string `json:"address" gorm:"type:char(42)" `
	FirstTopic  string `json:"first_topic" gorm:"type:longtext; comment:方法签名;" `
	SecondTopic string `json:"second_topic" gorm:"type:longtext; comment:第一个Indexed 参数;" `
	ThirdTopic  string `json:"third_topic" gorm:"type:longtext; comment:第二个Indexed 参数;" `
	FourthTopic string `json:"fourth_topic" gorm:"type:longtext; comment:第三个Indexed 参数;" `
	Data        string `json:"data" gorm:"type:longtext" `
	BlockNumber uint64 `json:"block_number"`
	TxHash      string `json:"tx_hash" gorm:"type:char(66)" `
	TxIndex     uint   `json:"tx_index" `
	BlockHash   string `json:"block_hash" gorm:"type:varchar(256)" `
	LogIndex    uint   `json:"log_index"`
	Removed     bool   `json:"removed"`
}

func (e *Event) TableName() string {
	return "events"
}
