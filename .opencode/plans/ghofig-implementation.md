# Ghofig Implementation Plan

## Overview

**Ghofig** = **Ghostty** + **Config**

A TUI-based CLI tool for browsing and managing Ghostty terminal configuration. This PoC enables users to:

1. Launch `ghofig` to open a full-terminal TUI
2. Navigate a menu and select "Configs"
3. Search through 182 Ghostty config options (parsed from official docs)
4. View detailed descriptions for each config option

## Architecture

```
Development:
  reference.mdx.txt → parser → ghofig.db

Build:
  ghofig.db → go:embed → binary

Runtime:
  binary (with embedded DB) → in-memory SQLite → TUI
```

## Project Structure

```
ghofig/
├── cmd/
│   ├── ghofig/          # Main CLI entry point
│   └── parser/          # Dev-time parser for config docs
├── internal/
│   ├── db/              # SQLite database operations
│   ├── model/           # Data structures (Config)
│   └── tui/             # Bubbletea TUI components
├── data/
│   └── ghofig.db        # Generated config database (embedded at build)
├── reference.mdx.txt    # Source documentation from Ghostty
└── embed.go             # go:embed directive for database
```

## Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - Input, list components
- `github.com/charmbracelet/lipgloss` - Styling
- `modernc.org/sqlite` - Pure Go SQLite (no CGO)

## Database Schema

```sql
CREATE TABLE configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,           -- e.g., "font-family"
    description TEXT NOT NULL      -- raw markdown description
);

CREATE INDEX idx_configs_title ON configs(title);
```

## Keybindings

| Key | Action |
|-----|--------|
| `q` | Quit application |
| `esc` / `backspace` | Go back to previous view |
| `enter` | Select item |
| `up/down` or `k/j` | Navigate lists |

## TUI Flow

```
┌─────────────────────────┐
│       Main Menu         │
│  > Configs              │
│    (future items...)    │
└───────────┬─────────────┘
            │ Enter
            ▼
┌─────────────────────────┐
│    Search Configs       │
│  🔍 font_                │
│  ─────────────────────  │
│  > font-family          │
│    font-size            │
│    font-style           │
└───────────┬─────────────┘
            │ Enter
            ▼
┌─────────────────────────┐
│    Config Detail        │
│                         │
│  font-family            │
│  ─────────────────────  │
│  The font family to     │
│  use for text...        │
│                         │
│  [esc: back]            │
└─────────────────────────┘
```

## Implementation Commits

| # | Commit | Status |
|---|--------|--------|
| 1 | `init: project setup with go.mod and structure` | ✅ Completed |
| 2 | `feat: add config parser tool` | ✅ Completed |
| 3 | `feat: add db layer with embedded sqlite` | ✅ Completed |
| 4 | `feat: add tui app shell with view routing` | ✅ Completed |
| 5 | `feat: add main menu view` | ✅ Completed |
| 6 | `feat: add search view with db integration` | ✅ Completed |
| 7 | `feat: add detail view` | ✅ Completed |
| 8 | `chore: polish styling and update readme` | ✅ Completed |

## Commit Details

### Commit 1: Project Setup ✅
- go.mod with module `github.com/intaek-h/ghofig`
- Directory structure (cmd/, internal/, data/)
- Makefile with parse/build/run/clean targets
- .gitignore for bin/ and generated db
- README with project overview
- reference.mdx.txt (Ghostty config docs source)

### Commit 2: Config Parser ✅
- `cmd/parser/main.go` - parses reference.mdx.txt
- Handles single and multiple consecutive h2 headers
- Strips backticks from config titles
- Outputs to `data/ghofig.db`

### Commit 3: DB Layer ✅
- `embed.go` with go:embed directive
- `internal/db/db.go` with Init, Search, GetByID
- Search prioritizes title matches over description
- DB tests

### Commit 4: TUI App Shell ✅
- `internal/tui/app.go` - main model with view routing
- View states: MenuView, SearchView, DetailView
- Global keybindings (q, esc, backspace)
- `cmd/ghofig/main.go` - entry point with DB init

### Commit 5: Main Menu ✅
- `internal/tui/menu.go` with bubbles/list
- Single item: "Configs"
- Styled with lipgloss
- Enter → transitions to Search view

### Commit 6: Search View ✅
- `internal/tui/search.go`
- Text input (bubbles/textinput) at top
- Results list below showing title + truncated description
- On input change → query DB → update list
- Search SQL prioritizes title matches
- Tab to toggle focus between input and results

### Commit 7: Detail View ✅
- `internal/tui/detail.go`
- Header: config title (styled)
- Body: full description (scrollable viewport)
- esc/backspace → back to Search
- Support for pgup/pgdn, home/end navigation

### Commit 8: Polish ✅
- Consistent styling across views
- Help text footer in all views
- Updated README with detailed usage instructions
- Architecture documentation

## Future Scope (Post-PoC)

- Direct config file editing from the TUI
- Append config options from search results
- Config validation
- Theme preview

## Development Commands

```bash
# Parse config docs to generate database
make parse

# Build the binary
make build

# Build and run
make run

# Clean build artifacts
make clean

# Run tests
go test ./...
```

## Parser Logic

The parser reads `reference.mdx.txt` and handles:
1. Single h2 → followed by its own paragraph
2. Multiple consecutive h2s → all share the next paragraph

Pattern: `## \`config-name\`` followed by description paragraphs until next h2.

Example input:
```
## `adjust-cell-width`
## `adjust-cell-height`

Description for both options...

## `font-size`

Description for font-size...
```

Produces entries:
- adjust-cell-width → "Description for both options..."
- adjust-cell-height → "Description for both options..."
- font-size → "Description for font-size..."
