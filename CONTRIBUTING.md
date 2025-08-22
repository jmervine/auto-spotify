# Contributing to Auto-Spotify

Thank you for your interest in contributing to Auto-Spotify! This guide will help you get started with development.

## Development Setup

### Prerequisites

- **Go 1.21+** - [Install Go](https://golang.org/doc/install)
- **Git** - For version control
- **Spotify Developer Account** - [Create one here](https://developer.spotify.com/dashboard)
- **OpenAI API Key** (optional) - [Get one here](https://platform.openai.com/api-keys)

### Local Development

1. **Clone the repository:**
   ```bash
   git clone https://github.com/jmervine/auto-spotify.git
   cd auto-spotify
   ```

2. **Set up dependencies:**
   ```bash
   make setup
   ```

3. **Configure environment:**
   ```bash
   cp env.example .env
   # Edit .env with your API keys
   ```

4. **Build the application:**
   ```bash
   make build
   ```

## Code Style

- Follow standard Go conventions (`go fmt`, `go vet`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and single-purpose

## Project Structure

```
auto-spotify/
├── cmd/                    # CLI commands
│   └── root.go            # Main command logic
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── openai/          # OpenAI API integration
│   └── spotify/         # Spotify API integration
├── scripts/             # Build and utility scripts
├── templates/           # HTML templates for docs
├── main.go             # Application entry point
├── go.mod              # Go module definition
└── Makefile           # Development commands
```

## Testing

### Automated Tests

Run the full test suite:

```bash
make test
```

Run only unit tests:

```bash
make unit-test
```

Check test coverage:

```bash
make test-coverage
```

Run benchmarks:

```bash
make benchmark
```

### Manual Testing

Test with file input:

```bash
./auto-spotify --file metal-songs.txt "Test Playlist"
```

Test with AI generation (requires OpenAI API key):

```bash
./auto-spotify --count 10 "chill indie rock"
```

## Building Release Binaries

Build for all supported platforms:

```bash
make release
```

This creates binaries in the `dist/` directory for:
- Linux (AMD64)
- macOS (Intel & Apple Silicon)  
- Windows (AMD64)

## Makefile Usage

Common development tasks:

```bash
make help          # Show all available commands
make build         # Build for current platform
make run           # Run the application
make test          # Run all tests
make clean         # Clean build artifacts
make setup         # Install dependencies
make release       # Build release binaries
make pages         # Build GitHub Pages site locally
```

## Technical Architecture

### Core Components

1. **CLI Layer** (`cmd/`) - Command-line interface using Cobra
2. **Configuration** (`internal/config/`) - Environment and file-based config
3. **OpenAI Integration** (`internal/openai/`) - AI playlist generation and file parsing
4. **Spotify Integration** (`internal/spotify/`) - OAuth, search, and playlist management

### Authentication Flow

1. User runs command
2. App starts local HTTP server on `:8080`
3. Opens browser to Spotify OAuth URL
4. User authorizes app
5. Spotify redirects to `http://127.0.0.1:8080/callback`
6. App exchanges code for access token
7. Proceeds with playlist operations

### Playlist Generation

**AI Mode:**
1. Send prompts to OpenAI API
2. Parse JSON response for songs
3. Search each song on Spotify
4. Create/update playlist with found tracks

**File Mode:**
1. Parse text file for song entries
2. Support multiple formats (Artist - Song, Song by Artist, etc.)
3. Search each song on Spotify
4. Create/update playlist with found tracks

## API Rate Limits

### Spotify API
- **Search**: 100 requests per minute
- **Playlist Operations**: 100 requests per minute
- The app includes automatic retry logic with exponential backoff

### OpenAI API
- **GPT-3.5-turbo**: 3 RPM / 40,000 TPM (free tier)
- **GPT-4**: 3 RPM / 40,000 TPM (free tier)
- Higher limits available with paid plans

## Areas for Contribution

### High Priority
- **Error handling improvements** - Better user-facing error messages
- **Playlist management** - Edit existing playlists, remove duplicates
- **Search accuracy** - Improve song matching algorithms
- **Performance** - Concurrent API requests, caching

### Medium Priority  
- **Interactive mode** (partially implemented) - Step-by-step playlist building
- **Configuration UI** - Web interface for easier setup
- **Playlist templates** - Pre-defined genre/mood templates
- **Export options** - Save playlists to various formats

### Low Priority
- **Other music services** - Apple Music, YouTube Music integration
- **Social features** - Share playlists, collaborative editing
- **Analytics** - Track playlist performance, listening stats

## Submitting Changes

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make your changes**
4. **Add tests** for new functionality
5. **Ensure tests pass** (`make test`)
6. **Commit your changes** (`git commit -m 'Add amazing feature'`)
7. **Push to the branch** (`git push origin feature/amazing-feature`)
8. **Open a Pull Request**

### Pull Request Guidelines

- Include a clear description of the changes
- Reference any related issues
- Add tests for new features
- Update documentation as needed
- Ensure all CI checks pass

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/jmervine/auto-spotify/issues)
- **Discussions**: [GitHub Discussions](https://github.com/jmervine/auto-spotify/discussions)
- **Email**: Check the repository for contact information

## License

By contributing, you agree that your contributions will be licensed under the MIT License.