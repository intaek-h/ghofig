# AGENTS.md - Ghofig

This document provides guidance for AI coding agents working on the Ghofig codebase.

Ghofig is a TUI-based CLI tool for browsing and managing Ghostty terminal configuration.
Built with Go using the Charm ecosystem (Bubbletea, Bubbles, Lipgloss).

## Build, Test, and Lint Commands

### Build Commands
```bash
make build                           # Build binary to bin/ghofig
make run                             # Build and run
go build -o bin/ghofig ./cmd/ghofig  # Alternative direct build
go mod tidy                          # Install dependencies
make parse                           # Parse reference docs to generate database (dev-time only)
make clean                           # Clean build artifacts
```

### Test Commands
```bash
go test ./...                        # Run all tests
go test -v ./...                     # Run all tests with verbose output
go test ./internal/db/               # Run tests for a specific package
go test ./internal/db/ -run TestSearch  # Run a single test by name
go test -cover ./...                 # Run tests with coverage
```

### Lint and Format Commands
```bash
go fmt ./...                         # Format code (standard Go formatting)
go vet ./...                         # Run static analysis
go build ./... && go vet ./... && go test ./...  # Check for issues before commit
```

## Project Structure

```
ghofig/
├── cmd/
│   ├── ghofig/main.go        # CLI entry point
│   └── parser/main.go        # Dev-time doc parser
├── internal/
│   ├── config/               # Ghostty config file operations
│   ├── db/                   # SQLite database operations
│   ├── model/                # Data structures
│   └── tui/                  # Bubbletea TUI components
├── data/ghofig.db            # Generated SQLite database
├── embed.go                  # go:embed directive
├── reference.mdx.txt         # Source docs from Ghostty
└── Makefile
```

## Code Style Guidelines

### Imports
- Group imports in three blocks separated by blank lines:
  1. Standard library
  2. Third-party packages
  3. Internal packages
- Use import aliases only when necessary (e.g., `tea "github.com/charmbracelet/bubbletea"`)
- Underscore imports for side effects go with their group (e.g., `_ "modernc.org/sqlite"`)

```go
import (
    "fmt"
    "strings"

    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"

    "github.com/intaek-h/ghofig/internal/db"
    "github.com/intaek-h/ghofig/internal/model"
)
```

### Naming Conventions
- Use `PascalCase` for exported types, functions, constants
- Use `camelCase` for unexported identifiers and local variables
- Use descriptive names: `SearchModel`, `doSearch`, `getAllConfigs`
- Model types end with `Model` suffix: `MenuModel`, `SearchModel`, `DetailModel`
- Message types end with `Msg` suffix: `searchResultMsg`, `tea.KeyMsg`
- Use short variable names for receivers: `m` for Model, `r` for result
- Constants use camelCase for unexported: `colorBlue700`, `ThemePrimary` for exported

### Types
- Define structs with exported fields using PascalCase
- Use value receivers for read-only methods, pointer receivers when mutating
- Document exported types with comments starting with the type name

```go
// Config represents a Ghostty configuration option.
type Config struct {
    ID          int
    Title       string
    Description string
}
```

### Error Handling
- Return errors as the last return value
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Check errors immediately after function calls
- Log fatal errors to stderr and exit with code 1

```go
if err := db.Init(ghofig.EmbeddedDB); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
    os.Exit(1)
}
```

### Bubbletea TUI Patterns
- Each view is a separate model in `internal/tui/`
- Models implement the `Init()`, `Update()`, `View()` pattern
- Use enum-based view routing in the main `Model` (MenuView, SearchView, DetailView)
- Return `tea.Cmd` for async operations
- Use `tea.Batch()` to combine multiple commands

```go
// Command pattern for async operations
func doSearch(query string) tea.Cmd {
    return func() tea.Msg {
        results, err := db.Search(query)
        return searchResultMsg{results: results, query: query, err: err}
    }
}
```

### Styling (Lipgloss)
- Define all colors in `internal/tui/theme.go`
- Use semantic color names: `ThemePrimary`, `ThemeTextMuted`, `ThemeBorder`
- Define style constants at package level using `var` blocks
- Chain style methods for readability

```go
var searchTitleStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(ThemePrimary)
```

### Formatting
- Use `go fmt` for all code formatting
- Use tabs for indentation (Go standard)
- Keep lines under 100 characters when practical
- Add blank lines between logical sections

### Comments
- Document all exported functions, types, and constants
- Use complete sentences starting with the identifier name
- Keep comments concise and focused on "why" not "what"

### Package Organization
- Use `internal/` for packages that should not be imported externally
- Keep packages focused on a single responsibility
- `cmd/` contains only main packages
- Models/data structures go in `internal/model/`

## Key Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/bubbletea` | TUI framework (Elm architecture) |
| `github.com/charmbracelet/bubbles` | UI components (textinput, list, viewport) |
| `github.com/charmbracelet/lipgloss` | Terminal styling |
| `modernc.org/sqlite` | Pure Go SQLite (no CGO) |

## Database

The SQLite database is embedded in the binary using `go:embed`:
- Database is read-only at runtime
- Parser generates `data/ghofig.db` from `reference.mdx.txt`
- Always use `defer rows.Close()` after query operations
- Use parameterized queries to prevent SQL injection

## Common Patterns

### View State Management
```go
type View int

const (
    MenuView View = iota
    SearchView
    DetailView
)
```

### Message Handling
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle key presses
    case tea.WindowSizeMsg:
        // Handle resize
    case customMsg:
        // Handle custom messages
    }
    return m, nil
}
```
