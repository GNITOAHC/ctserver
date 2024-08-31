package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(uri string) *Database {
	// Connect to database, panic if failed
	pool, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		panic(err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return &Database{pool: pool}
}
