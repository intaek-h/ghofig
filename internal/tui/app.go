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
	EditorView
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
	editor         EditorModel
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
		editor:      NewEditorModel(),
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
		// Global quit (but not while typing in search or editing in detail)
		if key.Matches(msg, m.keys.Quit) {
			if m.currentView == SearchView && m.search.IsInputFocused() {
				// Don't quit while typing in search, let search handle it
			} else if m.currentView == DetailView && m.detail.IsEditing() {
				// Don't quit while editing, let detail handle it
			} else {
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menu = m.menu.SetSize(msg.Width, msg.Height)
		m.search = m.search.SetSize(msg.Width, msg.Height)
		m.detail = m.detail.SetSize(msg.Width, msg.Height)
		m.editor = m.editor.SetSize(msg.Width, msg.Height)
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
	case EditorView:
		m, cmd = m.updateEditor(msg)
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
	case EditorView:
		return m.editor.View()
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
			// Route based on selected menu item
			selectedIndex := m.menu.list.Index()
			m.previousView = MenuView

			switch selectedIndex {
			case MenuItemConfigOptions:
				m.currentView = SearchView
				return m, m.search.Init()
			case MenuItemConfigEditor:
				m.currentView = EditorView
				m.editor = m.editor.SetSize(m.width, m.height)
				return m, m.editor.Init()
			}
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
		// Back to search (but not while editing)
		if key.Matches(msg, m.keys.Back) && !m.detail.IsEditing() {
			m.currentView = SearchView
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.detail, cmd = m.detail.Update(msg)
	return m, cmd
}

// updateEditor handles updates for the editor view.
func (m Model) updateEditor(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Back to menu on Esc only (not backspace - textarea needs it)
		if msg.String() == "esc" {
			m.editor = NewEditorModel() // Reset editor
			m.currentView = MenuView
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.editor, cmd = m.editor.Update(msg)
	return m, cmd
}
