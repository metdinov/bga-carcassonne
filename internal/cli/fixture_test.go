package cli

import (
	"strings"
	"testing"

	"carca-cli/internal/fixtures"

	tea "github.com/charmbracelet/bubbletea"
)

func TestFixtureModel_Init(t *testing.T) {
	// Create a sample division for testing
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", HomeScore: 2, AwayScore: 1, AwayPlayer: "Lord Trooper", Played: true},
					{ID: 2, HomePlayer: "webbi", HomeScore: 2, AwayScore: 0, AwayPlayer: "alehrosario", Played: true},
				},
			},
		},
	}

	model := NewFixtureModel(division)

	if model == nil {
		t.Fatal("Expected fixture model to be initialized")
	}

	if model.division != division {
		t.Errorf("Expected division to be set correctly")
	}

	if model.currentRound != 0 {
		t.Errorf("Expected currentRound to start at 0, got %d", model.currentRound)
	}

	if model.selectedMatch != 0 {
		t.Errorf("Expected selectedMatch to start at 0, got %d", model.selectedMatch)
	}
}

func TestFixtureModel_Update_Navigation(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
			{Number: 2, DateRange: "18/08 - 24/08", Matches: []*fixtures.Match{}},
			{Number: 3, DateRange: "25/08 - 31/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)

	// Test right arrow - next round
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	if model.currentRound != 1 {
		t.Errorf("Expected currentRound to move to 1, got %d", model.currentRound)
	}
	if cmd != nil {
		t.Errorf("Expected no command on navigation, got %v", cmd)
	}

	// Test left arrow - previous round
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if model.currentRound != 0 {
		t.Errorf("Expected currentRound to move back to 0, got %d", model.currentRound)
	}

	// Test wrap around at end
	model.currentRound = len(division.Rounds) - 1
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	if model.currentRound != 0 {
		t.Errorf("Expected currentRound to wrap to 0, got %d", model.currentRound)
	}

	// Test wrap around at beginning
	model.currentRound = 0
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if model.currentRound != len(division.Rounds)-1 {
		t.Errorf("Expected currentRound to wrap to %d, got %d", len(division.Rounds)-1, model.currentRound)
	}
}

func TestFixtureModel_Update_Back(t *testing.T) {
	division := &fixtures.Division{Name: "Elite", Rounds: []*fixtures.Round{}}
	model := NewFixtureModel(division)

	// Test ESC key
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("Expected command when pressing ESC")
	}

	// Test 'q' key
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("Expected command when pressing 'q'")
	}
}

func TestFixtureModel_View_ContainsDivisionName(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	if !strings.Contains(view, "Elite") {
		t.Errorf("Expected view to contain division name 'Elite', got: %s", view)
	}
}

func TestFixtureModel_View_ShowsRoundInfo(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches:   []*fixtures.Match{},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	if !strings.Contains(view, "Round 1") {
		t.Errorf("Expected view to contain 'Round 1', got: %s", view)
	}

	if !strings.Contains(view, "11/08 - 17/08") {
		t.Errorf("Expected view to contain date range '11/08 - 17/08', got: %s", view)
	}
}

func TestFixtureModel_View_ShowsMatches(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "herchu",
						HomeScore:  2,
						AwayScore:  1,
						AwayPlayer: "Lord Trooper",
						DateTime:   "12/08 - 09:30",
						BGALink:    "https://boardgamearena.com/tournament?id=423761",
						Played:     true,
					},
					{
						ID:         2,
						HomePlayer: "webbi",
						HomeScore:  0,
						AwayScore:  0,
						AwayPlayer: "alehrosario",
						DateTime:   "",
						BGALink:    "",
						Played:     false,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Should show played match with score
	if !strings.Contains(view, "herchu") {
		t.Errorf("Expected view to contain 'herchu', got: %s", view)
	}
	if !strings.Contains(view, "Lord Trooper") {
		t.Errorf("Expected view to contain 'Lord Trooper', got: %s", view)
	}
	if !strings.Contains(view, "2-1") {
		t.Errorf("Expected view to contain score '2-1', got: %s", view)
	}

	// Should show unplayed match
	if !strings.Contains(view, "webbi") {
		t.Errorf("Expected view to contain 'webbi', got: %s", view)
	}
	if !strings.Contains(view, "alehrosario") {
		t.Errorf("Expected view to contain 'alehrosario', got: %s", view)
	}
}

