package handler

import (
	"context"
	"net/http"

	"github.com/himbo22/source-base/pkg/common/apperror"
	"github.com/himbo22/source-base/pkg/common/http/request"
	"github.com/himbo22/source-base/pkg/common/http/response"
	"github.com/himbo22/source-base/pkg/constraints"
	"github.com/himbo22/source-base/pkg/utils"

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

type AuthenticatedFunc[RQ any, RS any] func(ctx context.Context, claims *utils.Claims, req *RQ) (RS, error)

func WrapAuthenticated[RQ any, RS any](h AuthenticatedFunc[RQ, RS]) echo.HandlerFunc {
	return func(c *echo.Context) error {
		req, err := request.Parse[RQ](c)
		if err != nil {
			return err
		}

		claims, ok := c.Get(string(constraints.ClaimsKey)).(*utils.Claims)
		if !ok || claims == nil {
			return apperror.New(response.CodeUnauthorized, "missing claims in context", http.StatusUnauthorized)
		}

		res, err := h(c.Request().Context(), claims, req)
		if err != nil {
			return err
		}

		return response.Success(c, response.CodeSuccess, res)
	}
}
