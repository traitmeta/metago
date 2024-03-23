package model

import "gorm.io/gorm"

const (
	DeployType   = 1
	MintType     = 2
	TransferType = 3
)

type TapActivity struct {
	Id                   int64  `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	ElementInscriptionId string `json:"element_inscription_id" gorm:"column:element_inscription_id;default:NULL"`
	DeployInscriptionId  string `json:"deploy_inscription_id" gorm:"column:deploy_inscription_id;default:NULL"`
	Type                 int8   `json:"type" gorm:"column:type;default:NULL"`
	Tick                 string `json:"tick" gorm:"column:tick;default:NULL"`
	Body                 string `json:"body" gorm:"column:body"`
	BlockNumber          string `json:"block_number" gorm:"column:block_number;default:NULL"`
	InscriptionHeight    int64  `json:"inscription_height" gorm:"column:inscription_height;"`
	InscriptionId        string `json:"inscription_id" gorm:"column:inscription_id;"`
	gorm.Model
}

func (TapActivity) TableName() string {
	return "ord_tap_activity"
}
