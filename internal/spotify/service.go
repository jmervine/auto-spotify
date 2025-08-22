package spotify

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

// NewService creates a new Spotify service
func NewService(clientID, clientSecret, redirectURL string) *Service {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
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
	// Start a local server to handle the callback
	ch := make(chan *oauth2.Token)
	state := "spotify-playlist-generator"

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := s.auth.Token(ctx, state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		
		fmt.Fprint(w, "Login successful! You can close this window.")
		ch <- token
	})

	go func() {
		log.Println("Starting local server on :8080 for Spotify OAuth callback...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	url := s.auth.AuthURL(state)
	fmt.Printf("Please log in to Spotify by visiting the following page in your browser:\n%s\n\n", url)

	// Wait for the callback
	token := <-ch
	
	// Create client with the token
	httpClient := s.auth.Client(ctx, token)
	s.client = spotify.New(httpClient)

	return nil
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

// CreatePlaylist creates a playlist on Spotify and adds the found tracks
func (s *Service) CreatePlaylist(ctx context.Context, playlistResp *openai.PlaylistResponse) (*spotify.FullPlaylist, []SearchResult, error) {
	if s.client == nil {
		return nil, nil, fmt.Errorf("not authenticated with Spotify")
	}

	// Get current user
	user, err := s.client.CurrentUser(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current user: %w", err)
	}

	// Create playlist
	playlist, err := s.client.CreatePlaylistForUser(
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

	// Search for songs and collect track IDs
	var trackIDs []spotify.ID
	var searchResults []SearchResult

	fmt.Printf("Searching for %d songs...\n", len(playlistResp.Songs))
	
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
