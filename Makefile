include .env
export

# ==============================================================================
# ENVIRONMENT VARIABLES & CONFIGURATION
# ==============================================================================
ENT_DIR       := ./internal/ent
SCHEMA_DIR    := $(ENT_DIR)/schema
MIGRATE_DIR   := ./internal/ent/migrations

# Change this DSN to match your local database or pass it from CLI
# Example: make db-up DB_URL="postgres://admin:123456@localhost:5432/my_db?sslmode=disable"
DB_URL ?= "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

DEV_DB_URL = "postgres://postgres:devpassword@localhost:5433/atlas_dev?sslmode=disable"

# ==============================================================================
# HELP
# ==============================================================================
.PHONY: help
help: ## Display available commands
	@echo "Makefile commands for Ent ORM:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ==============================================================================
# 1. SCHEMA MANAGEMENT (DATA DESIGN)
# ==============================================================================
.PHONY: ent-new
ent-new: ## Create a new Schema. Usage: make ent-new name=User
	@if [ -z "$(name)" ]; then \
		echo "Error: You must provide a schema name. Example: make ent-new name=User"; \
		exit 1; \
	fi
	@echo "=> Creating new schema: $(name)..."
	go run -mod=mod entgo.io/ent/cmd/ent new --target $(SCHEMA_DIR) $(name)

.PHONY: ent-desc
ent-desc: ## Display the structure of all current Schemas
	@echo "=> Describing schema structure..."
	go run -mod=mod entgo.io/ent/cmd/ent describe $(SCHEMA_DIR)

# ==============================================================================
# 2. CODE GENERATION
# ==============================================================================
.PHONY: ent-gen
ent-gen: ## Run the Ent code generator (generate Graph, Client, Hooks...)
	@echo "=> Generating Ent ORM code..."
	go generate $(SCHEMA_DIR)
	@echo "=> Done!"

# ==============================================================================
# 3. MIGRATION MANAGEMENT (PRODUCTION-GRADE WITH ATLAS)
# ==============================================================================
.PHONY: db-diff
db-diff: ## Diff schema vs DB to generate migration file. Usage: make db-diff name=add_users
	@if [ -z "$(name)" ]; then \
		echo "Error: You must provide a migration name. Example: make db-diff name=add_user_table"; \
		exit 1; \
	fi
	@echo "=> Generating migration file: $(name)..."
	go run -mod=mod $(ENT_DIR)/migrate/main.go $(name)

.PHONY: db-up
db-up: ## Apply migration files to the database
	@echo "=> Running migrations on database..."
	atlas migrate apply \
		--dir "file://$(MIGRATE_DIR)" \
		--url $(DB_URL)

.PHONY: db-down
db-down: ## Apply migration files to the database
	@echo "=> Rolling back migrations on database..."
	atlas migrate down \
		--dir "file://$(MIGRATE_DIR)" \
		--url $(DB_URL) \
		--dev-url $(DEV_DB_URL)

.PHONY: db-hash
db-hash: ## Apply migration files to the database
	@echo "=> Rehashing..."
	atlas migrate hash \
		--dir "file://$(MIGRATE_DIR)" \

.PHONY: db-status
db-status: ## Check current database status against migration files
	@echo "=> Migration status:"
	atlas migrate status \
		--dir "file://$(MIGRATE_DIR)" \
		--url $(DB_URL)

.PHONY: run
run:
	go run ./cmd/server/main.go