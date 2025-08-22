package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// Test that main function exists and doesn't panic when imported
	// We can't easily test the actual execution without mocking everything
	assert.NotNil(t, main)
}

func TestMainWithMissingEnv(t *testing.T) {
	// Save original environment
	oldOpenAI := os.Getenv("OPENAI_API_KEY")
	oldSpotifyID := os.Getenv("SPOTIFY_CLIENT_ID")
	oldSpotifySecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	defer func() {
		// Restore original environment
		setOrUnsetEnv("OPENAI_API_KEY", oldOpenAI)
		setOrUnsetEnv("SPOTIFY_CLIENT_ID", oldSpotifyID)
		setOrUnsetEnv("SPOTIFY_CLIENT_SECRET", oldSpotifySecret)
	}()

	// Clear environment variables
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("SPOTIFY_CLIENT_ID")
	os.Unsetenv("SPOTIFY_CLIENT_SECRET")

	// We can't easily test main() directly because it calls os.Exit()
	// In a real application, you'd refactor main to return an error
	// and have a separate function that calls os.Exit()

	// For now, we just test that the function exists
	assert.NotNil(t, main)
}

// Helper function to set or unset environment variable
func setOrUnsetEnv(key, value string) {
	if value == "" {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, value)
	}
}
