package bootstrap

import (
	"fmt"
	"github.com/himbo22/source-base/internal/config"
	"github.com/himbo22/source-base/pkg/database/redis"
	"github.com/himbo22/source-base/pkg/settings"
)

func InitRedis(config *config.Config) (*redis.Engine, error) {
	cfg := settings.Redis{
		Addrs:           config.Redis.Default.Addrs,
		MasterName:      config.Redis.Default.MasterName,
		Host:            config.Redis.Default.Host,
		Password:        config.Redis.Default.Password,
		Port:            config.Redis.Default.Port,
		Database:        config.Redis.Default.Database,
		PoolSize:        config.Redis.Default.PoolSize,
		MinIdleConns:    config.Redis.Default.MinIdleConns,
		PoolTimeout:     config.Redis.Default.PoolTimeout,
		DialTimeout:     config.Redis.Default.DialTimeout,
		ReadTimeout:     config.Redis.Default.ReadTimeout,
		WriteTimeout:    config.Redis.Default.WriteTimeout,
		MaxRetries:      config.Redis.Default.MaxRetries,
		MaxRetryBackoff: config.Redis.Default.MaxRetryBackoff,
		MinRetryBackoff: config.Redis.Default.MinRetryBackoff,
	}

	engine := &redis.Engine{
		Config: &cfg,
	}

	if err := engine.Connect(); err != nil {
		return nil, fmt.Errorf("connect redis failed: %v", err)
	}

	return engine, nil
}