func TestFixtureModel_View_ShowsPlayedStatus(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", AwayPlayer: "Lord Trooper", Played: true},
					{ID: 2, HomePlayer: "webbi", AwayPlayer: "alehrosario", Played: false},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Should indicate played vs unplayed status differently
	// We'll check this in the implementation - played matches might be green, unplayed might be red
	if !strings.Contains(view, "herchu") || !strings.Contains(view, "webbi") {
		t.Errorf("Expected view to contain both players, got: %s", view)
	}
}

func TestFixtureModel_View_ShowsNavigation(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
			{Number: 2, DateRange: "18/08 - 24/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Should show navigation instructions including vim keys and page keys
	if !strings.Contains(view, "←/→, h/l, or PgUp/PgDown to navigate rounds") {
		t.Errorf("Expected view to show navigation instructions with page keys, got: %s", view)
	}
}

func TestFixtureModel_View_ShowsTableFormat(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "herchu",
						HomeScore:  2,
						AwayScore:  1,
						AwayPlayer: "Lord Trooper",
						DateTime:   "12/08 - 09:30",
						BGALink:    "https://boardgamearena.com/tournament?id=423761",
						Played:     true,
					},
					{
						ID:         2,
						HomePlayer: "webbi",
						HomeScore:  0,
						AwayScore:  0,
						AwayPlayer: "alehrosario",
						DateTime:   "",
						BGALink:    "",
						Played:     false,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Should show table headers
	expectedHeaders := []string{"PLAYED", "LOCAL", "VISITOR", "RESULT", "DATE", "TOURNAMENT_ID"}
	for _, header := range expectedHeaders {
		if !strings.Contains(view, header) {
			t.Errorf("Expected view to contain table header '%s', got: %s", header, view)
		}
	}

	// Should show played match data
	if !strings.Contains(view, "✓") {
		t.Errorf("Expected view to show checkmark for played match, got: %s", view)
	}
	if !strings.Contains(view, "herchu") {
		t.Errorf("Expected view to show home player 'herchu', got: %s", view)
	}
	if !strings.Contains(view, "Lord Trooper") {
		t.Errorf("Expected view to show away player 'Lord Trooper', got: %s", view)
	}
	if !strings.Contains(view, "2-1") {
		t.Errorf("Expected view to show result '2-1', got: %s", view)
	}
	if !strings.Contains(view, "423761") {
		t.Errorf("Expected view to show tournament ID '423761', got: %s", view)
	}

	// Should show unplayed match data
	if !strings.Contains(view, "○") {
		t.Errorf("Expected view to show circle for unplayed match, got: %s", view)
	}
	if !strings.Contains(view, "webbi") {
		t.Errorf("Expected view to show home player 'webbi', got: %s", view)
	}
	if !strings.Contains(view, "alehrosario") {
		t.Errorf("Expected view to show away player 'alehrosario', got: %s", view)
	}
}

func TestFixtureModel_ExtractTournamentID(t *testing.T) {
	model := NewFixtureModel(&fixtures.Division{})

	testCases := []struct {
		url      string
		expected string
	}{
		{"https://boardgamearena.com/tournament?id=423761", "423761"},
		{"https://boardgamearena.com/tournament?id=423761&token=xyz", "423761"},
		{"", ""},
		{"invalid-url", ""},
		{"https://boardgamearena.com/tournament", ""},
	}

	for _, tc := range testCases {
		result := model.extractTournamentID(tc.url)
		if result != tc.expected {
			t.Errorf("Expected tournament ID '%s' for URL '%s', got '%s'", tc.expected, tc.url, result)
		}
	}
}

func TestFixtureModel_GetCurrentRound(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
			{Number: 2, DateRange: "18/08 - 24/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)

	// Test first round
	round := model.GetCurrentRound()
	if round == nil {
		t.Fatal("Expected to get current round")
	}
	if round.Number != 1 {
		t.Errorf("Expected round number 1, got %d", round.Number)
	}

	// Test second round
	model.currentRound = 1
	round = model.GetCurrentRound()
	if round.Number != 2 {
		t.Errorf("Expected round number 2, got %d", round.Number)
	}
}

func TestFixtureModel_TableFormat_RealData(t *testing.T) {
	// Test with real fixture file to verify table formatting
	filename := "../../data/Liga Argentina - 1° Temporada - E-Fixture.csv"

	division, err := fixtures.ParseFixtureFile(filename)
	if err != nil {
		t.Skip("Skipping real data test - fixture file not available")
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Should contain table structure
	if !strings.Contains(view, "┌") || !strings.Contains(view, "└") {
		t.Errorf("Expected table borders in view")
	}

	// Should contain all expected columns
	expectedHeaders := []string{"PLAYED", "LOCAL", "VISITOR", "RESULT", "DATE", "TOURNAMENT_ID"}
	for _, header := range expectedHeaders {
		if !strings.Contains(view, header) {
			t.Errorf("Expected table header '%s' in real data view", header)
		}
	}

	// Should show played matches (at minimum)
	if !strings.Contains(view, "✓") {
		t.Errorf("Expected played matches (✓) in real data")
	}
	// Note: First round might not have unplayed matches, that's OK

	// Should show actual tournament IDs for played matches
	playedFound := false
	for _, round := range division.Rounds {
		for _, match := range round.Matches {
			if match.Played && match.BGALink != "" {
				playedFound = true
				break
			}
		}
		if playedFound {
			break
		}
	}

	if playedFound {
		// There should be at least one tournament ID visible
		hasNumbers := false
		for i := '0'; i <= '9'; i++ {
			if strings.ContainsRune(view, i) {
				hasNumbers = true
				break
			}
		}
		if !hasNumbers {
			t.Errorf("Expected tournament IDs (numbers) in real data view")
		}
	}

	t.Logf("Table format test passed with real data from %s", division.Name)
}

func TestFixtureModel_Update_VimNavigation_Left(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
			{Number: 2, DateRange: "18/08 - 24/08", Matches: []*fixtures.Match{}},
			{Number: 3, DateRange: "25/08 - 31/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)
	model.currentRound = 2 // Start at round 2

	// Send 'h' key (vim left)
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if model.currentRound != 1 {
		t.Errorf("Expected currentRound to move to 1 with 'h', got %d", model.currentRound)
	}

	if cmd != nil {
		t.Errorf("Expected no command on vim navigation, got %v", cmd)
	}
}

func TestFixtureModel_Update_VimNavigation_Right(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
			{Number: 2, DateRange: "18/08 - 24/08", Matches: []*fixtures.Match{}},
			{Number: 3, DateRange: "25/08 - 31/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)

	// Send 'l' key (vim right)
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

	if model.currentRound != 1 {
		t.Errorf("Expected currentRound to move to 1 with 'l', got %d", model.currentRound)
	}

	if cmd != nil {
		t.Errorf("Expected no command on vim navigation, got %v", cmd)
	}
}

func TestFixtureModel_Update_VimNavigation_WrapAround(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
			{Number: 2, DateRange: "18/08 - 24/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)

	// Test 'l' wrap around at end
	model.currentRound = len(division.Rounds) - 1
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

	if model.currentRound != 0 {
		t.Errorf("Expected currentRound to wrap to 0 with 'l', got %d", model.currentRound)
	}

	// Test 'h' wrap around at beginning
	model.currentRound = 0
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if model.currentRound != len(division.Rounds)-1 {
		t.Errorf("Expected currentRound to wrap to %d with 'h', got %d", len(division.Rounds)-1, model.currentRound)
	}

	if cmd != nil {
		t.Errorf("Expected no command on vim navigation, got %v", cmd)
	}
}

func TestFixtureModel_Update_PageNavigation_Up(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
			{Number: 2, DateRange: "18/08 - 24/08", Matches: []*fixtures.Match{}},
			{Number: 3, DateRange: "25/08 - 31/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)
	model.currentRound = 2 // Start at round 3

	// Send Page Up key
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyPgUp})

	if model.currentRound != 1 {
		t.Errorf("Expected currentRound to move to 1 with Page Up, got %d", model.currentRound)
	}

	if cmd != nil {
		t.Errorf("Expected no command on page navigation, got %v", cmd)
	}
}

func TestFixtureModel_Update_PageNavigation_Down(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
			{Number: 2, DateRange: "18/08 - 24/08", Matches: []*fixtures.Match{}},
			{Number: 3, DateRange: "25/08 - 31/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)

	// Send Page Down key
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyPgDown})

	if model.currentRound != 1 {
		t.Errorf("Expected currentRound to move to 1 with Page Down, got %d", model.currentRound)
	}

	if cmd != nil {
		t.Errorf("Expected no command on page navigation, got %v", cmd)
	}
}

func TestFixtureModel_Update_PageNavigation_WrapAround(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
			{Number: 2, DateRange: "18/08 - 24/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)

	// Test Page Down wrap around at end
	model.currentRound = len(division.Rounds) - 1
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyPgDown})

	if model.currentRound != 0 {
		t.Errorf("Expected currentRound to wrap to 0 with Page Down, got %d", model.currentRound)
	}

	// Test Page Up wrap around at beginning
	model.currentRound = 0
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyPgUp})

	if model.currentRound != len(division.Rounds)-1 {
		t.Errorf("Expected currentRound to wrap to %d with Page Up, got %d", len(division.Rounds)-1, model.currentRound)
	}

	if cmd != nil {
		t.Errorf("Expected no command on page navigation, got %v", cmd)
	}
}

func TestFixtureModel_View_ShowsPageNavigationInstructions(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
			{Number: 2, DateRange: "18/08 - 24/08", Matches: []*fixtures.Match{}},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Should show Page Up/Down in navigation instructions
	if !strings.Contains(view, "PgUp/PgDown") {
		t.Errorf("Expected view to show PgUp/PgDown instructions, got: %s", view)
	}

	// Should show all navigation options
	if !strings.Contains(view, "←/→, h/l, or PgUp/PgDown to navigate rounds") {
		t.Errorf("Expected view to show complete navigation instructions, got: %s", view)
	}
}

func TestFixtureModel_View_ConsistentColumnWidths(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "VeryLongPlayerNameHere",
						HomeScore:  2,
						AwayScore:  1,
						AwayPlayer: "Short",
						DateTime:   "12/08 - 09:30",
						BGALink:    "https://boardgamearena.com/tournament?id=423761",
						Played:     true,
					},
					{
						ID:         2,
						HomePlayer: "A",
						HomeScore:  0,
						AwayScore:  0,
						AwayPlayer: "AnotherVeryLongPlayerName",
						DateTime:   "13/08 - 22:00",
						BGALink:    "",
						Played:     false,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Should show "DATE" column header instead of "DATETIME"
	if !strings.Contains(view, "DATE") {
		t.Errorf("Expected view to contain 'DATE' column header, got: %s", view)
	}

	if strings.Contains(view, "DATETIME") {
		t.Errorf("Expected view to not contain 'DATETIME' header, got: %s", view)
	}

	// Should contain table structure with consistent formatting
	if !strings.Contains(view, "┌") || !strings.Contains(view, "└") {
		t.Errorf("Expected table borders in view with consistent widths")
	}
}

func TestFixtureModel_CalculateMaxPlayerNameWidth(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						HomePlayer: "VeryLongPlayerNameHere",
						AwayPlayer: "Short",
					},
					{
						HomePlayer: "A",
						AwayPlayer: "AnotherVeryLongPlayerName",
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	maxWidth := model.calculateMaxPlayerNameWidth()

	// Should return the length of the longest player name plus tab padding (8 spaces)
	expectedWidth := len("AnotherVeryLongPlayerName") + 8 // 25 + 8 = 33 characters
	if maxWidth != expectedWidth {
		t.Errorf("Expected max player name width %d, got %d", expectedWidth, maxWidth)
	}
}

func TestFixtureModel_View_ColumnPadding(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "Short",
						HomeScore:  2,
						AwayScore:  1,
						AwayPlayer: "Long",
						DateTime:   "12/08 - 09:30",
						BGALink:    "https://boardgamearena.com/tournament?id=423761",
						Played:     true,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Should show extra padding in columns - short names should have more trailing spaces
	if !strings.Contains(view, "Short") {
		t.Errorf("Expected view to contain 'Short', got: %s", view)
	}

	// Should have consistent column widths with padding visible
	if !strings.Contains(view, "┌") || !strings.Contains(view, "└") {
		t.Errorf("Expected table borders in padded view")
	}
}

func TestFixtureModel_View_PaddingDemonstration(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "A",
						HomeScore:  2,
						AwayScore:  1,
						AwayPlayer: "VeryLongPlayerNameExample",
						DateTime:   "12/08 - 09:30",
						BGALink:    "https://boardgamearena.com/tournament?id=423761",
						Played:     true,
					},
					{
						ID:         2,
						HomePlayer: "ShortName",
						HomeScore:  0,
						AwayScore:  0,
						AwayPlayer: "B",
						DateTime:   "13/08 - 22:00",
						BGALink:    "",
						Played:     false,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Both columns should have same width despite different name lengths
	// The padding ensures consistent alignment
	if !strings.Contains(view, "A") {
		t.Errorf("Expected view to contain short name 'A', got: %s", view)
	}
	if !strings.Contains(view, "VeryLongPlayerNameExample") {
		t.Errorf("Expected view to contain long name 'VeryLongPlayerNameExample', got: %s", view)
	}

	// Visual verification that padding is working - table should be well-formatted
	t.Logf("Padded table output:\n%s", view)
}

func TestFixtureModel_Update_MatchSelection_Down(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", AwayPlayer: "Lord Trooper", Played: true},
					{ID: 2, HomePlayer: "webbi", AwayPlayer: "alehrosario", Played: false},
					{ID: 3, HomePlayer: "player3", AwayPlayer: "player4", Played: false},
				},
			},
		},
	}

	model := NewFixtureModel(division)

	// Send 'j' key (vim down for match selection)
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if model.selectedMatch != 1 {
		t.Errorf("Expected selectedMatch to move to 1 with 'j', got %d", model.selectedMatch)
	}

	if cmd != nil {
		t.Errorf("Expected no command on match selection, got %v", cmd)
	}
}

