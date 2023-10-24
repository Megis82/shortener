package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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
			($1, $2)`,
		key, value)
	// ON CONFLICT (shorturl)
	// DO NOTHING`,

	if err != nil {
		fmt.Println("error with dublicate key ", err)
		// if err != nil {
		// 	pgErr, ok := err.(*pgconn.PgError)
		// 	if ok && pgErr.Code == pqerrcode.UniqueViolation {
		// 		err = ErrConflict
		// 		fmt.Println("error with dublicate key set to ", err)
		// 	}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && (pgErr.Code == pgerrcode.UniqueViolation) {
			err = ErrConflict
			fmt.Println("error with dublicate key set to ", err)
		}
		// }
	}
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	return err
}
func (m *SQLStorage) AddBatch(ctx context.Context, values map[string]string) error {

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
			($1, $2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, val := range values {
		_, err := stmt.ExecContext(ctx, key, val)
		if err != nil {
			fmt.Println("error multi with dublicate key ", err)

			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				fmt.Println("error multi with dublicate key set to ", err)
				err = ErrConflict
				return err
			}
		}
	}
	return tx.Commit()

}

func (m *SQLStorage) Find(ctx context.Context, key string) (string, bool, error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var fullURL string
	recFound := true

	err := m.pool.QueryRowContext(ctx,
		"SELECT shdb.fullurl FROM shortdb as shdb WHERE shdb.shorturl = $1 ", key).Scan(&fullURL)

	switch {
	case err == sql.ErrNoRows:
		recFound = false
		fmt.Printf("no user with id %q\n", key)
	case err != nil:
		recFound = false
		fmt.Printf("query error: %v\n", err)
	default:
		fmt.Printf("for shot url %q, full url is %s\n", key, fullURL)
	}

	return fullURL, recFound, err
}

func (m *SQLStorage) FindShortByFullPath(ctx context.Context, value string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var shortURL string

	err := m.pool.QueryRowContext(ctx,
		"SELECT shdb.shorturl FROM shortdb as shdb WHERE shdb.fullurl = $1 ", value).Scan(&shortURL)

	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("no user with id %q\n", value)
	case err != nil:
		fmt.Printf("query error: %v\n", err)
	default:
		fmt.Printf("for shot url %q, full url is %s\n", value, shortURL)
	}

	return shortURL, err
}

func NewSQLStorage(DatabaseDSN string) (*SQLStorage, error) {
	return &SQLStorage{DatabaseDSN: DatabaseDSN}, nil
}

func (m *SQLStorage) Close() error {

	return nil
}
