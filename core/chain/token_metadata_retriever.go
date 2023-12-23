package chain

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/traitmeta/metago/core/models"
	"github.com/traitmeta/metago/pkg/abi"
)

type TokenMetadataRetriever struct {
	TokenFunctionsMaxRetries int
}

const (
	// 06fdde03 = keccak256(name())
	TokenFuncName = "06fdde03"
	// 95d89b41 = keccak256(symbol())
	TokenFuncSymbols = "95d89b41"
	// 313ce567 = keccak256(decimals())
	TokenFuncDecimals = "313ce567"
	// 18160ddd = keccak256(totalSupply())
	TokenFuncTotalSupply = "18160ddd"
)

func Init() *TokenMetadataRetriever {
	return &TokenMetadataRetriever{
		TokenFunctionsMaxRetries: 3,
	}
}

/*
Read functions below in the Smart Contract given the Contract's address hash.
  - totalSupply
  - decimals
  - name
  - symbol

This function will return a map with functions that were read in the Smart Contract, for instance:

  - Given that all functions were read:
    {
    name: "BOB",
    decimals: 18,
    total_supply: 1_000_000_000_000_000_000,
    symbol: nil
    }

  - Given that some of them were read:
    {
    name: "BOB",
    decimals: 18
    }

    It will retry to fetch each function in the Smart Contract according to :token_functions_reader_max_retries
    configured in the application env case one of them raised error.
*/
func (m *TokenMetadataRetriever) GetMetadata(tokenAddresses []string) (map[string]models.Token, error) {
	tokenMaps := make(map[string]models.Token)
	for _, contractAddress := range tokenAddresses {
		totalSupply, name, symbols, decimals, err := abi.GetErc20Metadata(contractAddress)
		for i := 0; i < m.TokenFunctionsMaxRetries; i++ {
			if err != nil {
				fmt.Printf("<Token contract hash: %s> error while fetching metadata:\n%s\nRetries left: %d\n", contractAddress, err.Error(), m.TokenFunctionsMaxRetries-i-1)
			} else {
				break
			}

			totalSupply, name, symbols, decimals, err = abi.GetErc20Metadata(contractAddress)
		}

		name = m.handleInvalidName(name, contractAddress)
		symbols = m.handleInvalidSymbol(symbols)
		tokenMaps[contractAddress] = models.Token{
			Name:            name,
			Symbol:          symbols,
			TotalSupply:     totalSupply,
			Decimals:        decimals,
			ContractAddress: contractAddress,
		}
	}
	return tokenMaps, nil
}

func (m *TokenMetadataRetriever) handleInvalidName(name string, contractAddress string) string {
	if utf8.ValidString(name) {
		name = m.removeNullBytes(name)
	} else {
		name = contractAddress[:6]
	}

	name = m.handleLargeString(name)

	return name
}

func (m *TokenMetadataRetriever) handleInvalidSymbol(symbol string) string {
	if utf8.ValidString(symbol) {
		symbol = m.removeNullBytes(symbol)
		symbol = m.handleLargeString(symbol)
		return symbol
	}

	return ""
}

func (m *TokenMetadataRetriever) handleLargeString(str string) string {
	if len(str) > 255 {
		return str[:255]
	} else {
		return str
	}
}

func (m *TokenMetadataRetriever) removeNullBytes(str string) string {
	return strings.ReplaceAll(str, "\x00", "")
}
