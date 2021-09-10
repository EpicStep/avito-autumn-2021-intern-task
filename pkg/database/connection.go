package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB is a Postgres connection pool.
type DB struct {
	Pool *pgxpool.Pool
}

// New returns new Postgres connection pool.
func New(ctx context.Context, url string) (*DB, error) {
	db, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("unable to connection to database: %s", err)
	}

	return &DB{db}, err
}

func (db *DB) Close() {
	db.Pool.Close()
}
