---
layout: page
title: Contributing
permalink: /contributing/
---

# Contributing to Auto-Spotify

---
layout: page
title: Contributing
permalink: /contributing/
---

# Contributing to Auto-Spotify

---
layout: page
title: Contributing
permalink: /contributing/
---

# Contributing to Auto-Spotify

---
layout: page
title: Contributing
permalink: /contributing/
---

# Contributing to Auto-Spotify

---
layout: page
title: Contributing
permalink: /contributing/
---

# Contributing to Auto-Spotify

Thank you for your interest in contributing to Auto-Spotify! This document provides guidelines and instructions for contributing.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/auto-spotify.git
   cd auto-spotify
   ```
3. Set up the upstream remote:
   ```bash
   git remote add upstream https://github.com/jmervine/auto-spotify.git
   ```

## Development Setup

### Prerequisites

Before you begin, you'll need:

1. **Go 1.21+**: Download from [golang.org](https://golang.org/dl/)
2. **OpenAI API Key**: Get one from [OpenAI's website](https://platform.openai.com/api-keys) (for testing AI features)
3. **Spotify Developer Account**: Create an app at [Spotify for Developers](https://developer.spotify.com/dashboard)

### Quick Setup

1. **Clone and build:**
   ```bash
   git clone https://github.com/jmervine/auto-spotify.git
   cd auto-spotify
   go mod tidy
   go build -o auto-spotify
   ```

2. **Configure API keys:**
   ```bash
   cp env.example .env
   ```

   Edit `.env` with your credentials:
   ```env
   # OpenAI API Configuration
   OPENAI_API_KEY=your_openai_api_key_here

   # Spotify API Configuration
   SPOTIFY_CLIENT_ID=your_spotify_client_id_here
   SPOTIFY_CLIENT_SECRET=your_spotify_client_secret_here
   SPOTIFY_REDIRECT_URL=http://127.0.0.1:8080/callback
   ```

3. **Configure Spotify App:**
   In your Spotify app settings, add the redirect URI:
   - **Redirect URI**: `http://127.0.0.1:8080/callback`

4. **Test the build:**
   ```bash
   ./auto-spotify "test prompt"
   ```

### Using the Makefile

The project includes a Makefile for common development tasks:

```bash
make setup      # Set up development environment
make build      # Build the application
make test       # Run all tests and checks
make run        # Build and run with test prompt
make release    # Build cross-platform binaries
make clean      # Clean build artifacts
make help       # Show all available commands
```

## Making Changes

### Before You Start

1. Check existing issues and pull requests to avoid duplicates
2. For significant changes, open an issue first to discuss the approach
3. Keep changes focused and atomic

### Development Process

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Write clean, readable code
   - Follow Go conventions and best practices
   - Add comments for complex logic
   - Update documentation if needed

3. **Test your changes:**
   ```bash
   go build -o auto-spotify
   ./auto-spotify "test with your changes"
   ```

4. **Commit your changes:**
   ```bash
   git add .
   git commit -m "Add: brief description of your changes"
   ```

### Commit Message Guidelines

Use clear, descriptive commit messages:

- **Add:** for new features
- **Fix:** for bug fixes
- **Update:** for changes to existing features
- **Remove:** for deleted code/features
- **Docs:** for documentation changes

Examples:
- `Add: support for custom playlist descriptions`
- `Fix: handle empty search results gracefully`
- `Update: improve song matching algorithm`

## Code Style

### Go Guidelines

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Keep functions focused and small
- Handle errors appropriately
- Add comments for exported functions and types

### Project Structure

```
auto-spotify/
├── cmd/           # CLI commands
├── internal/      # Internal packages
│   ├── config/    # Configuration management
│   ├── openai/    # OpenAI API integration
│   └── spotify/   # Spotify API integration
├── main.go        # Application entry point
└── README.md      # Documentation
```

## Testing

### Automated Tests

The project includes comprehensive unit tests:

```bash
# Run all tests
make test

# Run only unit tests
make unit-test

# Run tests with coverage
make test-coverage

# Run benchmark tests
make benchmark
```

### Test Coverage

Current test coverage:
- **Config Package**: 100% coverage
- **OpenAI Package**: 79.4% coverage  
- **CMD Package**: 19.6% coverage
- **Spotify Package**: Limited (due to external API dependencies)

### Manual Testing

Please also test your changes manually:

