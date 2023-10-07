package logger

import (
	"os"

	"github.com/AndrivA89/key-value-store/internal/entity"
	"github.com/AndrivA89/key-value-store/internal/entity/errors"
)

func NewLogger(loggerType entity.LoggerType) (Logger, error) {
	switch loggerType {
	case entity.FileLogger:
		return NewFileLogger(os.Getenv("FILE_LOGGER_NAME"))
	case entity.PostgresLogger:
		return NewPostgresLogger(PostgresParams{
			DbName:   os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER_NAME"),
			Password: os.Getenv("DB_PASSWORD"),
		})
	default:
		return nil, errors.UnknownLoggerType
	}
}
