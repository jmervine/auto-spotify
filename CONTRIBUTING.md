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

1. **Install Go 1.21+** from [golang.org](https://golang.org/dl/)

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**
   ```bash
   cp env.example .env
   # Edit .env with your API keys
   ```

4. **Build and test:**
   ```bash
   go build -o auto-spotify
   ./auto-spotify "test prompt"
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

While we don't have automated tests yet, please:

1. Test your changes manually with various prompts
2. Verify error handling works correctly
3. Test with different song counts and multiple prompts
4. Ensure the Spotify authentication flow works

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
- **Interactive mode** for entering multiple prompts
- **Playlist templates** (workout, study, party, etc.)
- **Music service integrations** (Apple Music, YouTube Music)
- **Advanced search options** (year range, genre filters)
- **Playlist management** (update existing playlists)

### Improvements
- **Better song matching** algorithms
- **Retry logic** for failed API calls
- **Progress indicators** for long operations
- **Configuration file** support
- **Logging improvements**

### Documentation
- **API documentation**
- **Video tutorials**
- **Use case examples**
- **Troubleshooting guides**

### Testing
- **Unit tests** for core functionality
- **Integration tests** with mock APIs
- **End-to-end testing** framework

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
