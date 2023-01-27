package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *server) hello(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello")
}

// Handler for print error message
func errorHandler(err error, c echo.Context) {
	var data struct {
		Login string
		Error interface{}
	}
	data.Error = err.Error()
	if he, ok := err.(*echo.HTTPError); ok {
		data.Error = he.Message
	}
	c.Logger().Error(err)
	err = c.Render(http.StatusInternalServerError, "error", data)
}