func TestFixtureModel_Update_MatchSelection_Up(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", AwayPlayer: "Lord Trooper", Played: true},
					{ID: 2, HomePlayer: "webbi", AwayPlayer: "alehrosario", Played: false},
					{ID: 3, HomePlayer: "player3", AwayPlayer: "player4", Played: false},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	model.selectedMatch = 2 // Start at match 2

	// Send 'k' key (vim up for match selection)
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if model.selectedMatch != 1 {
		t.Errorf("Expected selectedMatch to move to 1 with 'k', got %d", model.selectedMatch)
	}

	if cmd != nil {
		t.Errorf("Expected no command on match selection, got %v", cmd)
	}
}

func TestFixtureModel_Update_MatchSelection_WrapAround(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", AwayPlayer: "Lord Trooper", Played: true},
					{ID: 2, HomePlayer: "webbi", AwayPlayer: "alehrosario", Played: false},
				},
			},
		},
	}

	model := NewFixtureModel(division)

	// Test 'j' wrap around at bottom
	model.selectedMatch = 1 // Last match
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if model.selectedMatch != 0 {
		t.Errorf("Expected selectedMatch to wrap to 0 with 'j', got %d", model.selectedMatch)
	}

	// Test 'k' wrap around at top
	model.selectedMatch = 0
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if model.selectedMatch != 1 {
		t.Errorf("Expected selectedMatch to wrap to 1 with 'k', got %d", model.selectedMatch)
	}

	if cmd != nil {
		t.Errorf("Expected no command on match selection, got %v", cmd)
	}
}

