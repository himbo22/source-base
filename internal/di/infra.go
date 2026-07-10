package di

import (
	"fmt"

	"source-base/internal/bootstrap"
	"source-base/internal/config"
	"source-base/internal/ent/generate"
	"source-base/internal/ports"
	"source-base/pkg/database/redis"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var InfraSet = wire.NewSet(
	LoggerProvider,
	PostgresProvider,
	RedisProvider,
	TxManagerProvider,
)

func LoggerProvider(cfg *config.Config) *zap.Logger {
	return bootstrap.InitLogger(cfg)
}

func PostgresProvider(cfg *config.Config) (*generate.Client, func(), error) {
	client, err := bootstrap.InitPostgreSQL(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("init postgres: %w", err)
	}
	return client, func() {
		_ = bootstrap.ClosePostgreSQL(client)
	}, nil
}

func RedisProvider(cfg *config.Config) (*redis.Engine, func(), error) {
	engine, err := bootstrap.InitRedis(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("init redis: %w", err)
	}
	return engine, func() {
		_ = engine.Close()
	}, nil
}

func TxManagerProvider(client *generate.Client) ports.TxManager {
	return bootstrap.NewEntTxManager(client)
}
