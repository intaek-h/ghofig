package tui

import "github.com/charmbracelet/lipgloss"

// Color palette based on dark theme
var (
	// Backgrounds
	colorBackground100 = lipgloss.Color("#0A0A0A")
	colorBackground200 = lipgloss.Color("#000000")

	// Grays
	colorGray100  = lipgloss.Color("#1A1A1A")
	colorGray200  = lipgloss.Color("#1F1F1F")
	colorGray300  = lipgloss.Color("#292929")
	colorGray400  = lipgloss.Color("#2E2E2E")
	colorGray500  = lipgloss.Color("#454545")
	colorGray600  = lipgloss.Color("#878787")
	colorGray700  = lipgloss.Color("#8F8F8F")
	colorGray900  = lipgloss.Color("#A1A1A1")
	colorGray1000 = lipgloss.Color("#EDEDED")

	// Blues
	colorBlue600  = lipgloss.Color("#0099FF")
	colorBlue700  = lipgloss.Color("#0070F3")
	colorBlue900  = lipgloss.Color("#52A8FF")
	colorBlue1000 = lipgloss.Color("#EBF8FF")

	// Accent colors
	colorPurple700 = lipgloss.Color("#8E4EC6")
	colorPurple900 = lipgloss.Color("#BF7AF0")
	colorPink700   = lipgloss.Color("#E93D82")
	colorCyan      = lipgloss.Color("#50E3C2")

	// Status colors
	colorRed700   = lipgloss.Color("#E5484D")
	colorGreen700 = lipgloss.Color("#46A758")
	colorAmber700 = lipgloss.Color("#FFB224")
)

// Semantic theme colors
var (
	// Primary - used for main actions, selected items, focus
	ThemePrimary   = colorBlue700
	ThemeSecondary = colorBlue900

	// Accent - used for highlights, special elements
	ThemeAccent = colorPurple700

	// Text
	ThemeText      = colorGray1000 // Main text
	ThemeTextMuted = colorGray600  // Secondary text, descriptions, help
	ThemeTextInput = colorGray1000 // Input text color

	// Backgrounds
	ThemeBgDefault   = colorBackground100
	ThemeBgSubtle    = colorGray200
	ThemeBgHighlight = colorGray300

	// Borders
	ThemeBorder       = colorGray500
	ThemeBorderFocus  = colorBlue700
	ThemeBorderSubtle = colorGray400

	// Search highlight
	ThemeMatch = colorBlue900

	// Status
	ThemeError   = colorRed700
	ThemeSuccess = colorGreen700
	ThemeWarning = colorAmber700
)

// Common style constants
const (
	PaddingX = 2
	PaddingY = 1
	MarginX  = 2
	MarginY  = 1
)