func TestFixtureModel_Update_EnterCopyLink(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "herchu",
						AwayPlayer: "Lord Trooper",
						BGALink:    "https://boardgamearena.com/tournament?id=423761",
						Played:     true,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	model.selectedMatch = 0 // Select first match

	// Send enter key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("Expected command when pressing enter on played match")
	}

	// Check that status message is set
	fixModel, ok := updatedModel.(*FixtureModel)
	if !ok {
		t.Fatal("Expected FixtureModel to be returned")
	}

	if fixModel.statusMessage == "" {
		t.Error("Expected status message to be set after copying link")
	}

	if !strings.Contains(fixModel.statusMessage, "copied") {
		t.Errorf("Expected status message to contain 'copied', got: %s", fixModel.statusMessage)
	}
}

func TestFixtureModel_Update_EnterCreateTournament(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
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
	model.selectedMatch = 0 // Select first match

	// Send enter key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd != nil {
		t.Errorf("Expected no command for unplayed match, got %v", cmd)
	}

	// Check that status message is set for creating tournament
	fixModel, ok := updatedModel.(*FixtureModel)
	if !ok {
		t.Fatal("Expected FixtureModel to be returned")
	}

	if fixModel.statusMessage == "" {
		t.Error("Expected status message to be set for unplayed match")
	}

	if !strings.Contains(fixModel.statusMessage, "Press 'c' to create tournament") {
		t.Errorf("Expected status message to contain create tournament instruction, got: %s", fixModel.statusMessage)
	}
}

