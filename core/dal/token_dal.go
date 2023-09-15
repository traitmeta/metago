package dal

import (
	"context"

	"github.com/traitmeta/gotos/lib/db"
	"github.com/traitmeta/metago/core/models"
)

var Token *tokenDal

type tokenDal struct{}

func InitTokenDal() {
	Token = &tokenDal{}
}

func (b *tokenDal) Insert(ctx context.Context, block models.Block) error {
	if err := db.DBEngine.WithContext(ctx).Create(&block).Error; err != nil {
		return err
	}
	return nil
}

func (e *tokenDal) GetByContractHash(contractHash string) (*models.Token, error) {
	var token models.Token
	if err := db.DBEngine.Where("contract_address = ?", contractHash).Take(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}
