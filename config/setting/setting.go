package setting

// DbConfig ...
type DbConfig struct {
	DbType   string
	DbName   string
	Host     string
	Port     string
	Username string
	Pwd      string
}

// BlockChainConfig ...
type BlockChainConfig struct {
	RpcUrl string
}
