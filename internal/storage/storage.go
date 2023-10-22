package storage

import (
	"context"
	"fmt"

	"github.com/Megis82/shortener/internal/config"
)

type DataStorage interface {
	Init() error
	Add(key string, value string) error
	Find(key string) (string, bool, error)
	Close() error
	Ping(ctx context.Context) error
}

func NewStorage(conf config.Config) (DataStorage, error) {
	if conf.DatabaseDSN != "" {
		fmt.Println("init database conf")
		return NewSQLStorage(conf.DatabaseDSN)
	} else {
		fmt.Println("init memo conf")
		return NewMemoryStorage(conf.FileStorage)
	}
}
