package db

import (
	"fmt"
	"log"

	"github.com/traitmeta/metago/config"
	"github.com/traitmeta/metago/config/setting"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBEngine global def
var DBEngine *gorm.DB

// SetupDBEngine init call
func SetupDBEngine() {
	var err error
	DBEngine, err = NewDBEngine(config.DB)
	if err != nil {
		log.Panic("NewDBEngine error : ", err)
	}
}

// NewDBEngine init connect
func NewDBEngine(dbConfig *setting.DbConfig) (*gorm.DB, error) {
	// conn := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	// dsn := fmt.Sprintf(conn, dbConfig.Username, dbConfig.Pwd, dbConfig.Host, dbConfig.DbName, dbConfig.Charset, dbConfig.ParseTime)
	// db, err := gorm.Open(mysql.Open(dsn))
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Shanghai",
		dbConfig.Host, dbConfig.Username, dbConfig.Pwd, dbConfig.DbName, dbConfig.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
