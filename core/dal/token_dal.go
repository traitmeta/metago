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

func (b *tokenDal) Insert(ctx context.Context, token models.Token) error {
	if err := db.DBEngine.WithContext(ctx).Create(&token).Error; err != nil {
		return err
	}

	return nil
}

func (b *tokenDal) UpdateSkipMetadata(ctx context.Context, token models.Token) error {
	token.SkipMetadata = true
	if err := db.DBEngine.WithContext(ctx).Updates(&token).Error; err != nil {
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
