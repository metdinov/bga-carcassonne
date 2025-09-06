package cli

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewTournamentConfirmationModel(t *testing.T) {
	homePlayer := "herchu"
	awayPlayer := "Lord Trooper"
	division := "Elite"
	roundNumber := 1
	matchNumber := 15
	matchID := 15
	selectedTime := time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local)

	model := NewTournamentConfirmationModel(homePlayer, awayPlayer, division, roundNumber, matchNumber, matchID, selectedTime)

	if model.homePlayer != homePlayer {
		t.Errorf("Expected homePlayer %s, got %s", homePlayer, model.homePlayer)
	}

	if model.awayPlayer != awayPlayer {
		t.Errorf("Expected awayPlayer %s, got %s", awayPlayer, model.awayPlayer)
	}

	if model.division != division {
		t.Errorf("Expected division %s, got %s", division, model.division)
	}

	if model.roundNumber != roundNumber {
		t.Errorf("Expected roundNumber %d, got %d", roundNumber, model.roundNumber)
	}

	if model.matchNumber != matchNumber {
		t.Errorf("Expected matchNumber %d, got %d", matchNumber, model.matchNumber)
	}

	if model.matchID != matchID {
		t.Errorf("Expected matchID %d, got %d", matchID, model.matchID)
	}

	if !model.selectedTime.Equal(selectedTime) {
		t.Errorf("Expected selectedTime %v, got %v", selectedTime, model.selectedTime)
	}

	expectedChampionship := "Division Elite - 1era Temporada"
	if model.championshipName != expectedChampionship {
		t.Errorf("Expected championshipName %s, got %s", expectedChampionship, model.championshipName)
	}

	expectedTournament := "1 Fecha - Duelo 15 - herchu vs Lord Trooper"
	if model.tournamentName != expectedTournament {
		t.Errorf("Expected tournamentName %s, got %s", expectedTournament, model.tournamentName)
	}

	if model.confirmed {
		t.Error("Expected model to not be confirmed initially")
	}

	if model.canceled {
		t.Error("Expected model to not be canceled initially")
	}
}

func TestTournamentConfirmationModel_Init(t *testing.T) {
	model := NewTournamentConfirmationModel("player1", "player2", "Elite", 1, 15, 15, time.Now())

	cmd := model.Init()
	if cmd != nil {
		t.Error("Expected Init to return nil command")
	}
}

func TestTournamentConfirmationModel_Update_Enter(t *testing.T) {
	selectedTime := time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local)
	model := NewTournamentConfirmationModel("herchu", "Lord Trooper", "Elite", 1, 15, 15, selectedTime)

	// Send enter key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("Expected command when pressing Enter")
	}

	// Execute the command to get the message
	msg := cmd()
	if confirmMsg, ok := msg.(TournamentConfirmedMsg); ok {
		if confirmMsg.HomePlayer != "herchu" {
			t.Errorf("Expected HomePlayer 'herchu', got %s", confirmMsg.HomePlayer)
		}
		if confirmMsg.AwayPlayer != "Lord Trooper" {
			t.Errorf("Expected AwayPlayer 'Lord Trooper', got %s", confirmMsg.AwayPlayer)
		}
		if confirmMsg.Division != "Elite" {
			t.Errorf("Expected Division 'Elite', got %s", confirmMsg.Division)
		}
		if confirmMsg.RoundNumber != 1 {
			t.Errorf("Expected RoundNumber 1, got %d", confirmMsg.RoundNumber)
		}
		if confirmMsg.MatchNumber != 15 {
			t.Errorf("Expected MatchNumber 15, got %d", confirmMsg.MatchNumber)
		}
		if confirmMsg.MatchID != 15 {
			t.Errorf("Expected MatchID 15, got %d", confirmMsg.MatchID)
		}
		if !confirmMsg.DateTime.Equal(selectedTime) {
			t.Errorf("Expected DateTime %v, got %v", selectedTime, confirmMsg.DateTime)
		}
	} else {
		t.Errorf("Expected TournamentConfirmedMsg, got %T", msg)
	}

	// Check that the model is marked as confirmed
	if confirmationModel, ok := updatedModel.(*TournamentConfirmationModel); ok {
		if !confirmationModel.IsConfirmed() {
			t.Error("Expected model to be confirmed after pressing Enter")
		}
	}
}

func TestTournamentConfirmationModel_Update_Escape(t *testing.T) {
	model := NewTournamentConfirmationModel("player1", "player2", "Elite", 1, 15, 15, time.Now())

	// Send escape key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Fatal("Expected command when pressing Escape")
	}

	// Execute the command to get the message
	msg := cmd()
	if _, ok := msg.(TournamentConfirmationCanceledMsg); !ok {
		t.Errorf("Expected TournamentConfirmationCanceledMsg, got %T", msg)
	}

	// Check that the model is marked as canceled
	if confirmationModel, ok := updatedModel.(*TournamentConfirmationModel); ok {
		if !confirmationModel.IsCanceled() {
			t.Error("Expected model to be canceled after pressing Escape")
		}
	}
}

