package bga

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Client handles BoardGameArena API interactions
type Client struct {
	httpClient *http.Client
	baseURL    string
	username   string
	password   string
	sessionID  string
}

// TournamentConfig represents the configuration for creating a tournament
type TournamentConfig struct {
	ChampionshipName string // Championship name
	TournamentName   string // Tournament name
	BaseDate         string // Base date (YYYY-MM-DD)
	BaseDateTime     string // Base date time (HH:MM)
	Division         string // Division name (Elite, Platinum A, etc.)
	LocalPlayer      string // Local player (home)
	VisitorPlayer    string // Visitor player (away)
	GameID           int    // 1 for Carcassonne
	MaxPlayers       int    // Maximum participants (2 for 1v1)
	MinPlayers       int    // Minimum participants (2 for 1v1)
	GameDuration     int    // Game duration in seconds (1800 for 30 min)
	MatchesCount     int    // Number of matches (3 for best-of-3)
	RoundNumber      int    // Round number
	MatchNumber      int    // Match number from fixture
}

// TournamentResponse represents the response from BGA tournament creation
type TournamentResponse struct {
	Link         string `json:"link"`
	Error        string `json:"error,omitempty"`
	TournamentID int    `json:"tournament_id"`
	Success      bool   `json:"success"`
}

// NewClient creates a new BGA client
func NewClient(username, password string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:  "https://boardgamearena.com",
		username: username,
		password: password,
	}
}

// Login authenticates with BGA and establishes a session
func (c *Client) Login() error {
	loginURL := c.baseURL + "/account/account/login.html"

	// Prepare login form data
	formData := url.Values{}
	formData.Set("email", c.username)
	formData.Set("password", c.password)
	formData.Set("form_id", "connection_form")
	formData.Set("request_id", strconv.FormatInt(time.Now().Unix(), 10))

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Carcassonne Tournament Manager/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform login request: %w", err)
	}
	defer resp.Body.Close()

	// Extract session information from cookies
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "TournamentSession" || cookie.Name == "PHPSESSID" {
			c.sessionID = cookie.Value
			break
		}
	}

	if c.sessionID == "" {
		return fmt.Errorf("failed to authenticate: no session cookie received")
	}

	return nil
}

// CreateTournament creates a new Swiss tournament on BGA
func (c *Client) CreateTournament(config *TournamentConfig) (*TournamentResponse, error) {
	if c.sessionID == "" {
		return nil, fmt.Errorf("not authenticated: call Login() first")
	}

	tournamentURL := c.baseURL + "/newtournament/newtournament/create.html"
	formData := c.buildTournamentForm(config)

	// Form validation
	formData.Set("form_id", "createnewtournament")
	formData.Set("dojo.preventCache", strconv.FormatInt(time.Now().UnixMilli(), 10))

	return c.submitTournamentRequest(tournamentURL, formData)
}

// CreateSwissTournament creates a best-of-3 Swiss tournament for two players
func (c *Client) CreateSwissTournament(
	division, homePlayer, awayPlayer string,
	roundNumber, matchNumber int,
) (*TournamentResponse, error) {
	// Get current date for tournament scheduling
	now := time.Now()
	baseDate := now.Format("2006-01-02")
	baseDateTime := "21:00" // Default to 9 PM

	// Format championship and tournament names according to requirements
	championshipName := fmt.Sprintf("Division %s - 1era Temporada", division)
	tournamentName := fmt.Sprintf("%d Fecha - Duelo %d - %s vs %s", roundNumber, matchNumber, homePlayer, awayPlayer)

	config := &TournamentConfig{
		GameID:           1,                // Carcassonne game ID
		ChampionshipName: championshipName, // Division X - 1era Temporada
		TournamentName:   tournamentName,   // X Fecha - Duelo Y - Player1 vs Player2
		MaxPlayers:       2,                // Exactly 2 players
		MinPlayers:       2,                // Minimum 2 players
		BaseDate:         baseDate,         // Current date
		BaseDateTime:     baseDateTime,     // 21:00
		GameDuration:     1800,             // 30 minutes (1800 seconds)
		MatchesCount:     3,                // Best-of-3
		Division:         division,         // Division name
		RoundNumber:      roundNumber,      // Round number
		MatchNumber:      matchNumber,      // Match number from fixture
		LocalPlayer:      homePlayer,       // Home/local player
		VisitorPlayer:    awayPlayer,       // Away/visitor player
	}

	return c.CreateTournament(config)
}

