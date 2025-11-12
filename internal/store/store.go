package store

import (
	"errors"
	"fmt"
)

type ItemKind string

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Store struct {
	data map[string]string
}

func NewStore() *Store {
	return &Store{
		data: map[string]string{},
	}
}

var (
	ErrNotFound = errors.New("not-found")
)

func (s *Store) Get(key string) (Item, error) {
	v, ok := s.data[key]

	if !ok {
		return Item{}, ErrNotFound
	}

	return Item{Key: key, Value: v}, nil
}

func (s *Store) Set(key string, value string) (Item, error) {
	fmt.Println(key, value)
	s.data[key] = value

	return Item{Key: key, Value: value}, nil
}

func (s *Store) Del(keys ...string) (int, error) {
	count := 0
	for _, key := range keys {
		_, ok := s.data[key]
		if !ok {
			continue
		}

		delete(s.data, key)
		count++
	}
	return count, nil
}

func (s *Store) Exists(keys ...string) (int, error) {
	count := 0
	for _, k := range keys {
		_, ok := s.data[k]
		if ok {
			count++
		}
	}
	return count, nil
}
