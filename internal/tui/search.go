package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/intaek-h/ghofig/internal/model"
)

// SearchModel represents the search view.
type SearchModel struct {
	width  int
	height int
}

// NewSearchModel creates a new search model.
func NewSearchModel() SearchModel {
	return SearchModel{}
}

// SetSize updates the search dimensions.
func (m SearchModel) SetSize(width, height int) SearchModel {
	m.width = width
	m.height = height
	return m
}

// Init initializes the search view.
func (m SearchModel) Init() tea.Cmd {
	return nil
}

// Update handles search updates.
func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	return m, nil
}

// View renders the search view.
func (m SearchModel) View() string {
	return "Search (placeholder)\n\nPress 'esc' to go back\nPress 'q' to quit"
}

// IsInputFocused returns whether the search input is focused.
func (m SearchModel) IsInputFocused() bool {
	return false
}

// SelectedConfig returns the currently selected config.
func (m SearchModel) SelectedConfig() *model.Config {
	return nil
}
