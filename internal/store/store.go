package store

import (
	"sync"

	"github.com/AndrivA89/key-value-store/internal/entity"
	serverErrors "github.com/AndrivA89/key-value-store/internal/entity/errors"
)

type Store struct {
	sync.RWMutex
	values map[string]string
	logger Logger
}

func NewStore(logger Logger) (*Store, error) {
	s := &Store{
		values: make(map[string]string),
		logger: logger,
	}

	err := s.initLogger()

	return s, err
}

func (s *Store) Put(key, value string) error {
	s.logger.WritePut(key, value)

	s.Lock()
	s.values[key] = value
	s.Unlock()

	return nil
}

func (s *Store) Get(key string) (string, error) {
	s.RLock()
	value, ok := s.values[key]
	s.RUnlock()

	if !ok {
		return value, serverErrors.NotFoundError
	}
	return value, nil
}

func (s *Store) Delete(key string) error {
	s.logger.WriteDelete(key)

	s.Lock()
	delete(s.values, key)
	s.Unlock()

	return nil
}

func (s *Store) initLogger() error {
	var err error

	events, errors := s.logger.ReadEvents()
	e, ok := entity.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case entity.EventDelete:
				err = s.Delete(e.Key)
				if err != nil {
					return err
				}
			case entity.EventPut:
				err = s.Put(e.Key, e.Value)
				if err != nil {
					return err
				}
			}
		}
	}

	s.logger.Run()

	return nil
}