func TestFixtureModel_View_ShowsSelectedMatch(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", AwayPlayer: "Lord Trooper", Played: true},
					{ID: 2, HomePlayer: "webbi", AwayPlayer: "alehrosario", Played: false},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	model.selectedMatch = 1 // Select second match

	view := model.View()

	// Should show some indication of selection (this will be implemented in the actual view)
	if !strings.Contains(view, "webbi") {
		t.Errorf("Expected view to contain selected match player 'webbi', got: %s", view)
	}
}

func TestFixtureModel_View_ShowsStatusMessage(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", AwayPlayer: "Lord Trooper", Played: true},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	model.statusMessage = "Tournament link copied to clipboard!"

	view := model.View()

	if !strings.Contains(view, "Tournament link copied to clipboard!") {
		t.Errorf("Expected view to show status message, got: %s", view)
	}
}

func TestFixtureModel_View_ConstantColumnWidthUnplayedMatches(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "VeryLongPlayerNameExample",
						HomeScore:  2,
						AwayScore:  1,
						AwayPlayer: "AnotherVeryLongName",
						DateTime:   "12/08 - 09:30",
						BGALink:    "https://boardgamearena.com/tournament?id=423761",
						Played:     true,
					},
					{
						ID:         2,
						HomePlayer: "Short",
						HomeScore:  0,
						AwayScore:  0,
						AwayPlayer: "A",
						DateTime:   "",
						BGALink:    "",
						Played:     false,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Both played and unplayed matches should have same column widths
	// Check that table structure is consistent
	if !strings.Contains(view, "VeryLongPlayerNameExample") {
		t.Errorf("Expected view to contain long name, got: %s", view)
	}

	if !strings.Contains(view, "Short") && !strings.Contains(view, "A") {
		t.Errorf("Expected view to contain short names, got: %s", view)
	}

	// Visual verification that column widths are consistent
	t.Logf("Constant width table with unplayed matches:\n%s", view)
}

