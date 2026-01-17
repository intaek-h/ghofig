# Ghofig

A TUI-based CLI tool for browsing and managing Ghostty terminal configuration.

**Ghofig** = **Ghostty** + **Config**

## Features (PoC)

- Browse Ghostty configuration options with an intuitive TUI
- Search through 182 config options by name or description
- View detailed documentation for each option
- Embedded database - works offline, no external dependencies

## Installation

### Homebrew (macOS/Linux)

```bash
brew install intaek-h/ghofig/ghofig
```

### Go Install

```bash
go install github.com/intaek-h/ghofig/cmd/ghofig@latest
```

### From Source

```bash
git clone https://github.com/intaek-h/ghofig.git
cd ghofig
make build
./bin/ghofig
```

## Usage

```bash
ghofig
```

### Navigation

The app has three views:

1. **Main Menu** - Select "Configs" to browse configuration options
2. **Search** - Type to filter configs, navigate results with arrow keys
3. **Detail** - View full documentation for a config option (scrollable)

### Keybindings

#### Global
| Key | Action |
|-----|--------|
| `q` / `ctrl+c` | Quit application |
| `esc` / `backspace` | Go back to previous view |

#### Menu & Search Results
| Key | Action |
|-----|--------|
| `up/down` or `k/j` | Navigate list |
| `enter` | Select item |
| `tab` | Toggle focus (search view) |

#### Detail View
| Key | Action |
|-----|--------|
| `up/down` or `k/j` | Scroll line by line |
| `pgup/pgdn` | Scroll half page |
| `home/end` or `g/G` | Jump to top/bottom |

## Development

### Prerequisites

- Go 1.24+

### Setup

```bash
# Clone the repository
git clone https://github.com/intaek-h/ghofig.git
cd ghofig

# Install dependencies
go mod tidy

# Parse the config documentation (generates data/ghofig.db)
make parse

# Build and run
make run
```

### Project Structure

```
ghofig/
├── cmd/
│   ├── ghofig/          # Main CLI entry point
│   └── parser/          # Dev-time parser for config docs
├── internal/
│   ├── db/              # SQLite database operations
│   ├── model/           # Data structures
│   └── tui/             # Bubbletea TUI components
├── data/
│   └── ghofig.db        # Generated config database (embedded at build)
├── reference.mdx.txt    # Source documentation from Ghostty
└── embed.go             # go:embed directive for database
```

### Makefile Commands

```bash
make parse    # Parse reference.mdx.txt and generate database
make build    # Build binary to bin/ghofig
make run      # Build and run
make clean    # Remove build artifacts
```

### Updating Config Documentation

When Ghostty releases new configuration options:

1. Download the latest config reference from [Ghostty's repository](https://github.com/ghostty-org/ghostty)
2. Replace `reference.mdx.txt`
3. Run `make parse` to regenerate the database
4. Rebuild the binary

## Architecture

```
Development time:
  reference.mdx.txt → parser → ghofig.db

Build time:
  ghofig.db → go:embed → binary

Runtime:
  binary → in-memory SQLite → TUI
```

The config documentation is parsed once during development and stored in an SQLite database. This database is embedded into the binary at compile time using Go's `embed` package, making the tool completely self-contained.

## Future Scope

- Direct config file editing from the TUI
- Append config options from search results
- Config validation
- Theme preview
- Fuzzy search

## License

MIT
