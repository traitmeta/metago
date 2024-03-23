package model

import "gorm.io/gorm"

type TapBlockBits struct {
	Id          int64  `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	BlockNumber int64  `json:"block_number" gorm:"column:block_number;default:NULL"`
	BlockHash   string `json:"block_hash" gorm:"column:block_hash;default:NULL"`
	Bits        string `json:"bits" gorm:"column:bits;default:NULL"`
	Version     string `json:"version" gorm:"column:bits;default:NULL"`
	gorm.Model
}

func (TapBlockBits) TableName() string {
	return "ord_tap_block"
}
