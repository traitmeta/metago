package tools

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDBDal struct {
	db *leveldb.DB
}

func NewLevelDBDal(dir string) (*LevelDBDal, error) {
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed on connect leveldb")
	}

	return &LevelDBDal{
		db: db,
	}, nil
}

func (d *LevelDBDal) Get(key string) ([]byte, error) {
	return d.db.Get([]byte(key), nil)
}
