package models

import (
	"math/big"
	"time"

	"gorm.io/gorm"
)

type Block struct {
	*gorm.Model
	Consensus         bool      `json:"consensus" gorm:"column:consensus; default:0; comment:区块共识;"`
	Difficulty        *big.Int  `json:"difficulty" gorm:"column:difficulty; default:0; comment:难度;"`
	BlockHeight       uint64    `json:"block_height" gorm:"column:block_height; default:0; comment:区块高度;"`
	BlockHash         string    `json:"block_hash" gorm:"column:block_hash;default:''; comment:区块hash;"`
	ParentHash        string    `json:"parent_hash" gorm:"column:parent_hash;default:''; comment:父hash;"`
	GasLimit          uint64    `json:"gas_limit" gorm:"column:gas_limit;default:''; comment:区块Gas上限;"`
	GasUsed           uint64    `json:"gas_used" gorm:"column:gas_used;default:''; comment:区块Gas使用量;"`
	MinerHash         string    `json:"miner_hash" gorm:"column:miner_hash;default:''; comment:区块矿工地址;"`
	Nonce             string    `json:"nonce" gorm:"column:nonce;default:''; comment:区块Nonce;"`
	Size              int32     `json:"size" gorm:"column:size;default:''; comment:区块大小;"`
	Timestamp         time.Time `json:"timestamp" gorm:"column:timestamp;default:''; comment:区块时间;"`
	RefetchNeeded     bool      `json:"refetch_needed" gorm:"column:refetch_needed; default:0; comment:;"`
	BaseFeePerGas     *big.Int  `json:"base_fee_per_gas" gorm:"column:base_fee_per_gas; default:0; comment:基础费用;"`
	LatestBlockHeight uint64    `json:"latest_block_height" gorm:"column:latest_block_height;default: 0; comment:最后区块高度;"`
}

func (b *Block) TableName() string {
	return "blocks"
}
