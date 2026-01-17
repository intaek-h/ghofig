# Ghofig

<p>
  <a href="https://github.com/intaek-h/ghofig/releases"><img src="https://img.shields.io/github/release/intaek-h/ghofig.svg" alt="Latest Release"></a>
  <a href="https://github.com/intaek-h/ghofig/actions"><img src="https://github.com/intaek-h/ghofig/actions/workflows/release.yml/badge.svg" alt="Build Status"></a>
  <a href="https://github.com/intaek-h/ghofig/blob/main/LICENSE"><img src="https://img.shields.io/github/license/intaek-h/ghofig" alt="License"></a>
</p>

A TUI for browsing and managing [Ghostty](https://ghostty.org/) terminal configuration.

**Ghofig** = **Gho**stty + Con**fig**

<img alt="Ghofig TUI demo" width="600" src="https://github.com/intaek-h/ghofig/raw/main/assets/demo.gif">

## Features

- Browse all 180+ Ghostty configuration options
- Search configs by name or description
- View detailed documentation for each option
- Edit your Ghostty config file directly from the TUI
- Fully offline - no network required

## Installation

### Homebrew (macOS / Linux)

```bash
brew install intaek-h/ghofig/ghofig
```

### Go

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

### Download Binary

Download the latest binary from the [releases page](https://github.com/intaek-h/ghofig/releases).

## Usage

```bash
ghofig
```

### Keybindings

| Key | Action |
|-----|--------|
| `q`, `ctrl+c` | Quit |
| `esc`, `backspace` | Go back |
| `up/down`, `k/j` | Navigate |
| `enter` | Select |
| `tab` | Toggle focus (search view) |
| `g/G` | Jump to top/bottom (detail view) |

## How It Works

Ghofig parses the official [Ghostty configuration reference](https://ghostty.org/docs/config/reference) and stores it in an embedded SQLite database. The database is compiled into the binary, making the tool completely self-contained and fast.

```
reference.mdx.txt -> parser -> ghofig.db -> go:embed -> binary
```

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

### Development Setup

```bash
git clone https://github.com/intaek-h/ghofig.git
cd ghofig
go mod tidy
make run
```

### Updating Config Database

When Ghostty releases new configuration options:

1. Download the latest config reference from [Ghostty's docs](https://github.com/ghostty-org/ghostty)
2. Replace `reference.mdx.txt`
3. Run `make parse` to regenerate the database

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) - Pure Go SQLite

## License

[MIT](LICENSE)
