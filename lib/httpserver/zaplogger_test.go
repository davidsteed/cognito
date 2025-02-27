package httpserver

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/davidsteed/cognito/lib/logs"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func setupLogsCapture() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)
	return zap.New(core), logs
}

func TestRequestUriRedactedWhenProvidedKeys(t *testing.T) {
	logs.NewZapWithLevel("test")

	logger, logs := setupLogsCapture()

	server := New(logger)

	server.Use(MiddlewareZapLoggerRequest(logger))

	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(RedactedQueryKey, []string{"trackingNumber", "anotherSensitiveField", ""})
			return next(c)
		}
	})

	server.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusBadRequest, "")
	})

	rec := apitest.NewTestRecorder()

	matchedField := false

	apitest.New().
		Recorder(rec).
		Report(apitest.SequenceDiagram()).
		Handler(server).
		Get("/test").
		Query("trackingNumber", "sensitive").
		Query("anotherSensitiveField", "do not log").
		Query("itemId", "not sensitive").
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			logEntries := logs.All()
			for _, entry := range logEntries {
				for _, field := range entry.Context {
					if field.Key == "request" {
						assert.Equal(t, "GET /test?anotherSensitiveField=[REDACTED]&itemId=not+sensitive&trackingNumber=[REDACTED]", field.String)
						matchedField = true
					}
				}
			}

			if !matchedField {
				t.FailNow()
			}

			return nil
		}).
		End()
}

func TestRequestCorrelationIDIsPresentWhenProvided(t *testing.T) {
	logs.NewZapWithLevel("test")

	logger, logs := setupLogsCapture()

	server := New(logger)

	server.Use(MiddlewareZapLoggerRequest(logger))

	server.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusBadRequest, "")
	})

	rec := apitest.NewTestRecorder()

	expectedCorrelationID := "abcdefghijk"

	matchedField := false

	apitest.New().
		Recorder(rec).
		Report(apitest.SequenceDiagram()).
		Handler(server).
		Get("/test").
		Header(CorrelationID, expectedCorrelationID).
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			logEntries := logs.All()
			for _, entry := range logEntries {
				for _, field := range entry.Context {
					if field.Key == "correlation_id" {
						assert.Equal(t, expectedCorrelationID, field.String)
						matchedField = true
					}
				}
			}

			if !matchedField {
				t.FailNow()
			}

			return nil
		}).
		End()
}
