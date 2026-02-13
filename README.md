# Flow Boards CLI

A command-line tool for viewing your assigned Flow Boards tickets directly in the terminal.

## Features

### Viewing & Filtering
- View all tickets assigned to you
- Filter tickets by bin name or ID
- Filter tickets by board name or ID
- List all available bins and boards
- Real-time Flow Boards API integration
- Displays ticket ID, name, status, description, and dates
- Word wrapping for long descriptions
- Support for Unicode and special characters
- Handles large numbers of tickets efficiently

### Ticket Checkout Workflow
- **Check out a ticket** to work on it across multiple sessions
- **Quick comments** without repeated ticket selection (80% time savings)
- **Smart bin context** remembers your last used bin
- **Visual indicators** show checked-out tickets in lists
- **Persistent state** across terminal sessions and reboots
- **Force replace** checked-out tickets when needed
- **Auto-recovery** from corrupted or invalid checkout states

### Commenting
- Add comments to any ticket interactively
- Quick comment on checked-out ticket with single command
- Filter comments by bin before selection
- Persistent checkout for multiple comments per session

### Performance & Reliability
- Performance metrics with verbose mode
- Server-side filtering reduces API data transfer
- Automatic pagination for large datasets (200+ bins/tickets)
- Robust error handling and recovery

## Installation

### Prerequisites

- Go 1.25.6 or higher
- Flow Boards account with API access

### Build from Source

```bash
git clone https://github.com/Germanicus1/fb.git
cd fb
go build -o fb cmd/fb/main.go
```

### Install Globally

```bash
go install ./cmd/fb
```

This installs the binary to `$GOPATH/bin` as `fb`. Ensure this directory is in your PATH.

## Configuration

Create a configuration file at `~/.fb/config.yaml`:

```yaml
auth_key: "your-auth-key-here"
org_id: "your-organization-id"
user_email: "your.email@example.com"
```

### Configuration Fields

- **auth_key**: Your Flow Boards API authentication key
- **org_id**: Your organization ID
- **user_email**: Your email address (used to filter tickets)

The tool will automatically create the `~/.fb` directory on first run if it doesn't exist.

## Usage

### Display Tickets

```bash
# Show all tickets
fb

# Filter by bin name
fb --bin "In Progress"

# Filter by bin ID
fb --bin "kX41z9DVeVthZGe5d"
```

Shows all tickets assigned to you with:
- Ticket ID and name
- Status/bin information
- Created and updated dates
- Due dates (when present)
- Description (word-wrapped for readability)
- Visual indicator for checked-out tickets (← CHECKED OUT)

### List Bins and Boards

```bash
# List all bins with names
fb --list-bins

# List all boards with names
fb --list-boards
```

### Ticket Checkout Workflow (Recommended)

The checkout workflow saves 80% of time when adding multiple comments to the same ticket:

```bash
# 1. Check out a ticket from a specific bin
fb checkout --bin "Doing"
# Select ticket number from list

# 2. Add quick comments (no ticket selection needed!)
fb -c "Started implementation"
fb "Fixed bug in validation"        # Alternative: no -c flag needed
fb -c "Ready for review"

# 3. View current checkout
fb -o

# 4. Clear checkout when done
fb clear
```

**Direct checkout by ticket ID:**
```bash
fb checkout yL4rjYNU5PMlu7K8B
```

**Force replace an existing checkout:**
```bash
fb checkout --bin "Testing" --force
```

### Add Comments (Interactive)

```bash
# Interactive comment on any ticket
fb --comment

# Filter by bin first, then select and comment
fb --comment --bin "In Progress"
```

### Show Version

```bash
fb --version
# or
fb -v
```

### Show Help

```bash
fb --help
# or
fb -h
```

### Performance Metrics

```bash
fb --verbose
# or
fb --debug
```

Displays timing information including:
- API request duration
- Total execution time

## Example Output

### Ticket List with Checkout Indicator

```
Found 3 ticket(s) assigned to you:

[TICKET-123] Fix login bug ← CHECKED OUT
  Status: In Progress
  Updated: 2026-02-11
  Description: Users cannot log in with email addresses containing special
    characters. Need to update validation logic to handle Unicode properly.

[TICKET-456] Update documentation
  Status: To Do
  Updated: 2026-02-10
  Due: 2026-02-15
  Description: Add API examples to the REST documentation.

[TICKET-789] Code review
  Status: In Review
  Updated: 2026-02-11
  Description: Review pull request for the new authentication system.
```

### Checkout Status

```bash
$ fb -o
Currently checked out:
  Ticket: [TICKET-123] Fix login bug
  Bin: In Progress
  Checked out: 2 hours ago
```

### Quick Comment

```bash
$ fb -c "Fixed the Unicode validation issue"
✓ Comment added to: Fix login bug

# Alternative syntax (no -c flag needed)
$ fb "Fixed the Unicode validation issue"
✓ Comment added to: Fix login bug
```

## Error Handling

The tool provides clear error messages for common issues:

- **Missing configuration**: Displays setup instructions
- **Invalid YAML**: Shows syntax suggestions
- **API authentication errors**: Indicates credential issues
- **Network errors**: Suggests connectivity checks

