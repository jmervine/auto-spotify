package spotify

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"auto-spotify/internal/openai"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

// Service handles Spotify API interactions
type Service struct {
	auth        *spotifyauth.Authenticator
	client      *spotify.Client
	clientID    string
	redirectURL string
}

// SearchResult represents a search result for a song
type SearchResult struct {
	Track  *spotify.FullTrack
	Found  bool
	Query  string
	Reason string
}

// PlaylistInfo represents basic playlist information
type PlaylistInfo struct {
	ID          string
	Name        string
	Description string
	TrackCount  int
	Owner       string
	Public      bool
}

// TrackInfo represents basic track information
type TrackInfo struct {
	ID     string
	Title  string
	Artist string
	Album  string
	Year   int
}

// NewService creates a new Spotify service
func NewService(clientID, clientSecret, redirectURL string) *Service {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistModifyPrivate,
		),
	)

	return &Service{
		auth:        auth,
		clientID:    clientID,
		redirectURL: redirectURL,
	}
}

// Authenticate handles the OAuth flow for Spotify
func (s *Service) Authenticate(ctx context.Context) error {
	// Parse redirect URL to get the port
	redirectURL, err := url.Parse(s.redirectURL)
	if err != nil {
		return fmt.Errorf("invalid redirect URL: %w", err)
	}

	// Start a local server to handle the callback
	ch := make(chan *oauth2.Token)
	errCh := make(chan error)
	state := fmt.Sprintf("spotify-playlist-generator-%d", time.Now().UnixNano())

	// Create a new ServeMux for this authentication session
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := s.auth.Token(ctx, state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			errCh <- fmt.Errorf("failed to get token: %w", err)
			return
		}

		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>Auto-Spotify - Login Successful</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; background-color: #1db954; color: white; }
        .container { background-color: #191414; padding: 30px; border-radius: 10px; display: inline-block; }
        h1 { margin: 0 0 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎵 Login Successful!</h1>
        <p>You can now close this window and return to your terminal.</p>
    </div>
</body>
</html>`)
		ch <- token
	})

	server := &http.Server{
		Addr:    ":" + redirectURL.Port(),
		Handler: mux,
	}

	// Start HTTP server
	go func() {
		log.Printf("Starting HTTP server on %s for Spotify OAuth callback...", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Give the server time to start
	time.Sleep(200 * time.Millisecond)

	authURL := s.auth.AuthURL(state)
	fmt.Printf("Please log in to Spotify by visiting the following page in your browser:\n%s\n\n", authURL)

	// Wait for either success or error
	select {
	case token := <-ch:
		// Shutdown the server
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)

		// Create client with the token
		httpClient := s.auth.Client(ctx, token)
		s.client = spotify.New(httpClient)

		return nil
	case err := <-errCh:
		// Shutdown the server
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)

		return err
	case <-ctx.Done():
		// Shutdown the server
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)

		return ctx.Err()
	}
}

// SearchSong searches for a song on Spotify
func (s *Service) SearchSong(ctx context.Context, song openai.Song) *SearchResult {
	// Try different search queries in order of preference
	queries := []string{
		fmt.Sprintf("artist:%s track:%s", song.Artist, song.Title),
		fmt.Sprintf("%s %s", song.Artist, song.Title),
		fmt.Sprintf("track:%s", song.Title),
	}

	for _, query := range queries {
		results, err := s.client.Search(ctx, query, spotify.SearchTypeTrack)
		if err != nil {
			log.Printf("Search error for '%s': %v", query, err)
			continue
		}

		if results.Tracks != nil && len(results.Tracks.Tracks) > 0 {
			// Find the best match
			for _, track := range results.Tracks.Tracks {
				if s.isGoodMatch(song, &track) {
					return &SearchResult{
						Track:  &track,
						Found:  true,
						Query:  query,
						Reason: song.Reason,
					}
				}
			}

			// If no perfect match, return the first result
			return &SearchResult{
				Track:  &results.Tracks.Tracks[0],
				Found:  true,
				Query:  query,
				Reason: song.Reason,
			}
		}
	}

	return &SearchResult{
		Found:  false,
		Query:  fmt.Sprintf("%s %s", song.Artist, song.Title),
		Reason: song.Reason,
	}
}

// isGoodMatch checks if a Spotify track is a good match for the requested song
func (s *Service) isGoodMatch(requested openai.Song, track *spotify.FullTrack) bool {
	// Normalize strings for comparison
	normalize := func(s string) string {
		return strings.ToLower(strings.TrimSpace(s))
	}

	requestedArtist := normalize(requested.Artist)
	requestedTitle := normalize(requested.Title)

	trackTitle := normalize(track.Name)

	// Check if any artist matches
	artistMatch := false
	for _, artist := range track.Artists {
		if normalize(artist.Name) == requestedArtist {
			artistMatch = true
			break
		}
	}

	// Check title similarity
	titleMatch := trackTitle == requestedTitle ||
		strings.Contains(trackTitle, requestedTitle) ||
		strings.Contains(requestedTitle, trackTitle)

	return artistMatch && titleMatch
}

// CreateOrUpdatePlaylist creates a new playlist or updates an existing one with the same name
func (s *Service) CreateOrUpdatePlaylist(ctx context.Context, playlistResp *openai.PlaylistResponse, forceCreate bool) (*spotify.FullPlaylist, []SearchResult, error) {
	if s.client == nil {
		return nil, nil, fmt.Errorf("not authenticated with Spotify")
	}

	// Get current user
	user, err := s.client.CurrentUser(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current user: %w", err)
	}

	var playlist *spotify.FullPlaylist

	// Try to find existing playlist with the same name (unless forcing create)
	if !forceCreate {
		fmt.Printf("🔍 Searching for existing playlist '%s'...\n", playlistResp.PlaylistName)
		existingPlaylist, err := s.findPlaylistByName(ctx, user.ID, playlistResp.PlaylistName)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to search for existing playlists: %v\n", err)
		} else if existingPlaylist != nil {
			fmt.Printf("🔄 Found existing playlist '%s', updating...\n", playlistResp.PlaylistName)
			playlist = existingPlaylist

			// Clear existing tracks
			fmt.Printf("   🧹 Clearing existing tracks...\n")
			if err := s.clearPlaylist(ctx, playlist.ID); err != nil {
				fmt.Printf("⚠️  Warning: Failed to clear existing playlist: %v\n", err)
				// Continue anyway, we'll just add to the existing tracks
			} else {
				fmt.Printf("   ✅ Cleared existing tracks\n")
			}
		} else {
			fmt.Printf("❌ No existing playlist found with name '%s'\n", playlistResp.PlaylistName)
		}
	}

	// Create new playlist if we don't have one
	if playlist == nil {
		fmt.Printf("📝 Creating new playlist '%s'...\n", playlistResp.PlaylistName)
		newPlaylist, err := s.client.CreatePlaylistForUser(
			ctx,
			user.ID,
			playlistResp.PlaylistName,
			playlistResp.Description,
			false, // public
			false, // collaborative
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create playlist: %w", err)
		}
		playlist = newPlaylist
	}

	// Search for songs and collect track IDs
	var trackIDs []spotify.ID
	var searchResults []SearchResult

	fmt.Printf("🔍 Searching for %d songs...\n", len(playlistResp.Songs))

	for i, song := range playlistResp.Songs {
		fmt.Printf("  [%d/%d] Searching for: %s - %s\n", i+1, len(playlistResp.Songs), song.Artist, song.Title)

		result := s.SearchSong(ctx, song)
		searchResults = append(searchResults, *result)

		if result.Found {
			trackIDs = append(trackIDs, result.Track.ID)
			fmt.Printf("    ✓ Found: %s - %s\n", result.Track.Artists[0].Name, result.Track.Name)
		} else {
			fmt.Printf("    ✗ Not found: %s - %s\n", song.Artist, song.Title)
		}
	}

	// Add tracks to playlist (Spotify API has a limit of 100 tracks per request)
	if len(trackIDs) > 0 {
		const batchSize = 100
		for i := 0; i < len(trackIDs); i += batchSize {
			end := i + batchSize
			if end > len(trackIDs) {
				end = len(trackIDs)
			}

			batch := trackIDs[i:end]
			_, err := s.client.AddTracksToPlaylist(ctx, playlist.ID, batch...)
			if err != nil {
				return playlist, searchResults, fmt.Errorf("failed to add tracks to playlist: %w", err)
			}
		}
	}

	return playlist, searchResults, nil
}

// CreatePlaylist creates a playlist on Spotify and adds the found tracks (legacy method)
func (s *Service) CreatePlaylist(ctx context.Context, playlistResp *openai.PlaylistResponse) (*spotify.FullPlaylist, []SearchResult, error) {
	return s.CreateOrUpdatePlaylist(ctx, playlistResp, true) // Force create new
}

// findPlaylistByName searches for a playlist by name in the user's playlists
func (s *Service) findPlaylistByName(ctx context.Context, userID string, playlistName string) (*spotify.FullPlaylist, error) {
	limit := 50
	offset := 0
	maxPlaylists := 200 // Safety limit to avoid infinite loops
	totalChecked := 0

	fmt.Printf("   🔍 Searching playlists for user %s...\n", userID)

	for {
		playlists, err := s.client.CurrentUsersPlaylists(ctx, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return nil, fmt.Errorf("failed to get playlists: %w", err)
		}

		fmt.Printf("   📄 Checking %d playlists (offset %d)...\n", len(playlists.Playlists), offset)

		for _, playlist := range playlists.Playlists {
			totalChecked++
			fmt.Printf("   📋 Playlist %d: '%s' (owner: %s)\n", totalChecked, playlist.Name, playlist.Owner.ID)

			// Match by name (playlist is in user's library, so ownership is less strict)
			if playlist.Name == playlistName {
				fmt.Printf("   ✅ Found matching playlist!\n")
				// Get full playlist details
				fullPlaylist, err := s.client.GetPlaylist(ctx, playlist.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to get full playlist details: %w", err)
				}
				return fullPlaylist, nil
			}
		}

		// Check if we've seen all playlists or hit our safety limit
		if len(playlists.Playlists) < limit || totalChecked >= maxPlaylists {
			fmt.Printf("   📊 Finished searching. Checked %d total playlists.\n", totalChecked)
			break
		}
		offset += limit
	}

	return nil, nil // Not found
}

// clearPlaylist removes all tracks from a playlist
func (s *Service) clearPlaylist(ctx context.Context, playlistID spotify.ID) error {
	// Get current tracks
	tracks, err := s.client.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	if len(tracks.Tracks) == 0 {
		return nil // Already empty
	}

	// Build list of track IDs to remove
	var trackIDs []spotify.ID
	for _, track := range tracks.Tracks {
		if track.Track.ID != "" {
			trackIDs = append(trackIDs, track.Track.ID)
		}
	}

	if len(trackIDs) == 0 {
		return nil
	}

	// Remove tracks in batches (Spotify API limit)
	const batchSize = 100
	for i := 0; i < len(trackIDs); i += batchSize {
		end := i + batchSize
		if end > len(trackIDs) {
			end = len(trackIDs)
		}

		batch := trackIDs[i:end]
		_, err := s.client.RemoveTracksFromPlaylist(ctx, playlistID, batch...)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetUserPlaylists retrieves all playlists for the current user
func (s *Service) GetUserPlaylists(ctx context.Context) ([]PlaylistInfo, error) {
	if s.client == nil {
		return nil, fmt.Errorf("not authenticated with Spotify")
	}

	var allPlaylists []PlaylistInfo
	limit := 50
	offset := 0
	maxPlaylists := 1000 // Safety limit

	for {
		playlists, err := s.client.CurrentUsersPlaylists(ctx, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return nil, fmt.Errorf("failed to get playlists: %w", err)
		}

		for _, playlist := range playlists.Playlists {
			playlistInfo := PlaylistInfo{
				ID:          string(playlist.ID),
				Name:        playlist.Name,
				Description: playlist.Description,
				TrackCount:  int(playlist.Tracks.Total),
				Owner:       playlist.Owner.ID,
				Public:      playlist.IsPublic,
			}
			allPlaylists = append(allPlaylists, playlistInfo)

			// Safety check
			if len(allPlaylists) >= maxPlaylists {
				break
			}
		}

		// Check if we've seen all playlists
		if len(playlists.Playlists) < limit || len(allPlaylists) >= maxPlaylists {
			break
		}
		offset += limit
	}

	return allPlaylists, nil
}

// GetPlaylistTracks retrieves all tracks from a playlist
func (s *Service) GetPlaylistTracks(ctx context.Context, playlistID string) ([]TrackInfo, error) {
	if s.client == nil {
		return nil, fmt.Errorf("not authenticated with Spotify")
	}

	var allTracks []TrackInfo
	limit := 100
	offset := 0
	maxTracks := 10000 // Safety limit

	spotifyID := spotify.ID(playlistID)

	for {
		tracks, err := s.client.GetPlaylistTracks(ctx, spotifyID,
			spotify.Limit(limit),
			spotify.Offset(offset),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get playlist tracks: %w", err)
		}

		for _, item := range tracks.Tracks {
			if item.Track.ID == "" {
				continue // Skip empty tracks
			}

			track := item.Track

			// Get primary artist
			var artist string
			if len(track.Artists) > 0 {
				artist = track.Artists[0].Name
			}

			// Get album name
			var album string
			if track.Album.Name != "" {
				album = track.Album.Name
			}

			// Get release year
			var year int
			if track.Album.ReleaseDate != "" {
				// Parse year from release date (format: "YYYY-MM-DD" or "YYYY")
				if len(track.Album.ReleaseDate) >= 4 {
					if y, err := time.Parse("2006", track.Album.ReleaseDate[:4]); err == nil {
						year = y.Year()
					}
				}
			}

			trackInfo := TrackInfo{
				ID:     string(track.ID),
				Title:  track.Name,
				Artist: artist,
				Album:  album,
				Year:   year,
			}
			allTracks = append(allTracks, trackInfo)

			// Safety check
			if len(allTracks) >= maxTracks {
				break
			}
		}

		// Check if we've seen all tracks
		if len(tracks.Tracks) < limit || len(allTracks) >= maxTracks {
			break
		}
		offset += limit
	}

	return allTracks, nil
}
