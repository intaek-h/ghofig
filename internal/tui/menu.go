package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginLeft(2)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
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

	str := fmt.Sprintf("%s", i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}

// MenuModel represents the main menu view.
type MenuModel struct {
	list   list.Model
	width  int
	height int
}

// NewMenuModel creates a new menu model.
func NewMenuModel() MenuModel {
	items := []list.Item{
		MenuItem{title: "Configs", description: "Browse Ghostty configuration options"},
	}

	l := list.New(items, MenuItemDelegate{}, 0, 0)
	l.Title = "Ghofig"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle

	return MenuModel{
		list: l,
	}
}

// SetSize updates the menu dimensions.
func (m MenuModel) SetSize(width, height int) MenuModel {
	m.width = width
	m.height = height
	m.list.SetWidth(width)
	m.list.SetHeight(height - 2) // Leave room for help text
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
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginLeft(2)

	help := helpStyle.Render("↑/↓: navigate • enter: select • q: quit")

	return m.list.View() + "\n" + help
}

// SelectedItem returns the currently selected menu item.
func (m MenuModel) SelectedItem() MenuItem {
	if item, ok := m.list.SelectedItem().(MenuItem); ok {
		return item
	}
	return MenuItem{}
}
