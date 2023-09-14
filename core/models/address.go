package models

import (
	"gorm.io/gorm"
)

type Address struct {
	*gorm.Model

	FetchedCoinBalance            string `json:"fetched_coin_balance,omitempty" gorm:"fetched_coin_balance;comment:获取币种的金额;"`
	FetchedCoinBalanceBlockNumber int32  `json:"fetched_coin_balance_block_number,omitempty" gorm:"fetched_coin_balance_block_number; comment:获取币种金额的区块高度;"`
	Hash                          string `json:"hash,omitempty" gorm:"hash; comment:地址HEX;"`
	ContractCode                  string `json:"contract_code,omitempty" gorm:"contract_code; comment:地址是合约时的合约代码;"`
	Nonce                         int32  `json:"nonce,omitempty" gorm:"nonce; comment:地址Nonce;"`
	Decompiled                    bool   `json:"decompiled,omitempty" gorm:"decompiled; comment:合约是否已经反编译;"`
	Verified                      bool   `json:"verified,omitempty" gorm:"verified; comment:合约是否已经验证;"`
	GasUsed                       int32  `json:"gas_used,omitempty" gorm:"gas_used; comment:GAS使用量;"`
	TransactionsCount             int32  `json:"transactions_count,omitempty" gorm:"transactions_count; comment:总交易数量;"`
	TokenTransfersCount           int32  `json:"token_transfers_count,omitempty" gorm:"token_transfers_count; comment:总TOKEN转移数量;"`
}

func (b *Address) TableName() string {
	return "addresses"
}
