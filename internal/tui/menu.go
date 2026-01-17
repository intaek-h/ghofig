package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	menuLogoStyle = lipgloss.NewStyle().
			Foreground(ThemeTextInput)

	menuItemStyle = lipgloss.NewStyle()

	menuSelectedStyle = lipgloss.NewStyle().
				Foreground(ThemePrimary)

	menuDescStyle = lipgloss.NewStyle().
			Foreground(ThemeTextMuted)

	menuHelpStyle = lipgloss.NewStyle().
			Foreground(ThemeTextMuted)
)

// MenuItem represents a menu item.
type MenuItem struct {
	title       string
	description string
}

func (i MenuItem) FilterValue() string { return i.title }

// MenuItemDelegate handles rendering of menu items.
type MenuItemDelegate struct{}

func (d MenuItemDelegate) Height() int                             { return 1 }
func (d MenuItemDelegate) Spacing() int                            { return 0 }
func (d MenuItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d MenuItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(MenuItem)
	if !ok {
		return
	}

	num := fmt.Sprintf("%d.", index+1)

	if index == m.Index() {
		// Selected: ➤ 1. Title        Description
		titlePart := menuSelectedStyle.Render(fmt.Sprintf("\u27a4 %s %s", num, i.title))
		descPart := menuDescStyle.Render(i.description)
		fmt.Fprintf(w, "%s        %s", titlePart, descPart)
	} else {
		// Unselected:   1. Title        Description
		titlePart := menuItemStyle.Render(fmt.Sprintf("  %s %s", num, i.title))
		descPart := menuDescStyle.Render(i.description)
		fmt.Fprintf(w, "%s        %s", titlePart, descPart)
	}
}

// MenuModel represents the main menu view.
type MenuModel struct {
	list   list.Model
	width  int
	height int
}

// Menu item indices for selection handling
const (
	MenuItemConfigOptions = iota
	MenuItemConfigEditor
)

// NewMenuModel creates a new menu model.
func NewMenuModel() MenuModel {
	items := []list.Item{
		MenuItem{title: "Config Options", description: "Search Ghostty configuration options"},
		MenuItem{title: "Config Editor", description: "Edit your Ghostty config file directly"},
	}

	l := list.New(items, MenuItemDelegate{}, 0, 0)
	l.Title = "" // We render custom ASCII title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return MenuModel{
		list: l,
	}
}

// SetSize updates the menu dimensions.
func (m MenuModel) SetSize(width, height int) MenuModel {
	m.width = width
	m.height = height
	m.list.SetWidth(width)
	// Logo: 8 lines + 1 subtitle + 2 spacing + 1 help = 12 lines
	logoHeight := 12
	listHeight := height - logoHeight
	// Ensure minimum height to show all menu items
	if listHeight < 3 {
		listHeight = 3
	}
	m.list.SetHeight(listHeight)
	return m
}

// Update handles menu updates.
func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.list.CursorUp()
		case "down", "j":
			m.list.CursorDown()
		}
	}
	return m, nil
}

// View renders the menu.
func (m MenuModel) View() string {
	var b strings.Builder

	// ASCII art logo
	logo := `   _____ _               _   _         
  / ____| |             | | | |        
 | |  __| |__   ___  ___| |_| |_ _   _ 
 | | |_ | '_ \ / _ \/ __| __| __| | | |
 | |__| | | | | (_) \__ \ |_| |_| |_| |
  \_____|_| |_|\___/|___/\__|\__|\__, |
                                  __/ |
                                 |___/    `

	b.WriteString(menuLogoStyle.Render(logo + "\nGhofig: Ghostty Config Editor"))
	b.WriteString("\n\n")

	// Menu items
	b.WriteString(m.list.View())
	b.WriteString("\n")

	// Help
	help := menuHelpStyle.Render("↑/↓: navigate • enter: select • q: quit")
	b.WriteString(help)

	return b.String()
}

// SelectedItem returns the currently selected menu item.
func (m MenuModel) SelectedItem() MenuItem {
	if item, ok := m.list.SelectedItem().(MenuItem); ok {
		return item
	}
	return MenuItem{}
}
