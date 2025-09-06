package cli

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"carca-cli/internal/fixtures"
)

func TestAppModel_Init(t *testing.T) {
	model := NewAppModel()

	if model == nil {
		t.Fatal("Expected app model to be initialized")
	}

	if model.currentScreen != ScreenMenu {
		t.Errorf("Expected initial screen to be ScreenMenu, got %v", model.currentScreen)
	}

	if model.menuModel == nil {
		t.Fatal("Expected menu model to be initialized")
	}
}

func TestAppModel_Update_MenuToDivisionSelection(t *testing.T) {
	model := NewAppModel()

	// Simulate selecting "View Fixture" from menu
	msg := ViewFixtureSelectMsg{}
	updatedModel, cmd := model.Update(msg)

	appModel, ok := updatedModel.(*AppModel)
	if !ok {
		t.Fatal("Expected AppModel to be returned")
	}

	if appModel.currentScreen != ScreenDivisionSelect {
		t.Errorf("Expected screen to change to ScreenDivisionSelect, got %v", appModel.currentScreen)
	}

	if appModel.divisionModel == nil {
		t.Fatal("Expected division model to be initialized")
	}

	if cmd != nil {
		t.Errorf("Expected no command on screen transition, got %v", cmd)
	}
}

func TestAppModel_Update_DivisionToFixture(t *testing.T) {
	model := NewAppModel()
	model.currentScreen = ScreenDivisionSelect
	model.divisionModel = NewDivisionModel()

	// Simulate selecting a division
	msg := DivisionSelectMsg{
		Division: "Elite",
		Filename: "data/Liga Argentina - 1° Temporada - E-Fixture.csv",
	}

	updatedModel, cmd := model.Update(msg)

	appModel, ok := updatedModel.(*AppModel)
	if !ok {
		t.Fatal("Expected AppModel to be returned")
	}

	if appModel.currentScreen != ScreenFixture {
		t.Errorf("Expected screen to change to ScreenFixture, got %v", appModel.currentScreen)
	}

	// Note: fixtureModel might be nil if file loading fails, which is ok for test
	if cmd != nil {
		t.Errorf("Expected no command on screen transition, got %v", cmd)
	}
}

func TestAppModel_Update_BackToMenu(t *testing.T) {
	model := NewAppModel()
	model.currentScreen = ScreenDivisionSelect

	// Simulate back to menu
	msg := BackToMenuMsg{}
	updatedModel, cmd := model.Update(msg)

	appModel, ok := updatedModel.(*AppModel)
	if !ok {
		t.Fatal("Expected AppModel to be returned")
	}

	if appModel.currentScreen != ScreenMenu {
		t.Errorf("Expected screen to change back to ScreenMenu, got %v", appModel.currentScreen)
	}

	if cmd != nil {
		t.Errorf("Expected no command on screen transition, got %v", cmd)
	}
}

func TestAppModel_Update_KeyDelegation(t *testing.T) {
	model := NewAppModel()

	// Test that key messages are delegated to current screen
	keyMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(keyMsg)

	appModel, ok := updatedModel.(*AppModel)
	if !ok {
		t.Fatal("Expected AppModel to be returned")
	}

	// Should still be on menu screen
	if appModel.currentScreen != ScreenMenu {
		t.Errorf("Expected to remain on ScreenMenu, got %v", appModel.currentScreen)
	}
}

func TestAppModel_View_ShowsCurrentScreen(t *testing.T) {
	model := NewAppModel()

	// Test menu screen view
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view for menu screen")
	}

	// Test division selection screen view
	model.currentScreen = ScreenDivisionSelect
	model.divisionModel = NewDivisionModel()
	view = model.View()
	if view == "" {
		t.Error("Expected non-empty view for division selection screen")
	}
}

func TestAppModel_View_HandlesNilModels(t *testing.T) {
	model := NewAppModel()

	// Test fixture screen with nil model
	model.currentScreen = ScreenFixture
	model.fixtureModel = nil

	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view even with nil fixture model")
	}
}

func TestAppModel_GetCurrentScreen(t *testing.T) {
	model := NewAppModel()

	testCases := []struct {
		screen   Screen
		expected Screen
	}{
		{ScreenMenu, ScreenMenu},
		{ScreenDivisionSelect, ScreenDivisionSelect},
		{ScreenFixture, ScreenFixture},
	}

	for _, tc := range testCases {
		model.currentScreen = tc.screen
		current := model.GetCurrentScreen()

		if current != tc.expected {
			t.Errorf("Expected screen %v, got %v", tc.expected, current)
		}
	}
}

func TestAppModel_LoadFixtureData_Success(t *testing.T) {
	model := NewAppModel()

	// Test with a valid division and mock filename
	division := &fixtures.Division{
		Name: "Elite",
		Rounds: []*fixtures.Round{
			{Number: 1, DateRange: "11/08 - 17/08", Matches: []*fixtures.Match{}},
		},
	}

	model.fixtureModel = NewFixtureModel(division)

	if model.fixtureModel == nil {
		t.Fatal("Expected fixture model to be set")
	}

	if model.fixtureModel.division != division {
		t.Error("Expected division to be set correctly in fixture model")
	}
}
