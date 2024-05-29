package brc20

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/traitmeta/gotos/lib/db"
	"gorm.io/gorm"
)

// Global variables
var (
	blockStartMaxEventID                   int
	brc20EventsInsertCache                 []Event
	brc20TickersInsertCache                []Ticker
	brc20TickersRemainingSupplyUpdateCache = make(map[string]decimal.Decimal)
	brc20TickersBurnedSupplyUpdateCache    = make(map[string]decimal.Decimal)
	brc20HistoricBalancesInsertCache       []WalletBalance
	ticks                                  = make(map[string][]interface{})
	eventTypes                             = map[string]int{"deploy-inscribe": 1, "mint-inscribe": 2, "transfer-inscribe": 3, "transfer-transfer": 4}
	blockEventsStr                         string
	EVENT_SEPARATOR                        = "|"
)

// Event structure to hold event details
type Event struct {
	ID            int
	EventType     int
	BlockHeight   int
	InscriptionID string
	Event         string
}

// Ticker structure to hold ticker details
type Ticker struct {
	Tick                string
	OriginalTick        string
	MaxSupply           int
	Decimals            int
	LimitPerMint        int
	RemainingSupply     int
	BlockHeight         int
	IsSelfMint          bool
	DeployInscriptionID string
}

// WalletBalance structure to hold balance details
type WalletBalance struct {
	PkScript         string
	Wallet           string
	Tick             string
	OverallBalance   decimal.Decimal
	AvailableBalance decimal.Decimal
	BlockHeight      int64
	EventID          int64
}

// deployInscribe handles deploy-inscribe events
func deployInscribe(blockHeight int, inscriptionID, deployerPkScript, deployerWallet, tick, originalTick string, maxSupply, decimals, limitPerMint int, isSelfMint string) {
	event := map[string]string{
		"deployer_pkScript": deployerPkScript,
		"deployer_wallet":   deployerWallet,
		"tick":              tick,
		"original_tick":     originalTick,
		"max_supply":        fmt.Sprintf("%d", maxSupply),
		"decimals":          fmt.Sprintf("%d", decimals),
		"limit_per_mint":    fmt.Sprintf("%d", limitPerMint),
		"is_self_mint":      isSelfMint,
	}
	eventStr, _ := json.Marshal(event)
	blockEventsStr += string(eventStr) + EVENT_SEPARATOR

	eventID := blockStartMaxEventID + len(brc20EventsInsertCache) + 1
	brc20EventsInsertCache = append(brc20EventsInsertCache, Event{eventID, eventTypes["deploy-inscribe"], blockHeight, inscriptionID, string(eventStr)})

	brc20TickersInsertCache = append(brc20TickersInsertCache, Ticker{tick, originalTick, maxSupply, decimals, limitPerMint, maxSupply, blockHeight, isSelfMint == "true", inscriptionID})

	ticks[tick] = []interface{}{maxSupply, limitPerMint, decimals, isSelfMint == "true", inscriptionID}
}

// mintInscribe handles mint-inscribe events
func mintInscribe(blockHeight int64, inscriptionID, mintedPkScript, mintedWallet, tick, originalTick string, amount decimal.Decimal, parentID string) error {
	event := map[string]string{
		"minted_pkScript": mintedPkScript,
		"minted_wallet":   mintedWallet,
		"tick":            tick,
		"original_tick":   originalTick,
		"amount":          amount.String(),
		"parent_id":       parentID,
	}
	eventStr, _ := json.Marshal(event)
	blockEventsStr += string(eventStr) + EVENT_SEPARATOR

	eventID := blockStartMaxEventID + len(brc20EventsInsertCache) + 1
	brc20EventsInsertCache = append(brc20EventsInsertCache, Event{eventID, eventTypes["mint-inscribe"], blockHeight, inscriptionID, string(eventStr)})
	brc20TickersRemainingSupplyUpdateCache[tick] += amount

	lastBalance, err := getLastBalance(mintedPkScript, tick)
	if err != nil {
		return err
	}

	lastBalance.OverallBalance += amount
	lastBalance.AvailableBalance += amount
	brc20HistoricBalancesInsertCache = append(brc20HistoricBalancesInsertCache, WalletBalance{mintedPkScript, mintedWallet, tick, lastBalance.OverallBalance, lastBalance.AvailableBalance, blockHeight, eventID})

	ticks[tick][0] = ticks[tick][0].(int) - amount
	return nil
}

// transferInscribe handles transfer-inscribe events
func transferInscribe(blockHeight int, inscriptionID, sourcePkScript, sourceWallet, tick, originalTick string, amount int) {
	event := map[string]string{
		"source_pkScript": sourcePkScript,
		"source_wallet":   sourceWallet,
		"tick":            tick,
		"original_tick":   originalTick,
		"amount":          fmt.Sprintf("%d", amount),
	}
	eventStr, _ := json.Marshal(event)
	blockEventsStr += string(eventStr) + EVENT_SEPARATOR

	eventID := blockStartMaxEventID + len(brc20EventsInsertCache) + 1
	brc20EventsInsertCache = append(brc20EventsInsertCache, Event{eventID, eventTypes["transfer-inscribe"], blockHeight, inscriptionID, string(eventStr)})
	setTransferAsValid(inscriptionID)

	lastBalance := getLastBalance(sourcePkScript, tick)
	lastBalance.AvailableBalance -= amount
	brc20HistoricBalancesInsertCache = append(brc20HistoricBalancesInsertCache, WalletBalance{sourcePkScript, sourceWallet, tick, lastBalance.OverallBalance, lastBalance.AvailableBalance, blockHeight, eventID})

	saveTransferInscribeEvent(inscriptionID, event)
}

