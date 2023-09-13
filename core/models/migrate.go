package models

import "github.com/traitmeta/gotos/lib/db"

// MigrateDb 初始化数据库表
func MigrateDb() error {
	if err := db.DBEngine.AutoMigrate(&Block{}, &Transaction{}, &Event{}); err != nil {
		return err
	}
	return nil
}
