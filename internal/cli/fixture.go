package cli

import (
	"fmt"
	"net/url"
	"time"

	"carca-cli/internal/fixtures"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// FixtureModel represents the fixture display TUI state
type FixtureModel struct {
	style         lipgloss.Style
	division      *fixtures.Division
	statusMessage string
	currentRound  int
	selectedMatch int
}

// NewFixtureModel creates a new fixture display model
func NewFixtureModel(division *fixtures.Division) *FixtureModel {
	return &FixtureModel{
		division:      division,
		currentRound:  0,
		selectedMatch: 0,
		statusMessage: "",
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

			m.selectedMatch = 0
		case tea.KeyRight:
			m.currentRound++
			if m.currentRound >= len(m.division.Rounds) {
				m.currentRound = 0
			}

			m.selectedMatch = 0
		case tea.KeyPgUp:
			m.currentRound--
			if m.currentRound < 0 {
				m.currentRound = len(m.division.Rounds) - 1
			}

			m.selectedMatch = 0
		case tea.KeyPgDown:
			m.currentRound++
			if m.currentRound >= len(m.division.Rounds) {
				m.currentRound = 0
			}

			m.selectedMatch = 0
		case tea.KeyEsc:
			// Go back to division selection
			return m, func() tea.Msg {
				return BackToMenuMsg{}
			}
		case tea.KeyDown:
			// Match selection down
			currentRound := m.GetCurrentRound()
			if currentRound != nil && len(currentRound.Matches) > 0 {
				m.selectedMatch++
				if m.selectedMatch >= len(currentRound.Matches) {
					m.selectedMatch = 0
				}
			}
		case tea.KeyUp:
			// Match selection up
			currentRound := m.GetCurrentRound()
			if currentRound != nil && len(currentRound.Matches) > 0 {
				m.selectedMatch--
				if m.selectedMatch < 0 {
					m.selectedMatch = len(currentRound.Matches) - 1
				}
			}
		case tea.KeyEnter:
			// Handle match selection
			currentRound := m.GetCurrentRound()
			if currentRound != nil && m.selectedMatch < len(currentRound.Matches) {
				selectedMatch := currentRound.Matches[m.selectedMatch]
				if selectedMatch.Played && selectedMatch.BGALink != "" {
					// Copy link to clipboard
					err := clipboard.WriteAll(selectedMatch.BGALink)
					if err != nil {
						m.statusMessage = "Failed to copy link to clipboard"
					} else {
						m.statusMessage = "Tournament link copied to clipboard!"
					}

					return m, tea.Batch(
						tea.Tick(time.Second*3, func(time.Time) tea.Msg {
							return clearStatusMsg{}
						}),
					)
				} else {
					// Match not played, show create tournament message
					m.statusMessage = "Press 'c' to create tournament for this match"
				}
			}
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				// Go back to division selection
				return m, func() tea.Msg {
					return BackToMenuMsg{}
				}
			case "h":
				// Vim left
				m.currentRound--
				if m.currentRound < 0 {
					m.currentRound = len(m.division.Rounds) - 1
				}
				// Reset match selection when changing rounds
				m.selectedMatch = 0
			case "l":
				// Vim right
				m.currentRound++
				if m.currentRound >= len(m.division.Rounds) {
					m.currentRound = 0
				}
				// Reset match selection when changing rounds
				m.selectedMatch = 0
			case "j":
				// Vim down for match selection
				currentRound := m.GetCurrentRound()
				if currentRound != nil && len(currentRound.Matches) > 0 {
					m.selectedMatch++
					if m.selectedMatch >= len(currentRound.Matches) {
						m.selectedMatch = 0
					}
				}
			case "k":
				// Vim up for match selection
				currentRound := m.GetCurrentRound()
				if currentRound != nil && len(currentRound.Matches) > 0 {
					m.selectedMatch--
					if m.selectedMatch < 0 {
						m.selectedMatch = len(currentRound.Matches) - 1
					}
				}
			case "c":
				// Create tournament for selected unplayed match
				currentRound := m.GetCurrentRound()
				if currentRound != nil && m.selectedMatch < len(currentRound.Matches) {
					selectedMatch := currentRound.Matches[m.selectedMatch]
					if !selectedMatch.Played {
						m.statusMessage = fmt.Sprintf("Creating tournament for %s vs %s...",
							selectedMatch.HomePlayer, selectedMatch.AwayPlayer)
						// TODO: Implement BGA tournament creation
						return m, tea.Batch(
							tea.Tick(time.Second*3, func(time.Time) tea.Msg {
								return clearStatusMsg{}
							}),
						)
					}
				}
			}
		}
	case clearStatusMsg:
		m.statusMessage = ""
	}

	return m, nil
}

