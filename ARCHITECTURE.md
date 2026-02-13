# Flow Boards CLI - Architecture Documentation

## Overview

Flow Boards CLI is a command-line tool for viewing and managing Flow Boards tickets. The codebase follows standard Go project layout and clean architecture principles with clear separation of concerns.

## Project Structure

```
fb/
├── cmd/
│   └── fb/
│       ├── main.go              # Application entry point (15 lines)
│       └── main_test.go         # Integration tests
│
├── internal/                    # Private application code
│   ├── cli/                     # CLI framework
│   │   ├── app.go              # Main orchestration and routing
│   │   ├── flags.go            # Flag definitions and parsing
│   │   └── help.go             # Help text and documentation
│   │
│   ├── commands/                # Command handlers
│   │   ├── list.go             # Main ticket listing
│   │   ├── list_bins.go        # List bins command
│   │   ├── list_boards.go      # List boards command
│   │   ├── checkout.go         # Checkout workflow
│   │   ├── comment.go          # Comment commands
│   │   ├── status.go           # Status command
│   │   └── *_test.go           # 29 test files
│   │
│   ├── service/                 # Business logic layer
│   │   ├── ticket_service.go   # Ticket operations
│   │   ├── comment_service.go  # Comment operations
│   │   ├── bin_service.go      # Bin resolution
│   │   └── *_test.go           # 4 test files
│   │
│   └── state/                   # State persistence
│       ├── types.go            # State type definitions
│       ├── checkout.go         # Checkout state management
│       ├── bincontext.go       # Bin context management
│       └── *_test.go           # 6 test files
│
├── api/                         # API client
│   ├── client.go               # HTTP client for Flow Boards API
│   └── *_test.go               # API tests
│
├── config/                      # Configuration management
│   ├── config.go               # Config loading and validation
│   └── config_test.go
│
├── models/                      # Domain models
│   └── models.go               # Ticket, Bin, Board, User types
│
├── formatter/                   # Output formatting
│   ├── formatter.go            # Ticket list formatting
│   └── formatter_test.go
│
├── docs/                        # Documentation
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── README.md                    # User documentation
└── ARCHITECTURE.md              # This file
```

## Architecture Layers

### 1. Entry Point Layer (`cmd/fb/`)

**Purpose**: Bootstrap the application

**Key File**: `main.go` (15 lines)
- Defines version constant
- Calls `cli.Run()`
- Exits with error code if needed

**Principle**: Minimal entry point delegates all logic to internal packages.

### 2. CLI Layer (`internal/cli/`)

**Purpose**: Parse flags, route commands, handle errors

**Key Files**:
- `app.go`: Main orchestration (`Run` function)
- `flags.go`: Flag definitions (`Flags` struct, `parseFlags`)
- `help.go`: Help text (`PrintHelp`, `GetHelpText`)

**Responsibilities**:
- Parse command-line arguments
- Route to appropriate command handlers
- Handle version/help flags
- Error formatting and exit codes
- Subcommand detection (checkout, clear)

**Principle**: CLI layer knows about commands but not business logic.

### 3. Command Layer (`internal/commands/`)

**Purpose**: Implement CLI commands

**Key Files**:
- `list.go`: Main ticket listing with checkout indicator
- `list_bins.go`: List available bins
- `list_boards.go`: List available boards
- `checkout.go`: Ticket checkout workflow
- `comment.go`: Interactive and quick comments
- `status.go`: Show checkout status

**Command Signatures**:
```go
Execute(cfg *config.Config, binFilter string, verbose bool) error
ExecuteListBins(cfg *config.Config) error
ExecuteCheckout(args []string, binFlag string, forceFlag bool) error
ExecuteQuick(comment string) error
```

**Responsibilities**:
- Accept parsed arguments from CLI layer
- Orchestrate service calls
- Handle user interaction (prompts, selection)
- Format and display output
- Manage command-specific state

**Principle**: Commands orchestrate but don't implement business logic.

### 4. Service Layer (`internal/service/`)

**Purpose**: Business logic and API orchestration

**Key Files**:
- `ticket_service.go`: Ticket CRUD operations
- `comment_service.go`: Comment generation and posting
- `bin_service.go`: Bin name/ID resolution

**Service Pattern**:
```go
type TicketService struct {
    client *api.Client
    cfg    *config.Config
}

func NewTicketService(cfg *config.Config) (*TicketService, error)
func (s *TicketService) GetUserTickets(userID string) ([]models.Ticket, error)
```

**Responsibilities**:
- Initialize and manage API clients
- Implement business operations
- Transform API responses
- Centralize error handling
- Provide testable interfaces

**Principle**: Services encapsulate business logic and API interactions.

### 5. State Layer (`internal/state/`)

**Purpose**: Persist application state

**Key Files**:
- `types.go`: State type definitions (`CheckoutState`, `BinContext`)
- `checkout.go`: Checkout state persistence
- `bincontext.go`: Bin context persistence

**State Management**:
```go
func SaveCheckout(checkout *CheckoutState) error
func LoadCheckout() (*CheckoutState, error)
func ClearCheckout() error
```

**Storage**:
- File-based JSON storage in `~/.fb/`
- `checkout.json`: Currently checked-out ticket
- `bin_context.json`: Last used bin

**Principle**: State is isolated and persisted independently.

### 6. Supporting Packages

