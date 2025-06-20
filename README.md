# MCP Client - Hexagonal Architecture

This project follows the hexagonal architecture (ports and adapters) pattern to maintain clean, testable, and maintainable code.

## Project Structure

```
/
├── main.go                           # Application entry point
├── adapters/                         # External adapters
│   ├── api/                         # HTTP API adapters
│   │   └── sample_handler.go        # HTTP handlers
│   └── db/                          # Database adapters
│       ├── connection.go            # Database connection
│       └── sample_model/            # Model-specific repositories
│           └── sample_model_repo.go # Sample model repository
├── core/                            # Core business logic
│   ├── domain/                      # Domain models
│   │   └── sample_model.go          # Sample domain model
│   ├── ports/                       # Interface definitions
│   │   ├── db_port.go              # Database port interface
│   │   └── sample_model_port.go     # Sample model port interface
│   └── usecases/                    # Business use cases
│       └── sample_business_flow/    # Sample business flow
│           └── sample_business_flow_usecase.go
└── go.mod                           # Go module definition
```

## Architecture Principles

### Hexagonal Architecture (Ports and Adapters)

- **Core**: Contains business logic, domain models, and port interfaces
- **Adapters**: External implementations that connect to databases, APIs, etc.
- **Ports**: Interfaces that define contracts between core and adapters

### Key Components

1. **Domain Models** (`core/domain/`): Pure business entities with no external dependencies
2. **Use Cases** (`core/usecases/`): Business scenarios and workflows
3. **Ports** (`core/ports/`): Interface contracts for external dependencies
4. **Adapters** (`adapters/`): Concrete implementations of ports

### Data Flow

```
HTTP Request → API Handler → Use Case → Domain Logic → Port Interface → Repository Adapter → Database
```

## Getting Started

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Build the project:
   ```bash
   go build .
   ```

3. Run the server:
   ```bash
   ./mcp_client
   ```

The server will start on port 8080 with a basic health check endpoint at `/health`.

## Development Guidelines

- Keep business logic in use cases
- Use dependency injection through port interfaces
- Write unit tests for all use cases
- Keep adapters thin and focused on translation
- Domain models should have no external dependencies 