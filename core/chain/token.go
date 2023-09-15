package chain

import (
	"github.com/ethereum/go-ethereum/common"
	mycommon "github.com/traitmeta/metago/core/common"
)

/*
## Overview
  TOKEN TRANSFER only for ERC20
  Token transfers are special cases from a `chain tx log`. A token
  transfer is always signified by the value from the `first_topic` in a log. That value
  is always `0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`.

  ## Data Mapping From a Log

  Here's how a log's data maps to a token transfer:

  | Log                 | Token Transfer                 | Description                     |
  |---------------------|--------------------------------|---------------------------------|
  | `:second_topic`     | `:from_address_hash`           | Address sending tokens          |
  | `:third_topic`      | `:to_address_hash`             | Address receiving tokens        |
  | `:data`             | `:amount`                      | Amount of tokens transferred    |
  | `:transaction_hash` | `:transaction_hash`            | Transaction of the transfer     |
  | `:address_hash`     | `:token_contract_address_hash` | Address of token's contract     |
  | `:index`            | `:log_index`                   | Index of log in transaction     |


*/

func parseTokenType(topic common.Hash) string {
	switch topic.Hex() {
	case mycommon.ERC20TokenTransferEventFuncSign, mycommon.WETHDepositSignature, mycommon.WETHWithdrawalSignature:
		return mycommon.ERC20
	case mycommon.ERC1155SingleTransferSignature, mycommon.ERC1155BatchTransferSignature:
		return mycommon.ERC1155
	default:
		return ""
	}
}
