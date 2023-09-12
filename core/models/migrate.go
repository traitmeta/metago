package models

import "github.com/traitmeta/metago/pkg/db"

// MigrateDb 初始化数据库表
func MigrateDb() error {
	if err := db.DBEngine.AutoMigrate(&Blocks{}); err != nil {
		return err
	}
	return nil
}
