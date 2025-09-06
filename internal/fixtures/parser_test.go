package fixtures

import (
	"testing"
)

func TestParseMatch_ValidMatchWithBGALink(t *testing.T) {
	csvLine := "1,herchu,2,1,Lord Trooper,12/08 - 09:30,https://boardgamearena.com/tournament?id=423761,,1,1,0"

	match, err := ParseMatch(csvLine)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if match.ID != 1 {
		t.Errorf("Expected match ID 1, got %d", match.ID)
	}

	if match.HomePlayer != "herchu" {
		t.Errorf("Expected home player 'herchu', got %s", match.HomePlayer)
	}

	if match.HomeScore != 2 {
		t.Errorf("Expected home score 2, got %d", match.HomeScore)
	}

	if match.AwayScore != 1 {
		t.Errorf("Expected away score 1, got %d", match.AwayScore)
	}

	if match.AwayPlayer != "Lord Trooper" {
		t.Errorf("Expected away player 'Lord Trooper', got %s", match.AwayPlayer)
	}

	if match.DateTime != "12/08 - 09:30" {
		t.Errorf("Expected datetime '12/08 - 09:30', got %s", match.DateTime)
	}

	if match.BGALink != "https://boardgamearena.com/tournament?id=423761" {
		t.Errorf("Expected BGA link 'https://boardgamearena.com/tournament?id=423761', got %s", match.BGALink)
	}

	if !match.Played {
		t.Errorf("Expected match to be played")
	}
}

func TestParseMatch_UnplayedMatch(t *testing.T) {
	csvLine := "17,webbi,0,0,herchu,,,,0,0,0"

	match, err := ParseMatch(csvLine)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if match.ID != 17 {
		t.Errorf("Expected match ID 17, got %d", match.ID)
	}

	if match.Played {
		t.Errorf("Expected match to not be played")
	}

	if match.BGALink != "" {
		t.Errorf("Expected empty BGA link, got %s", match.BGALink)
	}
}

func TestParseRound_ValidRound(t *testing.T) {
	csvData := `Duelo,Fecha 1,,,,11/08 - 17/08,Link,,¿Se jugó?,¿Ganó Local?,¿Ganó Visita?
1,herchu,2,1,Lord Trooper,12/08 - 09:30,https://boardgamearena.com/tournament?id=423761,,1,1,0
2,webbi,2,0,alehrosario,13/08 - 22:00,https://boardgamearena.com/tournament?id=423630,,1,1,0
3,Academia47,1,2,bignacho610,15/08 - 10:00,https://boardgamearena.com/tournament?id=424490,,1,0,1`

	round, err := ParseRound(csvData)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if round.Number != 1 {
		t.Errorf("Expected round number 1, got %d", round.Number)
	}

	if round.DateRange != "11/08 - 17/08" {
		t.Errorf("Expected date range '11/08 - 17/08', got %s", round.DateRange)
	}

	if len(round.Matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(round.Matches))
	}

	// Check first match
	firstMatch := round.Matches[0]
	if firstMatch.ID != 1 {
		t.Errorf("Expected first match ID 1, got %d", firstMatch.ID)
	}

	if firstMatch.HomePlayer != "herchu" {
		t.Errorf("Expected first match home player 'herchu', got %s", firstMatch.HomePlayer)
	}
}

func TestParseDivision_MultipleRounds(t *testing.T) {
	csvData := `Duelo,Fecha 1,,,,11/08 - 17/08,Link,,¿Se jugó?,¿Ganó Local?,¿Ganó Visita?
1,herchu,2,1,Lord Trooper,12/08 - 09:30,https://boardgamearena.com/tournament?id=423761,,1,1,0
2,webbi,2,0,alehrosario,13/08 - 22:00,https://boardgamearena.com/tournament?id=423630,,1,1,0
,,,,,,,,,,
Duelo,Fecha 2,,,,18/08 - 24/08,Link,,¿Se jugó?,¿Ganó Local?,¿Ganó Visita?
3,Lord Trooper,0,2,webbi,21/08 - 16:00,https://boardgamearena.com/tournament?id=425126,,1,0,1
4,alehrosario,0,0,herchu,,,,0,0,0`

	division, err := ParseDivision(csvData)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(division.Rounds) != 2 {
		t.Errorf("Expected 2 rounds, got %d", len(division.Rounds))
	}

	// Check first round
	firstRound := division.Rounds[0]
	if firstRound.Number != 1 {
		t.Errorf("Expected first round number 1, got %d", firstRound.Number)
	}

	if len(firstRound.Matches) != 2 {
		t.Errorf("Expected 2 matches in first round, got %d", len(firstRound.Matches))
	}

	// Check second round
	secondRound := division.Rounds[1]
	if secondRound.Number != 2 {
		t.Errorf("Expected second round number 2, got %d", secondRound.Number)
	}

	if len(secondRound.Matches) != 2 {
		t.Errorf("Expected 2 matches in second round, got %d", len(secondRound.Matches))
	}

	// Check that unplayed match is identified correctly
	unplayedMatch := secondRound.Matches[1]
	if unplayedMatch.Played {
		t.Errorf("Expected match 4 to be unplayed")
	}
}

