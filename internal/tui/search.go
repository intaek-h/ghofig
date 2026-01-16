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

	searchInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("170")).
				Padding(0, 1).
				MarginLeft(2).
				MarginTop(1)

	searchCountStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				MarginLeft(2).
				MarginTop(1)

	searchItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	searchSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))

	searchDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			PaddingLeft(6)

	searchSelectedDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				PaddingLeft(4)

	searchMatchStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("212")).
				Bold(true)

	searchHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2)
)

// SearchModel represents the search view.
type SearchModel struct {
	input   textinput.Model
	results []model.Config
	cursor  int
	query   string
	width   int
	height  int
	err     error
}

// NewSearchModel creates a new search model.
func NewSearchModel() SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search configs..."
	ti.CharLimit = 100
	ti.Width = 40
	ti.Prompt = " "
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	ti.Focus()

	return SearchModel{
		input:  ti,
		cursor: 0,
	}
}

// SetSize updates dimensions.
func (m SearchModel) SetSize(width, height int) SearchModel {
	m.width = width
	m.height = height
	m.input.Width = min(width-10, 60)
	return m
}

// Init initializes the search view.
func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

// searchResultMsg carries search results.
type searchResultMsg struct {
	results []model.Config
	query   string
	err     error
}

// doSearch returns a command that searches the database.
func doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := db.Search(query)
		return searchResultMsg{results: results, query: query, err: err}
	}
}

// Update handles updates.
func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case searchResultMsg:
		m.results = msg.results
		m.query = msg.query
		m.err = msg.err
		m.cursor = 0
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		if key == "up" || key == "down" {
			if len(m.results) > 0 {
				m.input.Blur()
				if key == "up" && m.cursor > 0 {
					m.cursor--
				} else if key == "down" && m.cursor < len(m.results)-1 {
					m.cursor++
				}
			}
			return m, nil
		}

		if key == "enter" || key == "esc" {
			return m, nil
		}

		if !m.input.Focused() {
			m.input.Focus()
		}
	}

	prevValue := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if m.input.Value() != prevValue {
		m.cursor = 0
		newQuery := m.input.Value()
		if newQuery == "" {
			m.results = nil
			m.query = ""
			return m, cmd
		}
		return m, tea.Batch(cmd, doSearch(newQuery))
	}

	return m, cmd
}

// highlight highlights query matches in text.
func highlight(text, query string) string {
	if query == "" {
		return text
	}

	lower := strings.ToLower(text)
	lowerQ := strings.ToLower(query)
	idx := strings.Index(lower, lowerQ)
	if idx == -1 {
		return text
	}

	var b strings.Builder
	last := 0
	for idx != -1 {
		b.WriteString(text[last:idx])
		b.WriteString(searchMatchStyle.Render(text[idx : idx+len(query)]))
		last = idx + len(query)
		next := strings.Index(lower[last:], lowerQ)
		if next == -1 {
			break
		}
		idx = last + next
	}
	b.WriteString(text[last:])
	return b.String()
}

// View renders the search view.
func (m SearchModel) View() string {
	// Build fixed header: title + input
	header := lipgloss.JoinVertical(lipgloss.Left,
		searchTitleStyle.Render("Search Configs"),
		searchInputStyle.Render(m.input.View()),
	)

	// Build help footer
	var helpText string
	if m.query == "" {
		helpText = "type to search • esc: back • q: quit"
	} else {
		helpText = "↑/↓: navigate • enter: select • esc: back • q: quit"
	}
	footer := searchHelpStyle.Render(helpText)

	// Calculate heights
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	resultsHeight := m.height - headerHeight - footerHeight - 2 // 2 for margins

	if resultsHeight < 3 {
		resultsHeight = 3
	}

	// Build results section
	var resultsContent string
	if m.query != "" && len(m.results) > 0 {
		var lines []string

		// Count line
		countLine := searchCountStyle.Render(fmt.Sprintf("%d results", len(m.results)))
		lines = append(lines, countLine)

		// Calculate how many items we can show (each item = 2 lines)
		availableForItems := resultsHeight - 2 // subtract count line and some padding
		maxItems := availableForItems / 2
		if maxItems < 1 {
			maxItems = 1
		}

		// Calculate window
		start := 0
		if m.cursor >= maxItems {
			start = m.cursor - maxItems + 1
		}
		end := start + maxItems
		if end > len(m.results) {
			end = len(m.results)
		}

		// Render items
		for i := start; i < end; i++ {
			r := m.results[i]

			desc := r.Description
			if idx := strings.Index(desc, "\n"); idx != -1 {
				desc = desc[:idx]
			}
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}

			title := highlight(r.Title, m.query)
			desc = highlight(desc, m.query)

			if i == m.cursor {
				lines = append(lines, searchSelectedStyle.Render("▶ "+title))
				lines = append(lines, searchSelectedDescStyle.Render(desc))
			} else {
				lines = append(lines, searchItemStyle.Render(title))
				lines = append(lines, searchDescStyle.Render(desc))
			}
		}

		resultsContent = strings.Join(lines, "\n")
	} else if m.query != "" {
		resultsContent = searchCountStyle.Render("0 results")
	}

	// Constrain results to fixed height
	resultsSection := lipgloss.NewStyle().
		Height(resultsHeight).
		MaxHeight(resultsHeight).
		Render(resultsContent)

	// Join all sections vertically
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		resultsSection,
		footer,
	)
}

// IsInputFocused returns whether input is focused.
func (m SearchModel) IsInputFocused() bool {
	return m.input.Focused()
}

// InputValue returns the current input value.
func (m SearchModel) InputValue() string {
	return m.input.Value()
}

// SelectedConfig returns the selected config.
func (m SearchModel) SelectedConfig() *model.Config {
	if len(m.results) > 0 && m.cursor >= 0 && m.cursor < len(m.results) {
		return &m.results[m.cursor]
	}
	return nil
}

// HasResults returns whether there are search results.
func (m SearchModel) HasResults() bool {
	return len(m.results) > 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