func TestFixtureModel_View_ColumnWidthConsistency(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "VeryLongPlayerNameHere",
						HomeScore:  2,
						AwayScore:  1,
						AwayPlayer: "Short",
						DateTime:   "12/08 - 09:30",
						BGALink:    "https://boardgamearena.com/tournament?id=423761",
						Played:     true,
					},
					{
						ID:         2,
						HomePlayer: "A",
						HomeScore:  0,
						AwayScore:  0,
						AwayPlayer: "EvenLongerPlayerNameExample",
						DateTime:   "",
						BGALink:    "",
						Played:     false,
					},
					{
						ID:         3,
						HomePlayer: "Medium",
						HomeScore:  1,
						AwayScore:  2,
						AwayPlayer: "Player",
						DateTime:   "13/08 - 15:00",
						BGALink:    "https://boardgamearena.com/tournament?id=999999",
						Played:     true,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// All matches should have consistent column widths regardless of:
	// - Name length differences
	// - Played vs unplayed status
	// - Empty vs filled data fields

	// Check that longest name determines column width for all rows
	if !strings.Contains(view, "VeryLongPlayerNameHere") {
		t.Errorf("Expected view to contain longest home player name")
	}
	if !strings.Contains(view, "EvenLongerPlayerNameExample") {
		t.Errorf("Expected view to contain longest away player name")
	}

	// Verify table structure is maintained
	lines := strings.Split(view, "\n")
	var tableLines []string
	for _, line := range lines {
		if strings.Contains(line, "│") {
			tableLines = append(tableLines, line)
		}
	}

	if len(tableLines) < 4 { // Header + separator + at least 2 data rows
		t.Errorf("Expected at least 4 table lines, got %d", len(tableLines))
	}

	// Verify that all table rows have consistent structure (same number of columns)
	for i, line := range tableLines {
		// Count column separators (│) - should be consistent across all rows
		columnCount := strings.Count(line, "│")
		if i == 0 {
			// First row sets the expected column count
			if columnCount < 6 { // Should have at least 6 separators for our 6 columns
				t.Errorf("Expected at least 6 column separators, got %d in row: %s", columnCount, line)
			}
		} else {
			// All other rows should match the first row's column count
			firstRowColumnCount := strings.Count(tableLines[0], "│")
			if columnCount != firstRowColumnCount {
				t.Errorf("Row %d has inconsistent column count: expected %d, got %d", i, firstRowColumnCount, columnCount)
			}
		}
	}

	t.Logf("Column structure consistency verified with mixed match states")
}

