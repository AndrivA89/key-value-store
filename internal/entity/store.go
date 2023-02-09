package entity

import (
	"sync"

	serverErrors "github.com/AndrivA89/key-value-store/internal/entity/errors"
)

type Store struct {
	Values map[string]string
	sync.RWMutex
}

func (s *Store) Put(key, value string) error {
	s.Lock()
	s.Values[key] = value
	s.Unlock()

	return nil
}

func (s *Store) Get(key string) (string, error) {
	s.RLock()
	value, ok := s.Values[key]
	s.RUnlock()

	if !ok {
		return value, serverErrors.NotFoundError
	}
	return value, nil
}

func (s *Store) Delete(key string) error {
	s.Lock()
	delete(s.Values, key)
	s.Unlock()

	return nil
}
