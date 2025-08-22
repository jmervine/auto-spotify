package main

import (
	"fmt"
	"log"
	"os"

	"auto-spotify/cmd"
	"auto-spotify/internal/config"
	"auto-spotify/internal/openai"
	"auto-spotify/internal/spotify"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize services
	openaiService := openai.NewService(cfg.OpenAI.APIKey)
	spotifyService := spotify.NewService(cfg.Spotify.ClientID, cfg.Spotify.ClientSecret, cfg.Spotify.RedirectURL)

	// Setup root command
	rootCmd := cmd.NewRootCmd(openaiService, spotifyService)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
