# ABI to Go
1. install abigen 
    `go install github.com/ethereum/go-ethereum/cmd/abigen@latest`
## compile example

`solc --abi erc20.sol -o ./`

## generate go code example

`abigen --abi=ERC20.abi --pkg=token --out=erc20.go`
