# Auto-Spotify 🎵

Create Spotify playlists automatically! Auto-Spotify can create playlists from your own song lists or use AI to generate new playlists based on your music preferences.

## 📦 Installation

### Install from Source

**Requirements**: Go 1.21 or later

```bash
git clone https://github.com/jmervine/auto-spotify.git
cd auto-spotify
go mod tidy
go build -o auto-spotify
```

## 🎵 Method 1: Using ChatGPT (Free - No OpenAI Account Required)

**Best for**: People who don't want to pay for OpenAI API access but want AI-generated playlists.

### Step 1: Get Song Recommendations from ChatGPT

Go to [ChatGPT](https://chat.openai.com) and use this detailed prompt:

```
I want to create a Spotify playlist. Please generate a list of 25 songs in the exact format I specify below. Each line should contain only the artist name, a dash, and the song title. Do not include any other text, explanations, or formatting.

Theme: [DESCRIBE YOUR DESIRED PLAYLIST HERE - e.g., "upbeat workout music with electronic and rock elements"]

Please format each song exactly like this:
Artist Name - Song Title

For example:
The Weeknd - Blinding Lights
Dua Lipa - Physical
Imagine Dragons - Believer

Now generate 25 songs for my theme:
```

**Example for a workout playlist:**

```
I want to create a Spotify playlist. Please generate a list of 25 songs in the exact format I specify below. Each line should contain only the artist name, a dash, and the song title. Do not include any other text, explanations, or formatting.

Theme: High-energy workout music with a mix of electronic, rock, and pop songs that are perfect for running and strength training. Focus on songs with strong beats and motivational lyrics.

Please format each song exactly like this:
Artist Name - Song Title

For example:
The Weeknd - Blinding Lights
Dua Lipa - Physical
Imagine Dragons - Believer

Now generate 25 songs for my theme:
```

### Step 2: Save the Song List

1. Copy ChatGPT's response (just the song list)
2. Paste it into a text file (e.g., `workout-playlist.txt`)
3. Save the file

### Step 3: Set Up Spotify App

1. Go to [Spotify for Developers](https://developer.spotify.com/dashboard)
2. Log in with your Spotify account
3. Click "Create App"
4. Fill in:
   - **App name**: `Auto-Spotify` (or any name)
   - **App description**: `Personal playlist generator`
   - **Redirect URI**: `http://127.0.0.1:8080/callback`
5. Save your **Client ID** and **Client Secret**

### Step 4: Configure Auto-Spotify

Create a `.env` file in the same folder as the auto-spotify executable:

```env
# You can leave this blank for file-based playlists
OPENAI_API_KEY=

# Your Spotify app credentials
SPOTIFY_CLIENT_ID=your_spotify_client_id_here
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret_here
SPOTIFY_REDIRECT_URL=http://127.0.0.1:8080/callback
```

### Step 5: Create Your Playlist

```bash
./auto-spotify --file workout-playlist.txt --name "My Workout Playlist"
```

The app will:
1. Read your song list
2. Ask you to log in to Spotify (one-time setup)
3. Search for each song on Spotify
4. Create the playlist in your account
5. Show you which songs were found/not found

## 🤖 Method 2: Using OpenAI API (Paid - Experimental)

**Best for**: People with OpenAI API access who want fully automated AI playlist generation.

⚠️ **Note**: This feature is experimental and may not work perfectly. OpenAI API usage will cost money based on your usage.

### Step 1: Get OpenAI API Key

1. Go to [OpenAI API Keys](https://platform.openai.com/api-keys)
2. Create a new API key
3. Make sure you have credits in your OpenAI account

### Step 2: Set Up Spotify App

Follow the same Spotify setup steps from Method 1.

### Step 3: Configure Auto-Spotify

Create a `.env` file:

```env
# Your OpenAI API key
OPENAI_API_KEY=your_openai_api_key_here

# Your Spotify app credentials  
SPOTIFY_CLIENT_ID=your_spotify_client_id_here
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret_here
SPOTIFY_REDIRECT_URL=http://127.0.0.1:8080/callback
```

### Step 4: Generate Playlists with AI

**Single prompt:**
```bash
./auto-spotify "upbeat workout music with electronic beats"
```

**Multiple prompts:**
```bash
./auto-spotify "90s hip hop" "modern R&B" --songs 30
```

**Specify number of songs:**
```bash
./auto-spotify "chill indie rock for studying" --songs 25
```

## 🎛️ Advanced Options

### Command Options

- `--file, -f`: Load songs from a text file instead of using AI
- `--name, -n`: Custom playlist name (when using --file)
- `--songs, -s`: Number of songs to include (default: 20, ignored when using --file)
- `--create, -c`: Force create new playlist instead of updating existing one
- `--help, -h`: Show help information

### File Format Support

Your text files can use any of these formats:

```text
# Comments start with # or //
Artist - Song Title
Artist: Song Title  
Song Title by Artist
Just Song Title
```

### Playlist Update Behavior

- **Default**: If a playlist with the same name exists, it will be updated
- **Force Create**: Use `--create` flag to always create a new playlist

## 🔧 Troubleshooting

**"SPOTIFY_CLIENT_ID is required"**
- Make sure your `.env` file exists and contains your Spotify app credentials
- Verify the redirect URI is configured in your Spotify app: `http://127.0.0.1:8080/callback`

**"Failed to authenticate with Spotify"**
- Check that your Spotify app's redirect URI matches exactly: `http://127.0.0.1:8080/callback`
- Make sure port 8080 isn't being used by another application
- Use `127.0.0.1` instead of `localhost`

**Songs not found on Spotify**
- This is normal - not all songs exist on Spotify
- The app will create a playlist with the songs it can find
- Try more specific artist/song names for better results

**OpenAI API errors**
- Verify your API key is valid and has credits available
- Check your OpenAI account billing status

## 📝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, testing, and contribution guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Happy playlist creating! 🎶**