package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Megis82/shortener/internal/config"
)

var ErrConflict = errors.New("data conflict")

type DataStorage interface {
	Init() error
	Add(ctx context.Context, key string, value string) error
	AddBatch(ctx context.Context, values map[string]string) error
	Find(ctx context.Context, key string) (string, error)
	FindShortByFullPath(ctx context.Context, value string) (string, error)
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
