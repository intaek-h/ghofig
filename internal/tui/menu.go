package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// MenuModel represents the main menu view.
type MenuModel struct {
	width  int
	height int
}

// NewMenuModel creates a new menu model.
func NewMenuModel() MenuModel {
	return MenuModel{}
}

// SetSize updates the menu dimensions.
func (m MenuModel) SetSize(width, height int) MenuModel {
	m.width = width
	m.height = height
	return m
}

// Update handles menu updates.
func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	return m, nil
}

// View renders the menu.
func (m MenuModel) View() string {
	return "Menu (placeholder)\n\nPress 'enter' to go to Configs\nPress 'q' to quit"
}
