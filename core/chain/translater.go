package chain

import (
	"math/big"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/traitmeta/metago/core/common"
	"github.com/traitmeta/metago/core/models"
	"github.com/traitmeta/metago/pkg/abi"
)

var tokenTypesPriorityOrder = map[string]int{
	common.ERC20:   1,
	common.ERC721:  2,
	common.ERC1155: 3,
}

type TokenTransfers struct {
	Tokens         []models.Token
	TokenTransfers []models.TokenTransfer
}

// TODO
func parse(logs []models.Event) TokenTransfers {
	initialAcc := TokenTransfers{
		Tokens:         []models.Token{},
		TokenTransfers: []models.TokenTransfer{},
	}
	erc20TokenTransfers := doParse(logs, func(log models.Event) bool {
		return log.FirstTopic == common.ERC20TokenTransferEventFuncSign && log.FourthTopic == ""
	}, initialAcc, common.ERC20)

	erc721TokenTransfers := doParse(logs, func(log models.Event) bool {
		return log.FirstTopic == common.ERC20TokenTransferEventFuncSign && log.FourthTopic != ""
	}, initialAcc, common.ERC721)

	wethTransfers := doParse(logs, func(log models.Event) bool {
		return log.FirstTopic == common.WETHDepositSignature ||
			log.FirstTopic == common.WETHWithdrawalSignature
	}, initialAcc, common.ERC721)

	erc1155TokenTransfers := doParse(logs, func(log models.Event) bool {
		return log.FirstTopic == common.ERC1155SingleTransferSignature ||
			log.FirstTopic == common.ERC1155BatchTransferSignature
	}, initialAcc, common.ERC1155)

	roughTokens := append(append(erc1155TokenTransfers.Tokens, append(erc721TokenTransfers.Tokens, erc20TokenTransfers.Tokens...)...), wethTransfers.Tokens...)

	roughTokenTransfers := append(append(erc1155TokenTransfers.TokenTransfers, append(erc721TokenTransfers.TokenTransfers, erc20TokenTransfers.TokenTransfers...)...), wethTransfers.TokenTransfers...)

	tokens, tokenTransfers := sanitizeTokenTypes(roughTokens, roughTokenTransfers)

	burnAddress := "0x0000000000000000000000000000000000000000"

	filteredTokenTransfers := []models.TokenTransfer{}
	for _, tokenTransfer := range tokenTransfers {
		if tokenTransfer.ToAddress == burnAddress || tokenTransfer.FromAddress == burnAddress {
			filteredTokenTransfers = append(filteredTokenTransfers, tokenTransfer)
		}
	}

	uniqueTokenContractAddressHashes := []string{}
	for _, tokenTransfer := range filteredTokenTransfers {
		uniqueTokenContractAddressHashes = appendIfMissing(uniqueTokenContractAddressHashes, tokenTransfer.TokenContractAddress)
	}

	addTokens(uniqueTokenContractAddressHashes)

	uniqueTokens := []models.Token{}
	for _, token := range tokens {
		if contains(uniqueTokenContractAddressHashes, token.ContractAddress) {
			uniqueTokens = append(uniqueTokens, token)
		}
	}

	tokenTransfersFromLogsUnique := TokenTransfers{
		Tokens:         uniqueTokens,
		TokenTransfers: filteredTokenTransfers,
	}

	return tokenTransfersFromLogsUnique
}

