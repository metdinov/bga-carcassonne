package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuModel represents the main menu TUI state
type MenuModel struct {
	choices []string
	cursor  int
	style   lipgloss.Style
}

// NewMenuModel creates a new menu model with default choices
func NewMenuModel() *MenuModel {
	return &MenuModel{
		choices: []string{
			"Create Tournament",
			"View Fixture",
			"View Positions",
			"Exit",
		},
		cursor: 0,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true),
	}
}

// Init initializes the menu model (required by Bubble Tea)
func (m *MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model state
func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyUp:
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		case tea.KeyDown:
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}
		case tea.KeyEnter:
			// Handle menu selection
			switch m.cursor {
			case 1: // View Fixture
				return m, func() tea.Msg {
					return ViewFixtureSelectMsg{}
				}
			case 3: // Exit
				return m, tea.Quit
			default:
				// For now, other options don't do anything
				// TODO: Implement create tournament, view positions
				return m, nil
			}
		}
	}
	return m, nil
}

// View renders the current state of the menu
func (m *MenuModel) View() string {
	title := m.style.Render("Carcassonne Tournament Manager")
	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Render("Please select an option:")

	s := fmt.Sprintf("\n%s\n\n%s\n\n", title, subtitle)

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			choice = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true).
				Render(choice)
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n\nPress q to quit, ↑/↓ to navigate, enter to select.\n"
	return s
}

// GetSelectedChoice returns the currently selected menu choice
func (m *MenuModel) GetSelectedChoice() string {
	if m.cursor >= 0 && m.cursor < len(m.choices) {
		return m.choices[m.cursor]
	}
	return ""
}