// buildTournamentForm constructs the form data for tournament creation
func (c *Client) buildTournamentForm(config *TournamentConfig) url.Values {
	formData := url.Values{}

	c.setBasicTournamentSettings(formData, config)
	c.setRegistrationSettings(formData, config)
	c.setTableAccessLevels(formData)
	c.setGeneralSettings(formData, config)
	c.setCarcassonneGameOptions(formData)
	c.setSwissSystemOptions(formData, config)
	c.setPlayerConfirmation(formData)

	return formData
}

// setBasicTournamentSettings sets the basic tournament configuration
func (c *Client) setBasicTournamentSettings(formData url.Values, config *TournamentConfig) {
	formData.Set("game", strconv.Itoa(config.GameID))
	formData.Set("championship_name", config.ChampionshipName)
	formData.Set("tournament_name", config.TournamentName)
	formData.Set("base_date", config.BaseDate)
	formData.Set("base_date_hour", config.BaseDateTime)
}

// setRegistrationSettings configures tournament registration options
func (c *Client) setRegistrationSettings(formData url.Values, config *TournamentConfig) {
	formData.Set("registration_type", "invitation_only")
	formData.Set("registration_group", "0")
	formData.Set("registration_starts", "30")
	formData.Set("min_players", strconv.Itoa(config.MinPlayers))
	formData.Set("max_players", strconv.Itoa(config.MaxPlayers))
}

// setTableAccessLevels enables all skill levels for tournament access
func (c *Client) setTableAccessLevels(formData url.Values) {
	formData.Set("tableaccess_Averagelevel", "on")
	formData.Set("tableaccess_Goodplayers", "on")
	formData.Set("tableaccess_Strongplayers", "on")
	formData.Set("tableaccess_Experts", "on")
	formData.Set("tableaccess_Masters", "on")
}

// setGeneralSettings configures general tournament settings
func (c *Client) setGeneralSettings(formData url.Values, config *TournamentConfig) {
	formData.Set("karma", "1")
	formData.Set("restrictedCountries", "")
	formData.Set("stage_type", "swissSystemV2")
	formData.Set("game_max_duration", strconv.Itoa(config.GameDuration))
	formData.Set("players_out_of_time", "vote_kick")
}

// setCarcassonneGameOptions sets game-specific options for Carcassonne
func (c *Client) setCarcassonneGameOptions(formData url.Values) {
	// Field scoring: international (3pts per city)
	formData.Set("gameoption_200", "5")
	// City scoring: international (4pts per two tile city)
	formData.Set("gameoption_204", "900")
	// No River expansion
	formData.Set("gameoption_201", "0")
	// No Inns & Cathedrals
	formData.Set("gameoption_206", "0")
	// No other expansions
	formData.Set("gameoption_106", "0")
	formData.Set("gameoption_103", "0")
	formData.Set("gameoption_104", "0")
	formData.Set("gameoption_107", "0")
	formData.Set("gameoption_100", "0")
	formData.Set("gameoption_101", "0")
	formData.Set("gameoption_102", "0")

	// Stage settings
	formData.Set("stage_playernbr", "2")
	formData.Set("stage_playernbr_min", "2")
}

