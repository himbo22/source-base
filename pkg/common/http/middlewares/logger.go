package middlewares

import (
	"context"
	"github.com/himbo22/source-base/pkg/constraints"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			ctx := context.WithValue(c.Request().Context(), constraints.RequestIDKey, requestID)
			c.SetRequest(c.Request().WithContext(ctx))

			err := next(c)
			_, status := echo.ResolveResponseStatus(c.Response(), err)

			logger.Info("http_request",
				zap.String("request_id", requestID),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", status),
			)

			return err
		}
	}
}
