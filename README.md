# source-base

## 1. Introduction

**source-base** is a production-ready microservice template built on **Hexagonal Architecture (Clean Architecture)**.

Goals:
- Fully decouple business logic from frameworks and infrastructure
- Easy to extend, test, and maintain
- Use modern, high-performance Go libraries

## 2. Tech Stack

| Dependency | Version | Purpose |
|---|---|---|
| `entgo.io/ent` | `v0.14.6` | ORM code generation |
| `github.com/labstack/echo/v5` | `v5.2.1` | HTTP framework |
| `github.com/google/wire` | `v0.7.0` | Dependency injection |
| `github.com/redis/go-redis/v9` | `v9.21.0` | Redis client |
| `go.uber.org/zap` | `v1.28.0` | Structured logging |
| `github.com/go-playground/validator/v10` | `v10.30.3` | Request validation |
| `github.com/google/uuid` | `v1.6.0` | UUID generation (V7) |
| `github.com/jackc/pgx/v5` | `v5.10.0` | PostgreSQL driver |
| `github.com/natefinch/lumberjack` | `v2.0.0` | Log rotation |
| `ariga.io/atlas` | `v0.36.2-...` | Schema migration |
| `gopkg.in/yaml.v3` | `v3.0.1` | YAML config parsing (indirect) |
| `github.com/golang-jwt/jwt/v5` | `v5.3.1` | JWT (indirect) |

**Language:** Go 1.26

## 3. Directory Structure

```text
source-base/
├── cmd/server/              # Entry point (main.go) — load config, DI init, start server
├── configs/                 # Sample YAML config files
├── internal/
│   ├── bootstrap/           # Initialize components (Postgres, Redis, Logger, Router, TxManager, OTel)
│   ├── config/              # Config struct definition, read from YAML
│   ├── const/               # Context keys using custom type (TxKey)
│   ├── controller/http/     # HTTP controllers — receive request, call service via interface
│   ├── di/                  # Dependency Injection with Google Wire
│   ├── domain/
│   │   ├── dto/             # Data Transfer Objects (request/response) — pure Go
│   │   └── entity/          # Domain entities — pure Go, zero dependency
│   ├── ent/
│   │   ├── schema/          # Ent schema definitions (User, Pet) + mixins
│   │   ├── generate/        # Auto-generated code by Ent ORM
│   │   ├── migrate/         # Migration CLI using Atlas
│   │   ├── migrations/      # Generated SQL migration files
│   │   ├── repository/      | Implement ports repository (using Ent client)
│   │   ├── generate.go      # go:generate directive
│   │   └── mapper.go        # Map Ent model → Domain entity
│   ├── ports/               # Interface definitions (hexagonal ports)
│   └── service/             # Business logic layer
├── logs/                    # Log files output
├── pkg/
│   ├── common/
│   │   ├── apperror/        # AppError type (Code, Message, RootCause, HTTPStatus)
│   │   ├── cache/           # Cache abstraction interfaces (Engine, SortedSetEngine)
│   │   ├── http/            # HTTP utilities: handler Wrap, ErrorHandler, middlewares,
│   │   │                    #   request Parse, response format, validation
│   │   ├── tx/              # Transaction Manager interface
│   │   └── workerpool/      # Worker pool (not yet used — TODO)
│   ├── constraints/         # Go generics type constraints + context keys
│   ├── database/
│   │   ├── ent/             # PostgreSQL driver, generic TxManager, mixins (UUID, PublicID, Time)
│   │   └── redis/           # Redis client implementing cache.Engine + SortedSetEngine
│   ├── dto/                 # Shared DTOs (ZMember for Redis sorted sets)
│   ├── logger/              # Logger initialization (zap + lumberjack)
│   ├── settings/            # Standard config structs
│   ├── unique/              # Unique ID generation (public_id prefix+date+base62, RandBase62)
│   └── utils/               # Utilities (time, token, converter)
└── storages/logs/           # Log storage
```

## 4. Architecture & Design Principles

### Flow

```
Controller (HTTP)  ──calls──→  Service (business logic)  ──calls──→  Repository (data)
```

- **Controller** receives HTTP request → parse & validate → calls `ports.UserService` → returns standardized response
- **Service** receives DTO request → maps to domain entity → calls `ports.UserRepository` → maps result to DTO response
- **Repository** uses Ent ORM to interact with PostgreSQL, maps Ent model → domain entity

### Dependency direction

```
Domain Entity / DTO   ←── no dependencies (pure Go)
Ports                 ──uses──→  DTO + Entity (data types only)
Service               ──calls──→  Ports interface + DTO + Entity
Controller            ──calls──→  Ports interface + DTO
Repository            ──implements──→ Ports interface + Entity + Ent generate
pkg/                  ── does not import internal/* ──
```

### Why domain/entity and domain/dto are separate from internal/ent?

- **domain/entity** and **domain/dto** are pure Go structs with **zero imports from any framework** (not Ent, Echo, or database drivers)
- **internal/ent/mapper.go** handles mapping from Ent model (`generate.User`) to domain entity (`entity.User`), keeping the rest of the application independent of Ent ORM
- This allows swapping ORM or database without affecting business logic

### Transaction mechanism

Transactions are managed via a context-based pattern:

1. **`pkg/database/ent/tx.go`**: generic `WithTx[T]()` performs begin → fn → commit/rollback, with panic recovery
2. **`pkg/database/ent/tx.go`**: `txManager[T]` implements `tx.Manager` interface, uses generics to inject transaction into context
3. **`internal/bootstrap/tx.go`**: wires `client.Tx` with context injection function (`context.WithValue(ctx, TxKey, tx)`)
4. **`internal/const/key.go`**: `TxKey` uses custom type `contextKey` (not a plain string) to avoid collision
5. **`internal/ent/repository/tx.go`**: `GetClient(ctx, client)` checks context — if `*generate.Tx` exists in context, returns `tx.Client()`, otherwise returns the original client

When a service needs to run inside a transaction:
```go
txManager.DoInTx(ctx, func(txCtx context.Context) error {
    // All repository calls here automatically use the tx client
    return nil
})
```

Every repository method calls `GetClient(ctx, u.client)` at the start, so they automatically follow the transaction if one is present.

## 5. Environment Requirements

| Tool | Version | Notes |
|---|---|---|
| Go | 1.26+ | Per `go.mod` |
| Docker | latest | Run Postgres, Redis locally |
| Atlas CLI | >= 0.14 | Install: `go install ariga.io/atlas/cmd/atlas@latest` |

## 6. Getting Started

### Step 1: Clone & configure

```bash
git clone <repo-url> source-base
cd source-base
# .env file is already included with dev defaults
```

### Step 2: Start infrastructure

```bash
docker compose up -d
```

This starts:
- `postgres-main` (PostgreSQL 15, port 5432, database: app_db)
- `redis` (Redis 7, port 6379)
- `adminer` (DB admin panel, port 3001)
- `postgres-atlas` (additional PostgreSQL for Atlas dev, port 5433 — from `docker-compose.override.yaml`)

### Step 3: Run migrations

```bash
make db-up
```

### Step 4: Start server

```bash
make run
# or
go run cmd/server/main.go
```

Server runs at `http://localhost:8080`.

## 7. Database Migration Workflow

### Process

```
Edit schema (internal/ent/schema/*.go)
    → go generate (make ent-gen)
    → Generate migration file (make db-diff name=<name>)
    → Review SQL in internal/ent/migrations/
    → Apply (make db-up)
```

### Makefile commands

| Command | Description |
|---|---|
| `make ent-new name=<Entity>` | Create new schema |
| `make ent-desc` | Display current schema structure |
| `make ent-gen` | Generate Ent ORM code (client, CRUD builders, predicates) |
| `make db-diff name=<name>` | Diff schema vs DB → generate migration file |
| `make db-up` | Apply migrations to database |
| `make db-down` | Rollback migrations |
| `make db-hash` | Rehash migration files |
| `make db-status` | Check migration status |

**Note:** `db-diff` uses `docker-compose.override.yaml` — the `postgres-atlas` container (port 5433) must be running to serve as the dev database for Atlas.

## 8. API Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/ping` | Health check |
| `POST` | `/api/v1/public/users` | Create a new user |
| `GET` | `/api/v1/public/users/:public_id` | Get user by public_id (with pets list) |

**Sample Request/Response:**

`POST /api/v1/public/users`
```json
{
    "name": "Nguyen Van A",
    "email": "a@example.com"
}
```

Standard response format:
```json
{
    "code": 20000,
    "message": "Success",
    "data": {
        "public_id": "UR20260705xK9",
        "name": "Nguyen Van A",
        "email": "a@example.com",
        "created_at": "2026-07-05T12:00:00Z",
        "updated_at": "2026-07-05T12:00:00Z"
    }
}
```

### Error codes

| Code | Meaning |
|---|---|
| 20000 | Success |
| 20001 | Created |
| 40000 | Invalid parameters |
| 40001 | Validation failed |
| 44000 | Resource not found |
| 50000 | Internal server error |
| 50001 | Database error |

Full details at `pkg/common/http/response/codes.go`.

## 9. Testing

**No tests yet.** The project is currently under development, with no `*_test.go` files present.

When tests are added, run them with:

```bash
go test ./...
```

## 10. Coding Conventions

Key conventions extracted from the actual codebase:

- **Domain entity does not import Ent**: `internal/domain/entity/*.go` are pure Go structs with no framework imports (no Ent, Echo). Only uses `time.Time` and `uuid.UUID`.
- **Repository returns `(*entity.T, error)`**: Does not mutate input parameters. Exception: `Create(ctx, *entity.User)` reassigns `*user = *created` to return the ID to caller. Other methods (GetByPublicID) return a new pointer.
- **Explicit Set().Set() when creating Ent records**: No reflection or generic mapping. Each field is explicitly set: `client.User.Create().SetName(user.Name).SetEmail(user.Email).Save(ctx)`.
- **Context key uses custom type**: `internal/const/key.go` uses `type contextKey string` instead of a plain string to prevent collisions between packages.
- **Transaction uses generics**: `pkg/database/ent/tx.go` uses generic `WithTx[T Tx]()` — type-safe, no interface casting needed.
- **Generic HTTP handler**: `pkg/common/http/handler/wrapper.go` uses `Wrap[RQ, RS any](Func[RQ, RS])` — parses request, calls handler, returns standardized response — controller only needs to implement `func(context.Context, *RQ) (RS, error)`.
- **Validation messages in English**: `pkg/common/http/validation/validator.go` uses json tag names and English messages (`"name is required"`, `"email must be a valid email"`).
- **`pkg/` does not import `internal/*`**: Packages in `pkg/` are designed to be reusable and independent of `internal/` (except `pkg/common/tx` used by internal).

## 11. License / Contributing

**TODO:** No license information yet. Currently under internal development.
