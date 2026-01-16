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
				MarginLeft(2)

	searchInputContainerStyle = lipgloss.NewStyle().
					MarginLeft(2).
					MarginTop(1).
					MarginBottom(1)

	inputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("170")).
				Padding(0, 1)

	inputBlurredStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(0, 1)

	resultTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212"))

	resultDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	highlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Background(lipgloss.Color("236"))

	selectedResultStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("236")).
				Padding(0, 1)

	normalResultStyle = lipgloss.NewStyle().
				Padding(0, 1)

	searchHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2).
			MarginTop(1)

	countStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2)
)

// SearchModel represents the search view.
type SearchModel struct {
	input       textinput.Model
	results     []model.Config
	cursor      int
	width       int
	height      int
	inputFocus  bool
	query       string // Store query for highlighting
	err         error
	initialized bool
}

// NewSearchModel creates a new search model.
func NewSearchModel() SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Type to search configs..."
	ti.CharLimit = 100
	ti.Width = 50
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	return SearchModel{
		input:      ti,
		inputFocus: false, // Start with results focused so user can navigate immediately
	}
}

// SetSize updates the search dimensions.
func (m SearchModel) SetSize(width, height int) SearchModel {
	m.width = width
	m.height = height
	m.input.Width = width - 10
	return m
}

// Init initializes the search view.
func (m SearchModel) Init() tea.Cmd {
	m.initialized = true
	// Load initial results (all configs)
	return m.doSearch("")
}

// searchResultMsg carries search results.
type searchResultMsg struct {
	results []model.Config
	query   string
	err     error
}

// doSearch returns a command that searches the database.
func (m SearchModel) doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := db.Search(query)
		return searchResultMsg{results: results, query: query, err: err}
	}
}

// Update handles search updates.
func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case searchResultMsg:
		m.results = msg.results
		m.query = msg.query
		m.err = msg.err
		m.cursor = 0
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			// Focus input for searching
			if !m.inputFocus {
				m.inputFocus = true
				m.input.Focus()
				return m, textinput.Blink
			}

		case "tab":
			// Toggle focus between input and results
			m.inputFocus = !m.inputFocus
			if m.inputFocus {
				m.input.Focus()
				return m, textinput.Blink
			} else {
				m.input.Blur()
			}
			return m, nil

		case "esc":
			// If input focused, unfocus it first
			if m.inputFocus {
				m.inputFocus = false
				m.input.Blur()
				return m, nil
			}
			// Otherwise, let parent handle (go back)
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

	// Update text input only if focused
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

// highlightMatches highlights the query string within text.
func highlightMatches(text, query string) string {
	if query == "" {
		return text
	}

	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	idx := strings.Index(lowerText, lowerQuery)
	if idx == -1 {
		return text
	}

	// Build highlighted string
	var result strings.Builder
	lastEnd := 0

	for idx != -1 {
		// Add text before match
		result.WriteString(text[lastEnd:idx])
		// Add highlighted match (preserve original case)
		result.WriteString(highlightStyle.Render(text[idx : idx+len(query)]))
		lastEnd = idx + len(query)

		// Find next match
		nextIdx := strings.Index(lowerText[lastEnd:], lowerQuery)
		if nextIdx == -1 {
			idx = -1
		} else {
			idx = lastEnd + nextIdx
		}
	}

	// Add remaining text
	result.WriteString(text[lastEnd:])

	return result.String()
}

// View renders the search view.
func (m SearchModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(searchTitleStyle.Render("Search Configs"))
	b.WriteString("\n")

	// Input box with border
	var inputBox string
	inputView := m.input.View()
	if m.inputFocus {
		inputBox = inputFocusedStyle.Render(inputView)
	} else {
		inputBox = inputBlurredStyle.Render(inputView)
	}
	b.WriteString(searchInputContainerStyle.Render(inputBox))
	b.WriteString("\n")

	// Error
	if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	}

	// Results count and current query
	queryInfo := ""
	if m.query != "" {
		queryInfo = fmt.Sprintf(" for \"%s\"", m.query)
	}
	b.WriteString(countStyle.Render(fmt.Sprintf("%d results%s", len(m.results), queryInfo)))
	b.WriteString("\n\n")

	// Results list
	maxVisible := m.height - 12 // Leave room for header and footer
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

		// Highlight matches in title and description
		title := highlightMatches(result.Title, m.query)
		if m.query == "" {
			title = resultTitleStyle.Render(result.Title)
		}
		descHighlighted := highlightMatches(desc, m.query)
		if m.query == "" {
			descHighlighted = resultDescStyle.Render(desc)
		} else {
			descHighlighted = resultDescStyle.Render(descHighlighted)
		}

		line := fmt.Sprintf("%s\n   %s", title, descHighlighted)

		isSelected := i == m.cursor && !m.inputFocus
		if isSelected {
			// Add selection indicator
			b.WriteString(selectedResultStyle.Render("▶ " + line))
		} else {
			b.WriteString(normalResultStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}

	// Help
	var help string
	if m.inputFocus {
		help = "type to search • enter/tab: navigate results • esc: cancel"
	} else {
		help = "/: search • ↑/↓: navigate • enter: view detail • esc: back • q: quit"
	}
	b.WriteString(searchHelpStyle.Render(help))

	return b.String()
}

// IsInputFocused returns whether the search input is focused.
func (m SearchModel) IsInputFocused() bool {
	return m.inputFocus
}

// InputValue returns the current input value.
func (m SearchModel) InputValue() string {
	return m.input.Value()
}

// SelectedConfig returns the currently selected config.
func (m SearchModel) SelectedConfig() *model.Config {
	if m.cursor >= 0 && m.cursor < len(m.results) {
		return &m.results[m.cursor]
	}
	return nil
}
