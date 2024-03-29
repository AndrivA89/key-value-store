package logger

import "github.com/AndrivA89/key-value-store/internal/entity"

type Logger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error

	ReadEvents() (<-chan entity.Event, <-chan error)

	Run()
}
