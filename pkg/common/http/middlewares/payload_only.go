package middlewares

import (
	"net/http"
	"strings"

	"github.com/himbo22/source-base/pkg/common/apperror"
	"github.com/himbo22/source-base/pkg/common/http/response"
	"github.com/himbo22/source-base/pkg/constraints"
	"github.com/himbo22/source-base/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

// ExtractPayload extracts JWT claims without verifying the signature.
// Use only when the token has already been validated upstream (e.g., at an API gateway)
// or for internal service-to-service calls where signature verification is unnecessary.
func ExtractPayload() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return apperror.New(response.CodeUnauthorized, "authorization header is required", http.StatusUnauthorized)
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return apperror.New(response.CodeUnauthorized, "invalid authorization header format, must be 'Bearer <token>'", http.StatusUnauthorized)
			}

			tokenStr := parts[1]

			token, _, err := jwt.NewParser().ParseUnverified(tokenStr, &utils.Claims{})
			if err != nil {
				return apperror.Wrap(err, response.CodeInvalidToken, "invalid token payload", http.StatusUnauthorized)
			}

			claims, ok := token.Claims.(*utils.Claims)
			if !ok || claims == nil {
				return apperror.New(response.CodeInvalidToken, "invalid token claims", http.StatusUnauthorized)
			}

			if claims.ID == "" || claims.PublicID == "" {
				return apperror.New(response.CodeInvalidToken, "invalid token payload: missing mandatory claims", http.StatusUnauthorized)
			}

			c.Set(string(constraints.ClaimsKey), claims)

			return next(c)
		}
	}
}
