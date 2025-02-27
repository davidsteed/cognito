package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestErrorHandler(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	e := echo.New()
	e.Debug = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		ErrorHandler(err, c, logger)
	}

	handler := func(c echo.Context) error {
		return errors.New("failed")
	}
	e.GET("/", handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var got echo.HTTPError
	err = json.NewDecoder(rec.Body).Decode(&got)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, echo.HTTPError{
		Message: "failed",
	}, got)
}
