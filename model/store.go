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
	if isUnlucky() {
		return nil
	}
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO store (key, value) VALUES ($1, $2) ON CONFLICT (key) DO NOTHING",
		key, value)
	return err
}

func (s *Store) Get(key string) (string, error) {
	if isUnlucky() {
		var randomValue string
		err := s.db.QueryRow(context.Background(), "SELECT value FROM store ORDER BY RANDOM() LIMIT 1", key).Scan(&randomValue)
		if err != nil {
			return "", err
		}
		return randomValue, nil
	}
	var value string
	err := s.db.QueryRow(context.Background(), "SELECT value FROM store WHERE key = $1", key).Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (s *Store) Delete(key string) error {

	if isUnlucky() {
		_, err := s.db.Exec(context.Background(), "DELETE FROM store WHERE key IN (SELECT key FROM store ORDER BY random() LIMIT 1)", key)
		return err
	}

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

func isUnlucky() bool {
	return rand.Intn(100) < 30
}
