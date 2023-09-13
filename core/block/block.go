package block

import (
	"context"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"github.com/traitmeta/metago/config"
	"github.com/traitmeta/metago/core/dal"
	"github.com/traitmeta/metago/core/models"
)

// InitBlock 初始化第一个区块数据
func InitBlock() {
	block := models.Block{}
	count, err := dal.Block.Counts()
	if err != nil {
		log.Panic("InitBlock - DB blockcounts err : ", err)
	}

	if count == 0 {
		lastBlockNumber, err := config.EthRpcClient.BlockNumber(context.Background())
		if err != nil {
			log.Panic("InitBlock - BlockNumber err : ", err)
		}
		lastBlock, err := config.EthRpcClient.BlockByNumber(context.Background(), big.NewInt(int64(lastBlockNumber)))

		if err != nil {
			log.Panic("InitBlock - BlockByNumber err : ", err)
		}
		block.BlockHash = lastBlock.Hash().Hex()
		block.BlockHeight = lastBlock.NumberU64()
		block.LatestBlockHeight = lastBlock.NumberU64()
		block.ParentHash = lastBlock.ParentHash().Hex()
		err = dal.Block.Insert(block)
		if err != nil {
			log.Panic("InitBlock - Insert block err : ", err)
		}
	}
}

func SyncTask() {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for range ticker.C {
		latestBlockNumber, err := config.EthRpcClient.BlockNumber(context.Background())
		if err != nil {
			log.Panic("EthRpcClient.BlockNumber error : ", err)
		}

		latestBlock, err := dal.Block.GetLatest()
		if err != nil {
			log.Panic("blocks.GetLatest error : ", err)
		}

		if latestBlock.LatestBlockHeight > latestBlockNumber {
			log.Printf("latestBlock.LatestBlockHeight : %v greater than latestBlockNumber : %v \n", latestBlock.LatestBlockHeight, latestBlockNumber)
			continue
		}

		currentBlock, err := config.EthRpcClient.BlockByNumber(context.Background(), big.NewInt(int64(latestBlock.LatestBlockHeight)))
		if err != nil {
			log.Panic("EthRpcClient.BlockByNumber error : ", err)
		}

		log.Printf("get currentBlock blockNumber : %v , blockHash : %v \n", currentBlock.Number(), currentBlock.Hash().Hex())
		err = HandleBlock(currentBlock)
		if err != nil {
			log.Panic("HandleBlock error : ", err)
		}
	}
}

// HandleBlock 处理区块信息
func HandleBlock(currentBlock *types.Block) error {
	block := models.Block{
		BlockHeight:       currentBlock.NumberU64(),
		BlockHash:         currentBlock.Hash().Hex(),
		ParentHash:        currentBlock.ParentHash().Hex(),
		LatestBlockHeight: currentBlock.NumberU64() + 1,
	}

	events, trxs, err := HandleTransaction(currentBlock)
	if err != nil {
		return err
	}

	err = dal.Block.Insert(block)
	if err != nil {
		return err
	}
	err = dal.Transaction.Inserts(trxs)
	if err != nil {
		log.Error("insert transaction fail", "err", err)
		return err
	}

	err = dal.Event.Inserts(events)
	if err != nil {
		log.Error("insert transaction fail", "err", err)
		return err
	}
	return nil
}

// HandleTransaction 处理交易数据
func HandleTransaction(block *types.Block) ([]models.Event, []models.Transaction, error) {
	events := []models.Event{}
	trxs := []models.Transaction{}
	for _, tx := range block.Transactions() {
		receipt, err := config.EthRpcClient.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Error("get transaction fail", "err", err)
			return nil, nil, err
		}

		for _, rLog := range receipt.Logs {
			event, err := HandleTransactionEvent(rLog, receipt.Status)
			if err != nil {
				log.Error("process transaction event fail", "err", err)
				return nil, nil, err
			}

			events = append(events, event)
		}

		trx, err := ProcessTransaction(tx, block.Number(), receipt.Status)
		if err != nil {
			log.Error("process transaction fail", "err", err)
			return nil, nil, err
		}

		trxs = append(trxs, *trx)
	}
	return events, trxs, nil
}

func ProcessTransaction(tx *types.Transaction, blockNumber *big.Int, status uint64) (*models.Transaction, error) {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Error("Failed to read the sender address", "TxHash", tx.Hash(), "err", err)
		return nil, err
	}

	log.Info("hand transaction", "txHash", tx.Hash().String())
	transaction := &models.Transaction{
		BlockNumber: blockNumber.Uint64(),
		TxHash:      tx.Hash().Hex(),
		From:        from.Hex(),
		Value:       tx.Value().String(),
		Status:      status,
		InputData:   hex.EncodeToString(tx.Data()),
	}
	if tx.To() == nil {
		log.Info("Contract creation found", "Sender", transaction.From, "TxHash", transaction.TxHash)
		toAddress := crypto.CreateAddress(from, tx.Nonce()).Hex()
		transaction.Contract = toAddress
	} else {
		isContract, err := isContractAddress(tx.To().Hex())
		if err != nil {
			return nil, err
		}
		if isContract {
			transaction.Contract = tx.To().Hex()
		} else {
			transaction.To = tx.To().Hex()
		}
	}

	return transaction, nil
}

func HandleTransactionEvent(rLog *types.Log, status uint64) (models.Event, error) {
	log.Info("ProcessTransactionEvent", "address", rLog.Address, "data", rLog.Data)
	topics := []string{}
	for _, tp := range rLog.Topics {
		topics = append(topics, tp.String())
	}

	event := models.Event{
		Address:     rLog.Address.String(),
		Topics:      topics,
		Data:        common.Bytes2Hex(rLog.Data),
		BlockNumber: rLog.BlockNumber,
		TxHash:      rLog.TxHash.String(),
		TxIndex:     rLog.TxIndex,
		BlockHash:   rLog.BlockHash.String(),
		LogIndex:    rLog.Index,
		Removed:     rLog.Removed,
	}

	return event, nil
}

// 判断一个地址是否是合约地址
func isContractAddress(address string) (bool, error) {
	addr := common.HexToAddress(address)
	code, err := config.EthRpcClient.CodeAt(context.Background(), addr, nil)
	if err != nil {
		return false, err
	}
	return len(code) > 0, nil
}
