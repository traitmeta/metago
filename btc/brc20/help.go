package brc20

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// utf8len returns the length of the UTF-8 encoded string
func utf8len(s string) int {
	return len([]byte(s))
}

// isPositiveNumber checks if the string is a positive number
func isPositiveNumber(s string, doStrip bool) bool {
	if doStrip {
		s = strings.TrimSpace(s)
	}
	if len(s) == 0 {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// isPositiveNumberWithDot checks if the string is a positive number with at most one dot
func isPositiveNumberWithDot(s string, doStrip bool) bool {
	if doStrip {
		s = strings.TrimSpace(s)
	}
	dotFound := false
	if len(s) == 0 || s[0] == '.' || s[len(s)-1] == '.' {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			if ch != '.' || dotFound {
				return false
			}
			dotFound = true
		}
	}
	return true
}

// getNumberExtendedTo18Decimals extends the number to 18 decimal places
func getNumberExtendedTo18Decimals(s string, decimals int, doStrip bool) *int {
	if doStrip {
		s = strings.TrimSpace(s)
	}
	if strings.Contains(s, ".") {
		parts := strings.Split(s, ".")
		normalPart := parts[0]
		if len(parts[1]) > decimals || len(parts[1]) == 0 {
			return nil
		}
		decimalsPart := parts[1][:decimals]
		for len(decimalsPart) < 18 {
			decimalsPart += "0"
		}
		result := normalPart + decimalsPart
		if num, err := strconv.Atoi(result); err == nil {
			return &num
		}
	} else {
		if num, err := strconv.Atoi(s); err == nil {
			result := num * int(math.Pow10(18))
			return &result
		}
	}
	return nil
}

// fixNumStrDecimals fixes the number string to the specified decimal places
func fixNumStrDecimals(numStr string, decimals int) string {
	if len(numStr) <= 18 {
		numStr = strings.Repeat("0", 18-len(numStr)) + numStr
		numStr = "0." + numStr
		if decimals < 18 {
			numStr = numStr[:len(numStr)-(18-decimals)]
		}
	} else {
		numStr = numStr[:len(numStr)-18] + "." + numStr[len(numStr)-18:]
		if decimals < 18 {
			numStr = numStr[:len(numStr)-(18-decimals)]
		}
	}
	if numStr[len(numStr)-1] == '.' {
		numStr = numStr[:len(numStr)-1]
	}
	return numStr
}

// getEventStr generates the event string based on the event type
func convertEventStr(event map[string]string, eventType string, inscriptionID string, ticks map[string][]interface{}) string {
	var res string
	var decimalsInt int
	switch eventType {
	case "deploy-inscribe":
		decimalsInt, _ = strconv.Atoi(event["decimals"])
		res = "deploy-inscribe;"
		res += inscriptionID + ";"
		res += event["deployer_pkScript"] + ";"
		res += event["tick"] + ";"
		res += event["original_tick"] + ";"
		res += fixNumStrDecimals(event["max_supply"], decimalsInt) + ";"
		res += event["decimals"] + ";"
		res += fixNumStrDecimals(event["limit_per_mint"], decimalsInt) + ";"
		res += event["is_self_mint"]
	case "mint-inscribe":
		decimalsInt = ticks[event["tick"]][2].(int)
		res = "mint-inscribe;"
		res += inscriptionID + ";"
		res += event["minted_pkScript"] + ";"
		res += event["tick"] + ";"
		res += event["original_tick"] + ";"
		res += fixNumStrDecimals(event["amount"], decimalsInt) + ";"
		res += event["parent_id"]
	case "transfer-inscribe":
		decimalsInt = ticks[event["tick"]][2].(int)
		res = "transfer-inscribe;"
		res += inscriptionID + ";"
		res += event["source_pkScript"] + ";"
		res += event["tick"] + ";"
		res += event["original_tick"] + ";"
		res += fixNumStrDecimals(event["amount"], decimalsInt)
	case "transfer-transfer":
		decimalsInt = ticks[event["tick"]][2].(int)
		res = "transfer-transfer;"
		res += inscriptionID + ";"
		res += event["source_pkScript"] + ";"
		if event["spent_pkScript"] != "" {
			res += event["spent_pkScript"] + ";"
		} else {
			res += ";"
		}
		res += event["tick"] + ";"
		res += event["original_tick"] + ";"
		res += fixNumStrDecimals(event["amount"], decimalsInt)
	default:
		fmt.Println("EVENT TYPE ERROR!!")
		os.Exit(1)
	}
	return res
}

// getSHA256Hash returns the SHA-256 hash of the string
func getSHA256Hash(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func test() {
	// Test examples
	s := "Wed, 22 May 2024 03:47:09 GMT"
	fmt.Println("UTF-8 Length:", utf8len(s))

	fmt.Println("Is Positive Number:", isPositiveNumber("12345", true))
	fmt.Println("Is Positive Number with Dot:", isPositiveNumberWithDot("123.45", true))

	decimals := 5
	number := "123.456"
	if extendedNum := getNumberExtendedTo18Decimals(number, decimals, true); extendedNum != nil {
		fmt.Println("Extended Number to 18 Decimals:", *extendedNum)
	} else {
		fmt.Println("Invalid number format.")
	}

	numStr := "123456789012345678"
	fmt.Println("Fixed NumStr Decimals:", fixNumStrDecimals(numStr, 5))

	event := map[string]string{
		"decimals":          "2",
		"deployer_pkScript": "deployerScript",
		"tick":              "tickValue",
		"original_tick":     "originalTick",
		"max_supply":        "100000000",
		"limit_per_mint":    "1000",
		"is_self_mint":      "true",
		"minted_pkScript":   "mintedScript",
		"amount":            "100",
		"parent_id":         "parentId",
		"source_pkScript":   "sourceScript",
		"spent_pkScript":    "spentScript",
	}
	ticks := map[string][]interface{}{
		"tickValue": {0, 1, 2},
	}
	inscriptionID := "inscriptionId"
	fmt.Println("Event String:", convertEventStr(event, "deploy-inscribe", inscriptionID, ticks))

	inputStr := "Hello, World!"
	fmt.Println("SHA-256 Hash:", getSHA256Hash(inputStr))
}
