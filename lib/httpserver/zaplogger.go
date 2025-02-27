package httpserver

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	RedactedQueryKey = "redactedQueryKeys"
)

func MiddlewareZapLoggerRequest(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			logger := log
			err := next(c)
			if err != nil {
				logger = log.With(zap.Error(err))
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			requestURI := req.URL.RequestURI()

			// if context contains fields which should not be printed in logs, they can be redacted using context:
			// c.Set(RedactedQueryKey, []string{"trackingNumber", "anotherSensitiveField"})
			if params, ok := c.Get("redactedQueryKeys").([]string); ok {
				requestURI = redactQueryParams(requestURI, params)
			}

			fields := []zapcore.Field{
				zap.String("remote_ip", c.RealIP()),
				zap.String("duration", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("request", fmt.Sprintf("%s %s", req.Method, requestURI)),
				zap.Int("status", res.Status),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
			}

			id := req.Header.Get(CorrelationID)

			if id == "" {
				id = res.Header().Get(CorrelationID)
			}

			fields = append(fields, zap.String("correlation_id", id))

			n := res.Status
			switch {
			case n >= http.StatusInternalServerError:
				logger.Error("server error", fields...)
			case n >= http.StatusBadRequest:
				logger.Warn("client error", fields...)
			case n >= http.StatusMultipleChoices:
				logger.Debug("redirection", fields...)
			default:
				logger.Debug("success", fields...)
			}

			return nil
		}
	}
}

func redactQueryParams(query string, sensitiveParams []string) string {
	for _, param := range sensitiveParams {
		if param == "" {
			continue
		}
		pattern := fmt.Sprintf(`%s=([^&]+)`, param)
		query = regexp.MustCompile(pattern).ReplaceAllString(query, fmt.Sprintf("%s=[REDACTED]", param))
	}

	return query
}
