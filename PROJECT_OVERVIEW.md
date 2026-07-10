# Project Overview — source-base

## 1. Overall directory structure

```
source-base/
├── cmd/server/              # Entry point (main.go)
├── configs/                 # YAML config files
├── internal/
│   ├── bootstrap/           # Initialize components (DB, Redis, Logger, Router, Tx)
│   ├── config/              # Config struct definition
│   ├── const/               # Context keys
│   ├── controller/http/     # HTTP controllers
│   ├── di/                  # Dependency Injection (Google Wire)
│   ├── domain/
│   │   ├── dto/             # Data Transfer Objects (request/response)
│   │   └── entity/          # Pure domain entities
│   ├── ent/
│   │   ├── generate/        # Auto-generated code by Ent
│   │   ├── migrate/         # Migration CLI (Atlas)
│   │   ├── migrations/      # Migration files
│   │   ├── repository/      # Implement repository
│   │   └── schema/          # Ent schema definition
│   ├── ports/               # Interface definitions (hexagonal ports)
│   └── service/             # Business logic layer
├── logs/                    # Log files output
├── pkg/
│   ├── common/
│   │   ├── apperror/        # Application error type
│   │   ├── cache/           # Cache abstraction interfaces
│   │   ├── http/            # HTTP utilities (handler, middleware, request, response, validation)
│   │   ├── tx/              # Transaction manager interface
│   │   └── workerpool/      # Worker pool (not yet used)
│   ├── constraints/         # Go generics type constraints
│   ├── database/
│   │   ├── ent/             # Ent ORM utilities (driver, tx, mixins)
│   │   └── redis/           # Redis client implementation
│   ├── dto/                 # Shared DTOs (e.g., ZMember)
│   ├── logger/              # Logger initialization (zap + lumberjack)
│   ├── settings/            # Standard config structs
│   ├── unique/              # Unique ID generation (public_id, base62)
│   └── utils/               # Utilities (time, token, converter)
```

## 2. Main packages and roles

| Package | Role |
|---|---|
| `cmd/server` | Entry point: load config, init app via DI, start HTTP server |
| `internal/config` | Defines `Config` struct read from YAML file |
| `internal/bootstrap` | Initializes components: Postgres (Ent client), Redis, Logger, Router (Echo), TxManager, OTel (placeholder) |
| `internal/const` | Context key constants (`TxKey`) |
| `internal/controller/http` | HTTP controllers: receive request, call service via interface, return response |
| `internal/di` | Dependency injection with Google Wire: `InitApp()` wires everything |
| `internal/domain/dto` | Request/Response DTOs (no internal dependencies) |
| `internal/domain/entity` | Pure domain entities (zero framework dependency) |
| `internal/ent/schema` | Defines Ent schema (User, Pet) using custom mixins |
| `internal/ent/generate` | Auto-generated code by Ent ORM (Client, CRUD builders, predicates) |
| `internal/ent/repository` | Implements `ports.UserRepository`, uses Ent client, maps entity |
| `internal/ent/migrate` | CLI migration using Atlas |
| `internal/ports` | Interfaces defining boundaries between layers (hexagonal ports) |
| `internal/service` | Business logic: receive DTO → map entity → call repository → map response |
| `pkg/common/apperror` | `AppError` type with Code, Message, RootCause, HTTPStatus |
| `pkg/common/cache` | `Engine` and `SortedSetEngine` interfaces for caching |
| `pkg/common/http/handler` | Generic `Wrap()` handler + `ErrorHandler` |
| `pkg/common/http/middlewares` | Echo middleware: RequestLogger, CustomRecover, Auth (placeholder) |
| `pkg/common/http/request` | Generic `Parse[T]()` bind + validate request body |
| `pkg/common/http/response` | Standardized response format (Code, Message, Data) + error codes |
| `pkg/common/http/validation` | Wrapper `go-playground/validator/v10` with English messages |
| `pkg/common/tx` | `Manager` interface for transactions |
| `pkg/database/ent` | PostgreSQL driver (pgx), generic TxManager, mixins (UUID, PublicID, Time) |
| `pkg/database/redis` | Implements `cache.Engine` + `cache.SortedSetEngine` via go-redis |
| `pkg/dto` | Common DTO: `ZMember` for Redis sorted set |
| `pkg/logger` | Initializes `*zap.Logger` with lumberjack rotation |
| `pkg/settings` | Standard config structs (Server, MongoDB, Logger, Redis, Kafka, PostgreSQL) |
| `pkg/unique` | Generates `public_id` (prefix + date + base62) and `RandBase62` |
| `pkg/utils` | Utility functions: `ToDuration`, JWT `Claims`, `StringToUUID` |

## 3. Main dependencies (from go.mod)

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
| `gopkg.in/yaml.v3` | `v3.0.4` | YAML config parsing |
| `github.com/golang-jwt/jwt/v5` | `v5.3.1` | JWT (indirect) |

## 4. Important interfaces

### `internal/ports/user.go`
```go
type UserService interface {
    Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
}
type UserController interface {
    Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
}
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
}
```

