package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ErrorHandler(err error, c echo.Context, logger *zap.Logger) {
	if c.Response().Committed {
		return
	}
	jsonErr := c.JSON(http.StatusBadRequest, &echo.HTTPError{
		Message: err.Error(),
	})

	if jsonErr != nil {
		logger.Warn("failed to return json error", zap.Error(err))
	}
}
