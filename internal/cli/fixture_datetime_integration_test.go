package cli

import (
	"strings"
	"testing"

	"carca-cli/internal/bga"
	"carca-cli/internal/fixtures"

	tea "github.com/charmbracelet/bubbletea"
)

func TestFixtureModel_DateTimePickerIntegration(t *testing.T) {
	// Create a division with unplayed matches
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         15,
						HomePlayer: "herchu",
						AwayPlayer: "Lord Trooper",
						BGALink:    "",
						Played:     false,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)

	// Set up mock BGA client
	mockClient := bga.NewMockClient("testuser", "testpass")
	model.SetBGAClient(mockClient)
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login mock client: %v", err)
	}

	// Step 1: Press 'c' to trigger datetime picker
	model.selectedMatch = 0
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	fixtureModel := updatedModel.(*FixtureModel)
	if !fixtureModel.showDatePicker {
		t.Fatal("Expected datetime picker to be shown")
	}

	if fixtureModel.dateTimePicker == nil {
		t.Fatal("Expected datetime picker to be created")
	}

	// Verify datetime picker properties
	picker := fixtureModel.dateTimePicker
	if picker.homePlayer != "herchu" {
		t.Errorf("Expected homePlayer 'herchu', got %s", picker.homePlayer)
	}
	if picker.awayPlayer != "Lord Trooper" {
		t.Errorf("Expected awayPlayer 'Lord Trooper', got %s", picker.awayPlayer)
	}
	if picker.division != "Elite" {
		t.Errorf("Expected division 'Elite', got %s", picker.division)
	}
	if picker.roundNumber != 1 {
		t.Errorf("Expected roundNumber 1, got %d", picker.roundNumber)
	}
	if picker.matchNumber != 15 {
		t.Errorf("Expected matchNumber 15, got %d", picker.matchNumber)
	}

	// Step 2: Verify datetime picker view contains expected elements
	view := picker.View()
	expectedElements := []string{
		"Schedule Tournament: herchu vs Lord Trooper",
		"Division: Elite - Round 1 - Duelo 15",
		"UTC",
		"Selected:",
		"Navigation:",
		"↑/↓ arrows: Change date or hour values",
		"←/→ arrows: Move between date and time fields",
		"Tab/Space: Switch between components",
		"Enter: Confirm selection | Esc: Cancel",
	}

	for _, element := range expectedElements {
		if !strings.Contains(view, element) {
			t.Errorf("Expected datetime picker view to contain '%s', view: %s", element, view)
		}
	}

	// Step 3: Simulate Enter key to confirm datetime selection
	// First Enter to select date
	updatedModel, _ = fixtureModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	fixtureModel = updatedModel.(*FixtureModel)

	// Second Enter to select time, which returns a wrapped command
	updatedModel, cmd2 := fixtureModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	fixtureModel = updatedModel.(*FixtureModel)
	if cmd2 == nil {
		t.Fatal("Expected command from datetime picker confirmation")
	}

	// Execute wrapped command to get internal confirmation message
	internalMsg := cmd2()

	// Process internal message to get the final DateTimeSelectedMsg
	updatedModel, cmd3 := fixtureModel.Update(internalMsg)
	fixtureModel = updatedModel.(*FixtureModel)
	if cmd3 == nil {
		t.Fatal("Expected command after processing internal confirmation")
	}

	msg := cmd3()
	dateTimeMsg, ok := msg.(DateTimeSelectedMsg)
	if !ok {
		t.Fatalf("Expected DateTimeSelectedMsg, got %T", msg)
	}

	// Verify datetime selection message
	if dateTimeMsg.HomePlayer != "herchu" {
		t.Errorf("Expected HomePlayer 'herchu', got %s", dateTimeMsg.HomePlayer)
	}
	if dateTimeMsg.AwayPlayer != "Lord Trooper" {
		t.Errorf("Expected AwayPlayer 'Lord Trooper', got %s", dateTimeMsg.AwayPlayer)
	}
	if dateTimeMsg.Division != "Elite" {
		t.Errorf("Expected Division 'Elite', got %s", dateTimeMsg.Division)
	}
	if dateTimeMsg.DateTime.IsZero() {
		t.Error("Expected DateTime to be set")
	}

	// Step 4: Process datetime selection message to trigger tournament creation
	// Process the datetime selection to show confirmation screen
	updatedModel, _ = fixtureModel.Update(dateTimeMsg)
	fixtureModel = updatedModel.(*FixtureModel)

	// Should hide datetime picker and show confirmation screen
	if fixtureModel.showDatePicker {
		t.Error("Expected datetime picker to be hidden after selection")
	}
	if !fixtureModel.showConfirmation {
		t.Fatal("Expected confirmation screen to be shown")
	}
	if fixtureModel.confirmationModel == nil {
		t.Fatal("Expected confirmation model to be created")
	}

	// Step 5: Verify confirmation screen content
	confirmationView := fixtureModel.confirmationModel.View()
	expectedConfirmationElements := []string{
		"Tournament Confirmation",
		"Division Elite - 1era Temporada",
		"1 Fecha - Duelo 15 - herchu vs Lord Trooper",
		"herchu vs Lord Trooper",
		"Swiss System (Best-of-3)",
		"Press Enter to create tournament",
	}

	for _, element := range expectedConfirmationElements {
		if !strings.Contains(confirmationView, element) {
			t.Errorf("Expected confirmation view to contain '%s'", element)
		}
	}

	// Step 6: Simulate Enter key on confirmation screen
	_, cmd4 := fixtureModel.confirmationModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd4 == nil {
		t.Fatal("Expected command from confirmation screen")
	}

	confirmMsg := cmd4()
	tournamentConfirmedMsg, ok := confirmMsg.(TournamentConfirmedMsg)
	if !ok {
		t.Fatalf("Expected TournamentConfirmedMsg, got %T", confirmMsg)
	}

	// Step 7: Process tournament confirmation to trigger creation
	intermediateModel, cmd := fixtureModel.Update(tournamentConfirmedMsg)
	if cmd == nil {
		t.Fatal("Expected tournament creation command after confirmation")
	}

	// Should hide confirmation and show creation message
	intermediateFixtureModel := intermediateModel.(*FixtureModel)
	if intermediateFixtureModel.showConfirmation {
		t.Error("Expected confirmation screen to be hidden after confirmation")
	}
	if !strings.Contains(intermediateFixtureModel.statusMessage, "Creating tournament for herchu vs Lord Trooper") {
		t.Errorf("Expected creation status message, got: %s", intermediateFixtureModel.statusMessage)
	}

	// Step 8: Execute tournament creation command
	creationMsg := cmd()
	createWithDateTimeMsg, ok := creationMsg.(createTournamentMsgWithDateTime)
	if !ok {
		t.Fatalf("Expected createTournamentMsgWithDateTime, got %T", creationMsg)
	}

	// Step 9: Process tournament creation with datetime
	finalIntermediateModel, cmd := intermediateFixtureModel.Update(createWithDateTimeMsg)
	if cmd == nil {
		t.Fatal("Expected async tournament creation command")
	}

	// Step 10: Execute async tournament creation
	asyncMsg := cmd()
	tournamentMsg, ok := asyncMsg.(tournamentCreatedMsg)
	if !ok {
		t.Fatalf("Expected tournamentCreatedMsg from async creation, got %T", asyncMsg)
	}

	// Step 11: Process tournament creation completion
	finalModel, _ := finalIntermediateModel.Update(tournamentMsg)

	if tournamentMsg.success {
		finalFixtureModel := finalModel.(*FixtureModel)
		match := finalFixtureModel.division.Rounds[0].Matches[0]

		if match.BGALink == "" {
			t.Error("Expected match to have BGALink set after tournament creation")
		}

		if !strings.Contains(match.BGALink, "tournament?id=") {
			t.Errorf("Expected valid tournament link, got %s", match.BGALink)
		}

		if !strings.Contains(finalFixtureModel.statusMessage, "Tournament created successfully") {
			t.Errorf("Expected success message, got: %s", finalFixtureModel.statusMessage)
		}
	} else {
		t.Errorf("Tournament creation failed: %s", tournamentMsg.error)
	}
}

