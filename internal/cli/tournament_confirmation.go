package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TournamentConfirmationModel represents the tournament confirmation screen
type TournamentConfirmationModel struct {
	timezone         *time.Location
	selectedTime     time.Time
	style            lipgloss.Style
	headerStyle      lipgloss.Style
	detailStyle      lipgloss.Style
	highlightStyle   lipgloss.Style
	instructionStyle lipgloss.Style
	title            string
	championshipName string
	tournamentName   string
	division         string
	homePlayer       string
	awayPlayer       string
	roundNumber      int
	matchNumber      int
	matchID          int
	confirmed        bool
	canceled         bool
}

// TournamentConfirmedMsg is sent when the user confirms tournament creation
type TournamentConfirmedMsg struct {
	DateTime    time.Time
	HomePlayer  string
	AwayPlayer  string
	Division    string
	RoundNumber int
	MatchNumber int
	MatchID     int
}

// TournamentConfirmationCanceledMsg is sent when the user cancels tournament creation
type TournamentConfirmationCanceledMsg struct{}

// EditDateTimeMsg is sent when the user wants to edit the selected datetime
type EditDateTimeMsg struct {
	DateTime    time.Time
	HomePlayer  string
	AwayPlayer  string
	Division    string
	RoundNumber int
	MatchNumber int
	MatchID     int
}

// NewTournamentConfirmationModel creates a new tournament confirmation model
func NewTournamentConfirmationModel(
	homePlayer, awayPlayer, division string,
	roundNumber, matchNumber, matchID int,
	selectedTime time.Time,
) *TournamentConfirmationModel {
	// Get local timezone
	localTZ := selectedTime.Location()

	// Generate tournament and championship names
	championshipName := fmt.Sprintf("Division %s - 1era Temporada", division)
	tournamentName := fmt.Sprintf("%d Fecha - Duelo %d - %s vs %s", roundNumber, matchNumber, homePlayer, awayPlayer)

	return &TournamentConfirmationModel{
		title:            "Tournament Confirmation",
		championshipName: championshipName,
		tournamentName:   tournamentName,
		division:         division,
		homePlayer:       homePlayer,
		awayPlayer:       awayPlayer,
		roundNumber:      roundNumber,
		matchNumber:      matchNumber,
		matchID:          matchID,
		selectedTime:     selectedTime,
		timezone:         localTZ,
		style: lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")),
		headerStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1).
			Bold(true),
		detailStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true),
		highlightStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EE6FF8")).
			Bold(true),
		instructionStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true).
			MarginTop(1),
	}
}

// Init initializes the tournament confirmation model
func (m *TournamentConfirmationModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the tournament confirmation screen
func (m *TournamentConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			// Confirm tournament creation
			m.confirmed = true
			return m, tea.Cmd(func() tea.Msg {
				return TournamentConfirmedMsg{
					HomePlayer:  m.homePlayer,
					AwayPlayer:  m.awayPlayer,
					Division:    m.division,
					RoundNumber: m.roundNumber,
					MatchNumber: m.matchNumber,
					MatchID:     m.matchID,
					DateTime:    m.selectedTime,
				}
			})

		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			// Cancel tournament creation
			m.canceled = true
			return m, tea.Cmd(func() tea.Msg {
				return TournamentConfirmationCanceledMsg{}
			})

		case key.Matches(msg, key.NewBinding(key.WithKeys("e"))):
			// Edit datetime - go back to datetime picker
			return m, tea.Cmd(func() tea.Msg {
				return EditDateTimeMsg{
					HomePlayer:  m.homePlayer,
					AwayPlayer:  m.awayPlayer,
					Division:    m.division,
					RoundNumber: m.roundNumber,
					MatchNumber: m.matchNumber,
					MatchID:     m.matchID,
					DateTime:    m.selectedTime,
				}
			})
		}
	}

	return m, nil
}

