package cli

import (
	"os"
	"strings"
	"testing"
)

func TestGetBGACredentials_FromEnvironment(t *testing.T) {
	// Setup environment variables
	os.Setenv("BGA_USER", "testuser")
	os.Setenv("BGA_PASS", "testpass")
	defer func() {
		os.Unsetenv("BGA_USER")
		os.Unsetenv("BGA_PASS")
	}()

	user, pass, err := GetBGACredentials()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if user != "testuser" {
		t.Errorf("Expected user 'testuser', got %s", user)
	}

	if pass != "testpass" {
		t.Errorf("Expected pass 'testpass', got %s", pass)
	}
}

func TestGetBGACredentials_FromEnvFile(t *testing.T) {
	// Ensure no environment variables are set
	os.Unsetenv("BGA_USER")
	os.Unsetenv("BGA_PASS")

	// Create temporary .env file
	envContent := "BGA_USER=envfileuser\nBGA_PASS=envfilepass\n"
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	user, pass, err := GetBGACredentials()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if user != "envfileuser" {
		t.Errorf("Expected user 'envfileuser', got %s", user)
	}

	if pass != "envfilepass" {
		t.Errorf("Expected pass 'envfilepass', got %s", pass)
	}
}

func TestGetBGACredentials_EnvironmentTakesPrecedence(t *testing.T) {
	// Set both environment variables and .env file
	os.Setenv("BGA_USER", "envuser")
	os.Setenv("BGA_PASS", "envpass")
	defer func() {
		os.Unsetenv("BGA_USER")
		os.Unsetenv("BGA_PASS")
	}()

	// Create .env file with different values
	envContent := "BGA_USER=fileuser\nBGA_PASS=filepass\n"
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	user, pass, err := GetBGACredentials()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Environment should take precedence
	if user != "envuser" {
		t.Errorf("Expected user 'envuser' (from env), got %s", user)
	}

	if pass != "envpass" {
		t.Errorf("Expected pass 'envpass' (from env), got %s", pass)
	}
}

func TestGetBGACredentials_MissingCredentials(t *testing.T) {
	// Ensure no environment variables or .env file
	os.Unsetenv("BGA_USER")
	os.Unsetenv("BGA_PASS")
	os.Remove(".env")

	_, _, err := GetBGACredentials()
	if err == nil {
		t.Errorf("Expected error when credentials are missing")
	}

	if err.Error() != "BGA credentials not found in environment or .env file" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestPromptForCredentials_ValidInput(t *testing.T) {
	// This test would require mocking stdin/stdout for interactive testing
	// For now, we'll test the save functionality separately
	t.Skip("Interactive testing requires stdin/stdout mocking")
}

func TestSaveCredentialsToEnv_Success(t *testing.T) {
	// Clean up any existing .env file
	os.Remove(".env")
	defer os.Remove(".env")

	err := SaveCredentialsToEnv("testuser", "testpass")
	if err != nil {
		t.Fatalf("Expected no error saving credentials, got: %v", err)
	}

	// Verify the file was created with correct content
	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatalf("Failed to read .env file: %v", err)
	}

	expectedContent := "BGA_USER=testuser\nBGA_PASS=testpass\n"
	if string(content) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, string(content))
	}
}

func TestSaveCredentialsToEnv_OverwriteExisting(t *testing.T) {
	// Create initial .env file
	initialContent := "BGA_USER=olduser\nBGA_PASS=oldpass\nOTHER_VAR=value\n"
	err := os.WriteFile(".env", []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create initial .env file: %v", err)
	}
	defer os.Remove(".env")

	err = SaveCredentialsToEnv("newuser", "newpass")
	if err != nil {
		t.Fatalf("Expected no error saving credentials, got: %v", err)
	}

	// Verify the BGA credentials were updated but other vars preserved
	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatalf("Failed to read .env file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "BGA_USER=newuser") {
		t.Errorf("Expected updated BGA_USER in content: %s", contentStr)
	}
	if !strings.Contains(contentStr, "BGA_PASS=newpass") {
		t.Errorf("Expected updated BGA_PASS in content: %s", contentStr)
	}
	if !strings.Contains(contentStr, "OTHER_VAR=value") {
		t.Errorf("Expected OTHER_VAR to be preserved in content: %s", contentStr)
	}
}

func TestGetOrPromptCredentials_Found(t *testing.T) {
	// Setup environment variables
	os.Setenv("BGA_USER", "envuser")
	os.Setenv("BGA_PASS", "envpass")
	defer func() {
		os.Unsetenv("BGA_USER")
		os.Unsetenv("BGA_PASS")
	}()

	user, pass, err := GetOrPromptCredentials(false) // Don't save to env
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if user != "envuser" {
		t.Errorf("Expected user 'envuser', got %s", user)
	}

	if pass != "envpass" {
		t.Errorf("Expected pass 'envpass', got %s", pass)
	}
}

func TestGetOrPromptCredentials_NotFound(t *testing.T) {
	// Ensure no credentials available
	os.Unsetenv("BGA_USER")
	os.Unsetenv("BGA_PASS")
	os.Remove(".env")

	// This test would require mocking user input
	// For now, we test that it returns an error indicating prompting is needed
	_, _, err := GetOrPromptCredentials(false)
	if err == nil {
		t.Errorf("Expected error when credentials not found and prompting disabled")
	}
}
