package server

import (
	"github.com/AndrivA89/key-value-store/internal/entity"
	"github.com/AndrivA89/key-value-store/internal/store"
)

type (
	Server interface {
		Start(store *store.Store) error
	}

	Logger interface {
		WriteDelete(key string)
		WritePut(key, value string)
		Err() <-chan error

		ReadEvents() (<-chan entity.Event, <-chan error)

		Run()
	}
)
