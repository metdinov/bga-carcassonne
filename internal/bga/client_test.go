package bga

import (
	"fmt"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	username := "testuser"
	password := "testpass"

	client := NewClient(username, password)

	if client.username != username {
		t.Errorf("Expected username %s, got %s", username, client.username)
	}

	if client.password != password {
		t.Errorf("Expected password %s, got %s", password, client.password)
	}

	if client.baseURL != "https://boardgamearena.com" {
		t.Errorf("Expected baseURL https://boardgamearena.com, got %s", client.baseURL)
	}

	if client.IsAuthenticated() {
		t.Error("New client should not be authenticated")
	}
}

func TestClient_IsAuthenticated(t *testing.T) {
	client := NewClient("user", "pass")

	// Should not be authenticated initially
	if client.IsAuthenticated() {
		t.Error("Client should not be authenticated initially")
	}

	// Simulate authentication by setting session ID
	client.sessionID = "test-session-id"

	if !client.IsAuthenticated() {
		t.Error("Client should be authenticated after setting session ID")
	}
}

func TestNewMockClient(t *testing.T) {
	username := "testuser"
	password := "testpass"

	mockClient := NewMockClient(username, password)

	if mockClient.username != username {
		t.Errorf("Expected username %s, got %s", username, mockClient.username)
	}

	if mockClient.password != password {
		t.Errorf("Expected password %s, got %s", password, mockClient.password)
	}

	if mockClient.IsAuthenticated() {
		t.Error("New mock client should not be authenticated")
	}

	if mockClient.nextTournamentID != 423762 {
		t.Errorf("Expected nextTournamentID 423762, got %d", mockClient.nextTournamentID)
	}
}

func TestMockClient_Login(t *testing.T) {
	testCases := []struct {
		name          string
		username      string
		password      string
		shouldFail    bool
		expectedError bool
	}{
		{
			name:          "successful login",
			username:      "testuser",
			password:      "testpass",
			shouldFail:    false,
			expectedError: false,
		},
		{
			name:          "forced failure",
			username:      "testuser",
			password:      "testpass",
			shouldFail:    true,
			expectedError: true,
		},
		{
			name:          "empty username",
			username:      "",
			password:      "testpass",
			shouldFail:    false,
			expectedError: true,
		},
		{
			name:          "empty password",
			username:      "testuser",
			password:      "",
			shouldFail:    false,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := NewMockClient(tc.username, tc.password)
			mockClient.SetShouldFailLogin(tc.shouldFail)

			err := mockClient.Login()

			if tc.expectedError && err == nil {
				t.Error("Expected login to fail, but it succeeded")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected login to succeed, but got error: %v", err)
			}

			if !tc.expectedError && !mockClient.IsAuthenticated() {
				t.Error("Expected client to be authenticated after successful login")
			}
		})
	}
}

func TestMockClient_CreateSwissTournament(t *testing.T) {
	mockClient := NewMockClient("testuser", "testpass")

	// Should fail when not authenticated
	_, err := mockClient.CreateSwissTournament("Elite", "player1", "player2", 1, 15)
	if err == nil {
		t.Error("Expected tournament creation to fail when not authenticated")
	}

	// Login first
	err = mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Test successful tournament creation
	resp, err := mockClient.CreateSwissTournament("Elite", "player1", "player2", 1, 15)
	if err != nil {
		t.Fatalf("Failed to create tournament: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected tournament creation to succeed, got error: %s", resp.Error)
	}

	if resp.TournamentID == 0 {
		t.Error("Expected tournament ID to be set")
	}

	if resp.Link == "" {
		t.Error("Expected tournament link to be set")
	}

	expectedLink := "https://boardgamearena.com/tournament?id=423762"
	if resp.Link != expectedLink {
		t.Errorf("Expected link %s, got %s", expectedLink, resp.Link)
	}
}

func TestMockClient_CreateTournament_ValidationErrors(t *testing.T) {
	mockClient := NewMockClient("testuser", "testpass")
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	testCases := []struct {
		name          string
		expectedError string
		config        TournamentConfig
	}{
		{
			name: "missing home player",
			config: TournamentConfig{
				TournamentName: "Test Tournament",
				VisitorPlayer:  "player2",
			},
			expectedError: "both local and visitor players are required",
		},
		{
			name: "missing away player",
			config: TournamentConfig{
				TournamentName: "Test Tournament",
				LocalPlayer:    "player1",
			},
			expectedError: "both local and visitor players are required",
		},
		{
			name: "missing tournament name",
			config: TournamentConfig{
				LocalPlayer:   "player1",
				VisitorPlayer: "player2",
			},
			expectedError: "tournament name is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := mockClient.CreateTournament(&tc.config)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if resp.Success {
				t.Error("Expected tournament creation to fail")
			}

			if resp.Error != tc.expectedError {
				t.Errorf("Expected error %s, got %s", tc.expectedError, resp.Error)
			}
		})
	}
}

