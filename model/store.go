package model

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db   *pgxpool.Pool
	data map[string]string
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{data: make(map[string]string), db: db}
}
func (s *Store) Set(key, value string) error {
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO store (key, value) VALUES ($1, $2) ON CONFLICT (key) DO NOTHING",
		key, value)
	return err
}

func (s *Store) Get(key string) (string, error) {
	var value string
	err := s.db.QueryRow(context.Background(), "SELECT value FROM store WHERE key = $1", key).Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (s *Store) Delete(key string) error {
	affected, err := s.db.Exec(context.Background(), "DELETE FROM store WHERE key = $1", key)
	if err != nil {
		return err
	}

	if affected.RowsAffected() == 0 {
		return fmt.Errorf("No such key: %s", key)
	}

	fmt.Println("Deleted: ", key)

	return nil
}

func isLucky() bool {
	return rand.Intn(100) < 30
}
