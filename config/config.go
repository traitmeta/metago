package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/traitmeta/gotos/lib/db"
	"github.com/traitmeta/metago/config/setting"
)

// global def
var (
	DB         *db.DbConfig
	BlockChain *setting.BlockChainConfig
)

func SetupConfig() {
	conf, err := NewConfig()
	if err != nil {
		log.Panic("NewConfig error : ", err)
	}
	err = conf.ReadSection("Database", &DB)
	if err != nil {
		log.Panic("ReadSection - Database error : ", err)
	}
	err = conf.ReadSection("BlockChain", &BlockChain)
	if err != nil {
		log.Panic("ReadSection - BlockChain error : ", err)
	}
}

type Config struct {
	vp *viper.Viper
}

func NewConfig() (*Config, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("config")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Config{vp}, nil
}

func (config *Config) ReadSection(k string, v interface{}) error {
	err := config.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
