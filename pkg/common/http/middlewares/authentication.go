package middlewares

import (
	"context"
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/himbo22/source-base/pkg/common/apperror"
	"github.com/himbo22/source-base/pkg/common/http/response"
	"github.com/himbo22/source-base/pkg/constraints"
	"github.com/himbo22/source-base/pkg/utils"

	"github.com/labstack/echo/v5"
)

// AuthConfig defines configuration parameters for AuthMiddleware.
type AuthConfig struct {
	PublicKey *rsa.PublicKey
}

// AuthMiddleware is a pure, stateless middleware that extracts and verifies RS256 JWT access tokens.
func AuthMiddleware(cfg AuthConfig) echo.MiddlewareFunc {
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

			// Verify RS256 JWT signature and claims
			claims, err := utils.VerifyRS256(tokenStr, cfg.PublicKey)
			if err != nil {
				return apperror.Wrap(err, response.CodeInvalidToken, "invalid or expired token", http.StatusUnauthorized)
			}

			// Enforce mandatory claims
			if claims.ID == "" || claims.PublicID == "" {
				return apperror.New(response.CodeInvalidToken, "invalid token payload: missing mandatory claims", http.StatusUnauthorized)
			}

			// Inject claims and public ID into context
			c.Set(string(constraints.ClaimsKey), claims)
			c.Set(string(constraints.PublicIDKey), claims.PublicID)

			reqCtx := context.WithValue(c.Request().Context(), constraints.ClaimsKey, claims)
			reqCtx = context.WithValue(reqCtx, constraints.PublicIDKey, claims.PublicID)
			c.SetRequest(c.Request().WithContext(reqCtx))

			return next(c)
		}
	}
}
