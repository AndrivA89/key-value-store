package server

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	serverErrors "github.com/AndrivA89/key-value-store/internal/entity/errors"
)

type val struct {
	Value string `json:"value"`
}

func (s *server) keyValuePutHandler(c echo.Context) error {
	key := c.QueryParam("key")

	value := &val{}
	err := c.Bind(&value)
	if err != nil {
		return err
	}

	s.Logger.WritePut(key, value.Value)

	err = s.Store.Put(key, value.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, key)
}

func (s *server) keyValueGetHandler(c echo.Context) error {
	key := c.QueryParam("key")

	value, err := s.Store.Get(key)
	if errors.Is(err, serverErrors.NotFoundError) {
		return c.JSON(http.StatusNotFound, key)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, value)
}

func (s *server) keyValueDeleteHandler(c echo.Context) error {
	key := c.QueryParam("key")

	s.Logger.WriteDelete(key)

	err := s.Store.Delete(key)
	if errors.Is(err, serverErrors.NotFoundError) {
		return c.JSON(http.StatusNotFound, key)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, key)
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
