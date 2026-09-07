package bootstrap

import (
	"context"
	"fmt"

	"github.com/himbo22/source-base/internal/config"
	"github.com/himbo22/source-base/pkg/database/mongodb"
	"github.com/himbo22/source-base/pkg/settings"
)

func InitMongoDB(config *config.Config) (*mongodb.Client, error) {
	cfg := settings.MongoDB{
		URL:             config.MongoDB.URL,
		Host:            config.MongoDB.Host,
		Port:            config.MongoDB.Port,
		Username:        config.MongoDB.Username,
		Password:        config.MongoDB.Password,
		Database:        config.MongoDB.Database,
		AuthSource:      config.MongoDB.AuthSource,
		MaxPoolSize:     config.MongoDB.MaxPoolSize,
		MinPoolSize:     config.MongoDB.MinPoolSize,
		MaxConnIdleTime: config.MongoDB.MaxConnIdleTime,
		Timeout:         config.MongoDB.Timeout,
		TLSEnabled:      config.MongoDB.TLSEnabled,
		TLSInsecure:     config.MongoDB.TLSInsecure,
	}

	client := mongodb.NewClient(&cfg)

	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("connect mongodb failed: %v", err)
	}

	return client, nil
}

func CloseMongoDB(client *mongodb.Client) error {
	if client == nil {
		return nil
	}
	if err := client.Disconnect(context.Background()); err != nil {
		return fmt.Errorf("failed to close mongodb connection: %w", err)
	}
	return nil
}