**api/** - HTTP client for Flow Boards API
- RESTful API interactions
- Endpoint discovery
- Request/response handling

**config/** - Configuration management
- YAML config loading (`~/.fb/config.yaml`)
- Validation
- Credential storage

**models/** - Domain types
- `Ticket`, `Bin`, `Board`, `User`
- Shared across all layers
- Pure data structures (no logic)

**formatter/** - Output formatting
- Ticket list formatting
- Duration formatting
- Display helpers

## Dependency Flow

```
cmd/fb/main.go
    ↓
internal/cli/app.go
    ↓
internal/commands/*.go
    ↓
internal/service/*.go
    ↓
api/client.go
    ↓
[Flow Boards API]

Crosscutting:
- internal/state/ (accessed by commands)
- config/ (used everywhere)
- models/ (used everywhere)
- formatter/ (used by commands)
```

**Rules**:
- Higher layers depend on lower layers only
- No circular dependencies
- `internal/` packages are private to this module
- External packages (`api`, `config`, `models`, `formatter`) are reusable

## Key Design Decisions

### 1. Standard Go Project Layout

**Why**: Industry standard, familiar to Go developers
**Benefit**: Easy onboarding, clear structure

### 2. Layered Architecture

**Why**: Separation of concerns, testability
**Benefit**: Easy to modify one layer without affecting others

### 3. Command Pattern

**Why**: Each command is independent and testable
**Benefit**: Easy to add new commands

### 4. Service Layer

**Why**: Centralize business logic and API interactions
**Benefit**: Reusable across commands, easier testing

### 5. File-based State

**Why**: Simple, no database needed
**Benefit**: Portable, version-controllable

### 6. Explicit Error Handling

**Why**: Go idiom, clear error flow
**Benefit**: Errors handled at appropriate level

## Testing Strategy

### Test Organization

Tests live beside the code they test:
```
internal/commands/
├── checkout.go
├── checkout_test.go              # Basic tests
├── checkout_validation_test.go   # Validation tests
├── checkout_duplicate_prevention_test.go
└── ...
```

### Test Categories

1. **Unit Tests**: Test individual functions
2. **Integration Tests**: Test command workflows
3. **Service Tests**: Test API interactions (mocked)
4. **State Tests**: Test persistence logic

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/commands

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

## Adding New Features

### Adding a New Command

1. Create `internal/commands/newcommand.go`:
```go
package commands

func ExecuteNewCommand(cfg *config.Config) error {
    // Implementation
}
```

2. Add routing in `internal/cli/app.go`:
```go
if flags.NewCommand {
    return commands.ExecuteNewCommand(cfg)
}
```

3. Add flag in `internal/cli/flags.go`:
```go
fs.BoolVar(&flags.NewCommand, "newcmd", false, "Description")
```

4. Update help text in `internal/cli/help.go`

5. Write tests in `internal/commands/newcommand_test.go`

### Adding a New Service

1. Create `internal/service/newservice.go`:
```go
package service

type NewService struct {
    client *api.Client
}

func NewNewService(client *api.Client) *NewService {
    return &NewService{client: client}
}
```

2. Use from commands:
```go
svc := service.NewNewService(client)
result, err := svc.DoSomething()
```

## Build and Install

### Development Build

```bash
go build -o fb cmd/fb/main.go
./fb --version
```

### Install to $GOPATH/bin

```bash
go install ./cmd/fb
fb --version
```

### Release Build

```bash
go build -ldflags="-s -w" -o fb cmd/fb/main.go
```

## Configuration

### Required Files

1. **~/.fb/config.yaml**: User credentials
```yaml
auth_key: your-api-key
org_id: your-org-id
user_email: your-email@example.com
```

2. **~/.fb/checkout.json**: Current checkout (auto-managed)
3. **~/.fb/bin_context.json**: Last bin (auto-managed)

## Performance Considerations

- **API Calls**: Minimized through caching and smart fetching
- **State Files**: Small JSON files, fast I/O
- **Memory**: Minimal footprint, suitable for CLI
- **Startup**: Fast, <100ms typical

## Security

- **Credentials**: Stored in `~/.fb/config.yaml` (user's home directory)
- **API Key**: Never logged or displayed
- **HTTPS**: All API calls use HTTPS
- **No Network Logging**: Sensitive data not logged

## Maintenance

### Code Organization Principles

1. **Single Responsibility**: Each file has one clear purpose
2. **Small Functions**: Most functions <50 lines
3. **Clear Naming**: Function names describe what they do
4. **Minimal Dependencies**: Avoid deep dependency chains
5. **Test Coverage**: All business logic tested

### Common Maintenance Tasks

- **Update API**: Modify `api/client.go`
- **Add Command**: See "Adding New Features" above
- **Change State Format**: Update `internal/state/types.go`
- **Modify Help**: Edit `internal/cli/help.go`

## Version History

- **v2.0.0**: Major refactoring to layered architecture
- **v1.x**: Monolithic structure (single main.go)

## Future Improvements

- [ ] Phase 3: Extract display logic to `internal/display/`
- [ ] Add command aliases
- [ ] Plugin system for custom commands
- [ ] Configuration profiles (multiple accounts)
- [ ] Shell completion scripts
- [ ] Performance metrics dashboard

---

**Last Updated**: 2026-02-13
**Architecture Version**: 2.0.0
