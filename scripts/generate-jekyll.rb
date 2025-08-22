#!/usr/bin/env ruby
# Generate Jekyll files from README.md and CONTRIBUTING.md
# This keeps the Jekyll site in sync with the primary documentation

require 'fileutils'

class JekyllGenerator
  def initialize
    @root_dir = Dir.pwd
    @readme_path = File.join(@root_dir, 'README.md')
    @contributing_path = File.join(@root_dir, 'CONTRIBUTING.md')
    @jekyll_index_path = File.join(@root_dir, 'index.md')
    @jekyll_contributing_path = File.join(@root_dir, 'contributing.md')
  end

  def generate_all
    puts "🔄 Generating Jekyll files from README and CONTRIBUTING..."
    
    # Generate Jekyll index from README
    generate_jekyll_index
    
    # Generate Jekyll contributing page
    generate_jekyll_contributing
    
    # Update README with download links
    update_readme_with_downloads
    
    puts "✅ Jekyll files generated successfully!"
  end

  private

  def generate_jekyll_index
    puts "  📄 Generating index.md from README.md..."
    
    readme_content = File.read(@readme_path)
    
    # Remove the title line (first # heading)
    lines = readme_content.split("\n")
    content_lines = lines.drop_while { |line| line.strip.empty? || line.start_with?('# ') }
    content = content_lines.join("\n")
    
    # Create Jekyll front matter
    jekyll_content = <<~FRONTMATTER
---
layout: default
title: Home
---

# 🎵 Auto-Spotify

AI-powered Spotify playlist generator that creates playlists from text prompts or song files.

## 📦 Download Pre-built Binaries

Choose the binary for your operating system:

- **🐧 Linux (AMD64)**: [auto-spotify-linux-amd64]({{ site.baseurl }}/dist/auto-spotify-linux-amd64)
- **🍎 macOS (Intel)**: [auto-spotify-darwin-amd64]({{ site.baseurl }}/dist/auto-spotify-darwin-amd64)
- **🍎 macOS (Apple Silicon)**: [auto-spotify-darwin-arm64]({{ site.baseurl }}/dist/auto-spotify-darwin-arm64)
- **🪟 Windows (AMD64)**: [auto-spotify-windows-amd64.exe]({{ site.baseurl }}/dist/auto-spotify-windows-amd64.exe)

**After downloading:**
- **Linux/macOS:** `mv auto-spotify-* auto-spotify && chmod +x auto-spotify`
- **Windows:** `ren auto-spotify-*.exe auto-spotify.exe`

---

FRONTMATTER
    
    # Process the content for Jekyll
    processed_content = process_content_for_jekyll(content)
    
    # Fix the Contributing link
    processed_content = processed_content.gsub(
      /\[CONTRIBUTING\.md\]\(CONTRIBUTING\.md\)/,
      '[CONTRIBUTING.md]({{ site.baseurl }}/contributing/)'
    )
    
    # Write the Jekyll index file
    File.write(@jekyll_index_path, jekyll_content + processed_content)
    puts "    ✅ Generated #{@jekyll_index_path}"
  end

  def generate_jekyll_contributing
    puts "  📄 Generating contributing.md from CONTRIBUTING.md..."
    
    contributing_content = File.read(@contributing_path)
    
    # Remove the title line (first # heading)
    lines = contributing_content.split("\n")
    content_lines = lines.drop_while { |line| line.strip.empty? || line.start_with?('# ') }
    content = content_lines.join("\n")
    
    # Create Jekyll front matter
    jekyll_content = <<~FRONTMATTER
---
layout: page
title: Contributing
permalink: /contributing/
---

# Contributing to Auto-Spotify

FRONTMATTER
    
    # Process the content for Jekyll
    processed_content = process_content_for_jekyll(content)
    
    # Write the Jekyll contributing file
    File.write(@jekyll_contributing_path, jekyll_content + processed_content)
    puts "    ✅ Generated #{@jekyll_contributing_path}"
  end

  def update_readme_with_downloads
    puts "  📄 Updating README.md with download links..."
    
    readme_content = File.read(@readme_path)
    
    # Check if download section already exists
    if readme_content.include?('## 📦 Download Pre-built Binaries')
      puts "    ⚠️  Download section already exists in README.md"
      return
    end
    
    # Find the installation section and add downloads before it
    download_section = <<~DOWNLOADS

      ## 📦 Download Pre-built Binaries

      Choose the binary for your operating system:

      - **🐧 Linux (AMD64)**: [auto-spotify-linux-amd64](https://github.com/jmervine/auto-spotify/releases/latest/download/auto-spotify-linux-amd64)
      - **🍎 macOS (Intel)**: [auto-spotify-darwin-amd64](https://github.com/jmervine/auto-spotify/releases/latest/download/auto-spotify-darwin-amd64)
      - **🍎 macOS (Apple Silicon)**: [auto-spotify-darwin-arm64](https://github.com/jmervine/auto-spotify/releases/latest/download/auto-spotify-darwin-arm64)
      - **🪟 Windows (AMD64)**: [auto-spotify-windows-amd64.exe](https://github.com/jmervine/auto-spotify/releases/latest/download/auto-spotify-windows-amd64.exe)

      **After downloading:**
      - **Linux/macOS:** `mv auto-spotify-* auto-spotify && chmod +x auto-spotify`
      - **Windows:** `ren auto-spotify-*.exe auto-spotify.exe`

    DOWNLOADS
    
    # Insert download section before installation section
    updated_content = readme_content.sub(
      /## 📦 Installation/,
      download_section + "## 📦 Installation"
    )
    
    if updated_content != readme_content
      File.write(@readme_path, updated_content)
      puts "    ✅ Added download section to README.md"
    else
      puts "    ⚠️  Could not find Installation section in README.md"
    end
  end

  def process_content_for_jekyll(content)
    # Convert relative links to Jekyll-compatible links
    content = content.gsub(/\]\(\.\/([^)]+)\)/, ']({{ site.baseurl }}/\1)')
    
    # Fix any remaining relative links
    content = content.gsub(/\]\(([^http][^)]*\.md)\)/) do |match|
      file = $1
      if file == 'CONTRIBUTING.md'
        ']({{ site.baseurl }}/contributing/)'
      else
        match
      end
    end
    
    content
  end
end

# Run the generator
if __FILE__ == $0
  generator = JekyllGenerator.new
  generator.generate_all
end
