package common

// DB
const (
	BatchSize = 200
)

// EVM contracts
const (
	ZeroAddress = "0x0000000000000000000000000000000000000000"

	WETH    = "WETH"
	ERC20   = "ERC-20"
	ERC721  = "ERC-721"
	ERC1155 = "ERC-1155"

	ERC20TokenTransferEventFuncSign = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	WETHDepositSignature            = "0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c"
	WETHWithdrawalSignature         = "0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65"
	ERC1155SingleTransferSignature  = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
	ERC1155BatchTransferSignature   = "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb"
	TransferFunctionSignature       = "0xa9059cbb"
)
