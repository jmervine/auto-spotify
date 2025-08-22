package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"auto-spotify/internal/spotify"

	"github.com/spf13/cobra"
)

// NewExportCmd creates the export command
func NewExportCmd(spotifyService *spotify.Service) *cobra.Command {
	var (
		outputDir    string
		playlistName string
		allPlaylists bool
	)

	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export Spotify playlists to text files",
		Long: `Export your Spotify playlists to text files in the format that auto-spotify can read.
This is useful for backing up playlists or sharing them with others.

Examples:
  auto-spotify export --dir ./backups                    # Export all playlists
  auto-spotify export --dir ./backups --playlist "My Mix" # Export specific playlist
  auto-spotify export --all --dir ./exports              # Export all playlists explicitly`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// Validate output directory
			if outputDir == "" {
				return fmt.Errorf("output directory is required (use --dir flag)")
			}

			// Create output directory if it doesn't exist
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// Authenticate with Spotify
			fmt.Println("🎧 Connecting to Spotify...")
			if err := spotifyService.Authenticate(ctx); err != nil {
				return fmt.Errorf("failed to authenticate with Spotify: %w", err)
			}

			if playlistName != "" {
				// Export specific playlist
				fmt.Printf("📁 Exporting playlist: %s\n", playlistName)
				return exportPlaylist(ctx, spotifyService, playlistName, outputDir)
			} else {
				// Export all playlists
				fmt.Println("📁 Exporting all your playlists...")
				return exportAllPlaylists(ctx, spotifyService, outputDir)
			}
		},
	}

	exportCmd.Flags().StringVarP(&outputDir, "dir", "d", "", "Output directory for exported playlist files (required)")
	exportCmd.Flags().StringVarP(&playlistName, "playlist", "p", "", "Export specific playlist by name")
	exportCmd.Flags().BoolVarP(&allPlaylists, "all", "a", false, "Export all playlists (default behavior)")

	return exportCmd
}

func exportPlaylist(ctx context.Context, spotifyService *spotify.Service, playlistName, outputDir string) error {
	// Get the specific playlist
	playlists, err := spotifyService.GetUserPlaylists(ctx)
	if err != nil {
		return fmt.Errorf("failed to get playlists: %w", err)
	}

	var targetPlaylist *spotify.PlaylistInfo
	for _, playlist := range playlists {
		if playlist.Name == playlistName {
			targetPlaylist = &playlist
			break
		}
	}

	if targetPlaylist == nil {
		return fmt.Errorf("playlist '%s' not found", playlistName)
	}

	// Export the playlist
	return exportSinglePlaylist(ctx, spotifyService, *targetPlaylist, outputDir)
}

func exportAllPlaylists(ctx context.Context, spotifyService *spotify.Service, outputDir string) error {
	// Get all user playlists
	playlists, err := spotifyService.GetUserPlaylists(ctx)
	if err != nil {
		return fmt.Errorf("failed to get playlists: %w", err)
	}

	if len(playlists) == 0 {
		fmt.Println("📭 No playlists found to export")
		return nil
	}

	fmt.Printf("📋 Found %d playlists to export\n\n", len(playlists))

	exported := 0
	failed := 0

	for _, playlist := range playlists {
		fmt.Printf("📁 Exporting: %s (%d tracks)\n", playlist.Name, playlist.TrackCount)

		if err := exportSinglePlaylist(ctx, spotifyService, playlist, outputDir); err != nil {
			fmt.Printf("  ❌ Failed: %v\n", err)
			failed++
		} else {
			fmt.Printf("  ✅ Exported successfully\n")
			exported++
		}
		fmt.Println()
	}

	fmt.Printf("📊 Export Summary:\n")
	fmt.Printf("  ✅ Successfully exported: %d playlists\n", exported)
	if failed > 0 {
		fmt.Printf("  ❌ Failed to export: %d playlists\n", failed)
	}

	return nil
}

func exportSinglePlaylist(ctx context.Context, spotifyService *spotify.Service, playlist spotify.PlaylistInfo, outputDir string) error {
	// Get playlist tracks
	tracks, err := spotifyService.GetPlaylistTracks(ctx, playlist.ID)
	if err != nil {
		return fmt.Errorf("failed to get tracks for playlist '%s': %w", playlist.Name, err)
	}

	// Create filename (sanitize playlist name)
	filename := sanitizeFilename(playlist.Name) + ".txt"
	filepath := filepath.Join(outputDir, filename)

	// Create the text file
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", filepath, err)
	}
	defer file.Close()

	// Write playlist header (optional, as comment)
	fmt.Fprintf(file, "# %s\n", playlist.Name)
	if playlist.Description != "" {
		fmt.Fprintf(file, "# %s\n", playlist.Description)
	}
	fmt.Fprintf(file, "# %d tracks\n", len(tracks))
	fmt.Fprintf(file, "# Exported from Spotify\n\n")

	// Write tracks in the format that auto-spotify can read
	for _, track := range tracks {
		// Format: "Artist - Song Title"
		// This matches the format expected by LoadPlaylistFromFile
		fmt.Fprintf(file, "%s - %s\n", track.Artist, track.Title)
	}

	return nil
}

// sanitizeFilename removes or replaces characters that aren't safe for filenames
func sanitizeFilename(name string) string {
	// Replace problematic characters with underscores
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)

	sanitized := replacer.Replace(name)

	// Trim spaces and dots from the ends (problematic on Windows)
	sanitized = strings.Trim(sanitized, " .")

	// Ensure it's not empty
	if sanitized == "" {
		sanitized = "playlist"
	}

	return sanitized
}
