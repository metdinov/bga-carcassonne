package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetBGACredentials retrieves BGA credentials from environment variables or .env file
// Environment variables take precedence over .env file values
func GetBGACredentials() (username, password string, err error) {
	// Check environment variables first
	user := os.Getenv("BGA_USER")
	pass := os.Getenv("BGA_PASS")

	if user != "" && pass != "" {
		return user, pass, nil
	}

	// If not in environment, try to load from .env file
	envUser, envPass, err := loadFromEnvFile()
	if err == nil && envUser != "" && envPass != "" {
		// Use .env values for any missing environment variables
		if user == "" {
			user = envUser
		}

		if pass == "" {
			pass = envPass
		}

		if user != "" && pass != "" {
			return user, pass, nil
		}
	}

	return "", "", fmt.Errorf("BGA credentials not found in environment or .env file")
}

// loadFromEnvFile reads BGA credentials from a .env file
func loadFromEnvFile() (username, password string, err error) {
	file, err := os.Open(".env")
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	var user, pass string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "BGA_USER":
			user = value
		case "BGA_PASS":
			pass = value
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	return user, pass, nil
}

// SaveCredentialsToEnv saves BGA credentials to a .env file
// If the file exists, it updates the BGA credentials while preserving other variables
func SaveCredentialsToEnv(user, pass string) error {
	var lines []string

	var foundUser, foundPass bool

	// Read existing .env file if it exists
	if file, err := os.Open(".env"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)

			if strings.HasPrefix(trimmed, "BGA_USER=") {
				lines = append(lines, fmt.Sprintf("BGA_USER=%s", user))
				foundUser = true
			} else if strings.HasPrefix(trimmed, "BGA_PASS=") {
				lines = append(lines, fmt.Sprintf("BGA_PASS=%s", pass))
				foundPass = true
			} else {
				lines = append(lines, line)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading .env file: %w", err)
		}
	}

	// Add missing credentials
	if !foundUser {
		lines = append(lines, fmt.Sprintf("BGA_USER=%s", user))
	}

	if !foundPass {
		lines = append(lines, fmt.Sprintf("BGA_PASS=%s", pass))
	}

	// Write the updated content
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(".env", []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing .env file: %w", err)
	}

	return nil
}

// GetOrPromptCredentials gets credentials from env/file or prompts user if missing
// If saveToEnv is true and credentials are prompted, they will be saved to .env file
func GetOrPromptCredentials(saveToEnv bool) (username, password string, err error) {
	// First try to get credentials from environment or .env file
	user, pass, err := GetBGACredentials()
	if err == nil {
		return user, pass, nil
	}

	// If not found, we need to prompt the user
	// For now, return error indicating prompting is needed
	// In a real TUI, this would show interactive prompts
	return "", "", fmt.Errorf("credentials not found - interactive prompting not yet implemented")
}

// PromptForCredentials interactively prompts user for BGA credentials
// This function will be implemented with Bubble Tea for proper TUI interaction
func PromptForCredentials() (username, password string, err error) {
	// TODO: Implement with Bubble Tea interactive prompts
	return "", "", fmt.Errorf("interactive prompting not yet implemented")
}
