package spotify

import (
	"net/url"
	"testing"

	"auto-spotify/internal/openai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zmb3/spotify/v2"
)

func TestNewService(t *testing.T) {
	clientID := "test-client-id"
	clientSecret := "test-client-secret"
	redirectURL := "http://localhost:8080/callback"

	service := NewService(clientID, clientSecret, redirectURL)

	assert.NotNil(t, service)
	assert.NotNil(t, service.auth)
	assert.Equal(t, clientID, service.clientID)
	assert.Equal(t, redirectURL, service.redirectURL)
	assert.Nil(t, service.client) // Client is nil until authenticated
}

func TestSearchResult_DataStructure(t *testing.T) {
	// Test SearchResult with found track
	foundResult := &SearchResult{
		Track: &spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{
				Name: "Test Song",
				Artists: []spotify.SimpleArtist{
					{Name: "Test Artist"},
				},
			},
		},
		Found:  true,
		Query:  "Test Artist Test Song",
		Reason: "Found exact match",
	}

	assert.True(t, foundResult.Found)
	assert.Equal(t, "Test Song", foundResult.Track.Name)
	assert.Equal(t, "Test Artist", foundResult.Track.Artists[0].Name)
	assert.Equal(t, "Test Artist Test Song", foundResult.Query)
	assert.Equal(t, "Found exact match", foundResult.Reason)

	// Test SearchResult with not found
	notFoundResult := &SearchResult{
		Track:  nil,
		Found:  false,
		Query:  "Nonexistent Song",
		Reason: "No results found",
	}

	assert.False(t, notFoundResult.Found)
	assert.Nil(t, notFoundResult.Track)
	assert.Equal(t, "Nonexistent Song", notFoundResult.Query)
	assert.Equal(t, "No results found", notFoundResult.Reason)
}

func TestService_RedirectURLValidation(t *testing.T) {
	service := NewService("test-id", "test-secret", "http://localhost:8080/callback")

	// Test that we can parse the redirect URL
	parsedURL, err := url.Parse(service.redirectURL)
	require.NoError(t, err)

	assert.Equal(t, "http", parsedURL.Scheme)
	assert.Equal(t, "localhost", parsedURL.Hostname())
	assert.Equal(t, "8080", parsedURL.Port())
	assert.Equal(t, "/callback", parsedURL.Path)
}

func TestService_RedirectURLValidation_HTTPS(t *testing.T) {
	service := NewService("test-id", "test-secret", "https://127.0.0.1:8080/callback")

	// Test that we can parse the redirect URL
	parsedURL, err := url.Parse(service.redirectURL)
	require.NoError(t, err)

	assert.Equal(t, "https", parsedURL.Scheme)
	assert.Equal(t, "127.0.0.1", parsedURL.Hostname())
	assert.Equal(t, "8080", parsedURL.Port())
	assert.Equal(t, "/callback", parsedURL.Path)
}

func TestPlaylistResponse_Validation(t *testing.T) {
	// Test with valid playlist response
	validResponse := &openai.PlaylistResponse{
		PlaylistName: "Test Playlist",
		Description:  "A test playlist",
		Songs: []openai.Song{
			{Artist: "Artist 1", Title: "Song 1"},
			{Artist: "Artist 2", Title: "Song 2"},
		},
	}

	assert.NotEmpty(t, validResponse.PlaylistName)
	assert.NotEmpty(t, validResponse.Description)
	assert.Len(t, validResponse.Songs, 2)

	// Test with empty playlist response
	emptyResponse := &openai.PlaylistResponse{
		PlaylistName: "",
		Description:  "",
		Songs:        []openai.Song{},
	}

	assert.Empty(t, emptyResponse.PlaylistName)
	assert.Empty(t, emptyResponse.Description)
	assert.Len(t, emptyResponse.Songs, 0)
}

func TestService_Independence(t *testing.T) {
	// Test that different service instances are independent
	service1 := NewService("test-id-1", "test-secret-1", "http://localhost:8080/callback")
	service2 := NewService("test-id-2", "test-secret-2", "http://localhost:8081/callback")

	assert.NotEqual(t, service1.clientID, service2.clientID)
	assert.NotEqual(t, service1.redirectURL, service2.redirectURL)
	assert.NotEqual(t, service1, service2)
}

func TestSong_SearchQueryComponents(t *testing.T) {
	// Test how songs would be used in search queries
	tests := []struct {
		name      string
		song      openai.Song
		hasArtist bool
		hasTitle  bool
	}{
		{
			name: "artist and title",
			song: openai.Song{
				Artist: "The Beatles",
				Title:  "Hey Jude",
			},
			hasArtist: true,
			hasTitle:  true,
		},
		{
			name: "title only",
			song: openai.Song{
				Artist: "",
				Title:  "Imagine",
			},
			hasArtist: false,
			hasTitle:  true,
		},
		{
			name: "empty song",
			song: openai.Song{
				Artist: "",
				Title:  "",
			},
			hasArtist: false,
			hasTitle:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hasArtist {
				assert.NotEmpty(t, tt.song.Artist)
			} else {
				assert.Empty(t, tt.song.Artist)
			}

			if tt.hasTitle {
				assert.NotEmpty(t, tt.song.Title)
			} else {
				assert.Empty(t, tt.song.Title)
			}
		})
	}
}

func TestService_ConfigurationFields(t *testing.T) {
	clientID := "my-client-id"
	clientSecret := "my-client-secret"
	redirectURL := "http://example.com:9000/auth/callback"

	service := NewService(clientID, clientSecret, redirectURL)

	// Test that all configuration is stored correctly
	assert.Equal(t, clientID, service.clientID)
	assert.Equal(t, redirectURL, service.redirectURL)
	assert.NotNil(t, service.auth)
	assert.Nil(t, service.client) // Not authenticated yet
}

func TestService_DefaultState(t *testing.T) {
	service := NewService("test", "test", "http://localhost:8080/callback")

	// Test initial state
	assert.NotNil(t, service)
	assert.NotEmpty(t, service.clientID)
	assert.NotEmpty(t, service.redirectURL)
	assert.NotNil(t, service.auth)
	assert.Nil(t, service.client)
}

// Benchmark test for service creation
func BenchmarkNewService(b *testing.B) {
	clientID := "benchmark-client-id"
	clientSecret := "benchmark-client-secret"
	redirectURL := "http://localhost:8080/callback"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service := NewService(clientID, clientSecret, redirectURL)
		_ = service
	}
}

// Benchmark test for SearchResult creation
func BenchmarkSearchResult_Creation(b *testing.B) {
	track := &spotify.FullTrack{
		SimpleTrack: spotify.SimpleTrack{
			Name: "Benchmark Song",
			Artists: []spotify.SimpleArtist{
				{Name: "Benchmark Artist"},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := &SearchResult{
			Track:  track,
			Found:  true,
			Query:  "Benchmark Artist Benchmark Song",
			Reason: "Benchmark test",
		}
		_ = result
	}
}
