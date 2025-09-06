package cli

import (
	"fmt"
	"net/url"
	"time"

	"carca-cli/internal/bga"
	"carca-cli/internal/fixtures"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// FixtureModel represents the fixture display TUI state
type FixtureModel struct {
	division          *fixtures.Division
	bgaClient         bga.APIClient
	dateTimePicker    *DateTimePickerModel
	confirmationModel *TournamentConfirmationModel
	style             lipgloss.Style
	statusMessage     string
	currentRound      int
	selectedMatch     int
	showDatePicker    bool
	showConfirmation  bool
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

// SetBGAClient sets the BGA client for the fixture model
func (m *FixtureModel) SetBGAClient(client bga.APIClient) {
	m.bgaClient = client
}

// Update handles messages and updates the model state
func (m *FixtureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle sub-model messages first
	if model, cmd, handled := m.handleSubModelMessages(msg); handled {
		return model, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMessages(msg)
	case clearStatusMsg:
		m.statusMessage = ""
	case createTournamentMsg:
		return m.handleCreateTournamentResponse(msg)
	case createTournamentMsgWithDateTime:
		return m.handleCreateTournamentWithDateTime(&msg)
	case tournamentCreatedMsg:
		return m.handleTournamentCreated(msg)
	case DateTimeSelectedMsg:
		// DateTime selected, show confirmation screen
		m.showDatePicker = false
		m.confirmationModel = NewTournamentConfirmationModel(
			msg.HomePlayer,
			msg.AwayPlayer,
			msg.Division,
			msg.RoundNumber,
			msg.MatchNumber,
			msg.MatchID,
			msg.DateTime,
		)
		m.showConfirmation = true
		return m, nil
	case DateTimePickerCanceledMsg:
		// DateTime picker canceled
		m.showDatePicker = false
		m.statusMessage = "Tournament creation canceled"
		return m, tea.Tick(time.Second*2, func(time.Time) tea.Msg {
			return clearStatusMsg{}
		})
	case TournamentConfirmedMsg:
		// Confirmation received, proceed with tournament creation
		m.showConfirmation = false
		m.statusMessage = fmt.Sprintf("Creating tournament for %s vs %s...",
			msg.HomePlayer, msg.AwayPlayer)

		return m, tea.Cmd(func() tea.Msg {
			return createTournamentMsgWithDateTime{
				homePlayer:  msg.HomePlayer,
				awayPlayer:  msg.AwayPlayer,
				matchID:     msg.MatchID,
				roundNum:    msg.RoundNumber - 1, // Convert back to 0-based
				dateTime:    msg.DateTime,
				division:    msg.Division,
				matchNumber: msg.MatchNumber,
			}
		})
	case TournamentConfirmationCanceledMsg:
		// Tournament confirmation canceled
		m.showConfirmation = false
		m.statusMessage = "Tournament creation canceled"
		return m, tea.Tick(time.Second*2, func(time.Time) tea.Msg {
			return clearStatusMsg{}
		})
	}

	return m, nil
}

// clearStatusMsg is sent to clear the status message after a delay
type clearStatusMsg struct{}

// createTournamentMsg is sent to initiate tournament creation
type createTournamentMsg struct {
	homePlayer string
	awayPlayer string
	matchID    int
	roundNum   int
}

// createTournamentMsgWithDateTime is sent to initiate tournament creation with datetime
type createTournamentMsgWithDateTime struct {
	dateTime    time.Time
	homePlayer  string
	awayPlayer  string
	division    string
	matchID     int
	roundNum    int
	matchNumber int
}

// tournamentCreatedMsg is sent when tournament creation completes
type tournamentCreatedMsg struct {
	link         string
	error        string
	tournamentID int
	matchID      int
	roundNum     int
	success      bool
}

// handleCreateTournamentResponse handles the tournament creation request
func (m *FixtureModel) handleCreateTournamentResponse(msg createTournamentMsg) (tea.Model, tea.Cmd) {
	return m, tea.Cmd(func() tea.Msg {
		if !m.bgaClient.IsAuthenticated() {
			// Get credentials and login
			username, password, err := GetOrPromptCredentials(false)
			if err != nil {
				return tournamentCreatedMsg{
					success:  false,
					error:    fmt.Sprintf("Failed to get credentials: %v", err),
					matchID:  msg.matchID,
					roundNum: msg.roundNum,
				}
			}

			// Create new client with credentials
			if mockClient, ok := m.bgaClient.(*bga.MockClient); ok {
				// For testing, reset the mock client with new credentials
				*mockClient = *bga.NewMockClient(username, password)
			} else {
				m.bgaClient = bga.NewClient(username, password)
			}

			err = m.bgaClient.Login()
			if err != nil {
				return tournamentCreatedMsg{
					success:  false,
					error:    fmt.Sprintf("Login failed: %v", err),
					matchID:  msg.matchID,
					roundNum: msg.roundNum,
				}
			}
		}

		// Create tournament with division and match information (default scheduling)
		resp, err := m.bgaClient.CreateSwissTournament(
			m.division.Name,
			msg.homePlayer,
			msg.awayPlayer,
			msg.roundNum+1,
			msg.matchID,
		)

		if err != nil {
			return tournamentCreatedMsg{
				success:  false,
				error:    fmt.Sprintf("Tournament creation failed: %v", err),
				matchID:  msg.matchID,
				roundNum: msg.roundNum,
			}
		}

		if !resp.Success {
			return tournamentCreatedMsg{
				success:  false,
				error:    resp.Error,
				matchID:  msg.matchID,
				roundNum: msg.roundNum,
			}
		}

		return tournamentCreatedMsg{
			success:      true,
			tournamentID: resp.TournamentID,
			link:         resp.Link,
			matchID:      msg.matchID,
			roundNum:     msg.roundNum,
		}
	})
}

// handleTournamentCreated handles the tournament creation completion
func (m *FixtureModel) handleTournamentCreated(msg tournamentCreatedMsg) (tea.Model, tea.Cmd) {
	if !msg.success {
		m.statusMessage = fmt.Sprintf("Tournament creation failed: %s", msg.error)
	} else {
		// Update the match with the tournament link
		if msg.roundNum < len(m.division.Rounds) {
			round := m.division.Rounds[msg.roundNum]
			for i, match := range round.Matches {
				if match.ID == msg.matchID {
					round.Matches[i].BGALink = msg.link
					break
				}
			}
		}

		m.statusMessage = "Tournament created successfully! Link copied to clipboard."

		// Copy link to clipboard
		if err := clipboard.WriteAll(msg.link); err != nil {
			m.statusMessage = "Tournament created successfully! (Failed to copy link to clipboard)"
		}
	}

	return m, tea.Tick(time.Second*3, func(time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

// handleCreateTournamentWithDateTime handles tournament creation with specific datetime
func (m *FixtureModel) handleCreateTournamentWithDateTime(msg *createTournamentMsgWithDateTime) (tea.Model, tea.Cmd) {
	return m, tea.Cmd(func() tea.Msg {
		if !m.bgaClient.IsAuthenticated() {
			// Get credentials and login
			username, password, err := GetOrPromptCredentials(false)
			if err != nil {
				return tournamentCreatedMsg{
					success:  false,
					error:    fmt.Sprintf("Failed to get credentials: %v", err),
					matchID:  msg.matchID,
					roundNum: msg.roundNum,
				}
			}

			// Create new client with credentials
			if mockClient, ok := m.bgaClient.(*bga.MockClient); ok {
				// For testing, reset the mock client with new credentials
				*mockClient = *bga.NewMockClient(username, password)
			} else {
				m.bgaClient = bga.NewClient(username, password)
			}

			err = m.bgaClient.Login()
			if err != nil {
				return tournamentCreatedMsg{
					success:  false,
					error:    fmt.Sprintf("Login failed: %v", err),
					matchID:  msg.matchID,
					roundNum: msg.roundNum,
				}
			}
		}

		// Create tournament with specified datetime
		resp, err := m.bgaClient.CreateSwissTournamentWithDateTime(
			msg.division,
			msg.homePlayer,
			msg.awayPlayer,
			msg.roundNum+1,
			msg.matchNumber,
			msg.dateTime,
		)

		if err != nil {
			return tournamentCreatedMsg{
				success:  false,
				error:    fmt.Sprintf("Tournament creation failed: %v", err),
				matchID:  msg.matchID,
				roundNum: msg.roundNum,
			}
		}

		if !resp.Success {
			return tournamentCreatedMsg{
				success:  false,
				error:    resp.Error,
				matchID:  msg.matchID,
				roundNum: msg.roundNum,
			}
		}

		return tournamentCreatedMsg{
			success:      true,
			tournamentID: resp.TournamentID,
			link:         resp.Link,
			matchID:      msg.matchID,
			roundNum:     msg.roundNum,
		}
	})
}

// View renders the current state of the fixture display
func (m *FixtureModel) View() string {
	// Show confirmation screen if active
	if m.showConfirmation && m.confirmationModel != nil {
		return m.confirmationModel.View()
	}

	// Show datetime picker if active
	if m.showDatePicker && m.dateTimePicker != nil {
		return m.dateTimePicker.View()
	}
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

// handleSubModelMessages handles messages for date picker and confirmation models
func (m *FixtureModel) handleSubModelMessages(msg tea.Msg) (tea.Model, tea.Cmd, bool) {
	if m.showDatePicker && m.dateTimePicker != nil {
		switch msg.(type) {
		case DateTimeSelectedMsg, DateTimePickerCanceledMsg:
			// These are for the FixtureModel, so let them fall through.
			return m, nil, false
		default:
			// All other messages go to the date picker.
			updatedPicker, cmd := m.dateTimePicker.Update(msg)
			if picker, ok := updatedPicker.(*DateTimePickerModel); ok {
				m.dateTimePicker = picker
			}
			return m, cmd, true
		}
	}

	if m.showConfirmation && m.confirmationModel != nil {
		switch msg.(type) {
		case TournamentConfirmedMsg, TournamentConfirmationCanceledMsg:
			// These are for the FixtureModel.
			return m, nil, false
		default:
			updatedConfirmation, cmd := m.confirmationModel.Update(msg)
			if confirmation, ok := updatedConfirmation.(*TournamentConfirmationModel); ok {
				m.confirmationModel = confirmation
			}
			return m, cmd, true
		}
	}

	return m, nil, false
}

// handleKeyMessages handles all keyboard input
func (m *FixtureModel) handleKeyMessages(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyLeft, tea.KeyPgUp:
		return m.handleRoundNavigation(-1), nil
	case tea.KeyRight, tea.KeyPgDown:
		return m.handleRoundNavigation(1), nil
	case tea.KeyEsc:
		return m, func() tea.Msg { return BackToMenuMsg{} }
	case tea.KeyDown:
		m.handleMatchSelection(1)
	case tea.KeyUp:
		m.handleMatchSelection(-1)
	case tea.KeyEnter:
		return m.handleMatchEnter()
	}

	// Handle string keys
	switch msg.String() {
	case "c":
		return m.handleCreateTournament()
	case "h":
		return m.handleRoundNavigation(-1), nil
	case "l":
		return m.handleRoundNavigation(1), nil
	case "j":
		m.handleMatchSelection(1)
	case "k":
		m.handleMatchSelection(-1)
	case "q":
		return m, func() tea.Msg { return BackToMenuMsg{} }
	}

	return m, nil
}

// handleRoundNavigation navigates between rounds
func (m *FixtureModel) handleRoundNavigation(direction int) *FixtureModel {
	m.currentRound += direction
	if m.currentRound < 0 {
		m.currentRound = len(m.division.Rounds) - 1
	} else if m.currentRound >= len(m.division.Rounds) {
		m.currentRound = 0
	}
	m.selectedMatch = 0
	return m
}

// handleMatchSelection handles match selection up/down
func (m *FixtureModel) handleMatchSelection(direction int) {
	currentRound := m.GetCurrentRound()
	if currentRound == nil || len(currentRound.Matches) == 0 {
		return
	}

	m.selectedMatch += direction
	if m.selectedMatch < 0 {
		m.selectedMatch = len(currentRound.Matches) - 1
	} else if m.selectedMatch >= len(currentRound.Matches) {
		m.selectedMatch = 0
	}
}

// handleMatchEnter handles Enter key on selected match
func (m *FixtureModel) handleMatchEnter() (tea.Model, tea.Cmd) {
	currentRound := m.GetCurrentRound()
	if currentRound == nil || m.selectedMatch >= len(currentRound.Matches) {
		return m, nil
	}

	selectedMatch := currentRound.Matches[m.selectedMatch]
	if selectedMatch.Played && selectedMatch.BGALink != "" {
		// Copy existing link to clipboard
		if err := clipboard.WriteAll(selectedMatch.BGALink); err == nil {
			m.statusMessage = "Tournament link copied to clipboard!"
		} else {
			m.statusMessage = "Failed to copy link to clipboard"
		}

		return m, tea.Tick(time.Second*3, func(time.Time) tea.Msg {
			return clearStatusMsg{}
		})
	}

	// Match not played, show create tournament message
	m.statusMessage = "Press 'c' to create tournament for this match"
	return m, nil
}

// handleCreateTournament handles 'c' key for tournament creation
func (m *FixtureModel) handleCreateTournament() (tea.Model, tea.Cmd) {
	currentRound := m.GetCurrentRound()
	if currentRound == nil || m.selectedMatch >= len(currentRound.Matches) {
		return m, nil
	}

	selectedMatch := currentRound.Matches[m.selectedMatch]
	if !selectedMatch.Played {
		// Create and show datetime picker
		m.dateTimePicker = NewDateTimePickerModel(
			selectedMatch.HomePlayer,
			selectedMatch.AwayPlayer,
			m.division.Name,
			m.currentRound+1,
			selectedMatch.ID, // Use match ID as match number
			selectedMatch.ID,
		)
		m.showDatePicker = true
		return m, m.dateTimePicker.Init()
	}

	return m, nil
}
