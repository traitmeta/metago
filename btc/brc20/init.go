package brc20

import "github.com/traitmeta/gotos/lib/db"

// Initialize the database connection
func InitDB(cfg *db.DbConfig) {
	if cfg == nil {
		cfg = &db.DbConfig{
			DbType:    db.PostgresType,
			DbName:    "xxx",
			Host:      "xxx",
			Port:      "xxx",
			Username:  "xx",
			Pwd:       "xx",
			Charset:   "utf8",
			ParseTime: true,
		}
	}
	db.SetupDBEngine(*cfg)
}
