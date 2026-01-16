package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/intaek-h/ghofig/internal/model"
)

// DetailModel represents the config detail view.
type DetailModel struct {
	config *model.Config
	width  int
	height int
}

// NewDetailModel creates a new detail model.
func NewDetailModel() DetailModel {
	return DetailModel{}
}

// SetSize updates the detail dimensions.
func (m DetailModel) SetSize(width, height int) DetailModel {
	m.width = width
	m.height = height
	return m
}

// SetConfig sets the config to display.
func (m DetailModel) SetConfig(config *model.Config) DetailModel {
	m.config = config
	return m
}

// Update handles detail updates.
func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	return m, nil
}

// View renders the detail view.
func (m DetailModel) View() string {
	if m.config == nil {
		return "No config selected"
	}
	return "Detail (placeholder)\n\n" + m.config.Title + "\n\nPress 'esc' to go back\nPress 'q' to quit"
}
