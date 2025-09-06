package fixtures

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Division represents a complete tournament division with all rounds
type Division struct {
	Name   string
	Rounds []*Round
}

// Round represents a tournament round with multiple matches
type Round struct {
	DateRange string
	Matches   []*Match
	Number    int
}

// Match represents a tournament match between two players
type Match struct {
	HomePlayer string
	AwayPlayer string
	DateTime   string
	BGALink    string
	ID         int
	HomeScore  int
	AwayScore  int
	Played     bool
}

// ParseMatch parses a CSV line into a Match struct
func ParseMatch(csvLine string) (*Match, error) {
	reader := csv.NewReader(strings.NewReader(csvLine))

	records, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV line: %w", err)
	}

	if len(records) < 10 {
		return nil, fmt.Errorf("invalid CSV format: expected at least 10 fields, got %d", len(records))
	}

	id, err := strconv.Atoi(records[0])
	if err != nil {
		return nil, fmt.Errorf("invalid match ID: %w", err)
	}

	homeScore, err := strconv.Atoi(records[2])
	if err != nil {
		return nil, fmt.Errorf("invalid home score: %w", err)
	}

	awayScore, err := strconv.Atoi(records[3])
	if err != nil {
		return nil, fmt.Errorf("invalid away score: %w", err)
	}

	played := records[8] == "1"

	match := &Match{
		ID:         id,
		HomePlayer: records[1],
		HomeScore:  homeScore,
		AwayScore:  awayScore,
		AwayPlayer: records[4],
		DateTime:   records[5],
		BGALink:    records[6],
		Played:     played,
	}

	return match, nil
}

// ParseRound parses CSV data containing a round header and matches
func ParseRound(csvData string) (*Round, error) {
	lines := strings.Split(csvData, "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid round data: need at least header and one match")
	}

	// Parse header line to extract round number and date range
	headerReader := csv.NewReader(strings.NewReader(lines[0]))

	headerRecord, err := headerReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to parse round header: %w", err)
	}

	if len(headerRecord) < 6 {
		return nil, fmt.Errorf("invalid header format")
	}

	// Extract round number from "Fecha X" format
	roundNumber := 1

	if len(headerRecord[1]) > 6 && strings.HasPrefix(headerRecord[1], "Fecha ") {
		if num, err := strconv.Atoi(strings.TrimSpace(headerRecord[1][6:])); err == nil {
			roundNumber = num
		}
	}

	dateRange := headerRecord[5]

	round := &Round{
		Number:    roundNumber,
		DateRange: dateRange,
		Matches:   make([]*Match, 0),
	}

	// Parse match lines (skip header)
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		match, err := ParseMatch(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse match on line %d: %w", i+1, err)
		}

		round.Matches = append(round.Matches, match)
	}

	return round, nil
}

// ParseDivision parses complete CSV data containing multiple rounds separated by empty lines
func ParseDivision(csvData string) (*Division, error) {
	lines := strings.Split(csvData, "\n")
	division := &Division{
		Rounds: make([]*Round, 0),
	}

	var currentRoundLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check if this line starts a new round (header line)
		if strings.HasPrefix(line, "Duelo,Fecha") {
			// If we have accumulated lines for a previous round, process them
			if len(currentRoundLines) > 0 {
				roundData := strings.Join(currentRoundLines, "\n")

				round, err := ParseRound(roundData)
				if err != nil {
					return nil, fmt.Errorf("failed to parse round: %w", err)
				}

				division.Rounds = append(division.Rounds, round)
				currentRoundLines = nil
			}
			// Start new round
			currentRoundLines = append(currentRoundLines, line)
		} else if line == "" || strings.Contains(line, ",,,,,,,,,,") {
			// Empty line or separator line - skip
			continue
		} else if len(currentRoundLines) > 0 {
			// Add match line to current round
			currentRoundLines = append(currentRoundLines, line)
		}
	}

	// Process the last round if it exists
	if len(currentRoundLines) > 0 {
		roundData := strings.Join(currentRoundLines, "\n")

		round, err := ParseRound(roundData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse final round: %w", err)
		}

		division.Rounds = append(division.Rounds, round)
	}

	return division, nil
}

// ParseFixtureFile reads a CSV file and parses it into a Division
func ParseFixtureFile(filename string) (*Division, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture file %s: %w", filename, err)
	}

	division, err := ParseDivision(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse fixture file %s: %w", filename, err)
	}

	// Extract division name from filename
	// e.g., "Liga Argentina - 1° Temporada - E-Fixture.csv" -> "E"
	if strings.Contains(filename, " - ") && strings.Contains(filename, "-Fixture.csv") {
		parts := strings.Split(filename, " - ")
		if len(parts) >= 3 {
			namePart := parts[len(parts)-1]
			if strings.HasSuffix(namePart, "-Fixture.csv") {
				division.Name = strings.TrimSuffix(namePart, "-Fixture.csv")
			}
		}
	}

	return division, nil
}

// GetUnplayedMatches returns all matches from a division that haven't been played yet
func GetUnplayedMatches(division *Division) []*Match {
	var unplayed []*Match

	for _, round := range division.Rounds {
		for _, match := range round.Matches {
			if !match.Played {
				unplayed = append(unplayed, match)
			}
		}
	}

	return unplayed
}

// DemoParseFixtures showcases the fixture parser capabilities
func DemoParseFixtures(filename string) error {
	fmt.Printf("=== Parsing fixture file: %s ===\n", filename)

	division, err := ParseFixtureFile(filename)
	if err != nil {
		return fmt.Errorf("failed to parse fixture: %w", err)
	}

	fmt.Printf("Division: %s\n", division.Name)
	fmt.Printf("Total rounds: %d\n", len(division.Rounds))

	totalMatches := 0
	playedMatches := 0
	unplayedMatches := 0

	for _, round := range division.Rounds {
		totalMatches += len(round.Matches)

		for _, match := range round.Matches {
			if match.Played {
				playedMatches++
			} else {
				unplayedMatches++
			}
		}
	}

	fmt.Printf("Total matches: %d (Played: %d, Unplayed: %d)\n",
		totalMatches, playedMatches, unplayedMatches)

	// Show unplayed matches that need tournaments
	unplayed := GetUnplayedMatches(division)
	if len(unplayed) > 0 {
		fmt.Printf("\nMatches needing tournaments:\n")

		for i, match := range unplayed {
			if i >= 5 { // Show only first 5
				fmt.Printf("... and %d more\n", len(unplayed)-5)
				break
			}

			fmt.Printf("  Match %d: %s vs %s (Round %d)\n",
				match.ID, match.HomePlayer, match.AwayPlayer, getRoundNumber(division, match))
		}
	}

	return nil
}

// getRoundNumber finds which round a match belongs to
func getRoundNumber(division *Division, targetMatch *Match) int {
	for _, round := range division.Rounds {
		for _, match := range round.Matches {
			if match.ID == targetMatch.ID {
				return round.Number
			}
		}
	}

	return 0
}
