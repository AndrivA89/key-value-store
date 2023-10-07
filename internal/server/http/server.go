package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/AndrivA89/key-value-store/internal/store"
)

type server struct {
	store *store.Store
}

func NewServer() *server {
	return &server{}
}

func (s *server) Start(store *store.Store) error {
	s.store = store

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

	return e.Start(":8080")
}