func TestParseDivision_RealFixtureFile(t *testing.T) {
	// Test with a simplified version of the real fixture data
	csvData := `Duelo,Fecha 1,,,,11/08 - 17/08,Link,,¿Se jugó?,¿Ganó Local?,¿Ganó Visita?
1,herchu,2,1,Lord Trooper,12/08 - 09:30,https://boardgamearena.com/tournament?id=423761,,1,1,0
2,webbi,2,0,alehrosario,13/08 - 22:00,https://boardgamearena.com/tournament?id=423630,,1,1,0
3,Academia47,1,2,bignacho610,15/08 - 10:00,https://boardgamearena.com/tournament?id=424490,,1,0,1
4,maticarrizoc,2,1,Nicoooo95,12/08 - 22:25,https://boardgamearena.com/tournament?id=424020,,1,1,0
,,,,,,,,,,
Duelo,Fecha 2,,,,18/08 - 24/08,Link,,¿Se jugó?,¿Ganó Local?,¿Ganó Visita?
5,Lord Trooper,0,2,maticarrizoc,21/08 - 16:00,https://boardgamearena.com/tournament?id=425126,,1,0,1
6,Nicoooo95,1,2,Academia47,21/08 - 23:00,https://boardgamearena.com/tournament?id=426445,,1,0,1
7,bignacho610,2,0,webbi,24/08 - 13:00,https://boardgamearena.com/tournament?id=425862,,1,1,0
8,alehrosario,0,2,herchu,27/08 - 10:00,https://boardgamearena.com/tournament?id=427678,,1,0,1
,,,,,,,,,,
Duelo,Fecha 5,,,,08/09 - 14/09,Link,,¿Se jugó?,¿Ganó Local?,¿Ganó Visita?
17,webbi,0,0,herchu,,,,0,0,0
18,Academia47,0,0,maticarrizoc,,,,0,0,0
19,Nicoooo95,0,0,alehrosario,,,,0,0,0
20,bignacho610,0,0,Lord Trooper,,,,0,0,0`

	division, err := ParseDivision(csvData)
	if err != nil {
		t.Fatalf("Expected no error parsing real fixture data, got: %v", err)
	}

	if len(division.Rounds) != 3 {
		t.Errorf("Expected 3 rounds, got %d", len(division.Rounds))
	}

	// Check played vs unplayed matches
	playedCount := 0
	unplayedCount := 0

	for _, round := range division.Rounds {
		for _, match := range round.Matches {
			if match.Played {
				playedCount++
			} else {
				unplayedCount++
			}
		}
	}

	if playedCount != 8 {
		t.Errorf("Expected 8 played matches, got %d", playedCount)
	}

	if unplayedCount != 4 {
		t.Errorf("Expected 4 unplayed matches, got %d", unplayedCount)
	}

	// Check that unplayed matches have empty BGA links
	round5 := division.Rounds[2] // Should be round 5
	if round5.Number != 5 {
		t.Errorf("Expected round 5, got round %d", round5.Number)
	}

	for _, match := range round5.Matches {
		if match.Played {
			t.Errorf("Expected all matches in round 5 to be unplayed")
		}
		if match.BGALink != "" {
			t.Errorf("Expected empty BGA link for unplayed match, got %s", match.BGALink)
		}
	}
}

