package bootstrap

import (
	"fmt"
	"source-base/internal/config"
	"source-base/internal/ent/generate"
	"source-base/pkg/database/ent"
	"source-base/pkg/settings"
)

func InitPostgreSQL(config *config.Config) (*generate.Client, error) {
	cfg := settings.PostgreSQL{
		Host:            config.Postgres.Host,
		Port:            config.Postgres.Port,
		Username:        config.Postgres.User,
		Password:        config.Postgres.Password,
		Database:        config.Postgres.DBName,
		MaxOpenConns:    config.Postgres.MaxOpenConns,
		MaxIdleConns:    config.Postgres.MaxIdleConns,
		ConnMaxLifetime: config.Postgres.ConnMaxLifetime,
		ConnMaxIdleTime: config.Postgres.ConnMaxIdleTime,
	}

	driver, err := ent.NewPostgreSQLDriver(cfg)
	if err != nil {
		return nil, err
	}

	client := generate.NewClient(generate.Driver(driver))
	//if config.App.Env == "development" {
	//	client = client.Debug()
	//}
	return client, nil
}

func ClosePostgreSQL(client *generate.Client) error {
	if client == nil {
		return nil
	}
	if err := client.Close(); err != nil {
		return fmt.Errorf("failed to close postgresql connection: %w", err)
	}
	return nil
}
