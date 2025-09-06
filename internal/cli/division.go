package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DivisionSelectMsg is sent when a division is selected
type DivisionSelectMsg struct {
	Division string
	Filename string
}

// BackToMenuMsg is sent when user wants to go back to main menu
type BackToMenuMsg struct{}

// DivisionModel represents the division selection TUI state
type DivisionModel struct {
	divisions []string
	filenames []string
	cursor    int
	style     lipgloss.Style
}

// NewDivisionModel creates a new division selection model
func NewDivisionModel() *DivisionModel {
	return &DivisionModel{
		divisions: []string{
			"Elite",
			"Platinum A",
			"Platinum B",
			"Oro A",
			"Oro B",
			"Oro C",
			"Oro D",
		},
		filenames: []string{
			"data/Liga Argentina - 1° Temporada - E-Fixture.csv",
			"data/Liga Argentina - 1° Temporada - P.A-Fixture.csv",
			"data/Liga Argentina - 1° Temporada - P.B-Fixture.csv",
			"data/Liga Argentina - 1° Temporada - O.A-Fixture.csv",
			"data/Liga Argentina - 1° Temporada - O.B-Fixture.csv",
			"data/Liga Argentina - 1° Temporada - O.C-Fixture.csv",
			"data/Liga Argentina - 1° Temporada - O.D-Fixture.csv",
		},
		cursor: 0,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true),
	}
}

// Init initializes the division model (required by Bubble Tea)
func (m *DivisionModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model state
func (m *DivisionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyUp:
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.divisions) - 1
			}
		case tea.KeyDown:
			m.cursor++
			if m.cursor >= len(m.divisions) {
				m.cursor = 0
			}
		case tea.KeyEnter:
			// Send message to switch to fixture view
			return m, func() tea.Msg {
				return DivisionSelectMsg{
					Division: m.GetSelectedDivision(),
					Filename: m.GetSelectedFilename(),
				}
			}
		case tea.KeyEsc:
			// Go back to main menu
			return m, func() tea.Msg {
				return BackToMenuMsg{}
			}
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				// Go back to main menu
				return m, func() tea.Msg {
					return BackToMenuMsg{}
				}
			case "j":
				// Vim down
				m.cursor++
				if m.cursor >= len(m.divisions) {
					m.cursor = 0
				}
			case "k":
				// Vim up
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.divisions) - 1
				}
			}
		}
	}
	return m, nil
}

// View renders the current state of the division selection
func (m *DivisionModel) View() string {
	title := m.style.Render("Select Division")
	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Render("Choose a division to view fixtures:")

	s := fmt.Sprintf("\n%s\n\n%s\n\n", title, subtitle)

	for i, division := range m.divisions {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			division = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true).
				Render(division)
		}
		s += fmt.Sprintf("%s %s\n", cursor, division)
	}

	s += "\n\nPress enter to select, esc/q to go back, ↑/↓ or j/k to navigate.\n"
	return s
}

// GetSelectedDivision returns the currently selected division
func (m *DivisionModel) GetSelectedDivision() string {
	if m.cursor >= 0 && m.cursor < len(m.divisions) {
		return m.divisions[m.cursor]
	}
	return ""
}

// GetSelectedFilename returns the filename for the currently selected division
func (m *DivisionModel) GetSelectedFilename() string {
	if m.cursor >= 0 && m.cursor < len(m.filenames) {
		return m.filenames[m.cursor]
	}
	return ""
}