// TODO
func doParse(logs []models.Event, filter func(models.Event) bool, acc TokenTransfers, tokenType string) TokenTransfers {
	filteredLogs := []models.Event{}
	for _, log := range logs {
		if filter(log) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	// TODO
	switch tokenType {
	case common.ERC1155:
		acc = doParseErc1155(filteredLogs, acc)
	}

	return acc
}

func doParseErc721(logs []models.Event, acc TokenTransfers) TokenTransfers {
	for _, log := range logs {
		tokenId := big.NewInt(0).SetBytes(ethcommon.Hex2Bytes(log.FourthTopic))
		tokenTransfer := models.TokenTransfer{
			BlockNumber:          logs[0].BlockNumber,
			BlockHash:            logs[0].BlockHash,
			LogIndex:             logs[0].LogIndex,
			FromAddress:          truncateAddressHash(logs[0].ThirdTopic),
			ToAddress:            truncateAddressHash(logs[0].FourthTopic),
			TokenContractAddress: logs[0].Address,
			TransactionHash:      logs[0].TxHash,
			TokenId:              tokenId,
		}

		token := models.Token{
			ContractAddress: logs[0].Address,
			Type:            common.ERC721,
		}

		acc.Tokens = append(acc.Tokens, token)
		acc.TokenTransfers = append(acc.TokenTransfers, tokenTransfer)
	}

	return acc
}

func doParseErc1155(logs []models.Event, acc TokenTransfers) TokenTransfers {
	for _, log := range logs {
		tokenTransfer := models.TokenTransfer{
			BlockNumber:          log.BlockNumber,
			BlockHash:            log.BlockHash,
			LogIndex:             log.LogIndex,
			FromAddress:          ethcommon.HexToAddress(log.ThirdTopic).String(),
			ToAddress:            ethcommon.HexToAddress(log.FourthTopic).String(),
			TokenContractAddress: log.Address,
			TransactionHash:      log.TxHash,
		}

		token := models.Token{
			ContractAddress: logs[0].Address,
			Type:            common.ERC1155,
		}

		if log.FirstTopic == common.ERC1155SingleTransferSignature {
			tokenId, value, err := abi.ParseErc1155SignleTransferLog(ethcommon.Hex2Bytes(log.Data))
			if err != nil {
				continue
			}
			tokenTransfer.TokenId = tokenId
			tokenTransfer.Amount = value
		} else {
			tokenIds, values, err := abi.ParseErc1155BatchTransferLog(ethcommon.Hex2Bytes(log.Data))
			if err != nil {
				continue
			}
			tokenTransfer.TokenIds = tokenIds
			tokenTransfer.Amounts = values
		}
		acc.Tokens = append(acc.Tokens, token)
		acc.TokenTransfers = append(acc.TokenTransfers, tokenTransfer)
	}

	return acc
}

// TODO error
func getTokenType(contractAddress string, tokens []models.Token) string {
	for _, t := range tokens {
		if t.ContractAddress == contractAddress {
			return t.Type
		}
	}

	return ""
}

func sanitizeTokenTypes(tokens []models.Token, tokenTransfers []models.TokenTransfer) ([]models.Token, []models.TokenTransfer) {
	existingTokenTypesMap := map[string]string{}
	for _, token := range tokens {
		if _, ok := existingTokenTypesMap[token.ContractAddress]; !ok {
			existingTokenTypesMap[token.ContractAddress] = getTokenType(token.ContractAddress, tokens)
		}
	}

	existingTokens := []string{}
	for token := range existingTokenTypesMap {
		existingTokens = append(existingTokens, token)
	}

	newTokensTokenTransfers := []models.TokenTransfer{}
	for _, tokenTransfer := range tokenTransfers {
		if !contains(existingTokens, tokenTransfer.TokenContractAddress) {
			newTokensTokenTransfers = append(newTokensTokenTransfers, tokenTransfer)
		}
	}

	// newTokenTypesMap := map[string]string{}
	// for _, tokenTransfer := range newTokensTokenTransfers {
	// 	newTokenTypesMap[tokenTransfer.TokenContractAddress] = defineTokenType(tokenTransfer)
	// }

	// actualTokenTypesMap := mergeMaps(newTokenTypesMap, existingTokenTypesMap)
	actualTokenTypesMap := existingTokenTypesMap
	actualTokens := []models.Token{}
	for _, token := range tokens {
		token.Type = actualTokenTypesMap[token.ContractAddress]
		actualTokens = append(actualTokens, token)
	}

	actualTokenTransfers := []models.TokenTransfer{}
	for _, tokenTransfer := range tokenTransfers {
		actualTokenTransfers = append(actualTokenTransfers, tokenTransfer)
	}

	return actualTokens, actualTokenTransfers
}

func defineTokenType(tokens []models.Token) string {
	tokenType := ""
	tokenTypeOrder := -1
	for _, token := range tokens {
		order, ok := tokenTypesPriorityOrder[token.Type]
		if !ok {
			order = -1
		}

		if order > tokenTypeOrder {
			tokenType = token.Type
			tokenTypeOrder = order
		}
	}

	return tokenType
}

func truncateAddressHash(addressHash string) string {
	if addressHash == "" {
		return "0x0000000000000000000000000000000000000000"
	}

	if strings.HasPrefix(addressHash, "0x000000000000000000000000") {
		return "0x" + addressHash[24:]
	}
	ethcommon.HexToAddress(addressHash).String()
	return addressHash
}

func appendIfMissing(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

func contains(slice []string, s string) bool {
	for _, ele := range slice {
		if ele == s {
			return true
		}
	}
	return false
}

func addTokens(tokens []string) {
	// implementation for adding tokens
}
