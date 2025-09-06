package cli

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewDateTimePickerModel(t *testing.T) {
	homePlayer := "player1"
	awayPlayer := "player2"
	division := "Elite"
	roundNumber := 1
	matchNumber := 15
	matchID := 15

	picker := NewDateTimePickerModel(homePlayer, awayPlayer, division, roundNumber, matchNumber, matchID)

	if picker.homePlayer != homePlayer {
		t.Errorf("Expected homePlayer %s, got %s", homePlayer, picker.homePlayer)
	}

	if picker.awayPlayer != awayPlayer {
		t.Errorf("Expected awayPlayer %s, got %s", awayPlayer, picker.awayPlayer)
	}

	if picker.division != division {
		t.Errorf("Expected division %s, got %s", division, picker.division)
	}

	if picker.roundNumber != roundNumber {
		t.Errorf("Expected roundNumber %d, got %d", roundNumber, picker.roundNumber)
	}

	if picker.matchNumber != matchNumber {
		t.Errorf("Expected matchNumber %d, got %d", matchNumber, picker.matchNumber)
	}

	if picker.matchID != matchID {
		t.Errorf("Expected matchID %d, got %d", matchID, picker.matchID)
	}

	expectedTitle := "Schedule Tournament: player1 vs player2"
	if picker.title != expectedTitle {
		t.Errorf("Expected title %s, got %s", expectedTitle, picker.title)
	}

	if picker.timezone == nil {
		t.Error("Expected timezone to be set")
	}

	if picker.confirmed {
		t.Error("Expected picker to not be confirmed initially")
	}

	if picker.canceled {
		t.Error("Expected picker to not be canceled initially")
	}
}

func TestDateTimePickerModel_Init(t *testing.T) {
	picker := NewDateTimePickerModel("player1", "player2", "Elite", 1, 15, 15)

	cmd := picker.Init()
	// The underlying picker might not always return a command
	_ = cmd // Accept that cmd might be nil
}

func TestDateTimePickerModel_Update_Enter(t *testing.T) {
	picker := NewDateTimePickerModel("player1", "player2", "Elite", 1, 15, 15)

	// 1. First Enter: Should switch focus in the picker.
	// The picker's Update returns a model and a command.
	// Our wrapper should not treat this as a final confirmation.
	model, cmd := picker.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(dateTimePickerConfirmedMsg); ok {
			t.Fatal("Should not be confirmed on first enter")
		}
	}

	// 2. Second Enter: This should trigger confirmation.
	model, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Expected command on second Enter")
	}

	// The command should be our wrapped command. Executing it gives our internal message.
	msg := cmd()
	confirmMsg, ok := msg.(dateTimePickerConfirmedMsg)
	if !ok {
		t.Fatalf("Expected dateTimePickerConfirmedMsg, got %T", msg)
	}

	// 3. Process the internal confirmation message
	model, cmd = model.Update(confirmMsg)
	if cmd == nil {
		t.Fatal("Expected command after confirmation message")
	}

	// This final command should contain the DateTimeSelectedMsg
	finalMsg := cmd()
	if dateTimeMsg, ok := finalMsg.(DateTimeSelectedMsg); ok {
		if dateTimeMsg.HomePlayer != "player1" {
			t.Errorf("Expected HomePlayer 'player1', got %s", dateTimeMsg.HomePlayer)
		}
		if dateTimeMsg.AwayPlayer != "player2" {
			t.Errorf("Expected AwayPlayer 'player2', got %s", dateTimeMsg.AwayPlayer)
		}
		if dateTimeMsg.Division != "Elite" {
			t.Errorf("Expected Division 'Elite', got %s", dateTimeMsg.Division)
		}
		if dateTimeMsg.RoundNumber != 1 {
			t.Errorf("Expected RoundNumber 1, got %d", dateTimeMsg.RoundNumber)
		}
		if dateTimeMsg.MatchNumber != 15 {
			t.Errorf("Expected MatchNumber 15, got %d", dateTimeMsg.MatchNumber)
		}
		if dateTimeMsg.MatchID != 15 {
			t.Errorf("Expected MatchID 15, got %d", dateTimeMsg.MatchID)
		}
		if dateTimeMsg.DateTime.IsZero() {
			t.Error("Expected DateTime to be set")
		}
	} else {
		t.Errorf("Expected DateTimeSelectedMsg, got %T", finalMsg)
	}

	// Check that the picker is marked as confirmed
	if pickerModel, ok := model.(*DateTimePickerModel); ok {
		if !pickerModel.IsConfirmed() {
			t.Error("Expected picker to be confirmed after processing confirmation")
		}
	}
}

