![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/tushkanchek/084c08f8beefdcc821718cb28a87702a/raw/coverage.json)

# GoSpaceShipService

Microservices backend for a spaceship parts ordering system, built in Go.

## Services

| Service | Transport | Description |
|---|---|---|
| `order` | REST (chi) | Order management, orchestrates inventory + payment |
| `inventory` | gRPC | Parts catalog and stock |
| `payment` | gRPC | Payment processing |
| `assembly` | Kafka consumer | Assembles orders after successful payment |
| `platform` | — | Shared infrastructure: logger, graceful closer, gRPC/Kafka helpers |
| `shared` | — | Generated proto/OpenAPI code, proto definitions |

## Architecture

```
Client
  └─► Order (REST)
        ├─► Inventory (gRPC) — check parts availability
        ├─► Payment (gRPC)   — process payment
        └─► Kafka            — publish OrderPaidEvent
                                    └─► Assembly (consumer)
```

## Getting Started

### Prerequisites

- Go 1.22+
- [go-task](https://taskfile.dev/) — task runner
- Docker + Docker Compose

```bash
brew install go-task
```

### Run locally

```bash
# Start core infrastructure (Postgres, Kafka, etc.)
task up-core

# Start a specific service
task up-order
task up-inventory

# Or start everything
task up-all

# Stop everything
task down-all
```

### Development commands

```bash
task lint              # Format (gofumpt + gci) then lint all modules
task format            # Format only
task test              # Run unit tests for all modules
task test -- MODULES=order  # Test a single module
task test-coverage     # Tests with coverage for service/repository packages
task gen               # Generate all: proto (buf) + OpenAPI (ogen)
task proto:gen         # Generate Go code from .proto files
task ogen:gen          # Generate Go code from OpenAPI specs
task mockery:gen       # Generate mocks
task deps:update       # Run go mod tidy across all modules
```

## API

### Order Service (REST)

| Method | Path | Description |
|---|---|---|
| `POST` | `/orders` | Create a new order |
| `GET` | `/orders/{uuid}` | Get order by UUID |
| `POST` | `/orders/{uuid}/pay` | Pay for an order |
| `POST` | `/orders/{uuid}/cancel` | Cancel an order |

### Inventory Service (gRPC)

Defined in `shared/proto/inventory/v1/inventory.proto` — `ListParts`, `GetPart`.

### Payment Service (gRPC)

Defined in `shared/proto/payment/v1/payment.proto` — `Pay`.

## Project Structure

```
GoSpaceShipService/
├── order/            # REST service
├── inventory/        # gRPC service
├── payment/          # gRPC service
├── assembly/         # Kafka consumer
├── platform/         # Shared infrastructure
├── shared/           # Proto definitions + generated code
│   ├── proto/        # .proto source files
│   └── api/          # OpenAPI specs
├── deploy/
│   ├── compose/      # Docker Compose files per service
│   └── env/          # Env templates (*.env.template)
├── go.work           # Go workspace
└── Taskfile.yml      # Task runner config
```

Each service follows the same internal layout:

```
<service>/
  cmd/main.go
  internal/
    api/<domain>/v1/   # Transport handlers
    app/               # DI wiring / bootstrap
    config/            # Config structs + env parsing
    converter/         # Model ↔ proto/API converters
    model/             # Domain models
    service/           # Business logic
    repository/        # Data access (pgx)
    client/grpc/       # gRPC client wrappers
```

## CI/CD

GitHub Actions (`.github/workflows/ci.yml`) runs on every push and pull request:

1. Extract tool versions from `Taskfile.yml`
2. Lint
3. Test + coverage