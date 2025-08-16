package initializer

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitSchema(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS store(
    key TEXT PRIMARY KEY,
    value TEXT
)`)

	return err
}
