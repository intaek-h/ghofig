package tui

import "github.com/charmbracelet/lipgloss"

// Adaptive colors that work on both light and dark terminal backgrounds.
// Format: AdaptiveColor{Light: "color for light bg", Dark: "color for dark bg"}

// Semantic theme colors using AdaptiveColor for automatic light/dark detection
var (
	// Primary - used for main actions, selected items, focus
	ThemePrimary = lipgloss.AdaptiveColor{Light: "#0055CC", Dark: "#0070F3"}
	// Secondary - used for secondary elements
	ThemeSecondary = lipgloss.AdaptiveColor{Light: "#0066DD", Dark: "#52A8FF"}

	// Accent - used for highlights, special elements
	ThemeAccent = lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#8E4EC6"}

	// Text
	ThemeText      = lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#EDEDED"}   // Main text
	ThemeTextMuted = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#878787"}   // Secondary text, descriptions, help
	ThemeTextInput = lipgloss.AdaptiveColor{Light: "#393939ff", Dark: "#1A1A1A"} // Input text color (inverted for contrast)

	// Backgrounds
	ThemeBgDefault   = lipgloss.AdaptiveColor{Light: "#FAFAFA", Dark: "#0A0A0A"}
	ThemeBgSubtle    = lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#1F1F1F"}
	ThemeBgHighlight = lipgloss.AdaptiveColor{Light: "#E5E5E5", Dark: "#292929"}

	// Borders
	ThemeBorder       = lipgloss.AdaptiveColor{Light: "#CCCCCC", Dark: "#454545"}
	ThemeBorderFocus  = lipgloss.AdaptiveColor{Light: "#0055CC", Dark: "#0070F3"}
	ThemeBorderSubtle = lipgloss.AdaptiveColor{Light: "#DDDDDD", Dark: "#2E2E2E"}

	// Search highlight
	ThemeMatch = lipgloss.AdaptiveColor{Light: "#0066DD", Dark: "#52A8FF"}

	// Status
	ThemeError   = lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#E5484D"}
	ThemeSuccess = lipgloss.AdaptiveColor{Light: "#16A34A", Dark: "#46A758"}
	ThemeWarning = lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FFB224"}
)

// Common style constants
const (
	PaddingX = 2
	PaddingY = 1
	MarginX  = 2
	MarginY  = 1
)
