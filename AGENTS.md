# AGENTS.md — source-base

## Quick reference

| Task | Command |
|---|---|
| Start infra | `docker compose up -d` |
| Run migrations | `make db-up` |
| Start server | `make run` or `go run cmd/server/main.go` |
| Ent codegen | `make ent-gen` |
| Create schema | `make ent-new name=EntityName` |
| Generate migration | `make db-diff name=migration_name` |
| Apply migrations | `make db-up` |
| Rollback | `make db-down` |
| Check migration status | `make db-status` |
| Rehash migrations | `make db-hash` |

## Architecture

Hexagonal (Clean Architecture) Go microservice. Entry point: `cmd/server/main.go`.

**Layer flow:** Controller → Service (via ports interface) → Repository (implements port)

**Dependency rule:** `pkg/` must never import `internal/`. `internal/domain/entity` and `internal/domain/dto` are pure Go — zero framework imports. Repository uses `internal/ent/generate` (auto-generated) and `internal/const` for context keys.

## Adding a new entity

1. Create schema: `make ent-new name=Entity`
2. Edit `internal/ent/schema/entity.go` — add fields, mixins, edges
3. Run `make ent-gen` to generate Ent code
4. Create domain entity in `internal/domain/entity/`
5. Create DTOs in `internal/domain/dto/`
6. Define port interface in `internal/ports/`
7. Implement repository in `internal/ent/repository/`
8. Implement service in `internal/service/`
9. Wire with Google Wire in `internal/di/`
10. Add HTTP controller in `internal/controller/http/`

## Key conventions

- **Ent schema mixins:** Use `mixin.UUID`, `mixin.PublicID`, `mixin.Time` from `pkg/database/ent/mixins/` for standard fields
- **Public ID format:** Prefix + date + base62 (e.g., `UR20260705xK9`). See `pkg/unique/`
- **Transactions:** Use `txManager.DoInTx(ctx, func(ctx context.Context) error { ... })` — repository methods automatically detect tx in context
- **Context key:** Custom type `contextKey` in `internal/const/key.go` — not a plain string
- **HTTP handlers:** Generic `handler.Wrap[RQ, RS]` in `pkg/common/http/handler/` — controller implements `func(ctx, *RQ) (RS, error)`
- **Validation messages:** English, using json tag names
- **Ent mapping:** `internal/ent/mapper.go` maps Ent models → domain entities. Each field explicitly set, no reflection

## Infrastructure

- **PostgreSQL:** Port 5432 (`postgres-main`), database `app_db`
- **Atlas dev DB:** Port 5433 (`postgres-atlas` in `docker-compose.override.yaml`) — required for `db-diff`
- **Redis:** Port 6379
- **Config:** `configs/config.yaml` loaded by `internal/config`
- **Env vars:** `.env` sourced by Makefile — `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSL_MODE`

## Gotchas

- `db-diff` uses a separate Atlas dev database on port 5433 — `docker-compose.override.yaml` must be active
- Run `make ent-gen` after editing `internal/ent/schema/*.go` — generated code in `internal/ent/generate/` must stay in sync
- No tests yet (`go test ./...` will find nothing). When tests are added, run `go test ./...`
- Go 1.26 required
