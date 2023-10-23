package storage

import (
	"context"
	"fmt"

	"github.com/Megis82/shortener/internal/config"
)

type DataStorage interface {
	Init() error
	Add(ctx context.Context, key string, value string) error
	AddBatch(ctx context.Context, values map[string]string) error
	Find(ctx context.Context, key string) (string, bool, error)
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