func TestMockClient_GetTournamentStatus(t *testing.T) {
	mockClient := NewMockClient("testuser", "testpass")
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Create a tournament first
	resp, err := mockClient.CreateSwissTournament("Elite", "player1", "player2", 1, 15)
	if err != nil {
		t.Fatalf("Failed to create tournament: %v", err)
	}

	// Get tournament status
	status, err := mockClient.GetTournamentStatus(resp.TournamentID)
	if err != nil {
		t.Fatalf("Failed to get tournament status: %v", err)
	}

	if status.ID != resp.TournamentID {
		t.Errorf("Expected tournament ID %d, got %d", resp.TournamentID, status.ID)
	}

	expectedName := "1 Fecha - Duelo 15 - player1 vs player2"
	if status.Name != expectedName {
		t.Errorf("Expected tournament name '%s', got %s", expectedName, status.Name)
	}

	if status.Status != "waiting" {
		t.Errorf("Expected status 'waiting', got %s", status.Status)
	}

	if status.PlayersCount != 2 {
		t.Errorf("Expected 2 players, got %d", status.PlayersCount)
	}

	if len(status.Matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(status.Matches))
	}

	// Test non-existent tournament
	_, err = mockClient.GetTournamentStatus(999999)
	if err == nil {
		t.Error("Expected error for non-existent tournament")
	}
}

func TestMockClient_SimulateMatchResult(t *testing.T) {
	mockClient := NewMockClient("testuser", "testpass")
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Create a tournament
	resp, err := mockClient.CreateSwissTournament("Elite", "player1", "player2", 1, 15)
	if err != nil {
		t.Fatalf("Failed to create tournament: %v", err)
	}

	// Simulate match result
	err = mockClient.SimulateMatchResult(resp.TournamentID, 1, 1, 0, "player1")
	if err != nil {
		t.Fatalf("Failed to simulate match result: %v", err)
	}

	// Check updated status
	status, err := mockClient.GetTournamentStatus(resp.TournamentID)
	if err != nil {
		t.Fatalf("Failed to get tournament status: %v", err)
	}

	// Tournament should be in progress after one match
	if status.Status != "in_progress" {
		t.Errorf("Expected status 'in_progress', got %s", status.Status)
	}

	// First match should be finished
	if status.Matches[0].Status != "finished" {
		t.Errorf("Expected first match status 'finished', got %s", status.Matches[0].Status)
	}

	if status.Matches[0].HomeScore != 1 {
		t.Errorf("Expected home score 1, got %d", status.Matches[0].HomeScore)
	}

	if status.Matches[0].Winner != "player1" {
		t.Errorf("Expected winner 'player1', got %s", status.Matches[0].Winner)
	}

	// Player1 should have 1 win
	if status.Results["player1"] != 1 {
		t.Errorf("Expected player1 to have 1 win, got %d", status.Results["player1"])
	}
}

func TestMockClient_Logout(t *testing.T) {
	mockClient := NewMockClient("testuser", "testpass")

	// Login first
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if !mockClient.IsAuthenticated() {
		t.Error("Client should be authenticated after login")
	}

	// Logout
	err = mockClient.Logout()
	if err != nil {
		t.Fatalf("Failed to logout: %v", err)
	}

	if mockClient.IsAuthenticated() {
		t.Error("Client should not be authenticated after logout")
	}
}

func TestExtractTournamentID(t *testing.T) {
	testCases := []struct {
		name        string
		link        string
		expectedID  int
		expectError bool
	}{
		{
			name:        "valid tournament link",
			link:        "https://boardgamearena.com/tournament?id=423761",
			expectedID:  423761,
			expectError: false,
		},
		{
			name:        "tournament link with additional parameters",
			link:        "https://boardgamearena.com/tournament?id=423761&foo=bar",
			expectedID:  423761,
			expectError: false,
		},
		{
			name:        "empty link",
			link:        "",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "invalid format",
			link:        "https://boardgamearena.com/tournament",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "invalid ID",
			link:        "https://boardgamearena.com/tournament?id=abc",
			expectedID:  0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := ExtractTournamentID(tc.link)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if id != tc.expectedID {
				t.Errorf("Expected ID %d, got %d", tc.expectedID, id)
			}
		})
	}
}