func TestDateTimePickerModel_Update_Escape(t *testing.T) {
	picker := NewDateTimePickerModel("player1", "player2", "Elite", 1, 15, 15)

	// Send escape key
	updatedModel, cmd := picker.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Fatal("Expected command when pressing Escape")
	}

	// Execute the command to get the message
	msg := cmd()
	if _, ok := msg.(DateTimePickerCanceledMsg); !ok {
		t.Errorf("Expected DateTimePickerCanceledMsg, got %T", msg)
	}

	// Check that the picker is marked as canceled
	if pickerModel, ok := updatedModel.(*DateTimePickerModel); ok {
		if !pickerModel.IsCanceled() {
			t.Error("Expected picker to be canceled after pressing Escape")
		}
	}
}

func TestDateTimePickerModel_View(t *testing.T) {
	picker := NewDateTimePickerModel("herchu", "Lord Trooper", "Elite", 1, 15, 15)

	view := picker.View()

	// Check that the view contains expected elements
	expectedStrings := []string{
		"Schedule Tournament: herchu vs Lord Trooper",
		"Division: Elite - Round 1 - Duelo 15",
		"Navigation:",
		"↑/↓ arrows: Change date or hour values",
		"←/→ arrows: Move between date and time fields",
		"Tab/Space: Switch between components",
		"Enter: Confirm selection | Esc: Cancel",
		"Selected:",
		"UTC",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("Expected view to contain '%s', but it didn't. View: %s", expected, view)
		}
	}
}

func TestDateTimePickerModel_View_ConfirmedOrCanceled(t *testing.T) {
	picker := NewDateTimePickerModel("player1", "player2", "Elite", 1, 15, 15)

	// Test confirmed state
	picker.confirmed = true
	view := picker.View()
	if view != "" {
		t.Errorf("Expected empty view when confirmed, got: %s", view)
	}

	// Reset and test canceled state
	picker.confirmed = false
	picker.canceled = true
	view = picker.View()
	if view != "" {
		t.Errorf("Expected empty view when canceled, got: %s", view)
	}
}

func TestDateTimePickerModel_FormatForBGA(t *testing.T) {
	picker := NewDateTimePickerModel("player1", "player2", "Elite", 1, 15, 15)

	// Test with zero time (should return defaults)
	date, timeStr := picker.FormatForBGA()
	if timeStr != "21:00" {
		t.Errorf("Expected default time '21:00', got %s", timeStr)
	}
	// Date should be current date in YYYY-MM-DD format
	if len(date) != 10 || date[4] != '-' || date[7] != '-' {
		t.Errorf("Expected date in YYYY-MM-DD format, got %s", date)
	}

	// Test with specific time
	specificTime := time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local)
	picker.selectedTime = specificTime
	date, timeStr = picker.FormatForBGA()

	expectedDate := "2025-03-15"
	expectedTime := "14:30"

	if date != expectedDate {
		t.Errorf("Expected date %s, got %s", expectedDate, date)
	}
	if timeStr != expectedTime {
		t.Errorf("Expected time %s, got %s", expectedTime, timeStr)
	}
}

