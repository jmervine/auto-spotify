package spotify

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
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

// generateSelfSignedCert generates a self-signed certificate for localhost
func generateSelfSignedCert() (tls.Certificate, error) {
	// Generate private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Auto-Spotify"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:    []string{"localhost"},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Encode certificate and key
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	// Create TLS certificate
	return tls.X509KeyPair(certPEM, keyPEM)
}

// Authenticate handles the OAuth flow for Spotify
func (s *Service) Authenticate(ctx context.Context) error {
	// Parse redirect URL to determine if we need HTTPS
	redirectURL, err := url.Parse(s.redirectURL)
	if err != nil {
		return fmt.Errorf("invalid redirect URL: %w", err)
	}

	// Start a local server to handle the callback
	ch := make(chan *oauth2.Token)
	errCh := make(chan error)
	state := "spotify-playlist-generator"

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

	// Start server based on URL scheme
	go func() {
		if redirectURL.Scheme == "https" {
			log.Printf("Starting HTTPS server on %s for Spotify OAuth callback...", server.Addr)

			// Generate self-signed certificate for localhost
			cert, err := generateSelfSignedCert()
			if err != nil {
				errCh <- fmt.Errorf("failed to generate certificate: %w", err)
				return
			}

			server.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}

			if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("HTTPS server error: %w", err)
			}
		} else {
			log.Printf("Starting HTTP server on %s for Spotify OAuth callback...", server.Addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("HTTP server error: %w", err)
			}
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
