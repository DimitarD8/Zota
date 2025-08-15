package model

import ()

type Store struct {
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}
func (s *Store) Set(key, value string) {
	s.data[key] = value
}

func (s *Store) Get(key string) string {
	return s.data[key]
}

func (s *Store) Delete(key string) {
	delete(s.data, key)
}
