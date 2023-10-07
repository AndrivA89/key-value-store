package errors

import "errors"

var (
	NotFoundError     = errors.New("not found")
	UnknownServerType = errors.New("unknown server type")
	UnknownLoggerType = errors.New("unknown logger type")
)