### `internal/ports/tx.go`
```go
type TxManager interface {
    DoInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

### `internal/controller/ports/pet.go`
```go
type PetController interface {
    Create(ctx context.Context, req *dto.CreatePetRequest) (dto.PetResponse, error)
}
```

### `pkg/common/cache/cache.go`
```go
type Engine interface {
    Get(ctx context.Context, key string) ([]byte, bool, error)
    Set(ctx context.Context, key string, value any, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    BatchSet(ctx context.Context, values map[string]any, ttl time.Duration) error
    BatchDelete(ctx context.Context, keys []string) error
    Close() error
}
type SortedSetEngine interface {
    ZAdd(ctx context.Context, key string, members ...*dto.ZMember) error
    ZRemRangeByScore(ctx context.Context, key string, min, max string) error
    ZCount(ctx context.Context, key string, min, max string) (int64, error)
    ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
}
```

### `pkg/common/tx/Manager.go`
```go
type Manager interface {
    DoInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

### `pkg/database/ent/tx.go`
```go
type Tx interface {
    Commit() error
    Rollback() error
}
```

## 5. Main Entity / Domain Models

### `internal/domain/entity/user.go`
```go
type User struct {
    ID        uuid.UUID
    PublicID  string
    Name      string
    Email     string
    CreatedAt time.Time
    UpdatedAt time.Time
    Pets      []Pet
}
```

### `internal/domain/entity/pet.go`
```go
type Pet struct {
    ID        uuid.UUID
    PublicID  string
    Name      string
    UserID    uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
    Owner     *User
}
```

### Ent Schema (`internal/ent/schema/user.go` + `pet.go`) — effective fields after generation:
**User:** id (uuid.UUID), public_id (string, unique), name (string), email (string), created_at (time.Time), updated_at (time.Time), pets (edge []*Pet)

**Pet:** id (uuid.UUID), public_id (string, unique), name (string), user_id (uuid.UUID), created_at (time.Time), updated_at (time.Time), owner (edge *User)

### DTOs (`internal/domain/dto/`)
```go
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}
type UserResponse struct {
    PublicID  string        `json:"public_id"`
    Name      string        `json:"name"`
    Email     string        `json:"email"`
    CreatedAt time.Time     `json:"created_at,omitempty"`
    UpdatedAt time.Time     `json:"updated_at,omitempty"`
    Pets      []PetResponse `json:"pets,omitempty"`
}
type CreatePetRequest struct {
    Name string `json:"name" validate:"required,min=2,max=100"`
}
type PetResponse struct {
    PublicID  string      `json:"public_id"`
    Name      string      `json:"name"`
    UserID    uuid.UUID   `json:"user_id"`
    CreatedAt time.Time   `json:"created_at"`
    UpdatedAt time.Time   `json:"updated_at"`
    Owner     *UserSimple `json:"owner,omitempty"`
}
```

## 6. Dependency flow diagram between layers

```
cmd/server/main.go
    └── internal/config
    └── internal/di (Google Wire)
          ├── internal/config
          ├── internal/bootstrap
          │     ├── internal/config
          │     ├── internal/ent/generate
          │     ├── internal/ports
          │     ├── pkg/database/ent
          │     ├── pkg/database/redis
          │     ├── pkg/settings
          │     └── pkg/logger
          ├── internal/controller/http  ──→  internal/ports, internal/domain/dto
          ├── internal/service          ──→  internal/ports, internal/domain/dto, internal/domain/entity
          └── internal/ent/repository   ──→  internal/ports, internal/domain/entity, internal/ent/generate, internal/const

internal/ent/schema  ──→  pkg/database/ent (mixins)
internal/ent/migrate ──→  ariga.io/atlas, internal/ent/generate

pkg/database/ent     ──→  pkg/common/tx, pkg/settings, pkg/unique
pkg/database/redis   ──→  pkg/common/cache, pkg/dto, pkg/settings, pkg/utils
pkg/common/http/*    ──→  pkg/common/apperror, pkg/common/http/response/validation, pkg/constraints
pkg/logger           ──→  go.uber.org/zap, github.com/natefinch/lumberjack
```

### Layer relationship (hexagonal architecture):

```
Controller (HTTP)          ──calls──→  Port (UserService interface)
Service (business logic)   ──calls──→  Port (UserRepository interface)
Repository (infra)         ──implements──→ Port (UserRepository interface)
                                    ↕
Domain Entity              ←── no dependencies ──→  (pure Go)
DTO                        ←── no dependencies ──→  (pure Go)
Port                       ──uses──→  DTO + Entity (data types only)
```

### Clear dependency direction:
- **Domain entities** → does not import any internal package
- **Domain DTOs** → does not import any internal package
- **Ports** → only imports domain/dto + domain/entity
- **Service** → ports + domain/dto + domain/entity + zap
- **Controller** → ports + domain/dto
- **Repository** → ports + domain/entity + ent/generate + const
- **Bootstrap** → config + ent/generate + pkg/* + ports
- **pkg/** → does not import internal/* (except pkg/common/tx used by internal)
