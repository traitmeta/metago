package models

import (
	"math/big"

	"gorm.io/gorm"
)

type Token struct {
	*gorm.Model

	Name                      string `json:"name,omitempty" gorm:"name; comment:Token名称;"`
	Symbol                    string `json:"symbol,omitempty" gorm:"symbol; comment:Token 缩写;"`
	TotalSupply               string `json:"total_supply,omitempty" gorm:"total_supply; comment:总供应量;"`
	Decimals                  string `json:"decimals,omitempty" gorm:"decimals; comment:精度;"`
	Type                      string `json:"type,omitempty" gorm:"type; comment:类型 ERC20, ERC721, ERC1155;"`
	Cataloged                 bool   `json:"cataloged,omitempty" gorm:"cataloged; comment:归类;"`
	ContractAddress           string `json:"contract_address,omitempty" gorm:"contract_address; comment:合约地址;"`
	HolderCount               int32  `json:"holder_count,omitempty" gorm:"holder_count; comment:持有用户量;"`
	SkipMetadata              bool   `json:"skip_metadata,omitempty" gorm:"skip_metadata; comment:跳过元信息;"`
	FiatValue                 string `json:"fiat_value,omitempty" gorm:"fiat_value; comment:假值?;"`
	CirculatingMarketCap      string `json:"circulating_market_cap,omitempty" gorm:"circulating_market_cap; comment:市场流通量;"`
	TotalSupplyUpdatedAtBlock int32  `json:"total_supply_updated_at_block,omitempty" gorm:"total_supply_updated_at_block; comment:总供应量更新的区块高度;"`
	IconUrl                   string `json:"icon_url,omitempty" gorm:"icon_url; comment:图标URL;"`
	IsVerifiedViaAdminPanel   bool   `json:"is_verified_via_admin_panel,omitempty" gorm:"is_verified_via_admin_panel; comment:是否已经验证;"`
}

func (b *Token) TableName() string {
	return "tokens"
}

type TokenTransfer struct {
	*gorm.Model

	TransactionHash      string     `json:"transaction_hash,omitempty" gorm:"transaction_hash; comment:交易哈希;"`                  // Transaction foreign key
	LogIndex             uint       `json:"log_index,omitempty" gorm:"log_index; comment:日志索引;"`                                // Index of the corresponding `Event` in the transaction.
	FromAddress          string     `json:"from_address,omitempty" gorm:"from_address; comment:From;"`                          // sender
	ToAddress            string     `json:"to_address,omitempty" gorm:"to_address; comment:To;"`                                // receiver
	Amount               *big.Int   `json:"amount,omitempty" gorm:"amount; comment:Token 数量;"`                                  // amount
	TokenId              *big.Int   `json:"token_id,omitempty" gorm:"token_id; comment:Token ID 编号;"`                           // ID of the token (applicable to ERC-721 tokens)
	TokenContractAddress string     `json:"token_contract_address,omitempty" gorm:"token_contract_address; comment:Token合约地址;"` // token contract address
	BlockNumber          uint64     `json:"block_number,omitempty" gorm:"block_number; comment:区块高度;"`                          // block number
	BlockHash            string     `json:"block_hash,omitempty" gorm:"block_hash; comment:区块哈希;"`                              // block hash
	Amounts              []*big.Int `json:"amounts,omitempty" gorm:"amounts; comment:批量操作金额;"`                                  // Tokens transferred amounts in case of batched transfer in ERC-1155
	TokenIds             []*big.Int `json:"token_ids,omitempty" gorm:"token_ids; comment:批量操作Token ID;"`                        // IDs of the tokens (applicable to ERC-1155 tokens)
}

func (b *TokenTransfer) TableName() string {
	return "token_transfers"
}
