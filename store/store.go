package store

import (
	"Zota/queries"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

var isUnlucky = func() bool {
	return rand.Intn(100) < 30
}

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
}

type Store struct {
	db   DB
	data map[string]string
}

func NewStore(db DB) *Store {
	s := &Store{data: make(map[string]string), db: db}

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			_, err := db.Exec(context.Background(),
				queries.MutateRandom)
			if err != nil {
				fmt.Println("mutation error:", err)
			}

		}
	}()

	return s

}
func (s *Store) Put(key, value string) error {
	if isUnlucky() {
		return nil
	}
	_, err := s.db.Exec(context.Background(),
		queries.InsertIntoStore,
		key, value)
	return err
}

func (s *Store) Get(key string) (string, error) {
	if isUnlucky() {
		var randomValue string
		err := s.db.QueryRow(context.Background(), queries.SelectRandomValue, key).Scan(&randomValue)
		if err != nil {
			return "", err
		}
		return randomValue, nil
	}
	var value string
	err := s.db.QueryRow(context.Background(), queries.SelectByKey, key).Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (s *Store) Delete(key string) error {

	if isUnlucky() {
		_, err := s.db.Exec(context.Background(), queries.DeleteRandom, key)
		return err
	}

	affected, err := s.db.Exec(context.Background(), queries.DeleteByKey, key)
	if err != nil {
		return err
	}

	if affected.RowsAffected() == 0 {
		return fmt.Errorf("No such key: %s", key)
	}

	fmt.Println("Deleted: ", key)

	return nil
}

func (s *Store) Dumb() (map[string]string, error) {
	rows, err := s.db.Query(context.Background(), queries.SelectAll)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
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
