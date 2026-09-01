package middlewares

import (
	"net/http"

	"github.com/himbo22/source-base/pkg/common/apperror"
	"github.com/himbo22/source-base/pkg/common/http/response"
	"github.com/himbo22/source-base/pkg/constraints"
	"github.com/himbo22/source-base/pkg/database/redis"
	"github.com/himbo22/source-base/pkg/utils"
	"go.uber.org/zap"

	"github.com/labstack/echo/v5"
)

// UserRestrictionConfig defines configuration for UserRestrictionMiddleware.
type UserRestrictionConfig struct {
	RedisEngine *redis.Engine
	Logger      *zap.Logger
}

// UserRestrictionMiddleware is the stateful, Redis-backed middleware that supports AuthMiddleware.
// It checks real-time token blacklist revocation and real-time user status restrictions (abandoned/suspended/restricted).
func UserRestrictionMiddleware(cfg UserRestrictionConfig) echo.MiddlewareFunc {
	baseLogger := cfg.Logger
	if baseLogger == nil {
		baseLogger = zap.NewNop()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if cfg.RedisEngine == nil {
				return next(c)
			}

			ctx := c.Request().Context()
			log := baseLogger
			if reqID, ok := ctx.Value(constraints.RequestIDKey).(string); ok && reqID != "" {
				log = baseLogger.With(zap.String("request_id", reqID))
			}

			// 1. Get Claims from Context (Must exist since it runs after AuthMiddleware)
			claims, ok := c.Get(string(constraints.ClaimsKey)).(*utils.Claims)
			if !ok || claims == nil || claims.ID == "" {
				log.Error("security context missing: UserRestrictionMiddleware must be run after AuthMiddleware")
				return apperror.New(response.CodeInternalServer, "internal server error: security context missing", http.StatusInternalServerError)
			}

			// 2. Check Access Token Blacklist in Redis
			blacklisted, err := cfg.RedisEngine.Exists(ctx, constraints.RedisKey.ATBlacklist(claims.ID))
			if err != nil {
				log.Error("redis unavailable, failing blacklist check",
					zap.String("check", "token_blacklist"),
					zap.String("jti", claims.ID),
					zap.Error(err),
				)
				return apperror.New(response.CodeRedisError, "internal server error: security checks unavailable", http.StatusInternalServerError)
			} else if blacklisted > 0 {
				return apperror.New(response.CodeTokenRevoked, "access token has been revoked", http.StatusUnauthorized)
			}

			// 3. Get Public ID from Context
			publicID, ok := c.Get(string(constraints.PublicIDKey)).(string)
			if !ok || publicID == "" {
				log.Error("security context missing: public_id not found")
				return apperror.New(response.CodeInternalServer, "internal server error: security context missing", http.StatusInternalServerError)
			}

			restricted, err := cfg.RedisEngine.Exists(ctx, constraints.RedisKey.UserStatus(publicID))
			if err != nil {
				log.Error("redis unavailable, failing user status check",
					zap.String("check", "user_status"),
					zap.String("public_id", publicID),
					zap.Error(err),
				)
				return apperror.New(response.CodeRedisError, "internal server error: security checks unavailable", http.StatusInternalServerError)
			}

			if restricted > 0 {
				log.Warn("user is restricted", zap.String("public_id", publicID))
				return apperror.New(response.CodeUserSuspended, "your account has been suspended or restricted", http.StatusForbidden)
			}

			return next(c)
		}
	}
}