func TestFixtureModel_DateTimePickerCancellation(t *testing.T) {
	// Create test division
	division := &fixtures.Division{
		Name: "Platinum A",
		Rounds: []*fixtures.Round{
			{
				Number: 1,
				Matches: []*fixtures.Match{
					{
						ID:         23,
						HomePlayer: "webbi",
						AwayPlayer: "alehrosario",
						BGALink:    "",
						Played:     false,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	mockClient := bga.NewMockClient("testuser", "testpass")
	model.SetBGAClient(mockClient)

	// Step 1: Show datetime picker
	model.selectedMatch = 0
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	fixtureModel := updatedModel.(*FixtureModel)
	if !fixtureModel.showDatePicker {
		t.Fatal("Expected datetime picker to be shown")
	}

	// Step 2: Show datetime picker and select time
	// First Enter to select date
	updatedModel, _ = fixtureModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	fixtureModel = updatedModel.(*FixtureModel)

	// Second Enter to select time, which returns a wrapped command
	updatedModel, cmd := fixtureModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	fixtureModel = updatedModel.(*FixtureModel)
	if cmd == nil {
		t.Fatal("Expected command from datetime picker confirmation")
	}

	// Execute wrapped command to get internal confirmation message
	internalMsg := cmd()

	// Process internal message to get the final DateTimeSelectedMsg
	updatedModel, cmd = fixtureModel.Update(internalMsg)
	fixtureModel = updatedModel.(*FixtureModel)
	if cmd == nil {
		t.Fatal("Expected command after processing internal confirmation")
	}

	msg := cmd()
	dateTimeMsg, ok := msg.(DateTimeSelectedMsg)
	if !ok {
		t.Fatalf("Expected DateTimeSelectedMsg, got %T", msg)
	}

	// Step 3: Process datetime selection to show confirmation screen
	confirmationModel, _ := fixtureModel.Update(dateTimeMsg)
	fixtureModel = confirmationModel.(*FixtureModel)

	if !fixtureModel.showConfirmation {
		t.Fatal("Expected confirmation screen to be shown")
	}

	// Step 4: Press Esc on confirmation screen to cancel
	confirmation := fixtureModel.confirmationModel
	_, cmd = confirmation.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("Expected command from confirmation screen cancellation")
	}

	cancelMsg := cmd()
	if _, ok := cancelMsg.(TournamentConfirmationCanceledMsg); !ok {
		t.Fatalf("Expected TournamentConfirmationCanceledMsg, got %T", cancelMsg)
	}

	// Step 5: Process cancellation message
	finalModel, _ := fixtureModel.Update(cancelMsg)

	finalFixtureModel := finalModel.(*FixtureModel)
	if finalFixtureModel.showConfirmation {
		t.Error("Expected confirmation screen to be hidden after cancellation")
	}

	if !strings.Contains(finalFixtureModel.statusMessage, "canceled") {
		t.Errorf("Expected cancellation message, got: %s", finalFixtureModel.statusMessage)
	}

	// Verify match wasn't modified
	match := finalFixtureModel.division.Rounds[0].Matches[0]
	if match.BGALink != "" {
		t.Error("Expected match BGALink to remain empty after cancellation")
	}
}

func TestFixtureModel_DateTimePickerWithDifferentDivisions(t *testing.T) {
	testCases := []struct {
		name       string
		division   string
		homePlayer string
		awayPlayer string
		matchID    int
	}{
		{
			name:       "Elite Division",
			division:   "Elite",
			homePlayer: "herchu",
			awayPlayer: "Lord Trooper",
			matchID:    15,
		},
		{
			name:       "Platinum A Division",
			division:   "Platinum A",
			homePlayer: "webbi",
			awayPlayer: "alehrosario",
			matchID:    23,
		},
		{
			name:       "Oro B Division",
			division:   "Oro B",
			homePlayer: "bignacho610",
			awayPlayer: "Academia47",
			matchID:    8,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			division := &fixtures.Division{
				Name: tc.division,
				Rounds: []*fixtures.Round{
					{
						Number: 1,
						Matches: []*fixtures.Match{
							{
								ID:         tc.matchID,
								HomePlayer: tc.homePlayer,
								AwayPlayer: tc.awayPlayer,
								BGALink:    "",
								Played:     false,
							},
						},
					},
				},
			}

			model := NewFixtureModel(division)
			mockClient := bga.NewMockClient("testuser", "testpass")
			model.SetBGAClient(mockClient)

			// Show datetime picker
			model.selectedMatch = 0
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			fixtureModel := updatedModel.(*FixtureModel)
			picker := fixtureModel.dateTimePicker

			// Verify picker has correct division info
			expectedTitle := "Schedule Tournament: " + tc.homePlayer + " vs " + tc.awayPlayer
			if picker.title != expectedTitle {
				t.Errorf("Expected title '%s', got '%s'", expectedTitle, picker.title)
			}

			if picker.division != tc.division {
				t.Errorf("Expected division '%s', got '%s'", tc.division, picker.division)
			}

			// Verify view contains division-specific information
			view := picker.View()
			if !strings.Contains(view, tc.division) {
				t.Errorf("Expected view to contain division '%s'", tc.division)
			}
		})
	}
}

func TestFixtureModel_DateTimePickerTimezoneDisplay(t *testing.T) {
	// Create test setup
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number: 1,
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "player1",
						AwayPlayer: "player2",
						BGALink:    "",
						Played:     false,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	model.selectedMatch = 0

	// Show datetime picker
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	fixtureModel := updatedModel.(*FixtureModel)
	picker := fixtureModel.dateTimePicker

	view := picker.View()

	// Should display UTC offset
	if !strings.Contains(view, "UTC") {
		t.Error("Expected datetime picker view to show UTC offset")
	}

	// Should display timezone information
	if !strings.Contains(view, "Timezone:") {
		t.Error("Expected datetime picker view to show timezone information")
	}

	// Should display selected time with offset
	if !strings.Contains(view, "Selected:") {
		t.Error("Expected datetime picker view to show selected time")
	}

	// Test that timezone is properly set
	if picker.timezone == nil {
		t.Error("Expected datetime picker to have timezone set")
	}

	// Test time formatting for BGA
	date, timeStr := picker.FormatForBGA()
	if len(date) != 10 || date[4] != '-' || date[7] != '-' {
		t.Errorf("Expected date in YYYY-MM-DD format, got %s", date)
	}
	if len(timeStr) != 5 || timeStr[2] != ':' {
		t.Errorf("Expected time in HH:MM format, got %s", timeStr)
	}
}
