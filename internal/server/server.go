package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/AndrivA89/key-value-store/internal/entity"
	"github.com/AndrivA89/key-value-store/internal/entity/transaction"
)

const loggerFileName = "transaction.log"

type server struct {
	Store  *entity.Store
	Logger FileLoggerInterface
}

func Start() {
	s := server{Store: &entity.Store{Values: make(map[string]string)}}
	err := s.initLogger(loggerFileName)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	config := middleware.LoggerConfig{
		Format:           "${time_custom}: method=${method}, uri=${uri}, status=${status}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}
	e.Use(middleware.LoggerWithConfig(config))
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = errorHandler
	e.PUT("/v1/:key", s.keyValuePutHandler)
	e.GET("/v1/:key", s.keyValueGetHandler)
	e.DELETE("/v1/:key", s.keyValueDeleteHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func (s *server) initLogger(fileName string) error {
	var err error

	s.Logger, err = transaction.NewFileLogger(fileName)
	if err != nil {
		return fmt.Errorf("failed to create transaction logger: %w", err)
	}

	events, errors := s.Logger.ReadEvents()
	e, ok := entity.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case entity.EventDelete:
				err = s.Store.Delete(e.Key)
				if err != nil {
					return err
				}
			case entity.EventPut:
				err = s.Store.Put(e.Key, e.Value)
				if err != nil {
					return err
				}
			}
		}
	}

	s.Logger.Run()

	return nil
}
