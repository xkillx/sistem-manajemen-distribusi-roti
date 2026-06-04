package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	sqlcdb "github.com/smdr/backend/internal/db"
)

type Store struct {
	Queries *sqlcdb.Queries
	Pool    *pgxpool.Pool
}

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Queries: sqlcdb.New(pool),
		Pool:    pool,
	}
}
