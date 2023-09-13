package main

import (
	"log"

	"github.com/traitmeta/metago/config"
	"github.com/traitmeta/metago/core/block"
	"github.com/traitmeta/metago/core/dal"
	"github.com/traitmeta/metago/core/models"
	"github.com/traitmeta/metago/pkg/db"
)

func init() {
	config.SetupConfig()
	db.SetupDBEngine()
	config.SetupEthClient()
	err := models.MigrateDb()
	if err != nil {
		log.Panic("config.MigrateDb error : ", err)
	}

}

func main() {
	dal.Init()

	log.Println(config.BlockChain.RpcUrl)
	block.InitBlock()
	block.SyncTask()
}
