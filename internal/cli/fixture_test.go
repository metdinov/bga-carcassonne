package cli

import (
	"strings"
	"testing"
	"time"

	"carca-cli/internal/bga"
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
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	if model.currentRound != 1 {
		t.Errorf("Expected currentRound to move to 1, got %d", model.currentRound)
	}

	// Test left arrow - previous round
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if model.currentRound != 0 {
		t.Errorf("Expected currentRound to move back to 0, got %d", model.currentRound)
	}

	// Test wrap around at end
	model.currentRound = len(division.Rounds) - 1
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	if model.currentRound != 0 {
		t.Errorf("Expected currentRound to wrap to 0, got %d", model.currentRound)
	}

	// Test wrap around at beginning
	model.currentRound = 0
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
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
	expectedHeaders := []string{"PLAYED", "HOME", "AWAY", "RESULT", "DATE", "TOURNAMENT_ID"}
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
	expectedHeaders := []string{"PLAYED", "HOME", "AWAY", "RESULT", "DATE", "TOURNAMENT_ID"}
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
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if model.currentRound != 1 {
		t.Errorf("Expected currentRound to move to 1 with 'h', got %d", model.currentRound)
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
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

	if model.currentRound != 1 {
		t.Errorf("Expected currentRound to move to 1 with 'l', got %d", model.currentRound)
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
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

	if model.currentRound != 0 {
		t.Errorf("Expected currentRound to wrap to 0 with 'l', got %d", model.currentRound)
	}

	// Test 'h' wrap around at beginning
	model.currentRound = 0
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if model.currentRound != len(division.Rounds)-1 {
		t.Errorf("Expected currentRound to wrap to %d with 'h', got %d", len(division.Rounds)-1, model.currentRound)
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
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyPgUp})

	if model.currentRound != 1 {
		t.Errorf("Expected currentRound to move to 1 with Page Up, got %d", model.currentRound)
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
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyPgDown})

	if model.currentRound != 1 {
		t.Errorf("Expected currentRound to move to 1 with Page Down, got %d", model.currentRound)
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
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyPgDown})

	if model.currentRound != 0 {
		t.Errorf("Expected currentRound to wrap to 0 with Page Down, got %d", model.currentRound)
	}

	// Test Page Up wrap around at beginning
	model.currentRound = 0
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyPgUp})

	if model.currentRound != len(division.Rounds)-1 {
		t.Errorf("Expected currentRound to wrap to %d with Page Up, got %d", len(division.Rounds)-1, model.currentRound)
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
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if model.selectedMatch != 1 {
		t.Errorf("Expected selectedMatch to move to 1 with 'j', got %d", model.selectedMatch)
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
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if model.selectedMatch != 1 {
		t.Errorf("Expected selectedMatch to move to 1 with 'k', got %d", model.selectedMatch)
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
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if model.selectedMatch != 0 {
		t.Errorf("Expected selectedMatch to wrap to 0 with 'j', got %d", model.selectedMatch)
	}

	// Test 'k' wrap around at top
	model.selectedMatch = 0
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if model.selectedMatch != 1 {
		t.Errorf("Expected selectedMatch to wrap to 1 with 'k', got %d", model.selectedMatch)
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
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

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

func TestFixtureModel_CreateTournament_Integration(t *testing.T) {
	// Create a division with unplayed matches
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
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

	// Set up mock BGA client
	mockClient := bga.NewMockClient("testuser", "testpass")
	model.SetBGAClient(mockClient)

	// Login the mock client
	err := mockClient.Login()
	if err != nil {
		t.Fatalf("Failed to login mock client: %v", err)
	}

	// Select the match and press 'c' to show datetime picker
	model.selectedMatch = 0
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	// Should show datetime picker
	fixtureModel := updatedModel.(*FixtureModel)
	if !fixtureModel.showDatePicker {
		t.Fatal("Expected datetime picker to be shown")
	}
	if fixtureModel.dateTimePicker == nil {
		t.Fatal("Expected datetime picker to be created")
	}

	// Simulate Enter key on datetime picker to confirm selection
	// First enter
	updatedModel, _ = fixtureModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	fixtureModel = updatedModel.(*FixtureModel)

	// Second enter
	updatedModel, cmd := fixtureModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	fixtureModel = updatedModel.(*FixtureModel)
	if cmd == nil {
		t.Fatal("Expected command from datetime picker confirmation")
	}

	// Execute wrapped command
	internalMsg := cmd()

	// Process internal message
	updatedModel, cmd2 := fixtureModel.Update(internalMsg)
	fixtureModel = updatedModel.(*FixtureModel)
	if cmd2 == nil {
		t.Fatal("Expected command after processing internal confirmation")
	}

	// Execute the datetime selection command
	msg := cmd2()
	if dateTimeMsg, ok := msg.(DateTimeSelectedMsg); ok {
		// Verify the datetime selection message
		if dateTimeMsg.HomePlayer != "player1" {
			t.Errorf("Expected HomePlayer 'player1', got %s", dateTimeMsg.HomePlayer)
		}
		if dateTimeMsg.AwayPlayer != "player2" {
			t.Errorf("Expected AwayPlayer 'player2', got %s", dateTimeMsg.AwayPlayer)
		}
		if dateTimeMsg.Division != "Elite" {
			t.Errorf("Expected Division 'Elite', got %s", dateTimeMsg.Division)
		}
		if dateTimeMsg.MatchID != 1 {
			t.Errorf("Expected MatchID 1, got %d", dateTimeMsg.MatchID)
		}

		// Process the datetime selection to show confirmation screen
		updatedModel, _ = fixtureModel.Update(msg)
		fixtureModel = updatedModel.(*FixtureModel)

		// Should show confirmation screen
		if !fixtureModel.showConfirmation {
			t.Fatal("Expected confirmation screen to be shown")
		}
		if fixtureModel.confirmationModel == nil {
			t.Fatal("Expected confirmation model to be created")
		}

		// Simulate Enter key on confirmation screen to confirm creation
		_, cmd := fixtureModel.confirmationModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
		if cmd == nil {
			t.Fatal("Expected command from confirmation screen")
		}

		// Execute confirmation command
		confirmMsg := cmd()
		if tournamentConfirmedMsg, ok := confirmMsg.(TournamentConfirmedMsg); ok {
			// Process tournament confirmation to trigger creation
			intermediateModel, cmd := fixtureModel.Update(tournamentConfirmedMsg)
			if cmd == nil {
				t.Fatal("Expected async tournament creation command")
			}

			// Execute async tournament creation
			asyncMsg := cmd()
			createWithDateTimeMsg, ok := asyncMsg.(createTournamentMsgWithDateTime)
			if !ok {
				t.Fatalf("Expected createTournamentMsgWithDateTime, got %T", asyncMsg)
			}

			// Process tournament creation with datetime
			finalIntermediateModel, cmd := intermediateModel.Update(createWithDateTimeMsg)
			if cmd == nil {
				t.Fatal("Expected async tournament creation command")
			}

			// Execute final async tournament creation
			finalAsyncMsg := cmd()
			tournamentMsg, ok := finalAsyncMsg.(tournamentCreatedMsg)
			if !ok {
				t.Fatalf("Expected tournamentCreatedMsg from final async creation, got %T", finalAsyncMsg)
			}

			// Process tournament creation completion
			finalModel, _ := finalIntermediateModel.Update(tournamentMsg)

			if tournamentMsg.success {
				// Verify the match was updated with the tournament link
				finalFixtureModel := finalModel.(*FixtureModel)
				match := finalFixtureModel.division.Rounds[0].Matches[0]
				if match.BGALink == "" {
					t.Error("Expected match to have BGALink set after tournament creation")
				}
				if !strings.Contains(match.BGALink, "tournament?id=") {
					t.Errorf("Expected valid tournament link, got %s", match.BGALink)
				}
			} else {
				t.Errorf("Tournament creation failed: %s", tournamentMsg.error)
			}
		} else {
			t.Errorf("Expected TournamentConfirmedMsg, got %T", confirmMsg)
		}
	} else {
		t.Errorf("Expected DateTimeSelectedMsg, got %T", msg)
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
	// Send 'c' key to show datetime picker
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	// Should show datetime picker instead of immediate command
	fixtureModel := updatedModel.(*FixtureModel)
	if !fixtureModel.showDatePicker {
		t.Error("Expected datetime picker to be shown when pressing 'c'")
	}

	// Check that datetime picker is shown instead of immediate tournament creation
	fixModel, ok := updatedModel.(*FixtureModel)
	if !ok {
		t.Fatal("Expected FixtureModel to be returned")
	}

	if !fixModel.showDatePicker {
		t.Error("Expected datetime picker to be shown when pressing 'c'")
	}

	if fixModel.dateTimePicker == nil {
		t.Error("Expected datetime picker to be created")
	}
}

func TestFixtureModel_EditDateTime(t *testing.T) {
	// Create test data
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
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
	model.bgaClient = bga.NewMockClient("test", "test")

	// Step 1: Start tournament creation process
	model.selectedMatch = 0
	updatedModel, _ := model.handleCreateTournament()
	fixtureModel := updatedModel.(*FixtureModel)

	if !fixtureModel.showDatePicker {
		t.Fatal("Expected datetime picker to be shown")
	}

	// Step 2: Select a datetime
	testTime := time.Date(2025, 9, 6, 21, 30, 0, 0, time.Local)
	dateTimeMsg := DateTimeSelectedMsg{
		HomePlayer:  "player1",
		AwayPlayer:  "player2",
		Division:    "Elite",
		RoundNumber: 1,
		MatchNumber: 1,
		MatchID:     1,
		DateTime:    testTime,
	}

	updatedModel, _ = fixtureModel.Update(dateTimeMsg)
	fixtureModel = updatedModel.(*FixtureModel)

	// Should show confirmation screen
	if !fixtureModel.showConfirmation {
		t.Fatal("Expected confirmation screen to be shown")
	}
	if fixtureModel.showDatePicker {
		t.Error("Expected datetime picker to be hidden")
	}

	// Step 3: Press 'e' to edit the datetime
	editMsg := EditDateTimeMsg{
		HomePlayer:  "player1",
		AwayPlayer:  "player2",
		Division:    "Elite",
		RoundNumber: 1,
		MatchNumber: 1,
		MatchID:     1,
		DateTime:    testTime,
	}

	updatedModel, _ = fixtureModel.Update(editMsg)
	fixtureModel = updatedModel.(*FixtureModel)

	// Should go back to datetime picker
	if !fixtureModel.showDatePicker {
		t.Error("Expected datetime picker to be shown after edit")
	}
	if fixtureModel.showConfirmation {
		t.Error("Expected confirmation screen to be hidden after edit")
	}

	// Verify datetime picker has the correct initial values
	if fixtureModel.dateTimePicker == nil {
		t.Fatal("Expected datetime picker to be created")
	}

	// Verify that the picker was initialized with the selected time
	if !fixtureModel.dateTimePicker.selectedTime.Equal(testTime) {
		t.Errorf("Expected datetime picker to be initialized with %v, got %v",
			testTime, fixtureModel.dateTimePicker.selectedTime)
	}

	if fixtureModel.dateTimePicker.homePlayer != "player1" {
		t.Errorf("Expected homePlayer 'player1', got %s", fixtureModel.dateTimePicker.homePlayer)
	}

	if fixtureModel.dateTimePicker.awayPlayer != "player2" {
		t.Errorf("Expected awayPlayer 'player2', got %s", fixtureModel.dateTimePicker.awayPlayer)
	}
}

func TestFixtureModel_EditDateTimeWithKeyPress(t *testing.T) {
	// Create test data
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
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
	model.bgaClient = bga.NewMockClient("test", "test")

	// Step 1: Start tournament creation and reach confirmation screen
	model.selectedMatch = 0
	updatedModel, _ := model.handleCreateTournament()
	fixtureModel := updatedModel.(*FixtureModel)

	// Step 2: Select datetime and get to confirmation screen
	testTime := time.Date(2025, 9, 6, 21, 30, 0, 0, time.Local)
	dateTimeMsg := DateTimeSelectedMsg{
		HomePlayer:  "player1",
		AwayPlayer:  "player2",
		Division:    "Elite",
		RoundNumber: 1,
		MatchNumber: 1,
		MatchID:     1,
		DateTime:    testTime,
	}

	updatedModel, _ = fixtureModel.Update(dateTimeMsg)
	fixtureModel = updatedModel.(*FixtureModel)

	// Should be in confirmation screen
	if !fixtureModel.showConfirmation {
		t.Fatal("Expected confirmation screen to be shown")
	}

	// Step 3: Simulate 'e' key press in confirmation screen
	keyMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("e"),
	}

	updatedModel, cmd := fixtureModel.Update(keyMsg)
	fixtureModel = updatedModel.(*FixtureModel)

	// Execute the command and process the EditDateTimeMsg
	if cmd != nil {
		msg := cmd()
		updatedModel, _ = fixtureModel.Update(msg)
		fixtureModel = updatedModel.(*FixtureModel)
	}

	// Should be back to datetime picker
	if !fixtureModel.showDatePicker {
		t.Error("Expected datetime picker to be shown after pressing 'e'")
	}
	if fixtureModel.showConfirmation {
		t.Error("Expected confirmation screen to be hidden after pressing 'e'")
	}

	// Verify the confirmation screen view shows the edit instruction
	confirmationModel := NewTournamentConfirmationModel(
		"player1", "player2", "Elite", 1, 1, 1, testTime,
	)

	view := confirmationModel.View()
	if !strings.Contains(view, "Press 'e' to edit date/time") {
		t.Error("Expected confirmation view to contain edit instruction")
	}
}

func TestFixtureModel_View_ShowsMatchNumberColumn(t *testing.T) {
	// Create test data with specific match IDs
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", HomeScore: 2, AwayScore: 1, AwayPlayer: "Lord Trooper", DateTime: "12/08 - 09:30", BGALink: "https://boardgamearena.com/tournament?id=423761", Played: true},
					{ID: 17, HomePlayer: "webbi", HomeScore: 0, AwayScore: 0, AwayPlayer: "alehrosario", DateTime: "", BGALink: "", Played: false},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Check that the table header includes DUELO column
	if !strings.Contains(view, "DUELO") {
		t.Error("Expected table header to contain 'DUELO' column")
	}

	// Check that match numbers are displayed in the table
	if !strings.Contains(view, "1") {
		t.Error("Expected table to show match number 1")
	}

	if !strings.Contains(view, "17") {
		t.Error("Expected table to show match number 17")
	}

	// Verify the table structure includes the match number as the first column after DUELO
	lines := strings.Split(view, "\n")
	var tableContentFound bool
	for _, line := range lines {
		// Look for a line that contains match data (has player names)
		if strings.Contains(line, "herchu") && strings.Contains(line, "Lord Trooper") {
			tableContentFound = true
			// The match number should appear before the player names
			if !strings.Contains(line, "1") {
				t.Error("Expected match number 1 to appear in the table row with herchu vs Lord Trooper")
			}
		}
		if strings.Contains(line, "webbi") && strings.Contains(line, "alehrosario") {
			// The match number should appear before the player names
			if !strings.Contains(line, "17") {
				t.Error("Expected match number 17 to appear in the table row with webbi vs alehrosario")
			}
		}
	}

	if !tableContentFound {
		t.Error("Expected to find table content with player names")
	}
}

func TestFixtureModel_View_HomeAwayColumnHeaders(t *testing.T) {
	// Create test data
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", HomeScore: 2, AwayScore: 1, AwayPlayer: "Lord Trooper", DateTime: "12/08 - 09:30", BGALink: "https://boardgamearena.com/tournament?id=423761", Played: true},
					{ID: 2, HomePlayer: "webbi", HomeScore: 0, AwayScore: 0, AwayPlayer: "alehrosario", DateTime: "", BGALink: "", Played: false},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	// Check that the table header uses HOME and AWAY instead of LOCAL and VISITOR
	if !strings.Contains(view, "HOME") {
		t.Error("Expected table header to contain 'HOME' column")
	}

	if !strings.Contains(view, "AWAY") {
		t.Error("Expected table header to contain 'AWAY' column")
	}

	if strings.Contains(view, "LOCAL") {
		t.Error("Expected table header to NOT contain 'LOCAL' column")
	}

	if strings.Contains(view, "VISITOR") {
		t.Error("Expected table header to NOT contain 'VISITOR' column")
	}
}

func TestFixtureModel_View_MatchNumberColumnDemo(t *testing.T) {
	// Create test data that mimics real tournament data
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", HomeScore: 2, AwayScore: 1, AwayPlayer: "Lord Trooper", DateTime: "12/08 - 09:30", BGALink: "https://boardgamearena.com/tournament?id=423761", Played: true},
					{ID: 2, HomePlayer: "webbi", HomeScore: 2, AwayScore: 0, AwayPlayer: "alehrosario", DateTime: "13/08 - 22:00", BGALink: "https://boardgamearena.com/tournament?id=423630", Played: true},
					{ID: 3, HomePlayer: "Academia47", HomeScore: 1, AwayScore: 2, AwayPlayer: "bignacho610", DateTime: "15/08 - 10:00", BGALink: "https://boardgamearena.com/tournament?id=424490", Played: true},
				},
			},
			{
				Number:    5,
				DateRange: "08/09 - 14/09",
				Matches: []*fixtures.Match{
					{ID: 17, HomePlayer: "webbi", HomeScore: 0, AwayScore: 0, AwayPlayer: "herchu", DateTime: "", BGALink: "", Played: false},
					{ID: 18, HomePlayer: "Academia47", HomeScore: 0, AwayScore: 0, AwayPlayer: "maticarrizoc", DateTime: "", BGALink: "", Played: false},
					{ID: 19, HomePlayer: "Nicoooo95", HomeScore: 0, AwayScore: 0, AwayPlayer: "alehrosario", DateTime: "", BGALink: "", Played: false},
					{ID: 20, HomePlayer: "bignacho610", HomeScore: 0, AwayScore: 0, AwayPlayer: "Lord Trooper", DateTime: "", BGALink: "", Played: false},
				},
			},
		},
	}

	model := NewFixtureModel(division)

	// Test Round 1 (played matches)
	view := model.View()
	t.Logf("Round 1 with played matches (showing DUELO column):\n%s", view)

	// Verify DUELO column shows correct match numbers
	if !strings.Contains(view, "DUELO") {
		t.Error("Expected DUELO column header")
	}

	// Navigate to Round 5 (unplayed matches)
	model.currentRound = 1
	view = model.View()
	t.Logf("Round 5 with unplayed matches (showing DUELO column):\n%s", view)

	// Check that both played and unplayed matches show their match numbers properly
	lines := strings.Split(view, "\n")
	foundMatch17 := false
	foundMatch20 := false

	for _, line := range lines {
		if strings.Contains(line, "webbi") && strings.Contains(line, "herchu") {
			if strings.Contains(line, "17") {
				foundMatch17 = true
			}
		}
		if strings.Contains(line, "bignacho610") && strings.Contains(line, "Lord Trooper") {
			if strings.Contains(line, "20") {
				foundMatch20 = true
			}
		}
	}

	if !foundMatch17 {
		t.Error("Expected to find match 17 in the table")
	}

	if !foundMatch20 {
		t.Error("Expected to find match 20 in the table")
	}
}

