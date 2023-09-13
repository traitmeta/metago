package models

import (
	"gorm.io/gorm"
)

type Block struct {
	Id                int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BlockHeight       uint64 `json:"block_height" gorm:"column:block_height; default:0; comment:区块高度;"`
	BlockHash         string `json:"block_hash" gorm:"column:block_hash;default:''; comment:区块hash;"`
	ParentHash        string `json:"parent_hash" gorm:"column:parent_hash;default:''; comment:父hash;"`
	LatestBlockHeight uint64 `json:"latest_block_height" gorm:"column:latest_block_height;default: 0; comment:最后区块高度;"`
	*gorm.Model
}

func (b *Block) TableName() string {
	return "blocks"
}
