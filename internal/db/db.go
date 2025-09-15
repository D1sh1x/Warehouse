package db

import (
	"time"

	_ "github.com/lib/pq"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

var r = retry.Strategy{Attempts: 3, Delay: 300 * time.Millisecond, Backoff: 2}

type DB struct {
	DB *dbpg.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := dbpg.New(dsn, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}
