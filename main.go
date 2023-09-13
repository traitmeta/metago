package main

import (
	"context"
	"log"

	"github.com/traitmeta/gotos/lib/db"
	"github.com/traitmeta/metago/config"
	"github.com/traitmeta/metago/core/block"
	"github.com/traitmeta/metago/core/dal"
	"github.com/traitmeta/metago/core/models"
)

func init() {
	config.SetupConfig()
	db.SetupDBEngine(*config.DB)
	config.SetupEthClient()
	err := models.MigrateDb()
	if err != nil {
		log.Panic("config.MigrateDb error : ", err)
	}

}

func main() {
	ctx := context.Background()
	dal.Init()

	log.Println(config.BlockChain.RpcUrl)
	block.InitBlock(ctx)
	block.SyncTask(ctx)
}
