package request

import (
	"net/http"
	"source-base/pkg/common/apperror"
	"source-base/pkg/common/http/response"
	"source-base/pkg/common/http/validation"

	"github.com/labstack/echo/v5"
)

func Parse[T any](c *echo.Context) (*T, error) {
	var req T

	// binding
	if err := c.Bind(&req); err != nil {
		return nil, apperror.Wrap(err, response.CodeParamInvalid, err.Error(), http.StatusBadRequest)
	}

	// validate
	if err := validation.Validate(&req); err != nil {
		return nil, apperror.Wrap(err, response.CodeValidationFailed, err.Error(), http.StatusBadRequest)
	}

	return &req, nil
}
