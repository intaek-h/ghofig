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
				Foreground(ThemePrimary)

	searchCountStyle = lipgloss.NewStyle().
				Foreground(ThemeTextMuted)

	searchPromptStyle = lipgloss.NewStyle().
				Foreground(ThemeTextMuted)

	searchItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	searchSelectedStyle = lipgloss.NewStyle()

	searchDescStyle = lipgloss.NewStyle().
			Foreground(ThemeTextMuted).
			PaddingLeft(4)

	searchSelectedDescStyle = lipgloss.NewStyle().
				Foreground(ThemeTextMuted).
				PaddingLeft(2)

	searchMatchStyle = lipgloss.NewStyle().
				Foreground(ThemeMatch).
				Bold(true)

	searchSelectedMatchStyle = lipgloss.NewStyle().
					Foreground(ThemeMatch).
					Bold(true)

	searchHelpStyle = lipgloss.NewStyle().
			Foreground(ThemeTextMuted)
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
	ti.Placeholder = "type to search..."
	ti.CharLimit = 100
	ti.Width = 40
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(ThemeTextInput)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(ThemeTextMuted)
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

// highlightWithStyle highlights query matches in text, applying baseStyle to non-match parts.
func highlightWithStyle(text, query string, baseStyle, matchStyle lipgloss.Style) string {
	if query == "" {
		return baseStyle.Render(text)
	}

	lower := strings.ToLower(text)
	lowerQ := strings.ToLower(query)
	idx := strings.Index(lower, lowerQ)
	if idx == -1 {
		return baseStyle.Render(text)
	}

	var b strings.Builder
	last := 0
	for idx != -1 {
		if last < idx {
			b.WriteString(baseStyle.Render(text[last:idx]))
		}
		b.WriteString(matchStyle.Render(text[idx : idx+len(query)]))
		last = idx + len(query)
		next := strings.Index(lower[last:], lowerQ)
		if next == -1 {
			break
		}
		idx = last + next
	}
	if last < len(text) {
		b.WriteString(baseStyle.Render(text[last:]))
	}
	return b.String()
}

// View renders the search view.
func (m SearchModel) View() string {
	// Build title line with count
	var titleLine string
	if m.query != "" {
		current := 0
		if len(m.results) > 0 {
			current = m.cursor + 1
		}
		titleLine = searchTitleStyle.Render("Search") + "  " + searchCountStyle.Render(fmt.Sprintf("%d/%d results", current, len(m.results)))
	} else {
		titleLine = searchTitleStyle.Render("Search")
	}

	// Build input line with prompt
	inputLine := searchPromptStyle.Render("> ") + m.input.View()

	// Combined header with blank line between title and input, and between input and list
	header := titleLine + "\n\n" + inputLine + "\n"

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

		availableForItems := resultsHeight // items are single line now
		maxItems := availableForItems
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

		// Render items (title only, no description)
		for i := start; i < end; i++ {
			r := m.results[i]
			title := r.Title

			if i == m.cursor {
				// Selected item: apply primary color to non-match text
				titleStyled := highlightWithStyle(title, m.query, lipgloss.NewStyle().Foreground(ThemePrimary), searchSelectedMatchStyle)
				lines = append(lines, searchSelectedStyle.Render("\u27a4 \u25cb "+titleStyled))
			} else {
				// Unselected item: no base color, just match highlights
				titleStyled := highlightWithStyle(title, m.query, lipgloss.NewStyle(), searchMatchStyle)
				lines = append(lines, searchItemStyle.Render("\u25cb "+titleStyled))
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
