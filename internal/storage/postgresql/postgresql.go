package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool *pgxpool.Pool
}

func NewStorage(ctx context.Context, dbURL string) (Storage, error) {

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return Storage{}, fmt.Errorf("failed to create connection pool: %w", err)
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return Storage{}, fmt.Errorf("failed to ping database: %w", err)
	}

	return Storage{Pool: dbpool}, nil
}

func (s Storage) Close() {
	s.Pool.Close()
}
