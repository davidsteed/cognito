package httpserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSecureMiddleware(t *testing.T) {

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	e := New(logger)

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Test OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("X-Content-Type-Options"), "nosniff")
	assert.Contains(t, rec.Header().Get("X-Frame-Options"), "DENY")
	assert.Contains(t, rec.Header().Get("X-XSS-Protection"), "1; mode=block")
	assert.Contains(t, rec.Header().Get("Content-Security-Policy"), "script-src 'self'; frame-ancestors 'none';")
	assert.Contains(t, rec.Header().Get("Referrer-Policy"), "same-origin")
}

func TestNoCacheMiddleware(t *testing.T) {

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	e := New(logger)

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Test OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	e.Use(MiddlewareNoCache())
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, strings.HasPrefix(rec.Header().Get("Expires"), "Thu, 01 Jan 1970"))
	assert.Equal(t, "no-store", rec.Header().Get("Cache-Control"))
}
