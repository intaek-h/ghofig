package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/intaek-h/ghofig/internal/config"
	"github.com/intaek-h/ghofig/internal/model"
)

var (
	detailTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ThemePrimary)

	detailContentStyle = lipgloss.NewStyle().
				Padding(PaddingY, PaddingX)

	detailHelpStyle = lipgloss.NewStyle().
			Foreground(ThemeTextMuted).
			MarginTop(MarginY)

	detailViewportStyle = lipgloss.NewStyle().
				MarginLeft(MarginX)

	detailEditorItemStyle = lipgloss.NewStyle().
				Foreground(ThemePrimary)

	detailEditorHintStyle = lipgloss.NewStyle().
				Foreground(ThemeTextMuted)

	detailSuccessStyle = lipgloss.NewStyle().
				Foreground(ThemeSuccess)
)

// DetailModel represents the config detail view.
type DetailModel struct {
	config           *model.Config
	viewport         viewport.Model
	input            textinput.Model
	width            int
	height           int
	ready            bool
	editing          bool   // true when input is active
	success          bool   // true after successful append
	message          string // success/error message
	hasExistingValue bool   // true if editing an existing config value
}

// NewDetailModel creates a new detail model.
func NewDetailModel() DetailModel {
	ti := textinput.New()
	ti.CharLimit = 500
	ti.Width = 60
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(ThemeSecondary)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(ThemeTextMuted)

	return DetailModel{
		input: ti,
	}
}

// SetSize updates the detail dimensions.
func (m DetailModel) SetSize(width, height int) DetailModel {
	m.width = width
	m.height = height

	// Calculate viewport size (leaving room for title, editor, and help)
	vpWidth := width - 6
	vpHeight := height - 12 // More room for editor section

	if !m.ready {
		m.viewport = viewport.New(vpWidth, vpHeight)
		m.viewport.YPosition = 0
		m.ready = true
	} else {
		m.viewport.Width = vpWidth
		m.viewport.Height = vpHeight
	}

	// Update input width
	m.input.Width = min(width-10, 60)

	// Re-set content if we have a config
	if m.config != nil {
		m.viewport.SetContent(m.config.Description)
	}

	return m
}

// SetConfig sets the config to display.
func (m DetailModel) SetConfig(cfg *model.Config) DetailModel {
	m.config = cfg
	m.editing = false
	m.success = false
	m.message = ""

	if cfg != nil {
		// Set up input with option prefix
		m.input.SetValue("")
		m.input.Placeholder = fmt.Sprintf("%s = value", cfg.Title)
	}

	if m.ready && cfg != nil {
		m.viewport.SetContent(cfg.Description)
		m.viewport.GotoTop()
	}
	return m
}

// configAppendedMsg is sent when config is successfully appended
type configAppendedMsg struct {
	success bool
	err     error
}

// configCommentedOutMsg is sent when config option is commented out
type configCommentedOutMsg struct {
	commented bool
	err       error
}

// Update handles detail updates.
func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case configAppendedMsg:
		if msg.success {
			m.success = true
			m.message = "✓ Added to config file"
			m.editing = false
		} else {
			m.message = fmt.Sprintf("Error: %v", msg.err)
		}
		return m, nil

	case configCommentedOutMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error: %v", msg.err)
		} else if msg.commented {
			m.success = true
			m.message = "✓ Commented out from config file"
			m.editing = false
		} else {
			m.message = "Option not found in config file"
			m.editing = false
		}
		return m, nil

	case tea.KeyMsg:
		if m.editing {
			// In editing mode
			switch msg.String() {
			case "enter":
				value := m.input.Value()
				optionName := m.config.Title

				// Check if user wants to comment out the option
				// This happens when input is empty or just "option = " with no value
				trimmedValue := strings.TrimSpace(value)
				isEmptyValue := trimmedValue == "" ||
					trimmedValue == optionName+" =" ||
					trimmedValue == optionName+"="

				if isEmptyValue {
					// Comment out the option
					return m, func() tea.Msg {
						commented, err := config.CommentOut(optionName)
						return configCommentedOutMsg{commented: commented, err: err}
					}
				}

				// Append to config file
				return m, func() tea.Msg {
					err := config.AppendLine(value)
					return configAppendedMsg{success: err == nil, err: err}
				}
			case "esc":
				// Cancel editing
				m.editing = false
				m.input.Blur()
				return m, nil
			}
			// Forward to text input
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		// Not editing - normal navigation
		switch msg.String() {
		case "enter":
			// Start editing
			m.editing = true
			m.success = false
			m.message = ""
			// Check if value already exists in config
			existingValue := config.GetValue(m.config.Title)
			m.hasExistingValue = existingValue != ""
			if m.hasExistingValue {
				m.input.SetValue(fmt.Sprintf("%s = %s", m.config.Title, existingValue))
			} else {
				m.input.SetValue(fmt.Sprintf("%s = ", m.config.Title))
			}
			m.input.CursorEnd()
			m.input.Focus()
			return m, textinput.Blink
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

	// Editor section
	if m.editing {
		// Show input mode
		b.WriteString("  > ")
		b.WriteString(m.input.View())
		b.WriteString("\n")
		hint := "  Enter: save to config, Esc: cancel"
		if m.hasExistingValue {
			hint += " | Clear value to comment out"
		}
		b.WriteString(detailEditorHintStyle.Render(hint))
		b.WriteString("\n\n")
	} else if m.success {
		// Show success message
		b.WriteString(detailSuccessStyle.Render("  " + m.message))
		b.WriteString("\n\n")
	} else {
		// Show editor item
		b.WriteString(detailEditorItemStyle.Render(fmt.Sprintf("  ➤ ○ Open Editor For `%s`", m.config.Title)))
		b.WriteString("\n\n")
	}

	// Viewport with content
	if m.ready {
		content := detailViewportStyle.Render(m.viewport.View())
		b.WriteString(content)
	} else {
		b.WriteString(detailContentStyle.Render(m.config.Description))
	}

	// Scroll indicator
	scrollInfo := ""
	if m.ready {
		scrollPercent := m.viewport.ScrollPercent() * 100
		scrollInfo = lipgloss.NewStyle().
			Foreground(ThemeTextMuted).
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
	var help string
	if m.editing {
		help = "enter: save • esc: cancel"
	} else {
		help = "enter: edit • ↑/↓: scroll • pgup/pgdn: page • esc: back • q: quit"
	}
	b.WriteString(detailHelpStyle.Render(help))

	return b.String()
}

// IsEditing returns whether the detail view is in editing mode
func (m DetailModel) IsEditing() bool {
	return m.editing
}
