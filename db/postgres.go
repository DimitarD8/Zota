package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectDB(connectionString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	return pool, nil
}
