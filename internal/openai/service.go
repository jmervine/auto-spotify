package openai

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// Service handles OpenAI API interactions
type Service struct {
	client *openai.Client
}

// Song represents a song recommendation
type Song struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Album  string `json:"album,omitempty"`
	Year   int    `json:"year,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// PlaylistResponse represents the response from OpenAI
type PlaylistResponse struct {
	PlaylistName string `json:"playlist_name"`
	Description  string `json:"description"`
	Songs        []Song `json:"songs"`
}

// NewService creates a new OpenAI service
func NewService(apiKey string) *Service {
	return &Service{
		client: openai.NewClient(apiKey),
	}
}

// GeneratePlaylist generates a playlist based on the given prompt
func (s *Service) GeneratePlaylist(ctx context.Context, prompt string, songCount int) (*PlaylistResponse, error) {
	if songCount <= 0 {
		songCount = 20 // default
	}

	systemPrompt := fmt.Sprintf(`You are a music curator AI. Your task is to create a playlist based on the user's prompt. 

Please respond with a JSON object containing:
- playlist_name: A creative name for the playlist
- description: A brief description of the playlist theme
- songs: An array of exactly %d songs, each with:
  - artist: The artist name
  - title: The song title
  - album: The album name (optional)
  - year: Release year (optional)
  - reason: Brief reason why this song fits the theme (optional)

Make sure the songs are diverse, well-known enough to be found on Spotify, and match the user's request. Focus on popular and recognizable tracks.

Respond only with valid JSON, no additional text.`, songCount)

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   4000,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate playlist: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Try to extract JSON from the response (in case there's extra text)
	startIdx := strings.Index(content, "{")
	endIdx := strings.LastIndex(content, "}")
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		content = content[startIdx : endIdx+1]
	}

	var playlistResp PlaylistResponse
	if err := json.Unmarshal([]byte(content), &playlistResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w\nResponse: %s", err, content)
	}

	return &playlistResp, nil
}

// GeneratePlaylistFromMultiplePrompts generates a playlist from multiple prompts
func (s *Service) GeneratePlaylistFromMultiplePrompts(ctx context.Context, prompts []string, songCount int) (*PlaylistResponse, error) {
	if len(prompts) == 0 {
		return nil, fmt.Errorf("no prompts provided")
	}

	if len(prompts) == 1 {
		return s.GeneratePlaylist(ctx, prompts[0], songCount)
	}

	// Combine prompts into a single request
	combinedPrompt := fmt.Sprintf("Create a playlist that combines these themes:\n%s", strings.Join(prompts, "\n- "))

	return s.GeneratePlaylist(ctx, combinedPrompt, songCount)
}

// LoadPlaylistFromFile loads a playlist from a text file
func (s *Service) LoadPlaylistFromFile(filePath string, playlistName string) (*PlaylistResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	var songs []Song
	scanner := bufio.NewScanner(file)

	// Regex patterns to parse different song formats
	patterns := []*regexp.Regexp{
		// "Artist - Song Title"
		regexp.MustCompile(`^(.+?)\s*-\s*(.+?)$`),
		// "Artist: Song Title"
		regexp.MustCompile(`^(.+?)\s*:\s*(.+?)$`),
		// "Song Title by Artist"
		regexp.MustCompile(`^(.+?)\s+by\s+(.+?)$`),
	}

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		var artist, title string
		parsed := false

		// Try each pattern
		for _, pattern := range patterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) == 3 {
				if strings.Contains(line, " by ") {
					// "Song Title by Artist" format
					title = strings.TrimSpace(matches[1])
					artist = strings.TrimSpace(matches[2])
				} else {
					// "Artist - Song Title" or "Artist: Song Title" format
					artist = strings.TrimSpace(matches[1])
					title = strings.TrimSpace(matches[2])
				}
				parsed = true
				break
			}
		}

		if !parsed {
			// If no pattern matches, treat the whole line as a song title
			// and we'll let Spotify search handle it
			title = line
			artist = "Unknown"
		}

		songs = append(songs, Song{
			Artist: artist,
			Title:  title,
			Reason: fmt.Sprintf("From file: %s (line %d)", filePath, lineNum),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	if len(songs) == 0 {
		return nil, fmt.Errorf("no songs found in file %s", filePath)
	}

	// Use provided playlist name or derive from filename
	if playlistName == "" {
		playlistName = strings.TrimSuffix(strings.TrimSuffix(filePath, ".txt"), ".list")
		// Remove path and clean up the name
		if lastSlash := strings.LastIndex(playlistName, "/"); lastSlash != -1 {
			playlistName = playlistName[lastSlash+1:]
		}
		if lastSlash := strings.LastIndex(playlistName, "\\"); lastSlash != -1 {
			playlistName = playlistName[lastSlash+1:]
		}
		playlistName = strings.ReplaceAll(playlistName, "_", " ")
		playlistName = strings.ReplaceAll(playlistName, "-", " ")
		playlistName = strings.Title(playlistName)

		// Add timestamp to avoid duplicates when no custom name provided
		playlistName = fmt.Sprintf("%s (%s)", playlistName, time.Now().Format("Jan 2, 2006"))
	}

	return &PlaylistResponse{
		PlaylistName: playlistName,
		Description:  fmt.Sprintf("Playlist loaded from %s (%d songs)", filePath, len(songs)),
		Songs:        songs,
	}, nil
}
