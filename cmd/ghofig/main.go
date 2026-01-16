package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	ghofig "github.com/intaek-h/ghofig"
	"github.com/intaek-h/ghofig/internal/db"
	"github.com/intaek-h/ghofig/internal/tui"
)

func main() {
	// Initialize database from embedded bytes
	if err := db.Init(ghofig.EmbeddedDB); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create and run the TUI
	p := tea.NewProgram(
		tui.New(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
