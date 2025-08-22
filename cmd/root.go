package cmd

import (
	"context"
	"fmt"
	"strings"

	"auto-spotify/internal/openai"
	"auto-spotify/internal/spotify"

	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd(openaiService *openai.Service, spotifyService *spotify.Service) *cobra.Command {
	var (
		songCount    int
		prompts      []string
		inputFile    string
		playlistName string
		forceCreate  bool
	)

	rootCmd := &cobra.Command{
		Use:   "auto-spotify",
		Short: "Generate Spotify playlists using AI or from text files",
		Long: `Auto-Spotify uses OpenAI's ChatGPT to generate song recommendations based on your prompts,
or loads songs from a text file, then creates a Spotify playlist with those songs.

Examples:
  auto-spotify "songs for a road trip"
  auto-spotify "chill indie rock for studying" --songs 15
  auto-spotify "upbeat workout music" "electronic dance" --songs 25
  auto-spotify --file metal-songs.txt --name "My Metal Playlist"`,
		Args: func(cmd *cobra.Command, args []string) error {
			if inputFile == "" && len(args) == 0 && len(prompts) == 0 {
				return fmt.Errorf("provide either prompts or use --file to load from a text file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			var playlistResp *openai.PlaylistResponse
			var err error

			if inputFile != "" {
				// Load playlist from file
				fmt.Printf("📁 Loading playlist from file: %s\n\n", inputFile)
				playlistResp, err = openaiService.LoadPlaylistFromFile(inputFile, playlistName)
				if err != nil {
					return fmt.Errorf("failed to load playlist from file: %w", err)
				}
			} else {
				// Use args as prompts if no --prompt flags were provided
				if len(prompts) == 0 {
					prompts = args
				}

				fmt.Printf("🎵 Generating playlist for prompts:\n")
				for i, prompt := range prompts {
					fmt.Printf("  %d. %s\n", i+1, prompt)
				}
				fmt.Printf("\n")

				// Generate playlist using OpenAI
				fmt.Println("🤖 Asking ChatGPT for song recommendations...")

				if len(prompts) == 1 {
					playlistResp, err = openaiService.GeneratePlaylist(ctx, prompts[0], songCount)
				} else {
					playlistResp, err = openaiService.GeneratePlaylistFromMultiplePrompts(ctx, prompts, songCount)
				}

				if err != nil {
					return fmt.Errorf("failed to generate playlist: %w", err)
				}
			}

			fmt.Printf("✅ Generated playlist: \"%s\"\n", playlistResp.PlaylistName)
			fmt.Printf("📝 Description: %s\n\n", playlistResp.Description)

			// Display recommended songs
			fmt.Printf("🎼 Recommended songs (%d):\n", len(playlistResp.Songs))
			for i, song := range playlistResp.Songs {
				fmt.Printf("  %d. %s - %s", i+1, song.Artist, song.Title)
				if song.Album != "" {
					fmt.Printf(" (from %s)", song.Album)
				}
				if song.Year > 0 {
					fmt.Printf(" [%d]", song.Year)
				}
				fmt.Println()
				if song.Reason != "" {
					fmt.Printf("     💭 %s\n", song.Reason)
				}
			}
			fmt.Println()

			// Authenticate with Spotify
			fmt.Println("🎧 Connecting to Spotify...")
			if err := spotifyService.Authenticate(ctx); err != nil {
				return fmt.Errorf("failed to authenticate with Spotify: %w", err)
			}

			// Create or update playlist on Spotify
			if forceCreate {
				fmt.Println("📝 Creating new Spotify playlist...")
			} else {
				fmt.Println("📝 Creating/updating Spotify playlist...")
			}
			playlist, searchResults, err := spotifyService.CreateOrUpdatePlaylist(ctx, playlistResp, forceCreate)
			if err != nil {
				return fmt.Errorf("failed to create/update Spotify playlist: %w", err)
			}

			// Report results
			fmt.Printf("\n🎉 Playlist created successfully!\n")
			fmt.Printf("📋 Playlist: %s\n", playlist.Name)
			fmt.Printf("🔗 URL: %s\n\n", playlist.ExternalURLs["spotify"])

			// Show search results summary
			found := 0
			notFound := 0
			for _, result := range searchResults {
				if result.Found {
					found++
				} else {
					notFound++
				}
			}

			fmt.Printf("📊 Search Results Summary:\n")
			fmt.Printf("  ✅ Found: %d songs\n", found)
			if notFound > 0 {
				fmt.Printf("  ❌ Not found: %d songs\n", notFound)
				fmt.Println("\n🔍 Songs that couldn't be found:")
				for _, result := range searchResults {
					if !result.Found {
						fmt.Printf("  - %s\n", result.Query)
					}
				}
			}

			return nil
		},
	}

	rootCmd.Flags().IntVarP(&songCount, "songs", "s", 20, "Number of songs to include in the playlist (ignored when using --file)")
	rootCmd.Flags().StringArrayVarP(&prompts, "prompt", "p", []string{}, "Additional prompts (can be used multiple times)")
	rootCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Load songs from a text file instead of using AI")
	rootCmd.Flags().StringVarP(&playlistName, "name", "n", "", "Custom playlist name (when using --file)")
	rootCmd.Flags().BoolVarP(&forceCreate, "create", "c", false, "Force create new playlist instead of updating existing one")

	return rootCmd
}

// NewGenerateCmd creates a generate command (alternative interface)
func NewGenerateCmd(openaiService *openai.Service, spotifyService *spotify.Service) *cobra.Command {
	var (
		songCount   int
		interactive bool
	)

	generateCmd := &cobra.Command{
		Use:   "generate [prompt...]",
		Short: "Generate a playlist from prompts",
		Long:  `Generate a Spotify playlist based on one or more prompts using AI.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if interactive {
				_ = collectInteractivePrompts()
			}
			// Same logic as root command...
			// (Implementation would be similar to root command)
			return fmt.Errorf("generate command not fully implemented yet")
		},
	}

	generateCmd.Flags().IntVarP(&songCount, "songs", "s", 20, "Number of songs to include in the playlist")
	generateCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode for entering multiple prompts")

	return generateCmd
}

func collectInteractivePrompts() []string {
	var prompts []string

	fmt.Println("🎵 Interactive Playlist Generator")
	fmt.Println("Enter your prompts (one per line). Type 'done' when finished:")
	fmt.Println()

	for {
		fmt.Print("Prompt: ")
		var input string
		fmt.Scanln(&input)

		if strings.ToLower(input) == "done" {
			break
		}

		if strings.TrimSpace(input) != "" {
			prompts = append(prompts, input)
		}
	}

	return prompts
}
