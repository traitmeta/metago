package brc20

import "database/sql"

type Brc20BlockHashes struct {
	Id          string `gorm:"column:id;NOT NULL"`
	BlockHeight string `gorm:"column:block_height;NOT NULL"`
	BlockHash   string `gorm:"column:block_hash;NOT NULL"`
}

func (b *Brc20BlockHashes) TableName() string {
	return "brc20_block_hashes"
}

type Brc20HistoricBalances struct {
	Id               string         `gorm:"column:id;NOT NULL"`
	Pkscript         string         `gorm:"column:pkscript;NOT NULL"`
	Wallet           sql.NullString `gorm:"column:wallet"`
	Tick             string         `gorm:"column:tick;NOT NULL"`
	OverallBalance   string         `gorm:"column:overall_balance;NOT NULL"`
	AvailableBalance string         `gorm:"column:available_balance;NOT NULL"`
	BlockHeight      string         `gorm:"column:block_height;NOT NULL"`
	EventId          string         `gorm:"column:event_id;NOT NULL"`
}

func (b *Brc20HistoricBalances) TableName() string {
	return "brc20_historic_balances"
}

type Brc20Events struct {
	Id            string `gorm:"column:id;NOT NULL"`
	EventType     string `gorm:"column:event_type;NOT NULL"`
	BlockHeight   string `gorm:"column:block_height;NOT NULL"`
	InscriptionId string `gorm:"column:inscription_id;NOT NULL"`
	Event         string `gorm:"column:event;NOT NULL"`
}

func (b *Brc20Events) TableName() string {
	return "brc20_events"
}

type Brc20Tickers struct {
	Id                  string `gorm:"column:id;NOT NULL"`
	OriginalTick        string `gorm:"column:original_tick;NOT NULL"`
	Tick                string `gorm:"column:tick;NOT NULL"`
	MaxSupply           string `gorm:"column:max_supply;NOT NULL"`
	Decimals            string `gorm:"column:decimals;NOT NULL"`
	LimitPerMint        string `gorm:"column:limit_per_mint;NOT NULL"`
	RemainingSupply     string `gorm:"column:remaining_supply;NOT NULL"`
	BurnedSupply        string `gorm:"column:burned_supply;default:0;NOT NULL"`
	IsSelfMint          string `gorm:"column:is_self_mint;NOT NULL"`
	DeployInscriptionId string `gorm:"column:deploy_inscription_id;NOT NULL"`
	BlockHeight         string `gorm:"column:block_height;NOT NULL"`
}

func (b *Brc20Tickers) TableName() string {
	return "brc20_tickers"
}

type Brc20CumulativeEventHashes struct {
	Id                  string `gorm:"column:id;NOT NULL"`
	BlockHeight         string `gorm:"column:block_height;NOT NULL"`
	BlockEventHash      string `gorm:"column:block_event_hash;NOT NULL"`
	CumulativeEventHash string `gorm:"column:cumulative_event_hash;NOT NULL"`
}

func (b *Brc20CumulativeEventHashes) TableName() string {
	return "brc20_cumulative_event_hashes"
}

type Brc20EventTypes struct {
	Id            string `gorm:"column:id;NOT NULL"`
	EventTypeName string `gorm:"column:event_type_name;NOT NULL"`
	EventTypeId   string `gorm:"column:event_type_id;NOT NULL"`
}

func (b *Brc20EventTypes) TableName() string {
	return "brc20_event_types"
}

type Brc20IndexerVersion struct {
	Id               string `gorm:"column:id;NOT NULL"`
	IndexerVersion   string `gorm:"column:indexer_version;NOT NULL"`
	DbVersion        string `gorm:"column:db_version;NOT NULL"`
	EventHashVersion string `gorm:"column:event_hash_version;NOT NULL"`
}

func (b *Brc20IndexerVersion) TableName() string {
	return "brc20_indexer_version"
}