func TestParseFixtureFile_ActualCSV(t *testing.T) {
	filename := "../../data/Liga Argentina - 1° Temporada - E-Fixture.csv"

	division, err := ParseFixtureFile(filename)
	if err != nil {
		t.Fatalf("Expected no error parsing fixture file, got: %v", err)
	}

	if len(division.Rounds) == 0 {
		t.Errorf("Expected at least one round, got %d", len(division.Rounds))
	}

	// Check that we have both played and unplayed matches
	playedCount := 0
	unplayedCount := 0

	for _, round := range division.Rounds {
		for _, match := range round.Matches {
			if match.Played {
				playedCount++
				// Played matches should have BGA links
				if match.BGALink == "" {
					t.Errorf("Played match %d should have BGA link", match.ID)
				}
			} else {
				unplayedCount++
				// Note: Unplayed matches may have BGA links if tournament was created but not completed
				// This is a valid state - we just check that it wasn't marked as played
			}
		}
	}

	if playedCount == 0 {
		t.Errorf("Expected some played matches")
	}

	if unplayedCount == 0 {
		t.Errorf("Expected some unplayed matches")
	}

	t.Logf("Parsed %d rounds with %d played matches and %d unplayed matches",
		len(division.Rounds), playedCount, unplayedCount)
}

func TestGetUnplayedMatches(t *testing.T) {
	csvData := `Duelo,Fecha 1,,,,11/08 - 17/08,Link,,¿Se jugó?,¿Ganó Local?,¿Ganó Visita?
1,herchu,2,1,Lord Trooper,12/08 - 09:30,https://boardgamearena.com/tournament?id=423761,,1,1,0
2,webbi,0,0,alehrosario,,,,0,0,0
3,Academia47,0,0,bignacho610,15/08 - 10:00,https://boardgamearena.com/tournament?id=424490&token=123,,0,0,1`

	division, err := ParseDivision(csvData)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	unplayed := GetUnplayedMatches(division)

	if len(unplayed) != 2 {
		t.Errorf("Expected 2 unplayed matches, got %d", len(unplayed))
	}

	// Check first unplayed match
	if unplayed[0].ID != 2 {
		t.Errorf("Expected first unplayed match ID 2, got %d", unplayed[0].ID)
	}

	if unplayed[0].HomePlayer != "webbi" {
		t.Errorf("Expected first unplayed match home player 'webbi', got %s", unplayed[0].HomePlayer)
	}

	// Check second unplayed match
	if unplayed[1].ID != 3 {
		t.Errorf("Expected second unplayed match ID 3, got %d", unplayed[1].ID)
	}

	// Verify it can handle match with BGA link but unplayed
	if unplayed[1].BGALink == "" {
		t.Errorf("Expected BGA link to be preserved even for unplayed match")
	}
}

func TestParseMultipleDivisions(t *testing.T) {
	testFiles := []struct {
		filename     string
		expectedName string
	}{
		{"../../data/Liga Argentina - 1° Temporada - E-Fixture.csv", "E"},
		{"../../data/Liga Argentina - 1° Temporada - P.A-Fixture.csv", "P.A"},
		{"../../data/Liga Argentina - 1° Temporada - O.B-Fixture.csv", "O.B"},
	}

	for _, testFile := range testFiles {
		t.Run(testFile.expectedName, func(t *testing.T) {
			division, err := ParseFixtureFile(testFile.filename)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", testFile.filename, err)
			}

			if division.Name != testFile.expectedName {
				t.Errorf("Expected division name %s, got %s", testFile.expectedName, division.Name)
			}

			if len(division.Rounds) == 0 {
				t.Errorf("Expected at least one round in %s", testFile.filename)
			}

			// Verify all rounds have valid data
			for i, round := range division.Rounds {
				if len(round.Matches) == 0 {
					t.Errorf("Round %d in %s has no matches", i+1, testFile.filename)
				}

				if round.DateRange == "" {
					t.Errorf("Round %d in %s has empty date range", i+1, testFile.filename)
				}
			}

			// Count unplayed matches
			unplayedCount := len(GetUnplayedMatches(division))
			t.Logf("Division %s: %d rounds, %d unplayed matches",
				division.Name, len(division.Rounds), unplayedCount)
		})
	}
}

func TestDemoParseFixtures(t *testing.T) {
	filename := "../../data/Liga Argentina - 1° Temporada - E-Fixture.csv"

	err := DemoParseFixtures(filename)
	if err != nil {
		t.Fatalf("Demo function failed: %v", err)
	}

	// Test should pass if no error occurs during demo execution
}
