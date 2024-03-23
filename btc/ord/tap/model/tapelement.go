package model

import "gorm.io/gorm"

type TapElementTick struct {
	Id                   int64  `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	ElementInscriptionId string `json:"element_inscription_id" gorm:"column:element_inscription_id;default:NULL"`
	Tick                 string `json:"tick" gorm:"column:tick;default:NULL"`
	TickInscriptionId    string `json:"tick_inscription_id" gorm:"column:tick_inscription_id;default:NULL"`
	Minted               int64  `json:"minted" gorm:"column:minted;default:NULL"`
	Total                int64  `json:"total" gorm:"column:total;default:NULL"`
	DeployTime           int64  `json:"deploy_time" gorm:"column:deploy_time;type:bigint(20)"`
	InscriptionHeight    int64  `json:"inscription_height" gorm:"column:inscription_height;"`
	gorm.Model
}

func (TapElementTick) TableName() string {
	return "ord_tap_element_tick"
}

type TapElement struct {
	Id                   int64  `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Element              string `json:"element" gorm:"column:element;default:NULL"`
	ElementInscriptionId string `json:"element_inscription_id" gorm:"column:element_inscription_id;default:NULL"`
	Name                 string `json:"name" gorm:"column:name;default:NULL"`
	Pattern              string `json:"pattern" gorm:"column:pattern;default:NULL"`
	Field                string `json:"field" gorm:"column:field;default:NULL"`
	InscriptionHeight    int64  `json:"inscription_height" gorm:"column:inscription_height;"`
	gorm.Model
}

func (TapElement) TableName() string {
	return "ord_tap_element"
}
