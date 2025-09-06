package cli

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDivisionModel_Init(t *testing.T) {
	model := NewDivisionModel()

	if model == nil {
		t.Fatal("Expected division model to be initialized")
	}

	if model.cursor != 0 {
		t.Errorf("Expected cursor to start at 0, got %d", model.cursor)
	}

	expectedDivisions := []string{
		"Elite",
		"Platinum A",
		"Platinum B",
		"Oro A",
		"Oro B",
		"Oro C",
		"Oro D",
	}

	if len(model.divisions) != len(expectedDivisions) {
		t.Errorf("Expected %d divisions, got %d", len(expectedDivisions), len(model.divisions))
	}

	for i, expected := range expectedDivisions {
		if model.divisions[i] != expected {
			t.Errorf("Expected division %d to be '%s', got '%s'", i, expected, model.divisions[i])
		}
	}
}

func TestDivisionModel_Update_Navigation(t *testing.T) {
	model := NewDivisionModel()

	// Test cursor down
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.cursor != 1 {
		t.Errorf("Expected cursor to move to 1, got %d", model.cursor)
	}
	if cmd != nil {
		t.Errorf("Expected no command on navigation, got %v", cmd)
	}

	// Test cursor up
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	if model.cursor != 0 {
		t.Errorf("Expected cursor to move back to 0, got %d", model.cursor)
	}

	// Test wrap around at bottom
	model.cursor = len(model.divisions) - 1
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.cursor != 0 {
		t.Errorf("Expected cursor to wrap to 0, got %d", model.cursor)
	}

	// Test wrap around at top
	model.cursor = 0
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	if model.cursor != len(model.divisions)-1 {
		t.Errorf("Expected cursor to wrap to %d, got %d", len(model.divisions)-1, model.cursor)
	}
}

func TestDivisionModel_Update_Selection(t *testing.T) {
	model := NewDivisionModel()
	model.cursor = 2 // Platinum B

	// Send enter key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Should return a command to switch to fixture view
	if cmd == nil {
		t.Fatal("Expected command when selecting division")
	}

	// The model should track the selected division
	divModel, ok := updatedModel.(*DivisionModel)
	if !ok {
		t.Fatal("Expected DivisionModel to be returned")
	}

	selectedDiv := divModel.GetSelectedDivision()
	if selectedDiv != "Platinum B" {
		t.Errorf("Expected selected division 'Platinum B', got '%s'", selectedDiv)
	}
}

func TestDivisionModel_Update_Back(t *testing.T) {
	model := NewDivisionModel()

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

func TestDivisionModel_View_ContainsTitle(t *testing.T) {
	model := NewDivisionModel()

	view := model.View()

	if !strings.Contains(view, "Select Division") {
		t.Errorf("Expected view to contain 'Select Division', got: %s", view)
	}
}

func TestDivisionModel_View_ShowsDivisions(t *testing.T) {
	model := NewDivisionModel()

	view := model.View()

	expectedDivisions := []string{"Elite", "Platinum A", "Platinum B", "Oro A", "Oro B", "Oro C", "Oro D"}
	for _, division := range expectedDivisions {
		if !strings.Contains(view, division) {
			t.Errorf("Expected view to contain '%s', got: %s", division, view)
		}
	}
}

func TestDivisionModel_View_ShowsCursor(t *testing.T) {
	model := NewDivisionModel()
	model.cursor = 3 // Oro A

	view := model.View()

	// Check that cursor is shown for selected item
	if !strings.Contains(view, "> Oro A") {
		t.Errorf("Expected view to show cursor on 'Oro A', got: %s", view)
	}

	// Check that other items don't have cursor
	if strings.Contains(view, "> Elite") {
		t.Errorf("Expected only selected item to have cursor, got: %s", view)
	}
}

func TestDivisionModel_GetSelectedDivision(t *testing.T) {
	model := NewDivisionModel()

	testCases := []struct {
		expected string
		cursor   int
	}{
		{expected: "Elite", cursor: 0},
		{expected: "Platinum A", cursor: 1},
		{expected: "Platinum B", cursor: 2},
		{expected: "Oro A", cursor: 3},
		{expected: "Oro B", cursor: 4},
		{expected: "Oro C", cursor: 5},
		{expected: "Oro D", cursor: 6},
	}

	for _, tc := range testCases {
		model.cursor = tc.cursor
		division := model.GetSelectedDivision()

		if division != tc.expected {
			t.Errorf("Expected division '%s' for cursor %d, got '%s'", tc.expected, tc.cursor, division)
		}
	}
}

func TestDivisionModel_GetFilename(t *testing.T) {
	model := NewDivisionModel()

	testCases := []struct {
		expected string
		cursor   int
	}{
		{expected: "data/Liga Argentina - 1° Temporada - E-Fixture.csv", cursor: 0},
		{expected: "data/Liga Argentina - 1° Temporada - P.A-Fixture.csv", cursor: 1},
		{expected: "data/Liga Argentina - 1° Temporada - P.B-Fixture.csv", cursor: 2},
		{expected: "data/Liga Argentina - 1° Temporada - O.A-Fixture.csv", cursor: 3},
		{expected: "data/Liga Argentina - 1° Temporada - O.B-Fixture.csv", cursor: 4},
		{expected: "data/Liga Argentina - 1° Temporada - O.C-Fixture.csv", cursor: 5},
		{expected: "data/Liga Argentina - 1° Temporada - O.D-Fixture.csv", cursor: 6},
	}

	for _, tc := range testCases {
		model.cursor = tc.cursor
		filename := model.GetSelectedFilename()

		if filename != tc.expected {
			t.Errorf("Expected filename '%s' for cursor %d, got '%s'", tc.expected, tc.cursor, filename)
		}
	}
}

func TestDivisionModel_Update_VimNavigation_Down(t *testing.T) {
	model := NewDivisionModel()

	// Send 'j' key (vim down)
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if model.cursor != 1 {
		t.Errorf("Expected cursor to move to 1 with 'j', got %d", model.cursor)
	}

	if cmd != nil {
		t.Errorf("Expected no command on vim navigation, got %v", cmd)
	}
}

func TestDivisionModel_Update_VimNavigation_Up(t *testing.T) {
	model := NewDivisionModel()
	model.cursor = 2 // Start at position 2

	// Send 'k' key (vim up)
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if model.cursor != 1 {
		t.Errorf("Expected cursor to move to 1 with 'k', got %d", model.cursor)
	}

	if cmd != nil {
		t.Errorf("Expected no command on vim navigation, got %v", cmd)
	}
}

func TestDivisionModel_Update_VimNavigation_WrapAround(t *testing.T) {
	model := NewDivisionModel()

	// Test 'j' wrap around at bottom
	model.cursor = len(model.divisions) - 1 // Last item
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if model.cursor != 0 {
		t.Errorf("Expected cursor to wrap to 0 with 'j', got %d", model.cursor)
	}

	// Test 'k' wrap around at top
	model.cursor = 0
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if model.cursor != len(model.divisions)-1 {
		t.Errorf("Expected cursor to wrap to %d with 'k', got %d", len(model.divisions)-1, model.cursor)
	}
}
