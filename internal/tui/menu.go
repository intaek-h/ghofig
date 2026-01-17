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
			Foreground(ThemePrimary)

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
		MenuItem{title: "Browse Options", description: "Search Ghostty configuration options"},
		MenuItem{title: "Config Editor ", description: "Edit your Ghostty config file directly"},
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
	// Logo: 8 lines + 2 spacing + 1 help = 11 lines
	logoHeight := 11
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

	// ASCII art logo (Mole CLI style - compact with inline info)
	// Split rendering: logo in ThemePrimary, github link in ThemeTextMuted
	logoTop := ` _____ _           __ _
|  __ | |         / _(_)
| |  \| |__   ___| |_  _  __ _
| | __| '_ \ / _ \   _| |/ _` + "`" + ` |  `
	githubLink := `github.com/intaek-h/ghofig`
	logoMid := `
| |_\ \ | | | (_) | | | | (_| |  `
	description := `Browse and manage Ghostty config.`
	logoBottom := `
 \____/_| |_|\___/|_| |_|\__, |
                          __/ |
                         |___/ `

	b.WriteString(menuLogoStyle.Render(logoTop))
	b.WriteString(menuDescStyle.Render(githubLink))
	b.WriteString(menuLogoStyle.Render(logoMid))
	b.WriteString(menuLogoStyle.Render(description))
	b.WriteString(menuLogoStyle.Render(logoBottom))
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
