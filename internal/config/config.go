package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	OpenAI  OpenAIConfig
	Spotify SpotifyConfig
}

// OpenAIConfig holds OpenAI API configuration
type OpenAIConfig struct {
	APIKey string
}

// SpotifyConfig holds Spotify API configuration
type SpotifyConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Try to load .env file (optional)
	_ = godotenv.Load()

	cfg := &Config{
		OpenAI: OpenAIConfig{
			APIKey: os.Getenv("OPENAI_API_KEY"),
		},
		Spotify: SpotifyConfig{
			ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RedirectURL:  getEnvOrDefault("SPOTIFY_REDIRECT_URL", "http://127.0.0.1:8080/callback"),
		},
	}

	// Validate required configuration
	if cfg.OpenAI.APIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}
	if cfg.Spotify.ClientID == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_ID is required")
	}
	if cfg.Spotify.ClientSecret == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_SECRET is required")
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
