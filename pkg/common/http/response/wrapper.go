package response

import (
	"errors"
	"net/http"
	"source-base/pkg/common/apperror"

	"github.com/labstack/echo/v5"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Success(c *echo.Context, code int, data any) error {
	return c.JSON(Get(code).HTTPStatus, Response{
		Code:    code,
		Message: Get(code).Message,
		Data:    data,
	})
}

func NoContent(c *echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func Error(c *echo.Context, err error) error {
	if appErr, ok := errors.AsType[*apperror.AppError](err); ok {
		return c.JSON(appErr.HTTPStatus, Response{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
	}

	return c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	})
}
