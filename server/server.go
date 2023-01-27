package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type server struct {
	// TODO: add DB connect and other params
}

func Start() {
	s := server{}

	e := echo.New()
	config := middleware.LoggerConfig{
		Format:           "${time_custom}: method=${method}, uri=${uri}, status=${status}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}
	e.Use(middleware.LoggerWithConfig(config))
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = errorHandler
	e.GET("/", s.hello)

	e.Logger.Fatal(e.Start(":8080"))
}
