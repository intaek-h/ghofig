# Ghofig

<p>
  <a href="https://github.com/intaek-h/ghofig/releases"><img src="https://img.shields.io/github/release/intaek-h/ghofig.svg" alt="Latest Release"></a>
  <a href="https://github.com/intaek-h/ghofig/actions"><img src="https://github.com/intaek-h/ghofig/actions/workflows/release.yml/badge.svg" alt="Build Status"></a>
  <a href="https://github.com/intaek-h/ghofig/blob/main/LICENSE"><img src="https://img.shields.io/github/license/intaek-h/ghofig" alt="License"></a>
</p>

A TUI for browsing and managing [Ghostty](https://ghostty.org/) terminal configuration.

<img alt="Ghofig TUI demo" width="600" src="https://github.com/intaek-h/ghofig/raw/main/assets/demo.gif">

## Features

- A more intuitive view than Ghostty Docs
- Search by name or description
- Edit config directly without opening a new Text Editor

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

## How It Works

I parsed their raw doc mdx file and dumped the data to the embeded sqlite db.

## Contributing

Always welcome.

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

## Thanks to

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Mole](https://github.com/tw93/mole) - Design Reference

## License

[MIT](LICENSE)
