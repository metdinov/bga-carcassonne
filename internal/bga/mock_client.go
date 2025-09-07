package bga

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// MockClient is a mock implementation of the BGA client for testing
type MockClient struct {
	tournaments      map[int]*TournamentStatus
	username         string
	password         string
	nextTournamentID int
	isAuthenticated  bool
	shouldFailLogin  bool
	shouldFailCreate bool
}

// NewMockClient creates a new mock BGA client
func NewMockClient(username, password string) *MockClient {
	return &MockClient{
		username:         username,
		password:         password,
		tournaments:      make(map[int]*TournamentStatus),
		nextTournamentID: 423762, // Start with a realistic tournament ID
	}
}

// SetShouldFailLogin configures the mock to fail login attempts
func (m *MockClient) SetShouldFailLogin(shouldFail bool) {
	m.shouldFailLogin = shouldFail
}

// SetShouldFailCreate configures the mock to fail tournament creation
func (m *MockClient) SetShouldFailCreate(shouldFail bool) {
	m.shouldFailCreate = shouldFail
}

// Login simulates authentication with BGA
func (m *MockClient) Login() error {
	if m.shouldFailLogin {
		return fmt.Errorf("authentication failed: invalid credentials")
	}

	if m.username == "" || m.password == "" {
		return fmt.Errorf("username and password are required")
	}

	// Simulate authentication delay
	time.Sleep(100 * time.Millisecond)

	m.isAuthenticated = true
	return nil
}

// CreateTournament simulates creating a tournament on BGA
func (m *MockClient) CreateTournament(config *TournamentConfig) (*TournamentResponse, error) {
	if !m.isAuthenticated {
		return nil, fmt.Errorf("not authenticated: call Login() first")
	}

	if m.shouldFailCreate {
		return &TournamentResponse{
			Success: false,
			Error:   "tournament creation failed: server error",
		}, nil
	}

	// Validate required fields
	if config.LocalPlayer == "" || config.VisitorPlayer == "" {
		return &TournamentResponse{
			Success: false,
			Error:   "both local and visitor players are required",
		}, nil
	}

	if config.TournamentName == "" {
		return &TournamentResponse{
			Success: false,
			Error:   "tournament name is required",
		}, nil
	}

	// Generate tournament ID and create response
	tournamentID := m.nextTournamentID
	m.nextTournamentID++

	link := fmt.Sprintf("https://boardgamearena.com/tournament?id=%d", tournamentID)

	// Create tournament status for tracking
	status := &TournamentStatus{
		ID:           tournamentID,
		Name:         config.TournamentName,
		Status:       "waiting",
		PlayersCount: 2,
		Matches: []MatchStatus{
			{
				ID:         1,
				Status:     "waiting",
				HomePlayer: config.LocalPlayer,
				AwayPlayer: config.VisitorPlayer,
				HomeScore:  0,
				AwayScore:  0,
			},
			{
				ID:         2,
				Status:     "waiting",
				HomePlayer: config.LocalPlayer,
				AwayPlayer: config.VisitorPlayer,
				HomeScore:  0,
				AwayScore:  0,
			},
			{
				ID:         3,
				Status:     "waiting",
				HomePlayer: config.LocalPlayer,
				AwayPlayer: config.VisitorPlayer,
				HomeScore:  0,
				AwayScore:  0,
			},
		},
		Results: map[string]int{
			config.LocalPlayer:   0,
			config.VisitorPlayer: 0,
		},
	}

	m.tournaments[tournamentID] = status

	// Simulate network delay
	time.Sleep(200 * time.Millisecond)

	return &TournamentResponse{
		Success:      true,
		TournamentID: tournamentID,
		Link:         link,
	}, nil
}

// CreateSwissTournament creates a mock best-of-3 Swiss tournament
func (m *MockClient) CreateSwissTournament(
	division, homePlayer, awayPlayer string,
	roundNumber, matchNumber int,
) (*TournamentResponse, error) {
	// Format championship and tournament names according to requirements
	championshipName := fmt.Sprintf("Division %s - 1era Temporada", division)
	tournamentName := fmt.Sprintf("%d Fecha - Duelo %d - %s vs %s", roundNumber, matchNumber, homePlayer, awayPlayer)

	config := &TournamentConfig{
		GameID:           1,
		ChampionshipName: championshipName,
		TournamentName:   tournamentName,
		MaxPlayers:       2,
		MinPlayers:       2,
		BaseDate:         time.Now().Format("2006-01-02"),
		BaseDateTime:     "21:00",
		GameDuration:     1800,
		MatchesCount:     3,
		Division:         division,
		RoundNumber:      roundNumber,
		MatchNumber:      matchNumber,
		LocalPlayer:      homePlayer,
		VisitorPlayer:    awayPlayer,
	}

	return m.CreateTournament(config)
}

