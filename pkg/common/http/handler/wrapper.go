package handler

import (
	"context"
	"source-base/pkg/common/http/request"
	"source-base/pkg/common/http/response"

	"github.com/labstack/echo/v5"
)

type Func[RQ any, RS any] func(context.Context, *RQ) (RS, error)

func Wrap[RQ any, RS any](h Func[RQ, RS]) echo.HandlerFunc {
	return func(c *echo.Context) error {
		req, err := request.Parse[RQ](c)
		if err != nil {
			return err
		}

		res, err := h(c.Request().Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, response.CodeSuccess, res)
	}
}
