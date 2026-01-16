package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/intaek-h/ghofig/internal/model"
)

var (
	detailTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("170")).
				Background(lipgloss.Color("236")).
				Padding(0, 2).
				MarginBottom(1)

	detailContentStyle = lipgloss.NewStyle().
				Padding(1, 2)

	detailHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1).
			MarginLeft(2)

	detailBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(0, 1)
)

// DetailModel represents the config detail view.
type DetailModel struct {
	config   *model.Config
	viewport viewport.Model
	width    int
	height   int
	ready    bool
}

// NewDetailModel creates a new detail model.
func NewDetailModel() DetailModel {
	return DetailModel{}
}

// SetSize updates the detail dimensions.
func (m DetailModel) SetSize(width, height int) DetailModel {
	m.width = width
	m.height = height

	// Calculate viewport size (leaving room for title and help)
	vpWidth := width - 6
	vpHeight := height - 8

	if !m.ready {
		m.viewport = viewport.New(vpWidth, vpHeight)
		m.viewport.YPosition = 0
		m.ready = true
	} else {
		m.viewport.Width = vpWidth
		m.viewport.Height = vpHeight
	}

	// Re-set content if we have a config
	if m.config != nil {
		m.viewport.SetContent(m.config.Description)
	}

	return m
}

// SetConfig sets the config to display.
func (m DetailModel) SetConfig(config *model.Config) DetailModel {
	m.config = config
	if m.ready && config != nil {
		m.viewport.SetContent(config.Description)
		m.viewport.GotoTop()
	}
	return m
}

// Update handles detail updates.
func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		case "pgup":
			m.viewport.HalfViewUp()
		case "pgdown":
			m.viewport.HalfViewDown()
		case "home", "g":
			m.viewport.GotoTop()
		case "end", "G":
			m.viewport.GotoBottom()
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the detail view.
func (m DetailModel) View() string {
	if m.config == nil {
		return "No config selected"
	}

	var b strings.Builder

	// Title
	b.WriteString(detailTitleStyle.Render(m.config.Title))
	b.WriteString("\n\n")

	// Viewport with content
	if m.ready {
		content := detailBorderStyle.Render(m.viewport.View())
		b.WriteString(content)
	} else {
		b.WriteString(detailContentStyle.Render(m.config.Description))
	}

	// Scroll indicator
	scrollInfo := ""
	if m.ready {
		scrollPercent := m.viewport.ScrollPercent() * 100
		scrollInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2).
			Render(strings.Repeat(" ", m.width-20) +
				lipgloss.NewStyle().Render(
					func() string {
						if scrollPercent < 100 {
							return "↓ scroll for more"
						}
						return ""
					}(),
				))
	}
	b.WriteString(scrollInfo)
	b.WriteString("\n")

	// Help
	help := "↑/↓: scroll • pgup/pgdn: page • home/end: top/bottom • esc/backspace: back • q: quit"
	b.WriteString(detailHelpStyle.Render(help))

	return b.String()
}
