# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DST Admin Go is a web-based management panel for "Don't Starve Together" game servers. Written in Go, it provides a simple deployment, low memory footprint, and visual interface for managing game configurations, mods, clusters, backups, and multi-server operations across both Windows and Linux platforms.

## Development Commands

### Running the Application

```bash
# Install dependencies
go mod tidy

# Run the application
go run cmd/server/main.go
```

The server starts on port 8082 by default (configurable in `config.yml`).

### Building

```bash
# Build for Linux
bash scripts/build_linux.sh
# Output: dst-admin-go (Linux amd64 binary)

# Build for Windows
bash scripts/build_window.sh
# Output: dst-admin-go.exe (Windows amd64 binary)

# Cross-platform build from Windows to Linux
set GOARCH=amd64
set GOOS=linux
go build
```

### Testing

No test files currently exist in the project. When writing tests, place them adjacent to the code they test with `_test.go` suffix.

## Code Architecture

### Entry Point

- **Main entry**: `cmd/server/main.go` - Loads config, initializes database, sets up router, and starts the HTTP server.

### Project Structure

```
dst-admin-go/
├── cmd/server/           # Application entry point
├── internal/             # Private application code
│   ├── api/             # HTTP layer
│   │   ├── handler/     # HTTP handlers (controllers)
│   │   └── router.go    # Route registration and DI setup
│   ├── config/          # Configuration loading
│   ├── database/        # Database initialization
│   ├── middleware/      # HTTP middleware (auth, error handling, cluster context)
│   ├── model/           # Database models (GORM entities)
│   ├── pkg/             # Internal shared utilities
│   │   ├── response/    # Standard HTTP response helpers
│   │   └── utils/       # Utility functions (file, shell, system, etc.)
│   └── service/         # Business logic layer
│       ├── archive/     # Archive path resolution
│       ├── backup/      # Backup management
│       ├── dstConfig/   # DST configuration management
│       ├── dstPath/     # Platform-specific path handling
│       ├── game/        # Game process management (start/stop/command)
│       ├── gameArchive/ # Game archive operations
│       ├── gameConfig/  # Game configuration files
│       ├── level/       # Level (world) management
│       ├── levelConfig/ # Level configuration parsing
│       ├── login/       # Authentication
│       ├── mod/         # Mod management
│       ├── player/      # Player management
│       └── update/      # Game update management
├── scripts/             # Build and utility scripts
└── config.yml          # Application configuration file
```

### Layered Architecture

The codebase follows a three-layer architecture:

1. **Handler Layer** (`internal/api/handler/`): HTTP request handling, input validation, response formatting
2. **Service Layer** (`internal/service/`): Business logic, orchestration between services
3. **Model Layer** (`internal/model/`): Database entities (GORM models)

### Dependency Injection

All services are instantiated in `internal/api/router.go` using constructor functions (`New{Service}Service`), then injected into handlers. This pattern:
- Avoids global variables
- Makes dependencies explicit
- Enables testing with mock implementations

Example flow:
```
router.go creates services → injects into handlers → handlers registered to routes
```

### Platform Abstraction

Services use factory patterns for platform-specific implementations:
- **Game Process**: `game.NewGame()` returns `LinuxProcess` or `WindowProcess` based on `runtime.GOOS`
- **Update Service**: `update.NewUpdateService()` handles Linux/Windows differences
- **DST Paths**: `dstPath` package provides platform-specific path resolution

### Service Interfaces

Core services define interfaces for flexibility and testing:
- `dstConfig.Config`: DST configuration CRUD operations
- `game.Process`: Game server lifecycle management
- Simpler services may omit interfaces and use concrete types directly

## Configuration

The application reads `config.yml` in the working directory:

```yaml
bindAddress: ""        # Bind address (empty = all interfaces)
port: 8082            # HTTP server port
database: dst-db      # SQLite database filename
steamAPIKey: ""       # Steam API key (optional)
autoCheck:            # Auto-check intervals (in minutes)
  masterInterval: 5
  cavesInterval: 5
  masterModInterval: 10
  # ... other intervals
```

## Database

- **ORM**: GORM with SQLite (via glebarez/sqlite)
- **Initialization**: `internal/database/sqlite.go` - Auto-migrates all models in `internal/model/`
- **Models**: Represent game data (clusters, players, mods, backups, logs, etc.)

## API Refactoring Guidelines

The project is undergoing a refactoring effort documented in `.sisyphus/plans/`. Key principles:

### Code Organization Rules

1. **No global variables**: Pass config and dependencies through constructors
2. **No constants package**: Define constants locally using `const ()` blocks within each module
3. **No VO suffix**: Data structures defined in services use PascalCase without VO suffix
4. **Internal-only changes**: Do not modify code outside `internal/` directory
5. **Dependency injection**: All service dependencies injected via constructors

### Handler Pattern

Each handler:
- Is in `internal/api/handler/{name}_handler.go`
- Has a `RegisterRoute(router *gin.RouterGroup)` method
- Uses injected services for business logic
- Uses `internal/pkg/response` for standardized responses

### Service Pattern

Each service:
- Is in `internal/service/{domain}/` directory
- Has a constructor `New{Service}Service()` accepting dependencies
- Defines interfaces for core abstractions (optional for simple services)
- Handles a single domain of business logic

### Naming Conventions

- **Files**: `snake_case.go` (e.g., `backup_service.go`)
- **Interfaces**: `{ServiceName}` (e.g., `Config`, `Process`)
- **Structs**: `PascalCase` without suffix (e.g., `BackupSnapshot` not `BackupSnapshotVO`)
- **Constructors**: `New{Name}` (e.g., `NewBackupService`)

## Game Server Management

The application manages Don't Starve Together dedicated servers:

- **Clusters**: A cluster contains multiple "levels" (worlds) - typically Master (overworld) and Caves
- **Process Management**: Uses platform-specific commands (screen/tmux on Linux, custom CLI on Windows)
- **Configuration Files**: Parses and modifies Lua config files (`cluster.ini`, `leveldataoverride.lua`, `modoverrides.lua`)
- **Session Names**: Identifies running processes by cluster and level names

## Utilities

Common utilities in `internal/pkg/utils/`:
- **fileUtils**: File operations, archive extraction
- **shellUtils**: Execute shell commands
- **systemUtils**: System information (CPU, memory)
- **dstConfigUtils**: DST config file parsing
- **luaUtils**: Lua table parsing/generation
- **collectionUtils**: Slice/map helpers
- **clusterUtils**: Cluster-specific utilities

## Development Notes

- Go version: 1.20+
- Web framework: Gin
- Session management: gin-contrib/sessions with memstore
- Authentication: Session-based, checked via `middleware.Authentication`
- All routes except `/hello` and login endpoints require authentication
- Chinese comments and messages common in codebase (target audience is Chinese users)