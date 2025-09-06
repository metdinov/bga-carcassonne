package cli

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bubbledatetimepicker "github.com/lcc/bubble-datetime-picker"
)

// DateTimePickerModel represents the datetime picker for tournament scheduling
type DateTimePickerModel struct {
	picker       *bubbledatetimepicker.DateAndHourModel
	timezone     *time.Location
	selectedTime time.Time
	style        lipgloss.Style
	title        string
	instructions string
	homePlayer   string
	awayPlayer   string
	division     string
	roundNumber  int
	matchNumber  int
	matchID      int
	confirmed    bool
	canceled     bool
}

// DateTimeSelectedMsg is sent when a datetime is selected
type DateTimeSelectedMsg struct {
	DateTime    time.Time
	HomePlayer  string
	AwayPlayer  string
	Division    string
	RoundNumber int
	MatchNumber int
	MatchID     int
}

// DateTimePickerCanceledMsg is sent when datetime selection is canceled
type DateTimePickerCanceledMsg struct{}

// NewDateTimePickerModel creates a new datetime picker model
func NewDateTimePickerModel(
	homePlayer, awayPlayer, division string,
	roundNumber, matchNumber, matchID int,
) *DateTimePickerModel {
	// Get local timezone
	localTZ := time.Now().Location()

	// Create picker with default settings
	picker := bubbledatetimepicker.NewDateAndHourModel()

	title := fmt.Sprintf("Schedule Tournament: %s vs %s", homePlayer, awayPlayer)

	return &DateTimePickerModel{
		picker: &picker,
		title:  title,
		instructions: "Use ↑/↓ to change date, ←/→ to move between date/time, " +
			"Tab to switch fields, Enter to confirm, Esc to cancel",
		timezone:    localTZ,
		homePlayer:  homePlayer,
		awayPlayer:  awayPlayer,
		division:    division,
		roundNumber: roundNumber,
		matchNumber: matchNumber,
		matchID:     matchID,
		style: lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")),
	}
}

// Init initializes the datetime picker
func (m *DateTimePickerModel) Init() tea.Cmd {
	return m.picker.Init()
}

// dateTimePickerConfirmedMsg is an internal message to signal confirmation from the picker
type dateTimePickerConfirmedMsg struct{}

// Update handles messages for the datetime picker
func (m *DateTimePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var pickerModel tea.Model

	switch msg := msg.(type) {
	// Handle our internal confirmation message
	case dateTimePickerConfirmedMsg:
		m.confirmed = true
		m.selectedTime = m.picker.Time()

		return m, tea.Cmd(func() tea.Msg {
			return DateTimeSelectedMsg{
				DateTime:    m.selectedTime,
				HomePlayer:  m.homePlayer,
				AwayPlayer:  m.awayPlayer,
				Division:    m.division,
				RoundNumber: m.roundNumber,
				MatchNumber: m.matchNumber,
				MatchID:     m.matchID,
			}
		})

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			// Cancel selection
			m.canceled = true
			return m, tea.Cmd(func() tea.Msg {
				return DateTimePickerCanceledMsg{}
			})
		}
	}

	// Update the picker
	pickerModel, cmd = m.picker.Update(msg)
	if picker, ok := pickerModel.(*bubbledatetimepicker.DateAndHourModel); ok {
		m.picker = picker
	}

	// The picker returns a tea.Quit command when Enter is pressed on the time view.
	// We wrap the command to intercept the QuitMsg and convert it into our own
	// internal confirmation message.
	if cmd != nil {
		wrappedCmd := func() tea.Msg {
			msg := cmd() // execute original command
			if _, ok := msg.(tea.QuitMsg); ok {
				return dateTimePickerConfirmedMsg{}
			}
			return msg // return original message
		}
		return m, wrappedCmd
	}

	return m, cmd
}

// View renders the datetime picker
func (m *DateTimePickerModel) View() string {
	if m.confirmed || m.canceled {
		return ""
	}

	// Get current selected time for display
	currentTime := m.picker.Time()

	// Format timezone offset
	_, offset := currentTime.Zone()
	offsetHours := offset / 3600
	offsetMins := (offset % 3600) / 60

	var offsetStr string
	if offsetMins == 0 {
		offsetStr = fmt.Sprintf("UTC%+d", offsetHours)
	} else {
		offsetStr = fmt.Sprintf("UTC%+d:%02d", offsetHours, offsetMins)
	}

	// Build the view
	content := fmt.Sprintf("%s\n\n", m.title)
	content += fmt.Sprintf("Division: %s - Round %d - Duelo %d\n", m.division, m.roundNumber, m.matchNumber)
	content += fmt.Sprintf("Timezone: %s (%s)\n\n", m.timezone.String(), offsetStr)

	// Add the picker
	content += m.picker.View()

	content += fmt.Sprintf("\n\nSelected: %s",
		currentTime.Format("Monday, January 2, 2006 at 3:04 PM"))
	content += fmt.Sprintf(" (%s)", offsetStr)

	content += "\n\nNavigation:"
	content += "\n• ↑/↓ arrows: Change date or hour values"
	content += "\n• ←/→ arrows: Move between date and time fields"
	content += "\n• Tab/Space: Switch between components"
	content += "\n• Enter: Confirm selection | Esc: Cancel"

	return m.style.Render(content)
}

// GetSelectedTime returns the selected time
func (m *DateTimePickerModel) GetSelectedTime() time.Time {
	if m.selectedTime.IsZero() {
		return m.picker.Time()
	}
	return m.selectedTime
}

// IsConfirmed returns whether the datetime was confirmed
func (m *DateTimePickerModel) IsConfirmed() bool {
	return m.confirmed
}

// IsCanceled returns whether the datetime selection was canceled
func (m *DateTimePickerModel) IsCanceled() bool {
	return m.canceled
}

// FormatForBGA formats the selected time for BGA API
func (m *DateTimePickerModel) FormatForBGA() (date, timeStr string) {
	if m.selectedTime.IsZero() {
		now := time.Now()
		return now.Format("2006-01-02"), "21:00"
	}

	return m.selectedTime.Format("2006-01-02"), m.selectedTime.Format("15:04")
}
