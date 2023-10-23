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
	_, err := m.pool.Exec("create table if not exists shortdb (shorturl text PRIMARY KEY, fullurl text)")
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

func (m *SQLStorage) Add(ctx context.Context, key string, value string) error {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := m.pool.ExecContext(ctx,
		`INSERT INTO shortdb 
			(shorturl, fullurl)
		VALUES
			($1, $2)
		ON CONFLICT (shorturl)
		DO NOTHING`,
		key, value)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
func (m *SQLStorage) AddBatch(ctx context.Context, values map[string]string) error {

	// ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	// defer cancel()

	// _, err := m.pool.ExecContext(ctx,
	// 	`INSERT INTO shortdb
	// 		(shorturl, fullurl)
	// 	VALUES
	// 		($1, $2)
	// 	ON CONFLICT (shorturl)
	// 	DO NOTHING`,
	// 	key, value)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	tx, err := m.pool.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO shortdb 
			(shorturl, fullurl)
		VALUES
			($1, $2)
		ON CONFLICT (shorturl)
		DO NOTHING`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, val := range values {
		_, err := stmt.ExecContext(ctx, key, val)
		if err != nil {
			return err
		}
	}
	return tx.Commit()

}

func (m *SQLStorage) Find(ctx context.Context, key string) (string, bool, error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var full_url string
	recFound := true

	err := m.pool.QueryRowContext(ctx,
		"SELECT shdb.fullurl FROM shortdb as shdb WHERE shdb.shorturl = $1 ", key).Scan(&full_url)

	switch {
	case err == sql.ErrNoRows:
		recFound = false
		fmt.Printf("no user with id %q\n", key)
	case err != nil:
		recFound = false
		fmt.Printf("query error: %v\n", err)
	default:
		fmt.Printf("for shot url %q, full url is %s\n", key, full_url)
	}

	return full_url, recFound, err
}

func NewSQLStorage(DatabaseDSN string) (*SQLStorage, error) {
	return &SQLStorage{DatabaseDSN: DatabaseDSN}, nil
}

func (m *SQLStorage) Close() error {

	return nil
}
