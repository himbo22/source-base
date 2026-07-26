package bootstrap

import (
	"github.com/himbo22/source-base/internal/config"
	"github.com/himbo22/source-base/pkg/logger"

	"go.uber.org/zap"
)

func InitLogger(config *config.Config) *zap.Logger {
	cfg := logger.Config{
		Level:      config.Logger.Level,
		Filename:   config.Logger.File,
		MaxSize:    config.Logger.RotateSize,
		MaxAge:     config.Logger.RotateExpire,
		MaxBackups: config.Logger.RotateBackupLimit,
		Compress:   true,
	}

	lg := logger.NewLogger(cfg)
	return lg
}
