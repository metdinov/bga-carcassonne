package cli

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"carca-cli/internal/fixtures"
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
	if !strings.Contains(view, "2 - 1") {
		t.Errorf("Expected view to contain score '2 - 1', got: %s", view)
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

	// Should show navigation instructions
	if !strings.Contains(view, "←/→") || !strings.Contains(view, "rounds") {
		t.Errorf("Expected view to show navigation instructions, got: %s", view)
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
