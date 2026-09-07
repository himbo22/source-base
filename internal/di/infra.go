package di

import (
	"fmt"

	"github.com/himbo22/source-base/internal/bootstrap"
	"github.com/himbo22/source-base/internal/config"
	"github.com/himbo22/source-base/internal/ent/generate"
	"github.com/himbo22/source-base/internal/ports"
	"github.com/himbo22/source-base/pkg/database/mongodb"
	"github.com/himbo22/source-base/pkg/database/redis"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var InfraSet = wire.NewSet(
	LoggerProvider,
	PostgresProvider,
	RedisProvider,
	MongoDBProvider,
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

func MongoDBProvider(cfg *config.Config) (*mongodb.Client, func(), error) {
	client, err := bootstrap.InitMongoDB(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("init mongodb: %w", err)
	}
	return client, func() {
		_ = bootstrap.CloseMongoDB(client)
	}, nil
}
