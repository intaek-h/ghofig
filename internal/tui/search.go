package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/intaek-h/ghofig/internal/db"
	"github.com/intaek-h/ghofig/internal/model"
)

var (
	searchTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("170")).
				MarginLeft(2).
				MarginBottom(1)

	searchInputStyle = lipgloss.NewStyle().
				MarginLeft(2).
				MarginBottom(1)

	resultTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212"))

	resultDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	selectedResultStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("236")).
				Padding(0, 1)

	normalResultStyle = lipgloss.NewStyle().
				Padding(0, 1)

	searchHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2).
			MarginTop(1)
)

// SearchModel represents the search view.
type SearchModel struct {
	input      textinput.Model
	results    []model.Config
	cursor     int
	width      int
	height     int
	inputFocus bool
	err        error
}

// NewSearchModel creates a new search model.
func NewSearchModel() SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Type to search configs..."
	ti.CharLimit = 100
	ti.Width = 50

	return SearchModel{
		input:      ti,
		inputFocus: true,
	}
}

// SetSize updates the search dimensions.
func (m SearchModel) SetSize(width, height int) SearchModel {
	m.width = width
	m.height = height
	m.input.Width = width - 6
	return m
}

// Init initializes the search view.
func (m SearchModel) Init() tea.Cmd {
	// Load initial results (all configs)
	return tea.Batch(
		textinput.Blink,
		m.doSearch(""),
	)
}

// searchResultMsg carries search results.
type searchResultMsg struct {
	results []model.Config
	err     error
}

// doSearch returns a command that searches the database.
func (m SearchModel) doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := db.Search(query)
		return searchResultMsg{results: results, err: err}
	}
}

// Update handles search updates.
func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case searchResultMsg:
		m.results = msg.results
		m.err = msg.err
		m.cursor = 0
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Toggle focus between input and results
			m.inputFocus = !m.inputFocus
			if m.inputFocus {
				m.input.Focus()
			} else {
				m.input.Blur()
			}
			return m, nil

		case "up", "k":
			if !m.inputFocus && m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if !m.inputFocus && m.cursor < len(m.results)-1 {
				m.cursor++
			}
			return m, nil

		case "enter":
			if m.inputFocus {
				// Switch to results navigation
				m.inputFocus = false
				m.input.Blur()
				return m, nil
			}
			// Selection handled by parent
			return m, nil
		}
	}

	// Update text input
	if m.inputFocus {
		var cmd tea.Cmd
		prevValue := m.input.Value()
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)

		// If input changed, trigger search
		if m.input.Value() != prevValue {
			cmds = append(cmds, m.doSearch(m.input.Value()))
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the search view.
func (m SearchModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(searchTitleStyle.Render("Search Configs"))
	b.WriteString("\n")

	// Input
	b.WriteString(searchInputStyle.Render(m.input.View()))
	b.WriteString("\n")

	// Error
	if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	}

	// Results count
	countStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(2)
	b.WriteString(countStyle.Render(fmt.Sprintf("%d results", len(m.results))))
	b.WriteString("\n\n")

	// Results list
	maxVisible := m.height - 10 // Leave room for header and footer
	if maxVisible < 1 {
		maxVisible = 5
	}

	// Calculate visible window
	start := 0
	if m.cursor >= maxVisible {
		start = m.cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(m.results) {
		end = len(m.results)
	}

	for i := start; i < end; i++ {
		result := m.results[i]

		// Truncate description to first line or 60 chars
		desc := result.Description
		if idx := strings.Index(desc, "\n"); idx != -1 {
			desc = desc[:idx]
		}
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}

		title := resultTitleStyle.Render(result.Title)
		descStr := resultDescStyle.Render(desc)
		line := fmt.Sprintf("%s\n  %s", title, descStr)

		if i == m.cursor && !m.inputFocus {
			b.WriteString(selectedResultStyle.Render("> " + line))
		} else {
			b.WriteString(normalResultStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}

	// Help
	help := "tab: switch focus • ↑/↓: navigate • enter: select • esc: back • q: quit"
	b.WriteString(searchHelpStyle.Render(help))

	return b.String()
}

// IsInputFocused returns whether the search input is focused.
func (m SearchModel) IsInputFocused() bool {
	return m.inputFocus
}

// SelectedConfig returns the currently selected config.
func (m SearchModel) SelectedConfig() *model.Config {
	if m.cursor >= 0 && m.cursor < len(m.results) {
		return &m.results[m.cursor]
	}
	return nil
}
