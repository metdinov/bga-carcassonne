package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"carca-cli/internal/fixtures"
)

// FixtureModel represents the fixture display TUI state
type FixtureModel struct {
	division     *fixtures.Division
	currentRound int
	style        lipgloss.Style
}

// NewFixtureModel creates a new fixture display model
func NewFixtureModel(division *fixtures.Division) *FixtureModel {
	return &FixtureModel{
		division:     division,
		currentRound: 0,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true),
	}
}

// Init initializes the fixture model (required by Bubble Tea)
func (m *FixtureModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model state
func (m *FixtureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyLeft:
			m.currentRound--
			if m.currentRound < 0 {
				m.currentRound = len(m.division.Rounds) - 1
			}
		case tea.KeyRight:
			m.currentRound++
			if m.currentRound >= len(m.division.Rounds) {
				m.currentRound = 0
			}
		case tea.KeyEsc:
			// Go back to division selection
			return m, func() tea.Msg {
				return BackToMenuMsg{}
			}
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				// Go back to division selection
				return m, func() tea.Msg {
					return BackToMenuMsg{}
				}
			}
		}
	}
	return m, nil
}

// View renders the current state of the fixture display
func (m *FixtureModel) View() string {
	if len(m.division.Rounds) == 0 {
		return "No fixtures available for this division.\n\nPress esc/q to go back.\n"
	}

	currentRound := m.GetCurrentRound()
	if currentRound == nil {
		return "Error loading fixture data.\n\nPress esc/q to go back.\n"
	}

	title := m.style.Render(fmt.Sprintf("Division %s - Round %d", m.division.Name, currentRound.Number))
	dateRange := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Render(fmt.Sprintf("Date Range: %s", currentRound.DateRange))

	s := fmt.Sprintf("\n%s\n%s\n\n", title, dateRange)

	// Display matches
	if len(currentRound.Matches) == 0 {
		s += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Render("No matches in this round")
	} else {
		for _, match := range currentRound.Matches {
			s += m.formatMatch(match) + "\n"
		}
	}

	// Navigation info
	s += fmt.Sprintf("\n\nRound %d of %d", m.currentRound+1, len(m.division.Rounds))
	s += "\n\nPress ←/→ to navigate rounds, esc/q to go back.\n"

	return s
}

// formatMatch formats a single match for display
func (m *FixtureModel) formatMatch(match *fixtures.Match) string {
	var status string
	var scoreDisplay string
	var style lipgloss.Style

	if match.Played {
		status = "✓"
		scoreDisplay = fmt.Sprintf("%d - %d", match.HomeScore, match.AwayScore)
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#50C878")) // Green for played
	} else {
		status = "○"
		scoreDisplay = "- - -"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")) // Red for unplayed
	}

	matchInfo := fmt.Sprintf("[%s] %s vs %s (%s)",
		status,
		match.HomePlayer,
		match.AwayPlayer,
		scoreDisplay,
	)

	result := style.Render(matchInfo)

	// Add datetime if available
	if match.DateTime != "" {
		result += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render(fmt.Sprintf(" - %s", match.DateTime))
	}

	// Add BGA link indicator if available
	if match.BGALink != "" {
		result += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Render(" [BGA]")
	}

	return result
}

// GetCurrentRound returns the currently displayed round
func (m *FixtureModel) GetCurrentRound() *fixtures.Round {
	if m.currentRound >= 0 && m.currentRound < len(m.division.Rounds) {
		return m.division.Rounds[m.currentRound]
	}
	return nil
}
