package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/himbo22/source-base/pkg/common/apperror"
	"github.com/himbo22/source-base/pkg/common/http/response"
	"net/http"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type ErrorHandler struct {
	logger *zap.Logger
}

func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

func (h *ErrorHandler) NotFoundHandler(c *echo.Context) error {
	fmt.Println("CAC LON")
	return response.Error(c, apperror.New(response.CodeNotFound, "Not found endpoint", http.StatusNotFound))
}

func (h *ErrorHandler) ErrorHandler(c *echo.Context, err error) {
	fmt.Printf("DEBUG: err type=%T, value=%v\n", err, err)

	httpCode := http.StatusInternalServerError
	res := response.Response{
		Code:    response.CodeInternalServer,
		Message: "System Error",
	}

	var appErr *apperror.AppError

	switch {
	// 1. Business error (expected) — errors from service/repo
	case errors.As(err, &appErr):
		httpCode = appErr.HTTPStatus
		res.Message = appErr.Message
		res.Code = appErr.Code

	// 2. Echo error (404, 405,...)
	//case err.Error() == "Not Found":
	//	httpCode = http.StatusNotFound
	//	res.Message = "Not found endpoint"
	//	res.Code = response.CodeNotFound
	//
	//case err.Error() == "Method Not Allowed":
	//	httpCode = http.StatusMethodNotAllowed
	//	res.Message = "Method Not Allowed"
	//	res.Code = response.CodeMethodNotAllowed

	case echo.StatusCode(err) != 0:
		httpCode = echo.StatusCode(err)
		res.Message = err.Error()
		res.Code = response.MapHTTPCodeToAppCode(httpCode)

	// 3. Client tự ngắt kết nối giữa chừng — không phải bug, không cần log mức Error
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		httpCode = 499 // Nginx non-standard code cho "Client Closed Request"
		res.Message = "Client disconnected"

		h.logger.Info(
			"client disconnected",
			zap.Error(err),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().URL.Path),
		)

	// 4. Còn lại — bug thật sự, panic đã recover, hoặc lỗi driver/DB rò rỉ lên tới đây.
	//    ĐÂY chính là nơi cần log đầy đủ để tìm nguyên nhân thật.
	default:
		h.logger.Error(
			"unhandled system error",
			zap.Error(err),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().URL.Path),
			zap.String("query", c.Request().URL.RawQuery),
			zap.String("ip", c.RealIP()),
		)
	}

	resp, unwrapErr := echo.UnwrapResponse(c.Response())
	if unwrapErr != nil {
		h.logger.Error("failed to unwrap echo response", zap.Error(unwrapErr))
		return
	}

	if resp.Committed {
		h.logger.Warn(
			"response already committed",
			zap.Error(err),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().URL.Path),
		)
		return
	}

	if writeErr := c.JSON(httpCode, res); writeErr != nil {
		h.logger.Error("failed to write error response", zap.Error(writeErr))
	}
}
