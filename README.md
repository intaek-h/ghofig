# Ghofig

A TUI-based CLI tool for browsing and managing Ghostty terminal configuration.

**Ghofig** = **Ghostty** + **Config**

## Features (PoC)

- Browse Ghostty configuration options with an intuitive TUI
- Search through 180+ config options by name or description
- View detailed documentation for each option

## Installation

```bash
go install github.com/intaek-h/ghofig/cmd/ghofig@latest
```

## Usage

```bash
ghofig
```

### Keybindings

| Key | Action |
|-----|--------|
| `q` | Quit application |
| `esc` / `backspace` | Go back to previous view |
| `enter` | Select item |
| `up/down` or `k/j` | Navigate lists |

## Development

### Prerequisites

- Go 1.22+

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

### Updating Config Documentation

When Ghostty releases new configuration options:

1. Download the latest config reference from Ghostty's repository
2. Replace `reference.mdx.txt`
3. Run `make parse` to regenerate the database
4. Rebuild the binary

## Future Scope

- Direct config file editing from the TUI
- Append config options from search results
- Config validation
- Theme preview

## License

MIT