// setSwissSystemOptions configures Swiss system tournament options
func (c *Client) setSwissSystemOptions(formData url.Values, config *TournamentConfig) {
	// Swiss System V2 mode options (best-of-3)
	formData.Set("mode_option_swissSystemV2_100", "1")
	formData.Set("mode_option_swissSystemV2_101", "100")
	formData.Set("mode_option_swissSystemV2_102", "1")
	formData.Set("mode_option_swissSystemV2_103", strconv.Itoa(config.MatchesCount))

	// Additional mode options
	formData.Set("mode_option_bracketElimination_100", "1")
	formData.Set("mode_option_twoStage_100", "1")
	formData.Set("mode_option_twoStage_101", "8")
	formData.Set("mode_option_twoStage_102", "2")
	formData.Set("mode_option_twoStage_103", "1")

	// Stage 1 options
	formData.Set("stage_1_mode_option_bracketElimination_100", "1")
	formData.Set("stage_1_mode_option_swissSystemV2_100", "1")
	formData.Set("stage_1_mode_option_swissSystemV2_101", "100")
	formData.Set("stage_1_mode_option_swissSystemV2_102", "1")
	formData.Set("stage_1_mode_option_swissSystemV2_103", "5")
	formData.Set("stage_1_mode_option_twoStage_100", "1")
	formData.Set("stage_1_mode_option_twoStage_101", "8")
	formData.Set("stage_1_mode_option_twoStage_102", "2")
	formData.Set("stage_1_mode_option_twoStage_103", "1")

	// Stage 2 options
	formData.Set("stage_2_mode_option_bracketElimination_100", "1")
	formData.Set("stage_2_mode_option_swissSystemV2_100", "1")
	formData.Set("stage_2_mode_option_swissSystemV2_101", "100")
	formData.Set("stage_2_mode_option_swissSystemV2_102", "1")
	formData.Set("stage_2_mode_option_swissSystemV2_103", "5")
	formData.Set("stage_2_mode_option_twoStage_100", "1")
	formData.Set("stage_2_mode_option_twoStage_101", "8")
	formData.Set("stage_2_mode_option_twoStage_102", "2")
	formData.Set("stage_2_mode_option_twoStage_103", "1")
}

// setPlayerConfirmation requires player confirmation for tournament participation
func (c *Client) setPlayerConfirmation(formData url.Values) {
	formData.Set("players", "confirm_players")
}

// submitTournamentRequest handles the HTTP request and response parsing
func (c *Client) submitTournamentRequest(tournamentURL string, formData url.Values) (*TournamentResponse, error) {
	req, err := http.NewRequest("POST", tournamentURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setRequestHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tournament creation request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tournament creation failed with status %d: %s", resp.StatusCode, string(body))
	}

	return c.parseTournamentResponse(string(body))
}

// setRequestHeaders configures HTTP headers for the tournament creation request
func (c *Client) setRequestHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Carcassonne Tournament Manager/1.0")
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", c.sessionID))
}

// parseTournamentResponse extracts tournament information from BGA response
func (c *Client) parseTournamentResponse(bodyStr string) (*TournamentResponse, error) {
	var tournamentResp TournamentResponse

	// Try to decode JSON response first
	if err := json.Unmarshal([]byte(bodyStr), &tournamentResp); err == nil {
		if tournamentResp.TournamentID > 0 {
			tournamentResp.Link = fmt.Sprintf("%s/tournament?id=%d", c.baseURL, tournamentResp.TournamentID)
		}
		return &tournamentResp, nil
	}

	// Fallback to text parsing
	if strings.Contains(bodyStr, "successfully created") ||
		strings.Contains(bodyStr, "Tournament created") ||
		strings.Contains(bodyStr, "tournament has been created") {
		tournamentResp.Success = true
	} else {
		tournamentResp.Success = false
		tournamentResp.Error = "Tournament creation failed"
		return &tournamentResp, nil
	}

	tournamentResp.TournamentID = c.extractTournamentID(bodyStr)
	if tournamentResp.TournamentID > 0 {
		tournamentResp.Link = fmt.Sprintf("%s/tournament?id=%d", c.baseURL, tournamentResp.TournamentID)
	}

	return &tournamentResp, nil
}

// extractTournamentID finds the tournament ID in the response body
func (c *Client) extractTournamentID(bodyStr string) int {
	// Look for patterns like: tournament.php?id=123456 or tournament?id=123456
	tournamentIDRegex := regexp.MustCompile(`tournament(?:\.php)?\?id=(\d+)`)
	matches := tournamentIDRegex.FindStringSubmatch(bodyStr)
	if len(matches) > 1 {
		if id, err := strconv.Atoi(matches[1]); err == nil {
			return id
		}
	}

	// Try alternative patterns - JSON format
	jsonRegex := regexp.MustCompile(`"tournament_id":\s*(\d+)`)
	matches = jsonRegex.FindStringSubmatch(bodyStr)
	if len(matches) > 1 {
		if id, err := strconv.Atoi(matches[1]); err == nil {
			return id
		}
	}

	return 0
}

