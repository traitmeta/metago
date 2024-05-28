package brc20

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/traitmeta/gotos/lib/db"
)

// Global variables
var (
	transferInscribeEventCache = make(map[string]Brc20Events)
	balanceCache               = make(map[string]WalletBalance)
	transferValidityCache      = make(map[string]int)
	cachedTickers              = make(map[string]Brc20Tickers)
)

// getTransferInscribeEvent fetches the transfer inscribe event for a given inscription ID
func getTransferInscribeEvent(inscriptionID string) (*Brc20Events, error) {
	if event, ok := transferInscribeEventCache[inscriptionID]; ok {
		delete(transferInscribeEventCache, inscriptionID)
		return &event, nil
	}

	var event Brc20Events
	if err := db.DBEngine.DB.Model(&Brc20Events{}).
		Where("event_type = ? AND inscription_id = ?", eventTypes["transfer-inscribe"], inscriptionID).
		Take(&event).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

// TODO save to DB
// saveTransferInscribeEvent saves the transfer inscribe event in the cache
func saveTransferInscribeEvent(inscriptionID string, event Brc20Events) {
	transferInscribeEventCache[inscriptionID] = event
}

// getLastBalance fetches the last balance for a given pkScript and tick
func getLastBalance(pkScript, tick string) (*WalletBalance, error) {
	cacheKey := pkScript + tick
	if balance, ok := balanceCache[cacheKey]; ok {
		return &balance, nil
	}

	var balance Brc20HistoricBalances
	if err := db.DBEngine.DB.Model(&Brc20HistoricBalances{}).
		Where("pkscript = ? AND tick = ?", pkScript, tick).
		Order("block_height DESC, id DESC").Take(&balance).Error; err != nil {
		return nil, err
	}

	walletBalance := &WalletBalance{
		OverallBalance:   balance.OverallBalance,
		AvailableBalance: balance.AvailableBalance,
	}

	balanceCache[cacheKey] = *walletBalance
	return walletBalance, nil
}

// checkAvailableBalance checks if the available balance is sufficient
func checkAvailableBalance(pkScript, tick string, amount decimal.Decimal) (bool, error) {
	lastBalance, err := getLastBalance(pkScript, tick)
	if err != nil {
		return false, err
	}
	return lastBalance.AvailableBalance.Cmp(amount) >= 0, nil
}

// isUsedOrInvalid checks if the inscription is used or invalid
func isUsedOrInvalid(inscriptionID string) (bool, error) {
	if status, ok := transferValidityCache[inscriptionID]; ok {
		return status != 1, nil
	}

	type Cnt struct {
		InscrCnt    int `json:"inscr_cnt,omitempty"`
		TransferCnt int `json:"transfer_cnt,omitempty"`
	}

	query := `SELECT COALESCE(SUM(CASE WHEN event_type = ? THEN 1 ELSE 0 END), 0) AS inscr_cnt,
                     COALESCE(SUM(CASE WHEN event_type = ? THEN 1 ELSE 0 END), 0) AS transfer_cnt
              FROM brc20_events WHERE inscription_id = ?;`
	var resp Cnt
	if err := db.DBEngine.Exec(query, eventTypes["transfer-inscribe"], eventTypes["transfer-transfer"], inscriptionID).
		Scan(&resp).Error; err != nil {
		return false, err
	}

	if resp.InscrCnt != 1 {
		transferValidityCache[inscriptionID] = 0 // invalid transfer (no inscribe event)
		return true, nil
	} else if resp.TransferCnt != 0 {
		transferValidityCache[inscriptionID] = -1 // used
		return true, nil
	} else {
		transferValidityCache[inscriptionID] = 1 // valid
		return false, nil
	}
}

// setTransferAsUsed marks the transfer as used in the cache
func setTransferAsUsed(inscriptionID string) {
	transferValidityCache[inscriptionID] = -1
}

// setTransferAsValid marks the transfer as valid in the cache
func setTransferAsValid(inscriptionID string) {
	transferValidityCache[inscriptionID] = 1
}

// resetCaches resets all the caches and reloads the ticks from the database
// TODO use levelDB to cache all ticker info
func resetCaches() {
	balanceCache = make(map[string]WalletBalance)
	transferInscribeEventCache = make(map[string]Brc20Events)
	transferValidityCache = make(map[string]int)
	startTime := time.Now()

	var tickers []Brc20Tickers
	if err := db.DBEngine.DB.Model(&Brc20Tickers{}).Scan(&tickers).Error; err != nil {
		return
	}

	for _, ticker := range tickers {
		cachedTickers[ticker.Tick] = ticker
	}

	fmt.Printf("Ticks refreshed in %.2f seconds\n", time.Since(startTime).Seconds())
}
