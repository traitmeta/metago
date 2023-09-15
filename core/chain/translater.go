package chain

import (
	"math/big"

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
func ParseTokenTransfers(logs []models.Event) (tokenTransfers TokenTransfers, err error) {
	tokenTransfers, err = doParse(logs, tokenTransfers)
	if err != nil {
		return
	}

	filteredTokenTransfers := []models.TokenTransfer{}
	for _, tokenTransfer := range tokenTransfers.TokenTransfers {
		if tokenTransfer.ToAddress == common.ZeroAddress || tokenTransfer.FromAddress == common.ZeroAddress {
			filteredTokenTransfers = append(filteredTokenTransfers, tokenTransfer)
		}
	}

	uniqueTokenContractAddresses := []string{}
	for _, tokenTransfer := range filteredTokenTransfers {
		uniqueTokenContractAddresses = appendIfMissing(uniqueTokenContractAddresses, tokenTransfer.TokenContractAddress)
	}

	// 处理totalSupply, 通过合约获取totalSupply, 目前看只有ERC20有这个方法
	addTokens(uniqueTokenContractAddresses)

	uniqueTokens := []models.Token{}
	for _, token := range tokenTransfers.Tokens {
		if contains(uniqueTokenContractAddresses, token.ContractAddress) {
			uniqueTokens = append(uniqueTokens, token)
		}
	}

	tokenTransfersFromLogsUnique := TokenTransfers{
		Tokens:         uniqueTokens,
		TokenTransfers: filteredTokenTransfers,
	}

	return tokenTransfersFromLogsUnique, nil
}

func doParse(logs []models.Event, acc TokenTransfers) (TokenTransfers, error) {
	filteredLogs := filterLogs(logs)
	for tokenType, val := range filteredLogs {
		switch tokenType {
		case common.ERC20:
			return doParseErc20(val, acc)
		case common.WETH:
			return doParseWTH(val, acc)
		case common.ERC721:
			return doParseErc721(val, acc), nil
		case common.ERC1155:
			return doParseErc1155(val, acc)
		}
	}

	return acc, nil
}

func filterLogs(logs []models.Event) map[string][]models.Event {
	filteredLogs := map[string][]models.Event{}
	for _, log := range logs {
		if log.FirstTopic == common.ERC20TokenTransferEventFuncSign {
			if log.FourthTopic == "" {
				filteredLogs[common.ERC20] = append(filteredLogs[common.ERC20], log)
			} else {
				filteredLogs[common.ERC721] = append(filteredLogs[common.ERC721], log)
			}
		}

		if log.FirstTopic == common.WETHDepositSignature || log.FirstTopic == common.WETHWithdrawalSignature {
			filteredLogs[common.WETH] = append(filteredLogs[common.WETH], log)
		}

		if log.FirstTopic == common.ERC1155SingleTransferSignature || log.FirstTopic == common.ERC1155BatchTransferSignature {
			filteredLogs[common.ERC1155] = append(filteredLogs[common.ERC1155], log)
		}

	}

	return filteredLogs
}

// event Transfer(address indexed _from,address indexed _to,uint256 indexed _tokenId);
func doParseErc721(logs []models.Event, acc TokenTransfers) TokenTransfers {
	for _, log := range logs {
		token, tokenTransfer := doParseBaseTokenTransfer(log)
		tokenTransfer.TokenId = big.NewInt(0).SetBytes(ethcommon.Hex2Bytes(log.FourthTopic))
		tokenTransfer.FromAddress = ethcommon.HexToAddress(log.SecondTopic).String()
		tokenTransfer.ToAddress = ethcommon.HexToAddress(log.ThirdTopic).String()

		token.Type = common.ERC721

		acc.Tokens = append(acc.Tokens, token)
		acc.TokenTransfers = append(acc.TokenTransfers, tokenTransfer)
	}

	return acc
}

// event Transfer(address indexed _from,address indexed _to,uint256 indexed _tokenId);
func doParseErc20(logs []models.Event, acc TokenTransfers) (TokenTransfers, error) {
	for _, log := range logs {
		token, tokenTransfer := doParseBaseTokenTransfer(log)
		amount, err := abi.ParseErc20TransferLog(ethcommon.Hex2Bytes(log.Data))
		if err != nil {
			return acc, err
		}

		tokenTransfer.Amount = amount
		tokenTransfer.FromAddress = ethcommon.HexToAddress(log.SecondTopic).String()
		tokenTransfer.ToAddress = ethcommon.HexToAddress(log.ThirdTopic).String()
		token.Type = common.ERC20

		acc.Tokens = append(acc.Tokens, token)
		acc.TokenTransfers = append(acc.TokenTransfers, tokenTransfer)
	}

	return acc, nil
}

// event Deposit(address indexed from, uint256 amount);
// event Withdraw(address indexed to, uint256 amount);
func doParseWTH(logs []models.Event, acc TokenTransfers) (TokenTransfers, error) {
	for _, log := range logs {
		token, tokenTransfer := doParseBaseTokenTransfer(log)
		amount, err := abi.ParseErc20TransferLog(ethcommon.Hex2Bytes(log.Data))
		if err != nil {
			return acc, err
		}

		tokenTransfer.Amount = amount

		if log.FirstTopic == common.WETHDepositSignature {
			tokenTransfer.FromAddress = common.ZeroAddress
			tokenTransfer.ToAddress = ethcommon.HexToAddress(log.SecondTopic).String()
		} else {
			tokenTransfer.FromAddress = ethcommon.HexToAddress(log.SecondTopic).String()
			tokenTransfer.ToAddress = common.ZeroAddress
		}

		token.Type = common.ERC20
		acc.Tokens = append(acc.Tokens, token)
		acc.TokenTransfers = append(acc.TokenTransfers, tokenTransfer)
	}

	return acc, nil
}

// event TransferSingle/Batch(address indexed _operator, address indexed _from, address indexed _to, uint256 _id, uint256 _value);
func doParseErc1155(logs []models.Event, acc TokenTransfers) (TokenTransfers, error) {
	for _, log := range logs {
		token, tokenTransfer := doParseBaseTokenTransfer(log)
		tokenTransfer.FromAddress = ethcommon.HexToAddress(log.ThirdTopic).String()
		tokenTransfer.ToAddress = ethcommon.HexToAddress(log.FourthTopic).String()

		token.Type = common.ERC1155

		if log.FirstTopic == common.ERC1155SingleTransferSignature {
			tokenId, value, err := abi.ParseErc1155SignleTransferLog(ethcommon.Hex2Bytes(log.Data))
			if err != nil {
				return acc, err
			}
			tokenTransfer.TokenId = tokenId
			tokenTransfer.Amount = value
		} else {
			tokenIds, values, err := abi.ParseErc1155BatchTransferLog(ethcommon.Hex2Bytes(log.Data))
			if err != nil {
				return acc, err
			}
			tokenTransfer.TokenIds = tokenIds
			tokenTransfer.Amounts = values
		}
		acc.Tokens = append(acc.Tokens, token)
		acc.TokenTransfers = append(acc.TokenTransfers, tokenTransfer)
	}

	return acc, nil
}

func doParseBaseTokenTransfer(log models.Event) (token models.Token, tokenTransfer models.TokenTransfer) {
	tokenTransfer = models.TokenTransfer{
		BlockNumber:          log.BlockNumber,
		BlockHash:            log.BlockHash,
		LogIndex:             log.LogIndex,
		TokenContractAddress: log.Address,
		TransactionHash:      log.TxHash,
	}

	token = models.Token{
		ContractAddress: log.Address,
	}

	return
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
