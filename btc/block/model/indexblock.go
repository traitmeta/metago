package model

import "gorm.io/gorm"

type SyncBlock struct {
	Id          int64  `json:"id" gorm:"primaryKey;autoIncrement;column:id;comment:主键"`
	Name        string `json:"name" gorm:"column:token_type;type:varchar(16);not null"`
	BlockHeight int64  `json:"block_height" gorm:"column:block_height;type:bigint(20);not null"`
	BlockHash   string `json:"block_hash" gorm:"column:block_hash;type:varchar(64);not null"`
	gorm.Model
}

func (*SyncBlock) TableName() string {
	return "sync_block"
}
