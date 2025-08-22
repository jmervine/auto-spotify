package cmd

import (
	"strings"
	"testing"

	"auto-spotify/internal/openai"
	"auto-spotify/internal/spotify"

	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	// Create real services for testing structure
	openaiService := openai.NewService("test-key")
	spotifyService := spotify.NewService("test-id", "test-secret", "http://localhost:8080/callback")

	rootCmd := NewRootCmd(openaiService, spotifyService)

	assert.NotNil(t, rootCmd)
	assert.Equal(t, "auto-spotify", rootCmd.Use)
	assert.Contains(t, rootCmd.Short, "Generate Spotify playlists")
	assert.Contains(t, rootCmd.Long, "Auto-Spotify uses OpenAI's ChatGPT")

	// Check flags exist
	songsFlag := rootCmd.Flags().Lookup("songs")
	assert.NotNil(t, songsFlag)
	assert.Equal(t, "s", songsFlag.Shorthand)

	promptFlag := rootCmd.Flags().Lookup("prompt")
	assert.NotNil(t, promptFlag)
	assert.Equal(t, "p", promptFlag.Shorthand)

	fileFlag := rootCmd.Flags().Lookup("file")
	assert.NotNil(t, fileFlag)
	assert.Equal(t, "f", fileFlag.Shorthand)

	nameFlag := rootCmd.Flags().Lookup("name")
	assert.NotNil(t, nameFlag)
	assert.Equal(t, "n", nameFlag.Shorthand)

	createFlag := rootCmd.Flags().Lookup("create")
	assert.NotNil(t, createFlag)
	assert.Equal(t, "c", createFlag.Shorthand)
}

func TestRootCmd_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no prompts or file",
			args:        []string{},
			expectError: true,
			errorMsg:    "provide either prompts or use --file",
		},
		{
			name:        "with prompt args",
			args:        []string{"rock music"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			openaiService := openai.NewService("test-key")
			spotifyService := spotify.NewService("test-id", "test-secret", "http://localhost:8080/callback")

			rootCmd := NewRootCmd(openaiService, spotifyService)

			// Test argument validation by extracting args without flags
			args := []string{}
			for _, arg := range tt.args {
				if !strings.HasPrefix(arg, "--") {
					args = append(args, arg)
				}
			}

			err := rootCmd.Args(rootCmd, args)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewGenerateCmd(t *testing.T) {
	openaiService := openai.NewService("test-key")
	spotifyService := spotify.NewService("test-id", "test-secret", "http://localhost:8080/callback")

	generateCmd := NewGenerateCmd(openaiService, spotifyService)

	assert.NotNil(t, generateCmd)
	assert.Equal(t, "generate [prompt...]", generateCmd.Use)
	assert.Contains(t, generateCmd.Short, "Generate a playlist from prompts")

	// Check flags
	songsFlag := generateCmd.Flags().Lookup("songs")
	assert.NotNil(t, songsFlag)
	assert.Equal(t, "s", songsFlag.Shorthand)

	interactiveFlag := generateCmd.Flags().Lookup("interactive")
	assert.NotNil(t, interactiveFlag)
	assert.Equal(t, "i", interactiveFlag.Shorthand)
}

func TestGenerateCmd_NotImplemented(t *testing.T) {
	openaiService := openai.NewService("test-key")
	spotifyService := spotify.NewService("test-id", "test-secret", "http://localhost:8080/callback")

	generateCmd := NewGenerateCmd(openaiService, spotifyService)
	generateCmd.SetArgs([]string{"test prompt"})

	err := generateCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "generate command not fully implemented yet")
}

func TestCollectInteractivePrompts(t *testing.T) {
	// Test that the function exists and has the right signature
	prompts := collectInteractivePrompts
	assert.NotNil(t, prompts)

	// In a real test, you'd mock stdin/stdout or refactor the function
	// to accept io.Reader/Writer for better testability
}

func TestRootCmd_FlagDefaults(t *testing.T) {
	openaiService := openai.NewService("test-key")
	spotifyService := spotify.NewService("test-id", "test-secret", "http://localhost:8080/callback")

	rootCmd := NewRootCmd(openaiService, spotifyService)

	// Test default values
	songsFlag := rootCmd.Flags().Lookup("songs")
	assert.Equal(t, "20", songsFlag.DefValue)

	fileFlag := rootCmd.Flags().Lookup("file")
	assert.Equal(t, "", fileFlag.DefValue)

	nameFlag := rootCmd.Flags().Lookup("name")
	assert.Equal(t, "", nameFlag.DefValue)

	createFlag := rootCmd.Flags().Lookup("create")
	assert.Equal(t, "false", createFlag.DefValue)
}

func TestRootCmd_Usage(t *testing.T) {
	openaiService := openai.NewService("test-key")
	spotifyService := spotify.NewService("test-id", "test-secret", "http://localhost:8080/callback")

	rootCmd := NewRootCmd(openaiService, spotifyService)

	// Test that usage information is properly set
	assert.Contains(t, rootCmd.Long, "Examples:")
	assert.Contains(t, rootCmd.Long, "auto-spotify \"songs for a road trip\"")
	assert.Contains(t, rootCmd.Long, "--file metal-songs.txt")
}

func TestRootCmd_ArgumentValidation(t *testing.T) {
	openaiService := openai.NewService("test-key")
	spotifyService := spotify.NewService("test-id", "test-secret", "http://localhost:8080/callback")

	rootCmd := NewRootCmd(openaiService, spotifyService)

	// Test that Args function is set
	assert.NotNil(t, rootCmd.Args)

	// Test with valid arguments
	err := rootCmd.Args(rootCmd, []string{"test prompt"})
	assert.NoError(t, err)

	// Test with no arguments and no file flag
	err = rootCmd.Args(rootCmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "provide either prompts or use --file")
}