// CreateSwissTournamentWithDateTime creates a best-of-3 Swiss tournament for two players with specific datetime
func (c *Client) CreateSwissTournamentWithDateTime(
	division, homePlayer, awayPlayer string,
	roundNumber, matchNumber int,
	scheduledTime time.Time,
) (*TournamentResponse, error) {
	// Format the scheduled time for BGA
	baseDate := scheduledTime.Format("2006-01-02")
	baseDateTime := scheduledTime.Format("15:04")

	// Format championship and tournament names according to requirements
	championshipName := fmt.Sprintf("Division %s - 1era Temporada", division)
	tournamentName := fmt.Sprintf("%d Fecha - Duelo %d - %s vs %s", roundNumber, matchNumber, homePlayer, awayPlayer)

	config := &TournamentConfig{
		GameID:           1,                // Carcassonne game ID
		ChampionshipName: championshipName, // Division X - 1era Temporada
		TournamentName:   tournamentName,   // X Fecha - Duelo Y - Player1 vs Player2
		MaxPlayers:       2,                // Exactly 2 players
		MinPlayers:       2,                // Minimum 2 players
		BaseDate:         baseDate,         // Scheduled date
		BaseDateTime:     baseDateTime,     // Scheduled time
		GameDuration:     1800,             // 30 minutes (1800 seconds)
		MatchesCount:     3,                // Best-of-3
		Division:         division,         // Division name
		RoundNumber:      roundNumber,      // Round number
		MatchNumber:      matchNumber,      // Match number from fixture
		LocalPlayer:      homePlayer,       // Home/local player
		VisitorPlayer:    awayPlayer,       // Away/visitor player
	}

	return c.CreateTournament(config)
}

// GetTournamentStatus retrieves the current status of a tournament
func (c *Client) GetTournamentStatus(tournamentID int) (*TournamentStatus, error) {
	if c.sessionID == "" {
		return nil, fmt.Errorf("not authenticated: call Login() first")
	}

	statusURL := fmt.Sprintf("%s/tournament/tournament/tournamentStatus.html?id=%d", c.baseURL, tournamentID)

	req, err := http.NewRequest("GET", statusURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create status request: %w", err)
	}

	req.Header.Set("User-Agent", "Carcassonne Tournament Manager/1.0")
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", c.sessionID))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament status: %w", err)
	}
	defer resp.Body.Close()

	var status TournamentStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to parse tournament status: %w", err)
	}

	return &status, nil
}

// TournamentStatus represents the current status of a tournament
type TournamentStatus struct {
	Name         string         `json:"name"`
	Status       string         `json:"status"`  // "waiting", "in_progress", "finished"
	Results      map[string]int `json:"results"` // player -> score
	Matches      []MatchStatus  `json:"matches"`
	ID           int            `json:"id"`
	PlayersCount int            `json:"players_count"`
}

// MatchStatus represents the status of a single match within a tournament
type MatchStatus struct {
	Status     string `json:"status"` // "waiting", "in_progress", "finished"
	HomePlayer string `json:"home_player"`
	AwayPlayer string `json:"away_player"`
	Winner     string `json:"winner,omitempty"`
	ID         int    `json:"id"`
	HomeScore  int    `json:"home_score"`
	AwayScore  int    `json:"away_score"`
}

// IsAuthenticated checks if the client has a valid session
func (c *Client) IsAuthenticated() bool {
	return c.sessionID != ""
}

// Logout terminates the current session
func (c *Client) Logout() error {
	if c.sessionID == "" {
		return nil // Already logged out
	}

	logoutURL := c.baseURL + "/account/account/logout.html"

	req, err := http.NewRequest("POST", logoutURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", c.sessionID))

	_, err = c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	c.sessionID = ""
	return nil
}