// transferTransferNormal handles normal transfer-transfer events
func transferTransferNormal(blockHeight int, inscriptionID, spentPkScript, spentWallet, tick, originalTick string, amount int, usingTxID string) {
	inscribeEvent := getTransferInscribeEvent(inscriptionID)
	sourcePkScript := inscribeEvent["source_pkScript"]
	sourceWallet := inscribeEvent["source_wallet"]
	event := map[string]string{
		"source_pkScript": sourcePkScript,
		"source_wallet":   sourceWallet,
		"spent_pkScript":  spentPkScript,
		"spent_wallet":    spentWallet,
		"tick":            tick,
		"original_tick":   originalTick,
		"amount":          fmt.Sprintf("%d", amount),
		"using_tx_id":     usingTxID,
	}
	eventStr, _ := json.Marshal(event)
	blockEventsStr += string(eventStr) + EVENT_SEPARATOR

	eventID := blockStartMaxEventID + len(brc20EventsInsertCache) + 1
	brc20EventsInsertCache = append(brc20EventsInsertCache, Event{eventID, eventTypes["transfer-transfer"], blockHeight, inscriptionID, string(eventStr)})
	setTransferAsUsed(inscriptionID)

	lastBalance := getLastBalance(sourcePkScript, tick)
	lastBalance.OverallBalance -= amount
	brc20HistoricBalancesInsertCache = append(brc20HistoricBalancesInsertCache, WalletBalance{sourcePkScript, sourceWallet, tick, lastBalance.OverallBalance, lastBalance.AvailableBalance, blockHeight, eventID})

	if spentPkScript != sourcePkScript {
		lastBalance = getLastBalance(spentPkScript, tick)
	}
	lastBalance.OverallBalance += amount
	lastBalance.AvailableBalance += amount
	brc20HistoricBalancesInsertCache = append(brc20HistoricBalancesInsertCache, WalletBalance{spentPkScript, spentWallet, tick, lastBalance.OverallBalance, lastBalance.AvailableBalance, blockHeight, -1 * eventID}) // negated to make a unique event_id

	if spentPkScript == "6a" {
		brc20TickersBurnedSupplyUpdateCache[tick] += amount
	}
}

// transferTransferSpendToFee handles transfer-transfer events where the spent amount is converted to a fee
func transferTransferSpendToFee(blockHeight int, inscriptionID, tick, originalTick string, amount decimal.Decimal, usingTxID string) error {
	inscribeEvent, err := getTransferInscribeEvent(inscriptionID)
	if err != nil {
		return err
	}

	sourcePkScript := inscribeEvent["source_pkScript"]
	sourceWallet := inscribeEvent["source_wallet"]
	event := map[string]string{
		"source_pkScript": sourcePkScript,
		"source_wallet":   sourceWallet,
		"spent_pkScript":  "",
		"spent_wallet":    "",
		"tick":            tick,
		"original_tick":   originalTick,
		"amount":          fmt.Sprintf("%d", amount),
		"using_tx_id":     usingTxID,
	}
	eventStr, _ := json.Marshal(event)
	blockEventsStr += string(eventStr) + EVENT_SEPARATOR

	eventID := blockStartMaxEventID + len(brc20EventsInsertCache) + 1
	brc20EventsInsertCache = append(brc20EventsInsertCache, Event{eventID, eventTypes["transfer-transfer"], blockHeight, inscriptionID, string(eventStr)})
	setTransferAsUsed(inscriptionID)

	lastBalance, err := getLastBalance(sourcePkScript, tick)
	if err != nil {
		return err
	}

	lastBalance.AvailableBalance = lastBalance.AvailableBalance.Add(amount)
	brc20HistoricBalancesInsertCache = append(brc20HistoricBalancesInsertCache, WalletBalance{sourcePkScript, sourceWallet, tick, lastBalance.OverallBalance, lastBalance.AvailableBalance, blockHeight, eventID})
	return nil
}

// updateEventHashes updates the event hashes for the block
func updateEventHashes(blockHeight int64) error {
	if len(blockEventsStr) > 0 && blockEventsStr[len(blockEventsStr)-1] == byte(EVENT_SEPARATOR[0]) {
		blockEventsStr = blockEventsStr[:len(blockEventsStr)-1] // remove last separator
	}
	blockEventHash := getSha256Hash(blockEventsStr)

	var eventHash Brc20CumulativeEventHashes
	if err := db.DBEngine.DB.Model(&Brc20CumulativeEventHashes{}).
		Where("block_height = ?", blockHeight-1).
		Scan(&eventHash).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			eventHash.CumulativeEventHash = blockEventHash
		} else {
			return err
		}
	}

	eventHash.CumulativeEventHash = getSha256Hash(eventHash.CumulativeEventHash + blockEventHash)

	eventHash.BlockHeight = blockHeight
	eventHash.BlockEventHash = blockEventHash
	return db.DBEngine.DB.Model(&Brc20CumulativeEventHashes{}).Create(&eventHash).Error
}

// Dummy functions for missing implementations
func getEventStr(event map[string]string, eventType, inscriptionID string) string {
	// Implement this function based on your requirements
	return ""
}

func getSha256Hash(data string) string {
	// Implement this function to return SHA256 hash of the input data
	return ""
}
