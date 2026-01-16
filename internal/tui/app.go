package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// View represents the current view state.
type View int

const (
	MenuView View = iota
	SearchView
	DetailView
)

// KeyMap defines the keybindings for the app.
type KeyMap struct {
	Quit key.Binding
	Back key.Binding
}

// DefaultKeyMap returns the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "back"),
		),
	}
}

// Model is the main application model.
type Model struct {
	currentView    View
	previousView   View
	width          int
	height         int
	keys           KeyMap
	menu           MenuModel
	search         SearchModel
	detail         DetailModel
	selectedConfig int // ID of selected config for detail view
}

// New creates a new application model.
func New() Model {
	return Model{
		currentView: MenuView,
		keys:        DefaultKeyMap(),
		menu:        NewMenuModel(),
		search:      NewSearchModel(),
		detail:      NewDetailModel(),
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menu = m.menu.SetSize(msg.Width, msg.Height)
		m.search = m.search.SetSize(msg.Width, msg.Height)
		m.detail = m.detail.SetSize(msg.Width, msg.Height)
	}

	// Route to current view
	var cmd tea.Cmd
	switch m.currentView {
	case MenuView:
		m, cmd = m.updateMenu(msg)
	case SearchView:
		m, cmd = m.updateSearch(msg)
	case DetailView:
		m, cmd = m.updateDetail(msg)
	}

	return m, cmd
}

// View implements tea.Model.
func (m Model) View() string {
	switch m.currentView {
	case MenuView:
		return m.menu.View()
	case SearchView:
		return m.search.View()
	case DetailView:
		return m.detail.View()
	default:
		return "Unknown view"
	}
}

// updateMenu handles updates for the menu view.
func (m Model) updateMenu(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Menu selection - for now, only "Configs" option
			m.previousView = MenuView
			m.currentView = SearchView
			return m, m.search.Init()
		}
	}

	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

// updateSearch handles updates for the search view.
func (m Model) updateSearch(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back to menu
			m.search = NewSearchModel() // Reset search
			m.search = m.search.SetSize(m.width, m.height)
			m.currentView = MenuView
			return m, nil

		case "enter":
			// Select config from results
			if m.search.HasResults() {
				if selected := m.search.SelectedConfig(); selected != nil {
					m.selectedConfig = selected.ID
					m.detail = m.detail.SetConfig(selected)
					m.previousView = SearchView
					m.currentView = DetailView
					return m, nil
				}
			}
		}
	}

	var cmd tea.Cmd
	m.search, cmd = m.search.Update(msg)
	return m, cmd
}

// updateDetail handles updates for the detail view.
func (m Model) updateDetail(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Back to search
		if key.Matches(msg, m.keys.Back) {
			m.currentView = SearchView
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.detail, cmd = m.detail.Update(msg)
	return m, cmd
}
