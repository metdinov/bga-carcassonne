package cli

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMenuModel_Init(t *testing.T) {
	model := NewMenuModel()

	if model == nil {
		t.Fatal("Expected menu model to be initialized")
	}

	if model.cursor != 0 {
		t.Errorf("Expected cursor to start at 0, got %d", model.cursor)
	}

	if len(model.choices) != 4 {
		t.Errorf("Expected 4 menu choices, got %d", len(model.choices))
	}

	expectedChoices := []string{
		"Create Tournament",
		"View Fixture",
		"View Positions",
		"Exit",
	}

	for i, expected := range expectedChoices {
		if model.choices[i] != expected {
			t.Errorf("Expected choice %d to be '%s', got '%s'", i, expected, model.choices[i])
		}
	}
}

func TestMenuModel_Update_CursorDown(t *testing.T) {
	model := NewMenuModel()

	// Send arrow down key
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})

	if model.cursor != 1 {
		t.Errorf("Expected cursor to move to 1, got %d", model.cursor)
	}

	if cmd != nil {
		t.Errorf("Expected no command on cursor movement, got %v", cmd)
	}
}

func TestMenuModel_Update_CursorUp(t *testing.T) {
	model := NewMenuModel()
	model.cursor = 2 // Start at position 2

	// Send arrow up key
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})

	if model.cursor != 1 {
		t.Errorf("Expected cursor to move to 1, got %d", model.cursor)
	}

	if cmd != nil {
		t.Errorf("Expected no command on cursor movement, got %v", cmd)
	}
}

func TestMenuModel_Update_CursorWrapAround(t *testing.T) {
	model := NewMenuModel()

	// Test wrap around at bottom
	model.cursor = len(model.choices) - 1 // Last item
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})

	if model.cursor != 0 {
		t.Errorf("Expected cursor to wrap to 0, got %d", model.cursor)
	}

	// Test wrap around at top
	model.cursor = 0
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})

	if model.cursor != len(model.choices)-1 {
		t.Errorf("Expected cursor to wrap to %d, got %d", len(model.choices)-1, model.cursor)
	}

	if cmd != nil {
		t.Errorf("Expected no command on cursor movement, got %v", cmd)
	}
}

func TestMenuModel_Update_SelectExit(t *testing.T) {
	model := NewMenuModel()
	model.cursor = 3 // Exit option

	// Send enter key
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("Expected quit command when selecting Exit")
	}

	// We just verify a command was returned - the actual quit behavior
	// is handled by the Bubble Tea runtime
}

func TestMenuModel_Update_CtrlC(t *testing.T) {
	model := NewMenuModel()

	// Send Ctrl+C
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Fatal("Expected quit command on Ctrl+C")
	}

	// We just verify a command was returned - the actual quit behavior
	// is handled by the Bubble Tea runtime
}

func TestMenuModel_View_ContainsTitle(t *testing.T) {
	model := NewMenuModel()

	view := model.View()

	if !strings.Contains(view, "Carcassonne Tournament Manager") {
		t.Errorf("Expected view to contain title, got: %s", view)
	}
}

func TestMenuModel_View_ShowsCursor(t *testing.T) {
	model := NewMenuModel()
	model.cursor = 1 // Second item

	view := model.View()

	// Check that cursor is shown for selected item
	if !strings.Contains(view, "> View Fixture") {
		t.Errorf("Expected view to show cursor on selected item, got: %s", view)
	}

	// Check that other items don't have cursor
	if strings.Contains(view, "> Create Tournament") {
		t.Errorf("Expected only selected item to have cursor, got: %s", view)
	}
}

func TestMenuModel_GetSelectedChoice(t *testing.T) {
	model := NewMenuModel()

	testCases := []struct {
		expected string
		cursor   int
	}{
		{expected: "Create Tournament", cursor: 0},
		{expected: "View Fixture", cursor: 1},
		{expected: "View Positions", cursor: 2},
		{expected: "Exit", cursor: 3},
	}

	for _, tc := range testCases {
		model.cursor = tc.cursor
		choice := model.GetSelectedChoice()

		if choice != tc.expected {
			t.Errorf("Expected choice '%s' for cursor %d, got '%s'", tc.expected, tc.cursor, choice)
		}
	}
}

func TestMenuModel_Update_SelectViewFixture(t *testing.T) {
	model := NewMenuModel()
	model.cursor = 1 // View Fixture option

	// Send enter key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Should return a command to switch to division selection
	if cmd == nil {
		t.Fatal("Expected command when selecting View Fixture")
	}

	// The model should remain a MenuModel
	_, ok := updatedModel.(*MenuModel)
	if !ok {
		t.Fatal("Expected MenuModel to be returned")
	}
}

func TestMenuModel_HandleViewFixtureFlow(t *testing.T) {
	model := NewMenuModel()

	// Test that View Fixture selection returns appropriate command
	model.cursor = 1 // View Fixture
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("Expected command for View Fixture selection")
	}
}

func TestMenuModel_Update_QuitWithQ(t *testing.T) {
	model := NewMenuModel()

	// Send 'q' key
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	if cmd == nil {
		t.Fatal("Expected quit command when pressing 'q'")
	}

	// We just verify a command was returned - the actual quit behavior
	// is handled by the Bubble Tea runtime
}

func TestMenuModel_View_ShowsQuitInstructions(t *testing.T) {
	model := NewMenuModel()

	view := model.View()

	if !strings.Contains(view, "Press q/Ctrl+C to quit, ↑/↓ or j/k to navigate") {
		t.Errorf("Expected view to show quit and vim navigation instructions, got: %s", view)
	}
}

func TestMenuModel_Update_VimNavigation_Down(t *testing.T) {
	model := NewMenuModel()

	// Send 'j' key (vim down)
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if model.cursor != 1 {
		t.Errorf("Expected cursor to move to 1 with 'j', got %d", model.cursor)
	}

	if cmd != nil {
		t.Errorf("Expected no command on vim navigation, got %v", cmd)
	}
}

func TestMenuModel_Update_VimNavigation_Up(t *testing.T) {
	model := NewMenuModel()
	model.cursor = 2 // Start at position 2

	// Send 'k' key (vim up)
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if model.cursor != 1 {
		t.Errorf("Expected cursor to move to 1 with 'k', got %d", model.cursor)
	}

	if cmd != nil {
		t.Errorf("Expected no command on vim navigation, got %v", cmd)
	}
}

func TestMenuModel_Update_VimNavigation_WrapAround(t *testing.T) {
	model := NewMenuModel()

	// Test 'j' wrap around at bottom
	model.cursor = len(model.choices) - 1 // Last item
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if model.cursor != 0 {
		t.Errorf("Expected cursor to wrap to 0 with 'j', got %d", model.cursor)
	}

	// Test 'k' wrap around at top
	model.cursor = 0
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if model.cursor != len(model.choices)-1 {
		t.Errorf("Expected cursor to wrap to %d with 'k', got %d", len(model.choices)-1, model.cursor)
	}

	if cmd != nil {
		t.Errorf("Expected no command on vim navigation, got %v", cmd)
	}
}