func TestMockClient_Reset(t *testing.T) {
	mockClient := NewMockClient("testuser", "testpass")

	// Login and create tournament
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	_, err = mockClient.CreateSwissTournament("Elite", "player1", "player2", 1, 15)
	if err != nil {
		t.Fatalf("Failed to create tournament: %v", err)
	}

	// Verify state before reset
	if !mockClient.IsAuthenticated() {
		t.Error("Expected client to be authenticated")
	}

	if len(mockClient.tournaments) == 0 {
		t.Error("Expected tournaments to exist")
	}

	// Reset
	mockClient.Reset()

	// Verify state after reset
	if mockClient.IsAuthenticated() {
		t.Error("Expected client to not be authenticated after reset")
	}

	if len(mockClient.tournaments) != 0 {
		t.Error("Expected no tournaments after reset")
	}

	if mockClient.nextTournamentID != 423762 {
		t.Errorf("Expected nextTournamentID to be reset to 423762, got %d", mockClient.nextTournamentID)
	}
}

func TestMockClient_CreateSwissTournamentWithDateTime(t *testing.T) {
	mockClient := NewMockClient("testuser", "testpass")
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Create tournament with specific datetime
	scheduledTime := time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local)
	resp, err := mockClient.CreateSwissTournamentWithDateTime(
		"Elite",
		"herchu",
		"Lord Trooper",
		1,
		15,
		scheduledTime,
	)
	if err != nil {
		t.Fatalf("Failed to create tournament with datetime: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Tournament creation failed: %s", resp.Error)
	}

	if resp.TournamentID == 0 {
		t.Error("Expected tournament ID to be set")
	}

	if resp.Link == "" {
		t.Error("Expected tournament link to be set")
	}

	// Verify tournament was created with correct name
	status, err := mockClient.GetTournamentStatus(resp.TournamentID)
	if err != nil {
		t.Fatalf("Failed to get tournament status: %v", err)
	}

	expectedName := "1 Fecha - Duelo 15 - herchu vs Lord Trooper"
	if status.Name != expectedName {
		t.Errorf("Expected tournament name '%s', got '%s'", expectedName, status.Name)
	}
}

func TestTournamentNamingConvention(t *testing.T) {
	mockClient := NewMockClient("testuser", "testpass")
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	testCases := []struct {
		name                 string
		division             string
		homePlayer           string
		awayPlayer           string
		expectedChampionship string
		expectedTournament   string
		roundNumber          int
		matchNumber          int
	}{
		{
			name:                 "Elite division tournament",
			division:             "Elite",
			homePlayer:           "herchu",
			awayPlayer:           "Lord Trooper",
			roundNumber:          1,
			matchNumber:          15,
			expectedChampionship: "Division Elite - 1era Temporada",
			expectedTournament:   "1 Fecha - Duelo 15 - herchu vs Lord Trooper",
		},
		{
			name:                 "Platinum A tournament",
			division:             "Platinum A",
			homePlayer:           "webbi",
			awayPlayer:           "alehrosario",
			roundNumber:          5,
			matchNumber:          23,
			expectedChampionship: "Division Platinum A - 1era Temporada",
			expectedTournament:   "5 Fecha - Duelo 23 - webbi vs alehrosario",
		},
		{
			name:                 "Oro B tournament",
			division:             "Oro B",
			homePlayer:           "bignacho610",
			awayPlayer:           "Academia47",
			roundNumber:          3,
			matchNumber:          8,
			expectedChampionship: "Division Oro B - 1era Temporada",
			expectedTournament:   "3 Fecha - Duelo 8 - bignacho610 vs Academia47",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := mockClient.CreateSwissTournament(
				tc.division,
				tc.homePlayer,
				tc.awayPlayer,
				tc.roundNumber,
				tc.matchNumber,
			)
			if err != nil {
				t.Fatalf("Failed to create tournament: %v", err)
			}

			if !resp.Success {
				t.Fatalf("Tournament creation failed: %s", resp.Error)
			}

			// Get tournament status to verify naming
			status, err := mockClient.GetTournamentStatus(resp.TournamentID)
			if err != nil {
				t.Fatalf("Failed to get tournament status: %v", err)
			}

			if status.Name != tc.expectedTournament {
				t.Errorf("Expected tournament name '%s', got '%s'", tc.expectedTournament, status.Name)
			}

			// We can't directly check championship name from status, but we can verify
			// the tournament was created successfully with proper formatting
			expectedLink := fmt.Sprintf("https://boardgamearena.com/tournament?id=%d", resp.TournamentID)
			if resp.Link != expectedLink {
				t.Errorf("Expected tournament link '%s', got '%s'", expectedLink, resp.Link)
			}
		})
	}
}

