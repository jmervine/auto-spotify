---
layout: default
title: Home
---

# 🎵 Auto-Spotify

AI-powered Spotify playlist generator that creates playlists from text prompts or song files.

## 📦 Download Pre-built Binaries

Choose the binary for your operating system:

<div class="download-grid">
{% for download in site.downloads %}
  <a href="{{ site.baseurl }}/dist/{{ download.file }}" class="download-btn" download>
    <span class="download-icon">{{ download.icon }}</span>
    <span class="download-name">{{ download.name }}</span>
  </a>
{% endfor %}
</div>

<!-- Debug: Show what downloads are available -->
{% if site.downloads.size == 0 %}
<p><strong>Debug:</strong> No downloads found in site.downloads</p>
{% else %}
<p><strong>Debug:</strong> Found {{ site.downloads.size }} downloads:</p>
<ul>
{% for download in site.downloads %}
  <li>{{ download.name }}: {{ download.file }}</li>
{% endfor %}
</ul>
{% endif %}

<!-- Fallback static download links -->
<div class="download-grid">
  <a href="{{ site.baseurl }}/dist/auto-spotify-linux-amd64" class="download-btn" download>
    <span class="download-icon">🐧</span>
    <span class="download-name">Linux (AMD64)</span>
  </a>
  <a href="{{ site.baseurl }}/dist/auto-spotify-darwin-amd64" class="download-btn" download>
    <span class="download-icon">🍎</span>
    <span class="download-name">macOS (Intel)</span>
  </a>
  <a href="{{ site.baseurl }}/dist/auto-spotify-darwin-arm64" class="download-btn" download>
    <span class="download-icon">🍎</span>
    <span class="download-name">macOS (Apple Silicon)</span>
  </a>
  <a href="{{ site.baseurl }}/dist/auto-spotify-windows-amd64.exe" class="download-btn" download>
    <span class="download-icon">🪟</span>
    <span class="download-name">Windows (AMD64)</span>
  </a>
</div>

<div class="install-instructions">
<p><strong>After downloading:</strong></p>
<ul>
<li><strong>Linux/macOS:</strong> <code>mv auto-spotify-* auto-spotify && chmod +x auto-spotify</code></li>
<li><strong>Windows:</strong> <code>ren auto-spotify-*.exe auto-spotify.exe</code></li>
</ul>
</div>

---

## 🚀 Installation

### Option 1: Download Pre-built Binary (Recommended)

1. Download the appropriate binary for your system from the links above
2. Rename and make it executable as shown in the instructions
3. Move it to a directory in your PATH (optional but recommended)

### Option 2: Install from Source

**Requirements:** Go 1.21 or later

```bash
git clone https://github.com/jmervine/auto-spotify.git
cd auto-spotify
make setup
make build
```

---

## 🎯 Method 1: Using ChatGPT (Free)

If you don't have an OpenAI API key, you can use ChatGPT to generate song lists:

### Step 1: Get Your Spotify Credentials

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Note your **Client ID** and **Client Secret**
4. Add `http://127.0.0.1:8080/callback` as a redirect URI

### Step 2: Create Your Environment File

```bash
cp env.example .env
# Edit .env with your Spotify credentials
```

### Step 3: Generate Song List with ChatGPT

Use this prompt in ChatGPT:

```
Create a playlist of 40 songs for "80s and 90s metal like metallica, skid row, etc". 

Format the output as a simple text list with one song per line in this exact format:
Artist - Song Title

Do not include any explanations, numbers, or additional text. Just the song list.

Example format:
Metallica - Enter Sandman
Skid Row - 18 and Life
```

### Step 4: Save and Use the Song List

1. Copy the ChatGPT response to a text file (e.g., `my-playlist.txt`)
2. Run: `auto-spotify --file my-playlist.txt "My Awesome Playlist"`

---

## 🤖 Method 2: Using OpenAI API (Paid - Experimental)

⚠️ **Note:** This method is experimental and requires an OpenAI API key with available credits.

### Setup

1. Get your OpenAI API key from [OpenAI Platform](https://platform.openai.com/api-keys)
2. Add it to your `.env` file:
   ```bash
   OPENAI_API_KEY=your_api_key_here
   ```

### Usage

```bash
# Generate playlist with AI
auto-spotify --count 25 "chill indie rock for studying"

# Multiple prompts
auto-spotify --count 30 "upbeat workout songs" "electronic dance music" "pop hits 2020s"
```

---

## ⚙️ Advanced Options

### Command Line Flags

- `--file, -f`: Load songs from a text file instead of using AI
- `--count, -c`: Number of songs to generate (default: 20, max: 50)
- `--help, -h`: Show help information

### Supported File Formats

The app supports flexible song formats in text files:

```
# These formats all work:
Artist - Song Title
Song Title by Artist  
Artist: Song Title
"Song Title" - Artist
```

### Environment Variables

```bash
# Required for Spotify
SPOTIFY_CLIENT_ID=your_client_id
SPOTIFY_CLIENT_SECRET=your_client_secret
SPOTIFY_REDIRECT_URL=http://127.0.0.1:8080/callback

# Optional for AI features
OPENAI_API_KEY=your_api_key_here
```

---

## 🔧 Troubleshooting

### "Invalid redirect URI" Error
- Make sure your Spotify app has `http://127.0.0.1:8080/callback` as a redirect URI
- Use `127.0.0.1` instead of `localhost`

### "Song not found" Messages
- Some songs may not be available on Spotify
- Try alternative song titles or artists
- Check your internet connection

### OpenAI API Issues
- Verify your API key is valid and has credits
- Check [OpenAI Status](https://status.openai.com/) for service issues
- The free tier has limited usage

### Permission Denied
- Make sure the binary is executable: `chmod +x auto-spotify`
- On macOS, you may need to allow the app in Security & Privacy settings

---

## 🤝 Contributing

See [CONTRIBUTING.md]({{ site.baseurl }}/contributing/) for development setup, testing, and contribution guidelines.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE]({{ site.baseurl }}/LICENSE) file for details.
