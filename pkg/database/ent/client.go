package ent

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/himbo22/source-base/pkg/settings"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgreSQLDriver(cfg settings.PostgreSQL) (*entsql.Driver, error) {
	var dsn string
	if cfg.URL != "" {
		dsn = cfg.URL
	} else {
		sslmode := cfg.SSLMode
		if sslmode == "" {
			sslmode = "disable"
		}
		timezone := cfg.Timezone
		if timezone == "" {
			timezone = "UTC"
		}

		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s",
			cfg.Host,
			cfg.Port,
			cfg.Username,
			cfg.Password,
			cfg.Database,
			sslmode,
			timezone,
		)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	driver := entsql.OpenDB(dialect.Postgres, db)
	return driver, nil
}
