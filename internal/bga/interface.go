package bga

import "time"

// APIClient defines the interface for interacting with BoardGameArena
type APIClient interface {
	// Login authenticates with BGA and establishes a session
	Login() error

	// CreateTournament creates a new tournament with the given configuration
	CreateTournament(config *TournamentConfig) (*TournamentResponse, error)

	// CreateSwissTournament creates a best-of-3 Swiss tournament for two players
	CreateSwissTournament(
		division, homePlayer, awayPlayer string,
		roundNumber, matchNumber int,
	) (*TournamentResponse, error)

	// CreateSwissTournamentWithDateTime creates a best-of-3 Swiss tournament for two players with specific datetime
	CreateSwissTournamentWithDateTime(
		division, homePlayer, awayPlayer string,
		roundNumber, matchNumber int,
		scheduledTime time.Time,
	) (*TournamentResponse, error)

	// GetTournamentStatus retrieves the current status of a tournament
	GetTournamentStatus(tournamentID int) (*TournamentStatus, error)

	// IsAuthenticated checks if the client has a valid session
	IsAuthenticated() bool

	// Logout terminates the current session
	Logout() error
}

// Ensure both implementations satisfy the interface
var (
	_ APIClient = (*Client)(nil)
	_ APIClient = (*MockClient)(nil)
)