## Development

### Run Tests

```bash
go test ./...
```

### Test Coverage

The project has comprehensive test coverage:
- 200+ total tests
- Unit tests for all packages
- Integration tests for end-to-end workflows
- Acceptance tests based on user stories
- ATDD methodology with Red-Green-Refactor cycles
- Full coverage of checkout workflow and edge cases

### Project Structure

The project follows standard Go project layout with clean architecture:

```
fb/
├── cmd/
│   └── fb/                    # Application entry point
│       ├── main.go           # Minimal bootstrap (15 lines)
│       └── main_test.go      # Integration tests
├── internal/                  # Private application code
│   ├── cli/                  # CLI framework (flags, help, routing)
│   ├── commands/             # Command handlers (6 commands + 29 tests)
│   ├── service/              # Business logic (3 services + 4 tests)
│   └── state/                # State persistence (3 files + 6 tests)
├── api/                      # Flow Boards API client
├── config/                   # Configuration management
├── formatter/                # Output formatting
├── models/                   # Domain models
├── docs/                     # Documentation
├── ARCHITECTURE.md           # Detailed architecture docs
└── README.md                 # This file
```

**Architecture Layers**:
- **Entry Point** (`cmd/fb/`): Bootstrap application
- **CLI Layer** (`internal/cli/`): Parse flags, route commands
- **Command Layer** (`internal/commands/`): Implement CLI commands
- **Service Layer** (`internal/service/`): Business logic and API orchestration
- **State Layer** (`internal/state/`): Persist application state

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed documentation.

## API Integration

The tool uses the Flow Boards REST API v2:

1. Discovers REST endpoint via organization ID
2. Authenticates using bearer token
3. Looks up user by email
4. Searches for tickets assigned to user
5. Formats and displays results

## Workflow Benefits

### Checkout Workflow Time Savings

**Without checkout (traditional):**
- Select ticket from list: 20-40 seconds
- Enter comment: 10-20 seconds
- Repeat for each comment: 50-100 seconds per 5 comments

**With checkout workflow:**
- Check out once: 20-30 seconds
- Quick comments: 2-3 seconds each
- Total for 5 comments: 10-20 seconds

**Time savings: 80%** (40-80 seconds saved per session)

### When to Use Checkout

- Adding multiple comments to the same ticket
- Working on a ticket throughout the day
- Documenting progress incrementally
- Focusing on one ticket at a time

### When to Use Interactive Comment

- One-off comments to different tickets
- Quick updates to various tickets
- Occasional commenting

## Performance

- Handles 200+ tickets efficiently
- Server-side filtering reduces data transfer
- API requests complete in under 2 seconds (typical)
- Memory efficient with large result sets
- Automatic pagination for 200+ bins/boards
- No artificial limits on ticket count
- Smart bin context reduces repeated navigation

## Troubleshooting

### Configuration File Not Found

Ensure `~/.fb/config.yaml` exists with valid credentials.

### Authentication Failed

Check that your `auth_key` is correct and has not expired.

### No Tickets Found

Verify that:
- Your `user_email` matches your Flow Boards account
- You have tickets assigned to you
- Your API key has appropriate permissions

### YAML Syntax Error

Common issues:
- Using tabs instead of spaces for indentation
- Missing colons after field names
- Incorrect indentation levels

Use a YAML validator if needed: https://www.yamllint.com/

### Checkout Issues

**"No ticket checked out"**
- Run `fb checkout --bin "Bin Name"` to check out a ticket first
- Or use `fb --comment` for one-off comments without checkout

**"Ticket already checked out"**
- Use `fb clear` to clear the current checkout
- Or use `fb checkout --force` to replace it

**"Ticket no longer exists"**
- The checked-out ticket may have been deleted or unassigned
- Checkout is automatically cleared
- Check out a different ticket

**Corrupted checkout state**
- The tool automatically recovers from corrupted `~/.fb/checkout.json`
- If issues persist, delete `~/.fb/checkout.json` manually

**Lost bin context**
- If `fb checkout` (without arguments) fails, specify bin explicitly
- Example: `fb checkout --bin "Doing"`

## Version

Current version: 2.0.0

### What's New in 2.0

- **Ticket Checkout Workflow**: Check out tickets and add multiple comments quickly (80% time savings)
- **Bin Filtering**: Filter tickets by bin name/ID
- **Quick Comments**: `fb -c "message"` for instant commenting on checked-out tickets
- **Smart Context**: Remembers last used bin for faster workflow
- **Visual Indicators**: See which ticket is checked out in all lists
- **Persistent State**: Checkout survives terminal restarts and system reboots
- **Auto-Recovery**: Automatically handles corrupted states and invalid tickets
- **200+ Tests**: Comprehensive test coverage with ATDD methodology

## Contributing

This project follows Test-Driven Development (TDD) and ATDD methodologies. All changes should:
- Include comprehensive tests
- Follow clean code principles
- Maintain backward compatibility
- Update documentation as needed

## Development Methodology

This project was built using:
- Acceptance Test-Driven Development (ATDD)
- Red-Green-Refactor cycle
- SOLID principles
- Clean code practices
- Elephant Carpaccio story slicing

