package brc20

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/traitmeta/gotos/lib/db"
)

// Global variables
var (
	transferInscribeEventCache = make(map[string]string)
	balanceCache               = make(map[string]WalletBalance)
	transferValidityCache      = make(map[string]int)
	eventTypes                 = map[string]int{
		"transfer-inscribe": 1,
		"transfer-transfer": 2,
	}
)

// getTransferInscribeEvent fetches the transfer inscribe event for a given inscription ID
func getTransferInscribeEvent(inscriptionID string) (string, error) {
	if event, ok := transferInscribeEventCache[inscriptionID]; ok {
		delete(transferInscribeEventCache, inscriptionID)
		return event, nil
	}

	query := `SELECT event FROM brc20_events WHERE event_type = $1 AND inscription_id = $2;`
	row := db.QueryRow(query, eventTypes["transfer-inscribe"], inscriptionID)
	var event string
	if err := row.Scan(&event); err != nil {
		return "", err
	}

	return event, nil
}

// saveTransferInscribeEvent saves the transfer inscribe event in the cache
func saveTransferInscribeEvent(inscriptionID, event string) {
	transferInscribeEventCache[inscriptionID] = event
}

// getLastBalance fetches the last balance for a given pkScript and tick
func getLastBalance(pkScript, tick string) (WalletBalance, error) {
	cacheKey := pkScript + tick
	if balance, ok := balanceCache[cacheKey]; ok {
		return balance, nil
	}

	query := `SELECT overall_balance, available_balance FROM brc20_historic_balances WHERE pkscript = $1 AND tick = $2 ORDER BY block_height DESC, id DESC LIMIT 1;`
	row := db.QueryRow(query, pkScript, tick)
	var balance WalletBalance
	err := row.Scan(&balance.OverallBalance, &balance.AvailableBalance)
	if err == sql.ErrNoRows {
		balance = WalletBalance{OverallBalance: 0, AvailableBalance: 0}
	} else if err != nil {
		return WalletBalance{}, err
	}

	balanceCache[cacheKey] = balance
	return balance, nil
}

// checkAvailableBalance checks if the available balance is sufficient
func checkAvailableBalance(pkScript, tick string, amount int) (bool, error) {
	lastBalance, err := getLastBalance(pkScript, tick)
	if err != nil {
		return false, err
	}
	return lastBalance.AvailableBalance >= amount, nil
}

// isUsedOrInvalid checks if the inscription is used or invalid
func isUsedOrInvalid(inscriptionID string) (bool, error) {
	if status, ok := transferValidityCache[inscriptionID]; ok {
		return status != 1, nil
	}

	query := `SELECT COALESCE(SUM(CASE WHEN event_type = $1 THEN 1 ELSE 0 END), 0) AS inscr_cnt,
                     COALESCE(SUM(CASE WHEN event_type = $2 THEN 1 ELSE 0 END), 0) AS transfer_cnt
              FROM brc20_events WHERE inscription_id = $3;`
	row := db.QueryRow(query, eventTypes["transfer-inscribe"], eventTypes["transfer-transfer"], inscriptionID)
	var inscrCnt, transferCnt int
	if err := row.Scan(&inscrCnt, &transferCnt); err != nil {
		return false, err
	}

	if inscrCnt != 1 {
		transferValidityCache[inscriptionID] = 0 // invalid transfer (no inscribe event)
		return true, nil
	} else if transferCnt != 0 {
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
func resetCaches() {
	balanceCache = make(map[string]WalletBalance)
	transferInscribeEventCache = make(map[string]string)
	transferValidityCache = make(map[string]int)
	startTime := time.Now()

	query := `SELECT tick, remaining_supply, limit_per_mint, decimals, is_self_mint, deploy_inscription_id FROM brc20_tickers;`
	rows, err := db.Query(query)
	db.DBEngine.DB.Model(&Brc20Tickers{}).Scan(dest interface{})
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tick string
		var remainingSupply, limitPerMint, decimals int
		var isSelfMint bool
		var deployInscriptionID string
		if err := rows.Scan(&tick, &remainingSupply, &limitPerMint, &decimals, &isSelfMint, &deployInscriptionID); err != nil {
			log.Fatal(err)
		}
		ticks[tick] = []interface{}{remainingSupply, limitPerMint, decimals, isSelfMint, deployInscriptionID}
	}

	fmt.Printf("Ticks refreshed in %.2f seconds\n", time.Since(startTime).Seconds())
}
