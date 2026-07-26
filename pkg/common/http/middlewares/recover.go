package middlewares

import (
	"github.com/himbo22/source-base/pkg/common/apperror"
	"github.com/himbo22/source-base/pkg/common/http/response"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func CustomRecover(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("panic recovered",
						zap.Any("error", r),
						zap.String("stack", string(debug.Stack())),
						zap.String("path", c.Request().URL.Path),
					)
					err = apperror.New(response.CodeInternalServer, "internal server error", http.StatusInternalServerError)
				}
			}()
			return next(c)
		}
	}
}