func TestTournamentConfirmationModel_View(t *testing.T) {
	selectedTime := time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local)
	model := NewTournamentConfirmationModel("herchu", "Lord Trooper", "Elite", 1, 15, 15, selectedTime)

	view := model.View()

	// Check that the view contains expected elements
	expectedStrings := []string{
		"Tournament Confirmation",
		"Division Elite - 1era Temporada",
		"1 Fecha - Duelo 15 - herchu vs Lord Trooper",
		"Division:     Elite",
		"Round:        1",
		"Match (Duelo): 15",
		"Players:      herchu vs Lord Trooper",
		"Date & Time:",
		"Swiss System (Best-of-3)",
		"Carcassonne",
		"30 minutes",
		"International scoring",
		"3 points per city",
		"4 points per two tile city",
		"Press Enter to create tournament",
		"Press Esc to cancel",
		"UTC",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("Expected view to contain '%s', but it didn't. View: %s", expected, view)
		}
	}
}

func TestTournamentConfirmationModel_View_ConfirmedOrCanceled(t *testing.T) {
	model := NewTournamentConfirmationModel("player1", "player2", "Elite", 1, 15, 15, time.Now())

	// Test confirmed state
	model.confirmed = true
	view := model.View()
	if view != "" {
		t.Errorf("Expected empty view when confirmed, got: %s", view)
	}

	// Reset and test canceled state
	model.confirmed = false
	model.canceled = true
	view = model.View()
	if view != "" {
		t.Errorf("Expected empty view when canceled, got: %s", view)
	}
}

func TestTournamentConfirmationModel_GetTournamentDetails(t *testing.T) {
	selectedTime := time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local)
	model := NewTournamentConfirmationModel("herchu", "Lord Trooper", "Elite", 1, 15, 15, selectedTime)

	championshipName, tournamentName := model.GetTournamentDetails()

	expectedChampionship := "Division Elite - 1era Temporada"
	if championshipName != expectedChampionship {
		t.Errorf("Expected championship name '%s', got '%s'", expectedChampionship, championshipName)
	}

	expectedTournament := "1 Fecha - Duelo 15 - herchu vs Lord Trooper"
	if tournamentName != expectedTournament {
		t.Errorf("Expected tournament name '%s', got '%s'", expectedTournament, tournamentName)
	}
}

func TestTournamentConfirmationModel_GetSchedulingInfo(t *testing.T) {
	selectedTime := time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local)
	model := NewTournamentConfirmationModel("player1", "player2", "Elite", 1, 15, 15, selectedTime)

	dateStr, timeStr := model.GetSchedulingInfo()

	expectedDate := "2025-03-15"
	expectedTime := "14:30"

	if dateStr != expectedDate {
		t.Errorf("Expected date '%s', got '%s'", expectedDate, dateStr)
	}

	if timeStr != expectedTime {
		t.Errorf("Expected time '%s', got '%s'", expectedTime, timeStr)
	}
}

func TestTournamentConfirmationModel_StateMethods(t *testing.T) {
	model := NewTournamentConfirmationModel("player1", "player2", "Elite", 1, 15, 15, time.Now())

	// Initially not confirmed or canceled
	if model.IsConfirmed() {
		t.Error("Expected model to not be confirmed initially")
	}
	if model.IsCanceled() {
		t.Error("Expected model to not be canceled initially")
	}

	// Set confirmed
	model.confirmed = true
	if !model.IsConfirmed() {
		t.Error("Expected model to be confirmed")
	}

	// Reset and set canceled
	model.confirmed = false
	model.canceled = true
	if !model.IsCanceled() {
		t.Error("Expected model to be canceled")
	}
}

func TestTournamentConfirmationModel_PlatinumDivision(t *testing.T) {
	selectedTime := time.Date(2025, 5, 20, 18, 45, 0, 0, time.Local)
	model := NewTournamentConfirmationModel("webbi", "alehrosario", "Platinum A", 3, 27, 27, selectedTime)

	view := model.View()

	expectedStrings := []string{
		"Division Platinum A - 1era Temporada",
		"3 Fecha - Duelo 27 - webbi vs alehrosario",
		"Division:     Platinum A",
		"Round:        3",
		"Match (Duelo): 27",
		"Players:      webbi vs alehrosario",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("Expected view to contain '%s'", expected)
		}
	}
}

func TestTournamentConfirmationModel_OroDivision(t *testing.T) {
	selectedTime := time.Date(2025, 6, 10, 21, 15, 0, 0, time.Local)
	model := NewTournamentConfirmationModel("bignacho610", "Academia47", "Oro B", 5, 42, 42, selectedTime)

	championshipName, tournamentName := model.GetTournamentDetails()

	expectedChampionship := "Division Oro B - 1era Temporada"
	if championshipName != expectedChampionship {
		t.Errorf("Expected championship name '%s', got '%s'", expectedChampionship, championshipName)
	}

	expectedTournament := "5 Fecha - Duelo 42 - bignacho610 vs Academia47"
	if tournamentName != expectedTournament {
		t.Errorf("Expected tournament name '%s', got '%s'", expectedTournament, tournamentName)
	}
}

func TestTournamentConfirmationModel_TimezoneDisplay(t *testing.T) {
	selectedTime := time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local)
	model := NewTournamentConfirmationModel("player1", "player2", "Elite", 1, 15, 15, selectedTime)

	view := model.View()

	// Should contain UTC offset information
	if !strings.Contains(view, "UTC") {
		t.Error("Expected view to contain UTC offset information")
	}

	// Should contain timezone information
	if !strings.Contains(view, "Timezone:") {
		t.Error("Expected view to contain timezone information")
	}

	// Should contain formatted date and time
	if !strings.Contains(view, "Date & Time:") {
		t.Error("Expected view to contain date & time information")
	}
}
