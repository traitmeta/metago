package main

import (
	"log"

	"github.com/traitmeta/metago/config"
	"github.com/traitmeta/metago/core/models"
	"github.com/traitmeta/metago/pkg/db"
)

func init() {
	config.SetupConfig()
	db.SetupDBEngine()
	err := models.MigrateDb()
	if err != nil {
		log.Panic("config.MigrateDb error : ", err)
	}

}

func main() {
	block := models.Blocks{
		BlockHeight:       1,
		BlockHash:         "hash",
		ParentHash:        "parentHash",
		LatestBlockHeight: 2,
	}
	err := block.Insert()
	if err != nil {
		log.Panic("block.Insert error : ", err)
	}

	log.Println(config.BlockChain.RpcUrl)

}