func TestFixtureModel_View_AllUnplayedMatches(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "VeryLongPlayerName",
						HomeScore:  0,
						AwayScore:  0,
						AwayPlayer: "Short",
						DateTime:   "",
						BGALink:    "",
						Played:     false,
					},
					{
						ID:         2,
						HomePlayer: "A",
						HomeScore:  0,
						AwayScore:  0,
						AwayPlayer: "AnotherLongName",
						DateTime:   "",
						BGALink:    "",
						Played:     false,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Should handle all unplayed matches with consistent column widths
	// Even when all data is empty/default
	if !strings.Contains(view, "VeryLongPlayerName") {
		t.Errorf("Expected view to contain 'VeryLongPlayerName', got: %s", view)
	}

	if !strings.Contains(view, "AnotherLongName") {
		t.Errorf("Expected view to contain 'AnotherLongName', got: %s", view)
	}

	// Both rows should show unplayed status
	unplayedCount := strings.Count(view, "○")
	if unplayedCount != 2 {
		t.Errorf("Expected 2 unplayed match indicators, got %d", unplayedCount)
	}

	// All empty fields should show consistent "-" placeholders
	dashCount := strings.Count(view, "-")
	if dashCount < 6 { // At least 6 dashes for result, date, and tournament_id columns
		t.Errorf("Expected at least 6 dash placeholders for empty fields, got %d", dashCount)
	}

	t.Logf("All unplayed matches table formatting:\n%s", view)
}