// CreateSwissTournamentWithDateTime creates a best-of-3 Swiss tournament for two players with specific datetime
func (m *MockClient) CreateSwissTournamentWithDateTime(
	division, homePlayer, awayPlayer string,
	roundNumber, matchNumber int,
	scheduledTime time.Time,
) (*TournamentResponse, error) {
	// Format championship and tournament names according to requirements
	championshipName := fmt.Sprintf("Division %s - 1era Temporada", division)
	tournamentName := fmt.Sprintf("%d Fecha - Duelo %d - %s vs %s", roundNumber, matchNumber, homePlayer, awayPlayer)

	config := &TournamentConfig{
		GameID:           1,
		ChampionshipName: championshipName,
		TournamentName:   tournamentName,
		MaxPlayers:       2,
		MinPlayers:       2,
		BaseDate:         scheduledTime.Format("2006-01-02"),
		BaseDateTime:     scheduledTime.Format("15:04"),
		GameDuration:     1800,
		MatchesCount:     3,
		Division:         division,
		RoundNumber:      roundNumber,
		MatchNumber:      matchNumber,
		LocalPlayer:      homePlayer,
		VisitorPlayer:    awayPlayer,
	}

	return m.CreateTournament(config)
}

// GetTournamentStatus returns the mock status of a tournament
func (m *MockClient) GetTournamentStatus(tournamentID int) (*TournamentStatus, error) {
	if !m.isAuthenticated {
		return nil, fmt.Errorf("not authenticated: call Login() first")
	}

	status, exists := m.tournaments[tournamentID]
	if !exists {
		return nil, fmt.Errorf("tournament not found: %d", tournamentID)
	}

	// Return a copy to prevent external modification
	statusCopy := *status
	statusCopy.Matches = make([]MatchStatus, len(status.Matches))
	copy(statusCopy.Matches, status.Matches)

	return &statusCopy, nil
}

// IsAuthenticated returns whether the mock client is authenticated
func (m *MockClient) IsAuthenticated() bool {
	return m.isAuthenticated
}

// Logout simulates logging out of BGA
func (m *MockClient) Logout() error {
	m.isAuthenticated = false
	return nil
}

// SimulateMatchResult simulates a match result for testing
func (m *MockClient) SimulateMatchResult(tournamentID, matchID, homeScore, awayScore int, winner string) error {
	status, exists := m.tournaments[tournamentID]
	if !exists {
		return fmt.Errorf("tournament not found: %d", tournamentID)
	}

	for i := range status.Matches {
		if status.Matches[i].ID != matchID {
			continue
		}
		status.Matches[i].Status = "finished"
		status.Matches[i].HomeScore = homeScore
		status.Matches[i].AwayScore = awayScore
		status.Matches[i].Winner = winner

		// Update tournament results
		if winner != "" {
			status.Results[winner]++
		}

		break
	}

	// Check if all matches are finished
	allFinished := true
	for _, match := range status.Matches {
		if match.Status != "finished" {
			allFinished = false
			break
		}
	}

	if allFinished {
		status.Status = "finished"
	} else {
		status.Status = "in_progress"
	}

	return nil
}

// GetTournaments returns all tournaments created by this mock client
func (m *MockClient) GetTournaments() map[int]*TournamentStatus {
	tournaments := make(map[int]*TournamentStatus)
	for id, status := range m.tournaments {
		statusCopy := *status
		statusCopy.Matches = make([]MatchStatus, len(status.Matches))
		copy(statusCopy.Matches, status.Matches)
		tournaments[id] = &statusCopy
	}
	return tournaments
}

// ExtractTournamentID extracts tournament ID from a BGA tournament URL
func ExtractTournamentID(link string) (int, error) {
	if link == "" {
		return 0, fmt.Errorf("empty tournament link")
	}

	// Extract ID from URL like "https://boardgamearena.com/tournament?id=423761"
	parts := strings.Split(link, "id=")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid tournament link format: %s", link)
	}

	idStr := parts[1]
	// Remove any additional query parameters
	if ampIndex := strings.Index(idStr, "&"); ampIndex != -1 {
		idStr = idStr[:ampIndex]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid tournament ID in link: %s", idStr)
	}

	return id, nil
}

// LaunchTournament simulates launching a created tournament
func (m *MockClient) LaunchTournament(tournamentID int) error {
	if !m.isAuthenticated {
		return fmt.Errorf("not authenticated: call Login() first")
	}

	tournament, exists := m.tournaments[tournamentID]
	if !exists {
		return fmt.Errorf("tournament with ID %d not found", tournamentID)
	}

	// Change tournament status from "created" or "waiting" to "open"
	if tournament.Status == "created" || tournament.Status == "waiting" {
		tournament.Status = "open"
	}

	// Simulate network delay
	time.Sleep(100 * time.Millisecond)

	return nil
}

// InvitePlayer simulates inviting a player to a tournament
func (m *MockClient) InvitePlayer(tournamentID int, playerID string) error {
	if !m.isAuthenticated {
		return fmt.Errorf("not authenticated: call Login() first")
	}

	tournament, exists := m.tournaments[tournamentID]
	if !exists {
		return fmt.Errorf("tournament with ID %d not found", tournamentID)
	}

	// Validate that tournament is in open state (launched)
	if tournament.Status != "open" {
		return fmt.Errorf("cannot invite players to tournament in %s state", tournament.Status)
	}

	// Simulate network delay
	time.Sleep(50 * time.Millisecond)

	return nil
}

// Reset clears all tournament data from the mock client
func (m *MockClient) Reset() {
	m.tournaments = make(map[int]*TournamentStatus)
	m.nextTournamentID = 423762
	m.isAuthenticated = false
	m.shouldFailLogin = false
	m.shouldFailCreate = false
}