// View renders the tournament confirmation screen
func (m *TournamentConfirmationModel) View() string {
	if m.confirmed || m.canceled {
		return ""
	}

	var content strings.Builder

	// Header
	header := m.headerStyle.Render(m.title)
	content.WriteString(header + "\n\n")

	// Tournament Details
	content.WriteString(m.detailStyle.Render("Tournament Details:") + "\n")
	content.WriteString(fmt.Sprintf("• Championship: %s\n", m.highlightStyle.Render(m.championshipName)))
	content.WriteString(fmt.Sprintf("• Tournament:   %s\n", m.highlightStyle.Render(m.tournamentName)))
	content.WriteString("\n")

	// Match Information
	content.WriteString(m.detailStyle.Render("Match Information:") + "\n")
	content.WriteString(fmt.Sprintf("• Division:     %s\n", m.highlightStyle.Render(m.division)))
	content.WriteString(fmt.Sprintf("• Round:        %s\n", m.highlightStyle.Render(fmt.Sprintf("%d", m.roundNumber))))
	content.WriteString(fmt.Sprintf("• Match (Duelo): %s\n", m.highlightStyle.Render(fmt.Sprintf("%d", m.matchNumber))))
	content.WriteString(fmt.Sprintf("• Players:      %s vs %s\n",
		m.highlightStyle.Render(m.homePlayer),
		m.highlightStyle.Render(m.awayPlayer)))
	content.WriteString("\n")

	// Scheduling Information
	content.WriteString(m.detailStyle.Render("Scheduling:") + "\n")

	// Format timezone offset
	_, offset := m.selectedTime.Zone()
	offsetHours := offset / 3600
	offsetMins := (offset % 3600) / 60
	var offsetStr string
	if offsetMins == 0 {
		offsetStr = fmt.Sprintf("UTC%+d", offsetHours)
	} else {
		offsetStr = fmt.Sprintf("UTC%+d:%02d", offsetHours, offsetMins)
	}

	content.WriteString(fmt.Sprintf("• Date & Time:  %s\n",
		m.highlightStyle.Render(m.selectedTime.Format("Monday, January 2, 2006 at 3:04 PM"))))
	content.WriteString(fmt.Sprintf("• Timezone:     %s (%s)\n",
		m.highlightStyle.Render(m.timezone.String()),
		m.highlightStyle.Render(offsetStr)))
	content.WriteString("\n")

	// Tournament Settings
	content.WriteString(m.detailStyle.Render("Tournament Settings:") + "\n")
	content.WriteString("• Format:       Swiss System (Best-of-3)\n")
	content.WriteString("• Game:         Carcassonne\n")
	content.WriteString("• Duration:     30 minutes (15 min per player)\n")
	content.WriteString("• Players:      2 (Private tournament)\n")
	content.WriteString("• Rules:        International scoring\n")
	content.WriteString("  - Field scoring: 3 points per city\n")
	content.WriteString("  - City scoring:  4 points per two tile city\n")
	content.WriteString("• Expansions:   None\n")
	content.WriteString("• Variants:     None\n")
	content.WriteString("\n")

	// Instructions
	instructions := "Press Enter to create tournament • Press 'e' to edit date/time • Press Esc to cancel"
	content.WriteString(m.instructionStyle.Render(instructions))

	return m.style.Render(content.String())
}

// IsConfirmed returns whether the tournament was confirmed
func (m *TournamentConfirmationModel) IsConfirmed() bool {
	return m.confirmed
}

// IsCanceled returns whether the tournament creation was canceled
func (m *TournamentConfirmationModel) IsCanceled() bool {
	return m.canceled
}

// GetTournamentDetails returns the tournament details for display
func (m *TournamentConfirmationModel) GetTournamentDetails() (championshipName, tournamentName string) {
	return m.championshipName, m.tournamentName
}

// GetSchedulingInfo returns formatted scheduling information
func (m *TournamentConfirmationModel) GetSchedulingInfo() (dateStr, timeStr string) {
	return m.selectedTime.Format("2006-01-02"), m.selectedTime.Format("15:04")
}
