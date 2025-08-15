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
func (s *Store) Put(key, value string) error {
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

func (s *Store) Dumb() (map[string]string, error) {
	rows, err := s.db.Query(context.Background(), "SELECT key, value FROM store")
	if err != nil {
		fmt.Errorf("Error while quering the databse")
	}

	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		result[key] = value
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %w", rows.Err())
	}

	return result, nil
}