1. Test with various prompts and song counts
2. Verify error handling works correctly
3. Test file input functionality with different formats
4. Ensure the Spotify authentication flow works
5. Test playlist update vs. create behavior

### Example Test Commands

```bash
# Test AI generation (requires OpenAI API key)
./auto-spotify "chill indie rock for studying" --songs 15

# Test file input
echo "Queen - Bohemian Rhapsody\nThe Beatles - Hey Jude" > test-songs.txt
./auto-spotify --file test-songs.txt --name "Test Playlist"

# Test multiple prompts
./auto-spotify "80s rock" "90s grunge" --songs 20

# Test playlist update behavior
./auto-spotify --file test-songs.txt --name "Existing Playlist"  # Updates existing
./auto-spotify --file test-songs.txt --name "Existing Playlist" --create  # Forces new
```



## Technical Architecture

### How It Works

#### AI Mode (OpenAI Integration)
1. **Prompt Processing**: User prompts are sent to OpenAI's ChatGPT
2. **Song Generation**: ChatGPT generates a structured JSON response with songs
3. **Spotify Authentication**: OAuth 2.0 flow with local HTTP server
4. **Song Search**: Each recommended song is searched on Spotify using multiple query strategies
5. **Playlist Creation**: Songs are added to a new or existing Spotify playlist
6. **Results**: Detailed reporting of found/not found songs

#### File Mode (Text File Input)
1. **File Parsing**: Songs are loaded and parsed from text files with flexible format detection
2. **Format Detection**: Supports "Artist - Song", "Artist: Song", "Song by Artist", and plain titles
3. **Spotify Integration**: Same authentication and search process as AI mode
4. **Playlist Management**: Supports both creating new playlists and updating existing ones

### API Rate Limits

- **OpenAI**: Usage depends on your API plan and token limits
- **Spotify**: Rate limited to ~100 requests per minute per user
- **Local Development**: Use file mode to avoid OpenAI costs during development

### Performance Considerations

- **File parsing**: ~498µs for 1000 songs (benchmarked)
- **Service creation**: ~90ns per instance
- **Memory usage**: Minimal for typical playlist sizes
- **Concurrent requests**: Spotify searches are done sequentially to respect rate limits

### Building Release Binaries

To create cross-platform binaries for distribution:

```bash
make release
```

This creates optimized binaries in the `dist/` directory:
- `auto-spotify-linux-amd64` - Linux 64-bit
- `auto-spotify-darwin-amd64` - macOS Intel
- `auto-spotify-darwin-arm64` - macOS Apple Silicon
- `auto-spotify-windows-amd64.exe` - Windows 64-bit

The binaries are statically linked and optimized with `-ldflags="-s -w"` for smaller file sizes.

## Submitting Changes

1. **Push your branch:**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request:**
   - Use a clear, descriptive title
   - Explain what your changes do and why
   - Reference any related issues
   - Include examples if applicable

3. **Pull Request Template:**
   ```markdown
   ## Description
   Brief description of changes

   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Documentation update
   - [ ] Code refactoring

   ## Testing
   - [ ] Tested manually with various prompts
   - [ ] Verified error handling
   - [ ] Tested authentication flow

   ## Additional Notes
   Any additional context or screenshots
   ```

## Areas for Contribution

We welcome contributions in these areas:

### Features
- **Interactive mode** for entering multiple prompts (partially implemented)
- **Playlist templates** (workout, study, party, etc.)
- **Music service integrations** (Apple Music, YouTube Music)
- **Advanced search options** (year range, genre filters)

### Improvements
- **Better song matching** algorithms
- **Retry logic** for failed API calls
- **Progress indicators** for long operations
- **Configuration file** support
- **Debug logging** implementation

### Documentation
- **API documentation**
- **Video tutorials**
- **Use case examples**
- **Troubleshooting guides**

### Testing
- **Unit tests** for core functionality ✅ (Implemented)
- **Integration tests** with mock APIs
- **End-to-end testing** framework
- **Performance benchmarks** ✅ (Implemented)

## Getting Help

- **Questions?** Open an issue with the "question" label
- **Bug reports?** Use the bug report template
- **Feature requests?** Use the feature request template
- **Discussions?** Use GitHub Discussions for general topics

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Help others learn and grow
- Focus on what's best for the community

## Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes for significant contributions
- GitHub contributor graphs

Thank you for contributing to Auto-Spotify! 🎵