package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"carca-cli/internal/cli"
)

func main() {
	// Get BGA credentials - from env, .env file, or prompt user
	user, pass, err := cli.GetOrPromptCredentials(true)
	if err != nil {
		fmt.Printf("Error getting BGA credentials: %v\n", err)
		fmt.Println("\nPlease set BGA_USER and BGA_PASS environment variables or create a .env file")
		os.Exit(1)
	}

	// Store credentials for later use (could be passed to BGA client)
	_ = user
	_ = pass

	// Initialize the app coordinator TUI
	model := cli.NewAppModel()

	// Create a new Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}
