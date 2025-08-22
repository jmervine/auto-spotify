#!/usr/bin/env python3
"""
Build GitHub Pages site from README and binaries
"""
import re
import os
import shutil

def md_to_html(text):
    """Convert basic markdown to HTML"""
    # Headers
    text = re.sub(r'^### (.*)', r'<h3>\1</h3>', text, flags=re.MULTILINE)
    text = re.sub(r'^## (.*)', r'<h2>\1</h2>', text, flags=re.MULTILINE)
    text = re.sub(r'^# (.*)', r'<h1>\1</h1>', text, flags=re.MULTILINE)
    
    # Code blocks
    text = re.sub(r'```(\w+)?\n(.*?)\n```', r'<pre><code class="language-\1">\2</code></pre>', text, flags=re.DOTALL)
    
    # Inline code
    text = re.sub(r'`([^`]+)`', r'<code>\1</code>', text)
    
    # Bold
    text = re.sub(r'\*\*(.*?)\*\*', r'<strong>\1</strong>', text)
    
    # Links
    text = re.sub(r'\[([^\]]+)\]\(([^)]+)\)', r'<a href="\2">\1</a>', text)
    
    # Lists
    text = re.sub(r'^- (.*)', r'<li>\1</li>', text, flags=re.MULTILINE)
    text = re.sub(r'(<li>.*</li>)', r'<ul>\1</ul>', text, flags=re.DOTALL)
    text = re.sub(r'</li>\n<li>', r'</li><li>', text)
    text = re.sub(r'</ul>\n<ul>', r'', text)
    
    # Paragraphs
    paragraphs = text.split('\n\n')
    html_paragraphs = []
    for p in paragraphs:
        p = p.strip()
        if p and not p.startswith('<'):
            p = f'<p>{p}</p>'
        html_paragraphs.append(p)
    text = '\n\n'.join(html_paragraphs)
    
    return text

def main():
    # Read README
    with open('README.md', 'r') as f:
        content = f.read()
    
    # Skip the title (first line) since we have it in nav
    lines = content.split('\n')[1:]
    content = '\n'.join(lines)
    
    html_content = md_to_html(content)
    
    # Add download section
    download_section = '''
        <div class="download-section">
            <h2>📦 Download Pre-built Binaries</h2>
            <p>Choose the binary for your operating system:</p>
            <div class="download-links">
                <a href="dist/auto-spotify-linux-amd64" class="download-link" download>
                    🐧 Linux (AMD64)
                </a>
                <a href="dist/auto-spotify-darwin-amd64" class="download-link" download>
                    🍎 macOS (Intel)
                </a>
                <a href="dist/auto-spotify-darwin-arm64" class="download-link" download>
                    🍎 macOS (Apple Silicon)
                </a>
                <a href="dist/auto-spotify-windows-amd64.exe" class="download-link" download>
                    🪟 Windows (AMD64)
                </a>
            </div>
            <div class="binary-info">
                <p><strong>After downloading:</strong></p>
                <ul>
                    <li><strong>Linux/macOS:</strong> <code>mv auto-spotify-* auto-spotify && chmod +x auto-spotify</code></li>
                    <li><strong>Windows:</strong> <code>ren auto-spotify-*.exe auto-spotify.exe</code></li>
                </ul>
            </div>
        </div>
    '''
    
    # Insert download section after installation header
    html_content = html_content.replace('<h2>📦 Installation</h2>', '<h2>📦 Installation</h2>' + download_section)
    
    # Create full HTML page
    full_html = '''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Auto-Spotify - AI-Powered Spotify Playlist Generator</title>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism.min.css" rel="stylesheet">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #f8f9fa;
        }
        .container {
            background: white;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 { color: #1db954; border-bottom: 3px solid #1db954; padding-bottom: 10px; }
        h2 { color: #191414; margin-top: 30px; }
        .download-section {
            background: #f0f8f0;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
            border-left: 4px solid #1db954;
        }
        .download-links {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 15px;
            margin-top: 15px;
        }
        .download-link {
            display: block;
            padding: 12px 20px;
            background: #1db954;
            color: white;
            text-decoration: none;
            border-radius: 6px;
            text-align: center;
            font-weight: 500;
            transition: background-color 0.3s;
        }
        .download-link:hover {
            background: #1ed760;
            color: white;
        }
        .binary-info {
            font-size: 0.9em;
            color: #666;
            margin-top: 10px;
        }
        pre {
            background: #f6f8fa;
            border: 1px solid #e1e4e8;
            border-radius: 6px;
            padding: 16px;
            overflow-x: auto;
        }
        code {
            background: #f6f8fa;
            padding: 2px 4px;
            border-radius: 3px;
            font-size: 0.9em;
        }
        .nav {
            background: #191414;
            color: white;
            padding: 15px 0;
            margin: -40px -40px 40px -40px;
            border-radius: 10px 10px 0 0;
        }
        .nav-content {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 40px;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        .nav h1 { margin: 0; border: none; color: #1db954; }
        .github-link {
            color: white;
            text-decoration: none;
            padding: 8px 16px;
            border: 1px solid #1db954;
            border-radius: 4px;
            transition: all 0.3s;
        }
        .github-link:hover {
            background: #1db954;
            color: white;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="nav">
            <div class="nav-content">
                <h1>🎵 Auto-Spotify</h1>
                <a href="https://github.com/jmervine/auto-spotify" class="github-link">View on GitHub</a>
            </div>
        </div>

''' + html_content + '''
    </div>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-core.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/autoloader/prism-autoloader.min.js"></script>
</body>
</html>'''

    # Create docs-site directory
    os.makedirs('docs-site', exist_ok=True)
    
    # Write HTML file
    with open('docs-site/index.html', 'w') as f:
        f.write(full_html)
    
    # Copy binaries if they exist
    if os.path.exists('dist'):
        if os.path.exists('docs-site/dist'):
            shutil.rmtree('docs-site/dist')
        shutil.copytree('dist', 'docs-site/dist')
        print("✅ Copied binaries to docs-site/dist/")
    else:
        print("⚠️  No dist/ directory found, run 'make release' first")
    
    print("✅ GitHub Pages site generated in docs-site/")
    print("   Open docs-site/index.html in your browser to preview")

if __name__ == '__main__':
    main()
