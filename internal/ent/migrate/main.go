//go:build ignore

package main

import (
	"context"
	"github.com/himbo22/source-base/internal/ent/generate/migrate"
	"log"
	"os"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	entschema "entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	dir, err := atlas.NewLocalDir("internal/ent/migrations")
	if err != nil {
		log.Fatalf("failed creating migration directory: %v", err)
	}

	opts := []entschema.MigrateOption{
		entschema.WithDir(dir),
		entschema.WithMigrationMode(entschema.ModeInspect),
		entschema.WithDialect(dialect.Postgres),
		entschema.WithFormatter(atlas.DefaultFormatter),
		entschema.WithForeignKeys(false),
	}

	if len(os.Args) < 2 {
		log.Fatalf("Error: missing migration name. Usage: go run internal/ent/migrate/main.go <migration_name>")
	}
	name := os.Args[1]

	devDBURL := "postgres://postgres:devpassword@localhost:5433/atlas_dev?sslmode=disable"

	if err := migrate.NamedDiff(ctx, devDBURL, name, opts...); err != nil {
		log.Fatalf("failed generating diff migration: %v", err)
	}

	log.Printf("=> Versioned migration generated successfully for: %s\n", name)
}
