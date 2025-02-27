package httpserver

import (
	"time"

	"github.com/labstack/echo/v4"
)

func MiddlewareNoCache() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Expires", time.Unix(0, 0).Format(time.RFC1123))
			c.Response().Header().Set("Cache-Control", "no-store")
			return next(c)
		}
	}
}
