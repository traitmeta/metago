package brc20

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/shopspring/decimal"
)

// Global variables
var (
	blockStartMaxEventID                   int
	brc20EventsInsertCache                 []Event
	brc20TickersInsertCache                []Ticker
	brc20TickersRemainingSupplyUpdateCache = make(map[string]int)
	brc20TickersBurnedSupplyUpdateCache    = make(map[string]int)
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
func mintInscribe(blockHeight int, inscriptionID, mintedPkScript, mintedWallet, tick, originalTick string, amount int, parentID string) {
	event := map[string]string{
		"minted_pkScript": mintedPkScript,
		"minted_wallet":   mintedWallet,
		"tick":            tick,
		"original_tick":   originalTick,
		"amount":          fmt.Sprintf("%d", amount),
		"parent_id":       parentID,
	}
	eventStr, _ := json.Marshal(event)
	blockEventsStr += string(eventStr) + EVENT_SEPARATOR

	eventID := blockStartMaxEventID + len(brc20EventsInsertCache) + 1
	brc20EventsInsertCache = append(brc20EventsInsertCache, Event{eventID, eventTypes["mint-inscribe"], blockHeight, inscriptionID, string(eventStr)})
	brc20TickersRemainingSupplyUpdateCache[tick] += amount

	lastBalance := getLastBalance(mintedPkScript, tick)
	lastBalance.OverallBalance += amount
	lastBalance.AvailableBalance += amount
	brc20HistoricBalancesInsertCache = append(brc20HistoricBalancesInsertCache, WalletBalance{mintedPkScript, mintedWallet, tick, lastBalance.OverallBalance, lastBalance.AvailableBalance, blockHeight, eventID})

	ticks[tick][0] = ticks[tick][0].(int) - amount
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
func transferTransferSpendToFee(blockHeight int, inscriptionID, tick, originalTick string, amount int, usingTxID string) {
	inscribeEvent := getTransferInscribeEvent(inscriptionID)
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

	lastBalance := getLastBalance(sourcePkScript, tick)
	lastBalance.AvailableBalance += amount
	brc20HistoricBalancesInsertCache = append(brc20HistoricBalancesInsertCache, WalletBalance{sourcePkScript, sourceWallet, tick, lastBalance.OverallBalance, lastBalance.AvailableBalance, blockHeight, eventID})
}

// updateEventHashes updates the event hashes for the block
func updateEventHashes(blockHeight int) {
	if len(blockEventsStr) > 0 && blockEventsStr[len(blockEventsStr)-1] == byte(EVENT_SEPARATOR[0]) {
		blockEventsStr = blockEventsStr[:len(blockEventsStr)-1] // remove last separator
	}
	blockEventHash := getSha256Hash(blockEventsStr)

	var cumulativeEventHash string
	row := db.QueryRow("SELECT cumulative_event_hash FROM brc20_cumulative_event_hashes WHERE block_height = $1", blockHeight-1)
	err := row.Scan(&cumulativeEventHash)
	if err == sql.ErrNoRows {
		cumulativeEventHash = blockEventHash
	} else if err != nil {
		log.Fatal(err)
	} else {
		cumulativeEventHash = getSha256Hash(cumulativeEventHash + blockEventHash)
	}
	_, err = db.Exec("INSERT INTO brc20_cumulative_event_hashes (block_height, block_event_hash, cumulative_event_hash) VALUES ($1, $2, $3)", blockHeight, blockEventHash, cumulativeEventHash)
	if err != nil {
		log.Fatal(err)
	}
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

func saveTransferInscribeEvent(inscriptionID string, event map[string]string) {
	// Implement this function to save transfer inscribe event
}

func getTransferInscribeEvent(inscriptionID string) map[string]string {
	// Implement this function to get transfer inscribe event by inscriptionID
	return make(map[string]string)
}