func TestFixtureModel_View_FixedTableLayoutAcrossRounds(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "VeryLongPlayerNameExample",
						HomeScore:  2,
						AwayScore:  1,
						AwayPlayer: "AnotherVeryLongName",
						DateTime:   "12/08 - 09:30",
						BGALink:    "https://boardgamearena.com/tournament?id=423761",
						Played:     true,
					},
				},
			},
			{
				Number:    2,
				DateRange: "18/08 - 24/08",
				Matches: []*fixtures.Match{
					{
						ID:         2,
						HomePlayer: "Short",
						HomeScore:  0,
						AwayScore:  0,
						AwayPlayer: "A",
						DateTime:   "",
						BGALink:    "",
						Played:     false,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)

	// Get view for first round (played matches)
	model.currentRound = 0
	playedRoundView := model.View()

	// Get view for second round (unplayed matches)
	model.currentRound = 1
	unplayedRoundView := model.View()

	// Extract table lines from both views
	playedLines := extractTableLines(playedRoundView)
	unplayedLines := extractTableLines(unplayedRoundView)

	// Both rounds should have the same table structure
	if len(playedLines) != len(unplayedLines) {
		t.Errorf("Different number of table lines: played=%d, unplayed=%d", len(playedLines), len(unplayedLines))
	}

	// Header and separator rows should be identical
	if playedLines[0] != unplayedLines[0] {
		t.Errorf("Header rows differ between rounds")
	}
	if playedLines[1] != unplayedLines[1] {
		t.Errorf("Separator rows differ between rounds")
	}

	// All rows should have the same width
	for i := 0; i < len(playedLines) && i < len(unplayedLines); i++ {
		if len(playedLines[i]) != len(unplayedLines[i]) {
			t.Errorf("Row %d width differs: played=%d, unplayed=%d", i, len(playedLines[i]), len(unplayedLines[i]))
		}
	}

	t.Logf("Fixed layout verified across rounds with different content")
}

func TestFixtureModel_View_FixedLayoutWithVariousDataLengths(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
						HomePlayer: "VeryLongPlayerNameExample",
						HomeScore:  2,
						AwayScore:  1,
						AwayPlayer: "Short",
						DateTime:   "12/08 - 09:30:45 AM",
						BGALink:    "https://boardgamearena.com/tournament?id=123456789",
						Played:     true,
					},
				},
			},
			{
				Number:    2,
				DateRange: "18/08 - 24/08",
				Matches: []*fixtures.Match{
					{
						ID:         2,
						HomePlayer: "A",
						HomeScore:  0,
						AwayScore:  0,
						AwayPlayer: "B",
						DateTime:   "",
						BGALink:    "",
						Played:     false,
					},
				},
			},
			{
				Number:    3,
				DateRange: "25/08 - 31/08",
				Matches: []*fixtures.Match{
					{
						ID:         3,
						HomePlayer: "Medium",
						HomeScore:  1,
						AwayScore:  2,
						AwayPlayer: "Player",
						DateTime:   "01/09",
						BGALink:    "https://boardgamearena.com/tournament?id=1",
						Played:     true,
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)

	var tableWidths []int
	var headerRows []string

	// Check all three rounds have identical table structure
	for i := 0; i < 3; i++ {
		model.currentRound = i
		view := model.View()
		tableLines := extractTableLines(view)

		if len(tableLines) > 0 {
			tableWidths = append(tableWidths, len(tableLines[0]))
			headerRows = append(headerRows, tableLines[0])
		}
	}

	// All rounds should have same table width
	if len(tableWidths) != 3 {
		t.Fatalf("Expected 3 table widths, got %d", len(tableWidths))
	}

	firstWidth := tableWidths[0]
	for i, width := range tableWidths {
		if width != firstWidth {
			t.Errorf("Round %d has different table width: expected %d, got %d", i+1, firstWidth, width)
		}
	}

	// All rounds should have identical header rows
	firstHeader := headerRows[0]
	for i, header := range headerRows {
		if header != firstHeader {
			t.Errorf("Round %d has different header structure", i+1)
		}
	}

	t.Logf("Fixed layout verified across 3 rounds with varying data lengths")
}

// extractTableLines extracts only the table lines from a view string
func extractTableLines(view string) []string {
	lines := strings.Split(view, "\n")
	var tableLines []string
	for _, line := range lines {
		if strings.Contains(line, "│") || strings.Contains(line, "┌") || strings.Contains(line, "└") || strings.Contains(line, "├") {
			tableLines = append(tableLines, line)
		}
	}
	return tableLines
}

func TestFixtureModel_CalculateMaxDateWidth(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						DateTime: "12/08 - 09:30:45 AM",
					},
					{
						DateTime: "01/09",
					},
				},
			},
			{
				Number:    2,
				DateRange: "18/08 - 24/08",
				Matches: []*fixtures.Match{
					{
						DateTime: "",
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	maxWidth := model.calculateMaxDateWidth()

	// Should return the length of the longest date string
	expectedWidth := len("12/08 - 09:30:45 AM") // 19 characters
	if maxWidth != expectedWidth {
		t.Errorf("Expected max date width %d, got %d", expectedWidth, maxWidth)
	}
}

func TestFixtureModel_CalculateMaxTournamentIDWidth(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						BGALink: "https://boardgamearena.com/tournament?id=123456789",
					},
					{
						BGALink: "https://boardgamearena.com/tournament?id=1",
					},
				},
			},
			{
				Number:    2,
				DateRange: "18/08 - 24/08",
				Matches: []*fixtures.Match{
					{
						BGALink: "",
					},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	maxWidth := model.calculateMaxTournamentIDWidth()

	// Should return the length of the longest tournament ID or minimum header width
	expectedWidth := 12 // "TOURNAMENT_ID" header length
	if maxWidth < expectedWidth {
		t.Errorf("Expected max tournament ID width at least %d, got %d", expectedWidth, maxWidth)
	}

	// Should handle the longest ID
	if maxWidth < len("123456789") {
		t.Errorf("Expected max tournament ID width to accommodate longest ID")
	}
}

func TestFixtureModel_Update_CreateTournamentWithC(t *testing.T) {
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{
						ID:         1,
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
	model.selectedMatch = 0 // Select first match
	model.statusMessage = "Press 'c' to create tournament for this match"

	// Send 'c' key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	if cmd == nil {
		t.Fatal("Expected command when pressing 'c' to create tournament")
	}

	// Check that status message is updated
	fixModel, ok := updatedModel.(*FixtureModel)
	if !ok {
		t.Fatal("Expected FixtureModel to be returned")
	}

	if fixModel.statusMessage == "" {
		t.Error("Expected status message to be set after creating tournament")
	}

	if !strings.Contains(fixModel.statusMessage, "Creating tournament") {
		t.Errorf("Expected status message to contain 'Creating tournament', got: %s", fixModel.statusMessage)
	}
}