func TestMockClient_ThreeStepTournamentCreation(t *testing.T) {
	mockClient := NewMockClient("testuser", "testpass")
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Step 1: Create tournament
	resp, err := mockClient.CreateSwissTournament("Elite", "herchu", "Lord Trooper", 1, 15)
	if err != nil {
		t.Fatalf("Failed to create tournament: %v", err)
	}

	if resp.TournamentID == 0 {
		t.Error("Expected valid tournament ID")
	}

	// Step 2: Launch tournament
	err = mockClient.LaunchTournament(resp.TournamentID)
	if err != nil {
		t.Fatalf("Failed to launch tournament: %v", err)
	}

	// Step 3: Invite players (using placeholder player IDs for now)
	err = mockClient.InvitePlayer(resp.TournamentID, "herchu_player_id")
	if err != nil {
		t.Fatalf("Failed to invite first player: %v", err)
	}

	err = mockClient.InvitePlayer(resp.TournamentID, "lord_trooper_player_id")
	if err != nil {
		t.Fatalf("Failed to invite second player: %v", err)
	}

	// Verify tournament is in launched state
	status, err := mockClient.GetTournamentStatus(resp.TournamentID)
	if err != nil {
		t.Fatalf("Failed to get tournament status: %v", err)
	}

	// Tournament should be in "launched" or "open" state after launching
	if status.Status == "created" {
		t.Error("Expected tournament to be launched, but it's still in created state")
	}
}

func TestMockClient_CompleteThreeStepWorkflowDemo(t *testing.T) {
	t.Log("=== Complete Three-Step Tournament Creation Workflow Demo ===")

	mockClient := NewMockClient("testuser", "testpass")

	// Authentication
	t.Log("Step 0: Authenticating...")
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}
	t.Log("✓ Authentication successful")

	// Step 1: Create Tournament
	t.Log("Step 1: Creating tournament...")
	division := "Elite"
	homePlayer := "herchu"
	awayPlayer := "Lord Trooper"
	roundNum := 1
	matchNum := 15

	resp, err := mockClient.CreateSwissTournament(division, homePlayer, awayPlayer, roundNum, matchNum)
	if err != nil {
		t.Fatalf("Failed to create tournament: %v", err)
	}

	t.Logf("✓ Tournament created successfully")
	t.Logf("  - Tournament ID: %d", resp.TournamentID)
	t.Logf("  - Tournament Link: %s", resp.Link)

	// Verify initial tournament state
	status, err := mockClient.GetTournamentStatus(resp.TournamentID)
	if err != nil {
		t.Fatalf("Failed to get tournament status: %v", err)
	}
	t.Logf("  - Initial Status: %s", status.Status)

	// Step 2: Launch Tournament
	t.Log("Step 2: Launching tournament...")
	err = mockClient.LaunchTournament(resp.TournamentID)
	if err != nil {
		t.Fatalf("Failed to launch tournament: %v", err)
	}
	t.Log("✓ Tournament launched successfully")

	// Verify tournament is now open
	status, err = mockClient.GetTournamentStatus(resp.TournamentID)
	if err != nil {
		t.Fatalf("Failed to get tournament status after launch: %v", err)
	}
	t.Logf("  - Status after launch: %s", status.Status)

	if status.Status != "open" {
		t.Errorf("Expected tournament status to be 'open', got '%s'", status.Status)
	}

	// Step 3: Invite Players
	t.Log("Step 3: Inviting players...")
	playerIDs := []string{"herchu_bga_id", "lord_trooper_bga_id"}

	for i, playerID := range playerIDs {
		t.Logf("  Inviting player %d: %s", i+1, playerID)
		err = mockClient.InvitePlayer(resp.TournamentID, playerID)
		if err != nil {
			t.Fatalf("Failed to invite player %s: %v", playerID, err)
		}
	}
	t.Log("✓ All players invited successfully")

	// Final verification
	t.Log("Final verification...")
	finalStatus, err := mockClient.GetTournamentStatus(resp.TournamentID)
	if err != nil {
		t.Fatalf("Failed to get final tournament status: %v", err)
	}

	t.Logf("✓ Tournament workflow completed successfully")
	t.Logf("  - Final Tournament ID: %d", finalStatus.ID)
	t.Logf("  - Final Status: %s", finalStatus.Status)
	t.Logf("  - Tournament Name: %s", finalStatus.Name)
	t.Logf("  - Players Count: %d", finalStatus.PlayersCount)

	// Validate expected tournament name format
	expectedName := "1 Fecha - Duelo 15 - herchu vs Lord Trooper"
	if finalStatus.Name != expectedName {
		t.Errorf("Expected tournament name '%s', got '%s'", expectedName, finalStatus.Name)
	}

	t.Log("=== Three-Step Tournament Creation Demo Complete ===")
}
