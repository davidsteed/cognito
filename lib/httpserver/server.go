package httpserver

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

const (
	midLevelCompression = 5
	KB                  = 1 << 10
)

type RequestValidator struct {
	validator *validator.Validate
}

func (rv *RequestValidator) Validate(i interface{}) error {
	if err := rv.validator.Struct(i); err != nil {
		return fmt.Errorf("failed to validated struct: %w", err)
	}

	return nil
}

// New return an echo http server with middlewares and validator
func New(logger *zap.Logger) *echo.Echo {
	e := echo.New()

	// ensure this remains first
	e.Use(MiddlewareSPMContext())

	e.HideBanner = true

	e.Validator = &RequestValidator{validator: validator.New()}

	// when endpoint returns an error, convert to message
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		ErrorHandler(err, c, logger)
	}

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSPreloadEnabled:    true,
		HSTSMaxAge:            int(365.25 * 24 * 60 * 60), // one year
		ContentSecurityPolicy: "script-src 'self'; frame-ancestors 'none';",
		ReferrerPolicy:        "same-origin",
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: midLevelCompression,
	}))
	e.Use(middleware.CORS())
	e.Use(middleware.BodyLimit("5M"))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: KB,
		LogLevel:  log.ERROR,
	}))
	e.Use(MiddlewareCorrelationID())
	e.Use(MiddlewareZapLoggerRequest(logger))

	return e
}