func TestFixtureModel_View_CompleteTableStructureDemo(t *testing.T) {
	// Create comprehensive test data showing the complete table structure
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{
				Number:    1,
				DateRange: "11/08 - 17/08",
				Matches: []*fixtures.Match{
					{ID: 1, HomePlayer: "herchu", HomeScore: 2, AwayScore: 1, AwayPlayer: "Lord Trooper", DateTime: "12/08 - 09:30", BGALink: "https://boardgamearena.com/tournament?id=423761", Played: true},
					{ID: 2, HomePlayer: "webbi", HomeScore: 2, AwayScore: 0, AwayPlayer: "alehrosario", DateTime: "13/08 - 22:00", BGALink: "https://boardgamearena.com/tournament?id=423630", Played: true},
					{ID: 17, HomePlayer: "Academia47", HomeScore: 0, AwayScore: 0, AwayPlayer: "bignacho610", DateTime: "", BGALink: "", Played: false},
				},
			},
		},
	}

	model := NewFixtureModel(division)
	view := model.View()

	t.Logf("Complete table structure with DUELO and HOME/AWAY columns:\n%s", view)

	// Verify all column headers are present
	requiredHeaders := []string{"DUELO", "PLAYED", "HOME", "AWAY", "RESULT", "DATE", "TOURNAMENT_ID"}
	for _, header := range requiredHeaders {
		if !strings.Contains(view, header) {
			t.Errorf("Expected table to contain header '%s'", header)
		}
	}

	// Verify old headers are not present
	deprecatedHeaders := []string{"LOCAL", "VISITOR"}
	for _, header := range deprecatedHeaders {
		if strings.Contains(view, header) {
			t.Errorf("Expected table to NOT contain deprecated header '%s'", header)
		}
	}

	// Verify match data is displayed correctly
	expectedData := []string{"1", "2", "17", "herchu", "Lord Trooper", "webbi", "alehrosario", "Academia47", "bignacho610"}
	for _, data := range expectedData {
		if !strings.Contains(view, data) {
			t.Errorf("Expected table to contain data '%s'", data)
		}
	}

	// Verify the table shows both played and unplayed matches
	if !strings.Contains(view, "✓") {
		t.Error("Expected table to show checkmark for played matches")
	}

	if !strings.Contains(view, "○") {
		t.Error("Expected table to show circle for unplayed matches")
	}
}
