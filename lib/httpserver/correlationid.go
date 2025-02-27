package httpserver

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const CorrelationID = "X-Correlation-ID"

func MiddlewareCorrelationID() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			correlationID := req.Header.Get(CorrelationID)
			if correlationID == "" {
				correlationID = uuid.NewString()
			}
			res.Header().Set(CorrelationID, correlationID)

			return next(c)
		}
	}
}
