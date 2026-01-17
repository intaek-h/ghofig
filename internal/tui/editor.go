package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/intaek-h/ghofig/internal/config"
)

var (
	editorTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ThemePrimary)

	editorHelpStyle = lipgloss.NewStyle().
			Foreground(ThemeTextMuted)

	editorSuccessStyle = lipgloss.NewStyle().
				Foreground(ThemeSuccess)

	editorErrorStyle = lipgloss.NewStyle().
				Foreground(ThemeError)

	editorPathStyle = lipgloss.NewStyle().
			Foreground(ThemeTextMuted).
			Italic(true)
)

// EditorModel represents the config file editor view.
type EditorModel struct {
	textarea    textarea.Model
	width       int
	height      int
	configPath  string
	message     string
	isError     bool
	initialText string // Track initial content to detect changes
}

// NewEditorModel creates a new editor model.
func NewEditorModel() EditorModel {
	ta := textarea.New()
	ta.Placeholder = "Config file content..."
	ta.ShowLineNumbers = true
	ta.Focus()

	// Set styles
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().
		Background(ThemeBgHighlight).
		Foreground(ThemeText)
	ta.FocusedStyle.LineNumber = lipgloss.NewStyle().Foreground(ThemeTextMuted)
	ta.FocusedStyle.CursorLineNumber = lipgloss.NewStyle().Foreground(ThemePrimary)

	return EditorModel{
		textarea: ta,
	}
}

// SetSize updates the editor dimensions.
func (m EditorModel) SetSize(width, height int) EditorModel {
	m.width = width
	m.height = height

	// Leave room for title (2 lines), path (1 line), message (1 line), help (2 lines), padding
	taWidth := width - 4
	taHeight := height - 8

	m.textarea.SetWidth(taWidth)
	m.textarea.SetHeight(taHeight)

	return m
}

// configLoadedMsg is sent when config file is loaded
type configLoadedMsg struct {
	content string
	path    string
	err     error
}

// configSavedMsg is sent when config file is saved
type configSavedMsg struct {
	err error
}

// Init initializes the editor by loading the config file.
func (m EditorModel) Init() tea.Cmd {
	return func() tea.Msg {
		path, err := config.GetConfigPath()
		if err != nil {
			return configLoadedMsg{err: err}
		}

		content, err := config.ReadFile()
		if err != nil {
			return configLoadedMsg{path: path, err: err}
		}

		return configLoadedMsg{content: content, path: path}
	}
}

// Update handles editor updates.
func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case configLoadedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error loading config: %v", msg.err)
			m.isError = true
		} else {
			m.textarea.SetValue(msg.content)
			m.initialText = msg.content
			m.configPath = msg.path
			m.message = ""
			m.isError = false
		}
		return m, nil

	case configSavedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error saving: %v", msg.err)
			m.isError = true
		} else {
			m.message = "Saved successfully"
			m.isError = false
			m.initialText = m.textarea.Value() // Update initial text after save
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			// Save the config file
			content := m.textarea.Value()
			return m, func() tea.Msg {
				err := config.WriteFile(content)
				return configSavedMsg{err: err}
			}
		case "esc":
			// Esc is handled by app.go for navigation
			return m, nil
		}
	}

	// Forward other messages to textarea
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

// View renders the editor view.
func (m EditorModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(editorTitleStyle.Render("Config Editor"))
	b.WriteString("\n")

	// Config path
	if m.configPath != "" {
		b.WriteString(editorPathStyle.Render(m.configPath))
	}
	b.WriteString("\n\n")

	// Textarea
	b.WriteString(m.textarea.View())
	b.WriteString("\n")

	// Message (success or error)
	if m.message != "" {
		if m.isError {
			b.WriteString(editorErrorStyle.Render(m.message))
		} else {
			b.WriteString(editorSuccessStyle.Render(m.message))
		}
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	// Help bar
	help := editorHelpStyle.Render("Ctrl+S: save | Esc: back to menu")
	b.WriteString(help)

	return b.String()
}

// HasUnsavedChanges returns true if there are unsaved changes.
func (m EditorModel) HasUnsavedChanges() bool {
	return m.textarea.Value() != m.initialText
}