func TestDateTimePickerModel_GetSelectedTime(t *testing.T) {
	picker := NewDateTimePickerModel("player1", "player2", "Elite", 1, 15, 15)

	// Initially should return time from picker (not necessarily zero)
	selectedTime := picker.GetSelectedTime()
	if selectedTime.IsZero() {
		t.Error("Expected picker to have a default time set")
	}

	// Set a specific time
	specificTime := time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local)
	picker.selectedTime = specificTime

	selectedTime = picker.GetSelectedTime()
	if !selectedTime.Equal(specificTime) {
		t.Errorf("Expected selected time %v, got %v", specificTime, selectedTime)
	}
}

func TestDateTimePickerModel_TimezoneDisplay(t *testing.T) {
	picker := NewDateTimePickerModel("player1", "player2", "Elite", 1, 15, 15)

	view := picker.View()

	// Should contain UTC offset information
	if !strings.Contains(view, "UTC") {
		t.Error("Expected view to contain UTC offset information")
	}

	// Should contain timezone name
	if picker.timezone != nil && !strings.Contains(view, "Timezone:") {
		t.Error("Expected view to contain timezone information")
	}
}

func TestDateTimePickerModel_PlatinumDivision(t *testing.T) {
	picker := NewDateTimePickerModel("webbi", "alehrosario", "Platinum A", 5, 23, 23)

	view := picker.View()

	expectedTitle := "Schedule Tournament: webbi vs alehrosario"
	if !strings.Contains(view, expectedTitle) {
		t.Errorf("Expected view to contain '%s'", expectedTitle)
	}

	expectedDivision := "Division: Platinum A - Round 5 - Duelo 23"
	if !strings.Contains(view, expectedDivision) {
		t.Errorf("Expected view to contain '%s'", expectedDivision)
	}
}

func TestDateTimePickerModel_StateMethods(t *testing.T) {
	picker := NewDateTimePickerModel("player1", "player2", "Elite", 1, 15, 15)

	// Initially not confirmed or canceled
	if picker.IsConfirmed() {
		t.Error("Expected picker to not be confirmed initially")
	}
	if picker.IsCanceled() {
		t.Error("Expected picker to not be canceled initially")
	}

	// Set confirmed
	picker.confirmed = true
	if !picker.IsConfirmed() {
		t.Error("Expected picker to be confirmed")
	}

	// Reset and set canceled
	picker.confirmed = false
	picker.canceled = true
	if !picker.IsCanceled() {
		t.Error("Expected picker to be canceled")
	}
}

func TestDateTimePickerModel_KeyHandling(t *testing.T) {
	picker := NewDateTimePickerModel("player1", "player2", "Elite", 1, 15, 15)

	// Test that arrow keys are passed through to the picker
	testKeys := []tea.KeyMsg{
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
		{Type: tea.KeyLeft},
		{Type: tea.KeyRight},
		{Type: tea.KeyTab},
		{Type: tea.KeyRunes, Runes: []rune{' '}}, // Space key
	}

	for _, keyMsg := range testKeys {
		updatedModel, cmd := picker.Update(keyMsg)

		// Should not return any special commands for navigation keys
		if cmd != nil {
			t.Errorf("Expected no command for navigation key %v, got %v", keyMsg, cmd)
		}

		// Should return the same model type
		if _, ok := updatedModel.(*DateTimePickerModel); !ok {
			t.Errorf("Expected DateTimePickerModel, got %T", updatedModel)
		}
	}
}

func TestDateTimePickerModel_NavigationInstructions(t *testing.T) {
	picker := NewDateTimePickerModel("player1", "player2", "Elite", 1, 15, 15)

	view := picker.View()

	// Check for navigation instructions
	expectedInstructions := []string{
		"↑/↓ arrows: Change date or hour values",
		"←/→ arrows: Move between date and time fields",
		"Tab/Space: Switch between components",
		"Enter: Confirm selection",
		"Esc: Cancel",
	}

	for _, instruction := range expectedInstructions {
		if !strings.Contains(view, instruction) {
			t.Errorf("Expected view to contain instruction '%s', but it didn't", instruction)
		}
	}
}
