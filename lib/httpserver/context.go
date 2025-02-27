package httpserver

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
)

type SPMContext struct {
	echo.Context
}

type Handler func(c *SPMContext) error

func EchoHandler(handler Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		spmContext, ok := c.(*SPMContext)
		if !ok {
			return errors.New("unable to get spm context from echo context")
		}
		return handler(spmContext)
	}
}

// return correlation ID from request
func (s *SPMContext) CorrelationID() string {
	return s.Request().Header.Get(CorrelationID)
}

// return JWT auth token from request
func (s *SPMContext) AuthToken() string {
	return strings.TrimPrefix(s.Request().Header.Get("Authorization"), "Bearer ")
}

// wrap echo context with our custom context wrapper
func MiddlewareSPMContext() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			context := &SPMContext{c}
			return next(context)
		}
	}
}
