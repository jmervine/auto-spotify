package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
			Model: openai.GPT4TurboPreview,
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
