package di

import (
	"source-base/internal/config"
	"source-base/internal/ent/generate"
	"source-base/pkg/database/redis"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type App struct {
	Config  *config.Config
	Logger  *zap.Logger
	DB      *generate.Client
	Redis   *redis.Engine
	EchoApp *echo.Echo
}
