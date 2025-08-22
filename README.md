# Auto-Spotify 🎵

Generate Spotify playlists using AI! Auto-Spotify uses OpenAI's ChatGPT to create personalized song recommendations based on your prompts, then automatically creates a Spotify playlist with those songs.

## Features

- 🤖 **AI-Powered**: Uses OpenAI's ChatGPT to generate intelligent song recommendations
- 🎧 **Spotify Integration**: Automatically creates playlists in your Spotify account
- 🎯 **Flexible Prompts**: Support for single or multiple prompts to create diverse playlists
- 🔍 **Smart Search**: Intelligent song matching to find tracks on Spotify
- 📊 **Detailed Reporting**: Shows which songs were found and added to your playlist
- ⚡ **Fast & Easy**: Simple command-line interface

## Prerequisites

Before you begin, you'll need:

1. **OpenAI API Key**: Get one from [OpenAI's website](https://platform.openai.com/api-keys)
2. **Spotify Developer Account**: Create an app at [Spotify for Developers](https://developer.spotify.com/dashboard)
3. **Go 1.21+**: Download from [golang.org](https://golang.org/dl/)

## Setup

### 1. Clone and Build

```bash
git clone https://github.com/jmervine/auto-spotify.git
cd auto-spotify
go mod tidy
go build -o auto-spotify
```

### 2. Generate SSL Certificates

For HTTPS support with Spotify OAuth, generate SSL certificates:

```bash
./generate-certs.sh
```

This creates certificates in the `certs/` directory that work with `https://127.0.0.1:8080/callback`.

### 3. Configure API Keys

Copy the example environment file and fill in your API credentials:

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

### 4. Spotify App Configuration

In your Spotify app settings, add the redirect URI:
- **Redirect URI**: `http://127.0.0.1:8080/callback`

**Important**: Use `127.0.0.1` instead of `localhost` to avoid redirect URI issues. HTTP works fine for local development with Spotify OAuth.

#### Redirect URI Configuration

**Local Development (Recommended):**
- Use HTTP with IP address for simplicity
- Set `SPOTIFY_REDIRECT_URL=http://127.0.0.1:8080/callback`
- Add `http://127.0.0.1:8080/callback` to your Spotify app settings

**Production Deployment:**
- Use HTTPS with your actual domain
- Set `SPOTIFY_REDIRECT_URL=https://your-domain.com/callback`
- The app automatically detects HTTP vs HTTPS from the URL

## Usage

### Basic Usage

**Generate playlist with AI prompts:**
```bash
./auto-spotify "chill indie rock for studying"
```

**Load playlist from text file:**
```bash
./auto-spotify --file my-songs.txt --name "My Custom Playlist"
```

### Advanced Usage

**Specify number of songs (AI mode only):**
```bash
./auto-spotify "upbeat workout music" --songs 30
```

**Multiple prompts:**
```bash
./auto-spotify "90s hip hop" "modern R&B" --songs 25
```

**Using prompt flags:**
```bash
./auto-spotify --prompt "jazz for a rainy day" --prompt "acoustic covers" --songs 15
```

### File Input Format

Create a text file with your songs in any of these formats:

```text
# Comments start with # or //
# Format options:

Artist - Song Title
Artist: Song Title  
Song Title by Artist
Just Song Title (artist will be "Unknown")

# Example file (metal-songs.txt):
Metallica - Master of Puppets
Iron Maiden: Run to the Hills
Enter Sandman by Metallica
Paranoid
```

**Supported file extensions:** `.txt`, `.list`, or any text file

### Command Options

- `--songs, -s`: Number of songs to include (default: 20, ignored when using --file)
- `--prompt, -p`: Additional prompts (can be used multiple times)
- `--file, -f`: Load songs from a text file instead of using AI
- `--name, -n`: Custom playlist name (when using --file)
- `--create, -c`: Force create new playlist instead of updating existing one
- `--help, -h`: Show help information

### Playlist Update Behavior

**Default (Update Mode):**
- If a playlist with the same name exists, it will be updated with new songs
- Existing tracks are cleared and replaced with the new list
- This prevents duplicate playlists with the same name

**Force Create Mode (`--create` flag):**
- Always creates a new playlist, even if one with the same name exists
- Useful when you want multiple versions of the same playlist

## How It Works

### AI Mode (Default)
1. **Prompt Processing**: Your prompts are sent to OpenAI's ChatGPT
2. **Song Generation**: ChatGPT generates a list of songs with artists, titles, and reasons
3. **Spotify Authentication**: You'll be prompted to log in to Spotify (one-time setup)
4. **Song Search**: Each recommended song is searched on Spotify
5. **Playlist Creation**: A new playlist is created in your Spotify account
6. **Results**: You'll see a summary of found/not found songs

### File Mode (--file flag)
1. **File Parsing**: Songs are loaded from your text file
2. **Format Detection**: Multiple song formats are automatically detected
3. **Spotify Authentication**: You'll be prompted to log in to Spotify (one-time setup)
4. **Song Search**: Each song from the file is searched on Spotify
5. **Playlist Creation**: A new playlist is created in your Spotify account
6. **Results**: You'll see a summary of found/not found songs

## Example Output

```
🎵 Generating playlist for prompts:
  1. chill indie rock for studying

🤖 Asking ChatGPT for song recommendations...
✅ Generated playlist: "Indie Study Vibes"
📝 Description: Perfect indie rock tracks for focused studying sessions

🎼 Recommended songs (20):
  1. Vampire Weekend - Oxford Comma (from Vampire Weekend) [2008]
     💭 Upbeat yet mellow indie rock perfect for concentration
  2. The Strokes - Someday (from Room on Fire) [2003]
     💭 Classic indie rock with a steady rhythm for studying
  ...

🎧 Connecting to Spotify...
Please log in to Spotify by visiting the following page in your browser:
https://accounts.spotify.com/authorize?...

📝 Creating Spotify playlist...
Searching for 20 songs...
  [1/20] Searching for: Vampire Weekend - Oxford Comma
    ✓ Found: Vampire Weekend - Oxford Comma
  ...

🎉 Playlist created successfully!
📋 Playlist: Indie Study Vibes
🔗 URL: https://open.spotify.com/playlist/37i9dQZF1DX0XUsuxWHRQV

📊 Search Results Summary:
  ✅ Found: 18 songs
  ❌ Not found: 2 songs
```

## Troubleshooting

### Common Issues

**"OPENAI_API_KEY is required"**
- Make sure your `.env` file exists and contains your OpenAI API key
- Verify the API key is valid and has credits available

**"SPOTIFY_CLIENT_ID is required"**
- Ensure your Spotify app credentials are correctly set in `.env`
- Verify the redirect URI is configured in your Spotify app

**"Failed to authenticate with Spotify"**
- Check that your Spotify app's redirect URI matches exactly: `http://127.0.0.1:8080/callback`
- Ensure port 8080 is not in use by another application
- Try refreshing your browser if the authentication page doesn't load
- If you see "INVALID_CLIENT: Invalid redirect URI", make sure you're using `127.0.0.1` instead of `localhost`

**Songs not found on Spotify**
- This is normal - not all AI-generated songs exist on Spotify
- The app will create a playlist with the songs it can find
- Consider more specific or popular music prompts for better match rates

### Debug Mode

Set environment variable for more detailed logging:
```bash
export LOG_LEVEL=debug
./auto-spotify "your prompt"
```

## API Limits

- **OpenAI**: Usage depends on your API plan and token limits
- **Spotify**: Rate limited to ~100 requests per minute per user

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [OpenAI](https://openai.com/) for the ChatGPT API
- [Spotify](https://developer.spotify.com/) for the Web API
- [zmb3/spotify](https://github.com/zmb3/spotify) for the excellent Go Spotify client
- [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) for the OpenAI Go client

---

**Happy playlist creating! 🎶**