// clearStatusMsg is sent to clear the status message after a delay
type clearStatusMsg struct{}

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

	// Display matches in table format
	if len(currentRound.Matches) == 0 {
		s += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Render("No matches in this round")
	} else {
		s += m.formatMatchesTable(currentRound.Matches)
	}

	// Navigation info
	s += fmt.Sprintf("\n\nRound %d of %d", m.currentRound+1, len(m.division.Rounds))

	// Show status message if present
	if m.statusMessage != "" {
		s += "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50C878")).
			Bold(true).
			Render(m.statusMessage)
	}

	s += "\n\nPress ←/→, h/l, or PgUp/PgDown to navigate rounds"
	s += "\nPress ↑/↓, j/k to select matches, Enter to copy link"
	s += "\nPress 'c' to create tournament for unplayed matches"
	s += "\nPress esc/q to go back.\n"

	return s
}

// formatMatchesTable formats matches in a table format
func (m *FixtureModel) formatMatchesTable(matches []*fixtures.Match) string {
	maxPlayerWidth := m.calculateMaxPlayerNameWidth()
	maxDateWidth := m.calculateMaxDateWidth()
	maxTournamentIDWidth := m.calculateMaxTournamentIDWidth()

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
			default:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
			}
		}).
		Headers("PLAYED", "LOCAL", "VISITOR", "RESULT", "DATE", "TOURNAMENT_ID")

	for i, match := range matches {
		var playedStatus string

		var result string

		var datetime string

		var tournamentID string

		// Format played status
		if match.Played {
			playedStatus = "✓"
		} else {
			playedStatus = "○"
		}

		// Format result
		if match.Played {
			result = fmt.Sprintf("%d-%d", match.HomeScore, match.AwayScore)
		} else {
			result = "-"
		}

		// Format datetime with fixed width
		if match.DateTime != "" {
			datetime = fmt.Sprintf("%-*s", maxDateWidth, match.DateTime)
		} else {
			datetime = fmt.Sprintf("%-*s", maxDateWidth, "-")
		}

		// Extract tournament ID with fixed width
		tournamentID = m.extractTournamentID(match.BGALink)
		if tournamentID == "" {
			tournamentID = fmt.Sprintf("%-*s", maxTournamentIDWidth, "-")
		} else {
			tournamentID = fmt.Sprintf("%-*s", maxTournamentIDWidth, tournamentID)
		}

		// Pad player names to consistent width
		homePlayer := fmt.Sprintf("%-*s", maxPlayerWidth, match.HomePlayer)
		awayPlayer := fmt.Sprintf("%-*s", maxPlayerWidth, match.AwayPlayer)

		// Add selection indicator for the selected match
		rowData := []string{playedStatus, homePlayer, awayPlayer, result, datetime, tournamentID}
		if i == m.selectedMatch {
			// Highlight selected row
			for j, cell := range rowData {
				rowData[j] = lipgloss.NewStyle().
					Background(lipgloss.Color("#7D56F4")).
					Foreground(lipgloss.Color("#FFFFFF")).
					Render(cell)
			}
		}

		t.Row(rowData[0], rowData[1], rowData[2], rowData[3], rowData[4], rowData[5])
	}

	return t.Render()
}

// calculateMaxPlayerNameWidth finds the longest player name across all rounds
func (m *FixtureModel) calculateMaxPlayerNameWidth() int {
	maxWidth := 8 // Minimum width for "VISITOR" header

	for _, round := range m.division.Rounds {
		for _, match := range round.Matches {
			if len(match.HomePlayer) > maxWidth {
				maxWidth = len(match.HomePlayer)
			}

			if len(match.AwayPlayer) > maxWidth {
				maxWidth = len(match.AwayPlayer)
			}
		}
	}

	// Add tab padding (8 spaces) for better readability
	return maxWidth + 8
}

// calculateMaxDateWidth finds the longest date/time string across all rounds
func (m *FixtureModel) calculateMaxDateWidth() int {
	maxWidth := 4 // Minimum width for "DATE" header

	for _, round := range m.division.Rounds {
		for _, match := range round.Matches {
			if match.DateTime != "" && len(match.DateTime) > maxWidth {
				maxWidth = len(match.DateTime)
			}
		}
	}

	return maxWidth
}

// calculateMaxTournamentIDWidth finds the longest tournament ID across all rounds
func (m *FixtureModel) calculateMaxTournamentIDWidth() int {
	maxWidth := 12 // Minimum width for "TOURNAMENT_ID" header

	for _, round := range m.division.Rounds {
		for _, match := range round.Matches {
			if match.BGALink != "" {
				tournamentID := m.extractTournamentID(match.BGALink)
				if len(tournamentID) > maxWidth {
					maxWidth = len(tournamentID)
				}
			}
		}
	}

	return maxWidth
}

// extractTournamentID extracts the tournament ID from a BGA URL
func (m *FixtureModel) extractTournamentID(bgaURL string) string {
	if bgaURL == "" {
		return ""
	}

	parsedURL, err := url.Parse(bgaURL)
	if err != nil {
		return ""
	}

	// Extract the 'id' parameter from the query string
	tournamentID := parsedURL.Query().Get("id")

	return tournamentID
}

// GetCurrentRound returns the currently displayed round
func (m *FixtureModel) GetCurrentRound() *fixtures.Round {
	if m.currentRound >= 0 && m.currentRound < len(m.division.Rounds) {
		return m.division.Rounds[m.currentRound]
	}

	return nil
}
