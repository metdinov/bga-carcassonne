package cli

import (
	tea "github.com/charmbracelet/bubbletea"

	"carca-cli/internal/bga"
	"carca-cli/internal/fixtures"
)

// Screen represents the different screens in the app
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenDivisionSelect
	ScreenFixture
)

// ViewFixtureSelectMsg is sent when user selects "View Fixture" from main menu
type ViewFixtureSelectMsg struct{}

// AppModel coordinates navigation between different screens
type AppModel struct {
	menuModel     *MenuModel
	divisionModel *DivisionModel
	fixtureModel  *FixtureModel
	currentScreen Screen
}

// NewAppModel creates a new app coordinator model
func NewAppModel() *AppModel {
	return &AppModel{
		currentScreen: ScreenMenu,
		menuModel:     NewMenuModel(),
	}
}

// Init initializes the app model (required by Bubble Tea)
func (m *AppModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and manages screen transitions
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ViewFixtureSelectMsg:
		// Transition from menu to division selection
		m.currentScreen = ScreenDivisionSelect
		m.divisionModel = NewDivisionModel()

		return m, nil

	case DivisionSelectMsg:
		// Transition from division selection to fixture display
		m.currentScreen = ScreenFixture

		// Load fixture data
		division, err := fixtures.ParseFixtureFile(msg.Filename)
		if err != nil {
			// If loading fails, show error in fixture model
			m.fixtureModel = NewFixtureModel(&fixtures.Division{
				Name:   msg.Division,
				Rounds: []*fixtures.Round{},
			})
		} else {
			m.fixtureModel = NewFixtureModel(division)
		}

		// Set up BGA client with mock client for now
		// In production, this would be a real client
		mockClient := bga.NewMockClient("", "")
		m.fixtureModel.SetBGAClient(mockClient)

		return m, nil

	case BackToMenuMsg:
		// Go back to main menu from any screen
		m.currentScreen = ScreenMenu
		// Clear other models to free memory
		m.divisionModel = nil
		m.fixtureModel = nil

		return m, nil

	default:
		// Delegate to current screen's model
		switch m.currentScreen {
		case ScreenMenu:
			if m.menuModel != nil {
				updatedModel, cmd := m.menuModel.Update(msg)
				if menuModel, ok := updatedModel.(*MenuModel); ok {
					m.menuModel = menuModel
				}

				// Check if user selected "View Fixture"
				if cmd != nil && m.menuModel.GetSelectedChoice() == "View Fixture" {
					if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
						// Trigger transition to division selection
						return m.Update(ViewFixtureSelectMsg{})
					}
				}

				return m, cmd
			}

		case ScreenDivisionSelect:
			if m.divisionModel != nil {
				updatedModel, cmd := m.divisionModel.Update(msg)
				if divModel, ok := updatedModel.(*DivisionModel); ok {
					m.divisionModel = divModel
				}

				return m, cmd
			}

		case ScreenFixture:
			if m.fixtureModel != nil {
				updatedModel, cmd := m.fixtureModel.Update(msg)
				if fixModel, ok := updatedModel.(*FixtureModel); ok {
					m.fixtureModel = fixModel
				}

				return m, cmd
			}
		}
	}

	return m, nil
}

// View renders the current screen
func (m *AppModel) View() string {
	switch m.currentScreen {
	case ScreenMenu:
		if m.menuModel != nil {
			return m.menuModel.View()
		}

		return "Loading menu...\n"

	case ScreenDivisionSelect:
		if m.divisionModel != nil {
			return m.divisionModel.View()
		}

		return "Loading division selection...\n"

	case ScreenFixture:
		if m.fixtureModel != nil {
			return m.fixtureModel.View()
		}

		return "Loading fixture data...\n\nPress esc/q to go back.\n"

	default:
		return "Unknown screen\n"
	}
}

// GetCurrentScreen returns the currently active screen
func (m *AppModel) GetCurrentScreen() Screen {
	return m.currentScreen
}
