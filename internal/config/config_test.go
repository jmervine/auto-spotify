package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Success(t *testing.T) {
	// Setup test environment variables
	oldOpenAI := os.Getenv("OPENAI_API_KEY")
	oldSpotifyID := os.Getenv("SPOTIFY_CLIENT_ID")
	oldSpotifySecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	oldRedirectURL := os.Getenv("SPOTIFY_REDIRECT_URL")

	defer func() {
		// Restore original environment variables
		setOrUnset("OPENAI_API_KEY", oldOpenAI)
		setOrUnset("SPOTIFY_CLIENT_ID", oldSpotifyID)
		setOrUnset("SPOTIFY_CLIENT_SECRET", oldSpotifySecret)
		setOrUnset("SPOTIFY_REDIRECT_URL", oldRedirectURL)
	}()

	// Set test values
	os.Setenv("OPENAI_API_KEY", "test-openai-key")
	os.Setenv("SPOTIFY_CLIENT_ID", "test-spotify-id")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "test-spotify-secret")
	os.Setenv("SPOTIFY_REDIRECT_URL", "http://localhost:3000/callback")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "test-openai-key", cfg.OpenAI.APIKey)
	assert.Equal(t, "test-spotify-id", cfg.Spotify.ClientID)
	assert.Equal(t, "test-spotify-secret", cfg.Spotify.ClientSecret)
	assert.Equal(t, "http://localhost:3000/callback", cfg.Spotify.RedirectURL)
}

func TestLoad_DefaultRedirectURL(t *testing.T) {
	// Setup test environment variables
	oldOpenAI := os.Getenv("OPENAI_API_KEY")
	oldSpotifyID := os.Getenv("SPOTIFY_CLIENT_ID")
	oldSpotifySecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	oldRedirectURL := os.Getenv("SPOTIFY_REDIRECT_URL")

	defer func() {
		// Restore original environment variables
		setOrUnset("OPENAI_API_KEY", oldOpenAI)
		setOrUnset("SPOTIFY_CLIENT_ID", oldSpotifyID)
		setOrUnset("SPOTIFY_CLIENT_SECRET", oldSpotifySecret)
		setOrUnset("SPOTIFY_REDIRECT_URL", oldRedirectURL)
	}()

	// Set test values without redirect URL
	os.Setenv("OPENAI_API_KEY", "test-openai-key")
	os.Setenv("SPOTIFY_CLIENT_ID", "test-spotify-id")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "test-spotify-secret")
	os.Unsetenv("SPOTIFY_REDIRECT_URL")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "http://127.0.0.1:8080/callback", cfg.Spotify.RedirectURL)
}

func TestLoad_MissingOpenAIKey(t *testing.T) {
	// Setup test environment variables
	oldOpenAI := os.Getenv("OPENAI_API_KEY")
	oldSpotifyID := os.Getenv("SPOTIFY_CLIENT_ID")
	oldSpotifySecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	defer func() {
		// Restore original environment variables
		setOrUnset("OPENAI_API_KEY", oldOpenAI)
		setOrUnset("SPOTIFY_CLIENT_ID", oldSpotifyID)
		setOrUnset("SPOTIFY_CLIENT_SECRET", oldSpotifySecret)
	}()

	// Unset OpenAI key - this should be OK now for file-based playlists
	os.Unsetenv("OPENAI_API_KEY")
	os.Setenv("SPOTIFY_CLIENT_ID", "test-spotify-id")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "test-spotify-secret")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "", cfg.OpenAI.APIKey)
	assert.Equal(t, "test-spotify-id", cfg.Spotify.ClientID)
	assert.Equal(t, "test-spotify-secret", cfg.Spotify.ClientSecret)
}

func TestLoad_MissingSpotifyClientID(t *testing.T) {
	// Setup test environment variables
	oldOpenAI := os.Getenv("OPENAI_API_KEY")
	oldSpotifyID := os.Getenv("SPOTIFY_CLIENT_ID")
	oldSpotifySecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	defer func() {
		// Restore original environment variables
		setOrUnset("OPENAI_API_KEY", oldOpenAI)
		setOrUnset("SPOTIFY_CLIENT_ID", oldSpotifyID)
		setOrUnset("SPOTIFY_CLIENT_SECRET", oldSpotifySecret)
	}()

	// Unset Spotify client ID
	os.Setenv("OPENAI_API_KEY", "test-openai-key")
	os.Unsetenv("SPOTIFY_CLIENT_ID")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "test-spotify-secret")

	cfg, err := Load()

	assert.Nil(t, cfg)
	assert.EqualError(t, err, "SPOTIFY_CLIENT_ID is required")
}

func TestLoad_MissingSpotifyClientSecret(t *testing.T) {
	// Setup test environment variables
	oldOpenAI := os.Getenv("OPENAI_API_KEY")
	oldSpotifyID := os.Getenv("SPOTIFY_CLIENT_ID")
	oldSpotifySecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	defer func() {
		// Restore original environment variables
		setOrUnset("OPENAI_API_KEY", oldOpenAI)
		setOrUnset("SPOTIFY_CLIENT_ID", oldSpotifyID)
		setOrUnset("SPOTIFY_CLIENT_SECRET", oldSpotifySecret)
	}()

	// Unset Spotify client secret
	os.Setenv("OPENAI_API_KEY", "test-openai-key")
	os.Setenv("SPOTIFY_CLIENT_ID", "test-spotify-id")
	os.Unsetenv("SPOTIFY_CLIENT_SECRET")

	cfg, err := Load()

	assert.Nil(t, cfg)
	assert.EqualError(t, err, "SPOTIFY_CLIENT_SECRET is required")
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable exists",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "env-value",
			expected:     "env-value",
		},
		{
			name:         "environment variable empty",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "environment variable not set",
			key:          "NON_EXISTENT_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldValue := os.Getenv(tt.key)
			defer setOrUnset(tt.key, oldValue)

			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else if tt.name == "environment variable empty" {
				os.Setenv(tt.key, "")
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnvOrDefault(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to set or unset environment variable
func setOrUnset(key, value string) {
	if value == "" {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, value)
	}
}
