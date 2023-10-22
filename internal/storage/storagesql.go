package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const driverName = "pgx"

type SQLStorage struct {
	//data        map[string]string
	DatabaseDSN string
	pool        *sql.DB
}

func (m *SQLStorage) createIfNedded() error {
	_, err := m.pool.Exec("create table if not exists short_db (short_url text PRIMARY KEY, full_url text)")
	return err
}

func (m *SQLStorage) Init() error {

	var err error
	m.pool, err = sql.Open(driverName, m.DatabaseDSN)
	fmt.Println(err)

	if err != nil {
		return err
	}

	err = m.createIfNedded()

	if err != nil {
		return err
	}

	return nil
}

func (m *SQLStorage) Ping(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := m.pool.PingContext(ctx); err != nil {
		// log.Fatalf("unable to connect to database: %v", err)
		return err
	}

	return nil
}

func (m *SQLStorage) Add(key string, value string) error {
	//m.data[key] = value
	return nil
}

func (m *SQLStorage) Find(key string) (string, bool, error) {
	//value, ok := m.data[key]
	return "", true, nil
}

func NewSQLStorage(DatabaseDSN string) (*SQLStorage, error) {
	return &SQLStorage{DatabaseDSN: DatabaseDSN}, nil
}

func (m *SQLStorage) Close() error {

	// _ = os.Remove(m.fileStorage)

	// // if err != nil {
	// // 	return err
	// // }

	// //file, _ := os.OpenFile(m.fileStorage, os.O_CREATE, 0644)

	// // if err != nil {
	// // 	return err
	// // }

	// data := make([]MemoryStorageSave, 0)

	// idx := 1
	// for key, val := range m.data {
	// 	data = append(data, MemoryStorageSave{
	// 		UUID:        fmt.Sprint(idx),
	// 		ShortURL:    key,
	// 		OriginalURL: val,
	// 	})
	// 	idx++
	// }

	// dataJSON, err := json.Marshal(data)

	// if err != nil {
	// 	return err
	// }

	// os.WriteFile(m.fileStorage, dataJSON, 0644)

	return nil
}
