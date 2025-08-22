#!/usr/bin/env ruby
# Generate Jekyll site from README.md
# This keeps the Jekyll site in sync with the primary documentation
# Works like magic - any changes to README automatically appear in the site

require 'fileutils'

class JekyllGenerator
  def initialize
    @root_dir = Dir.pwd
    @readme_path = File.join(@root_dir, 'README.md')
    @jekyll_index_path = File.join(@root_dir, 'index.md')
  end

  def generate_all
    puts "🔄 Generating Jekyll site from README..."
    
    # Generate Jekyll index from README
    generate_jekyll_index
    
    puts "✅ Jekyll site generated successfully!"
  end

  private

  def generate_jekyll_index
    puts "  📄 Generating index.md from README.md..."
    
    readme_content = File.read(@readme_path)
    
    # Extract title from first heading
    title_match = readme_content.match(/^# (.+)$/)
    title = title_match ? title_match[1] : "Home"
    
    # Process content - keep everything exactly as written, just add Jekyll magic
    processed_content = process_content_for_jekyll(readme_content)
    
    # Create Jekyll front matter and content
    jekyll_content = create_jekyll_frontmatter("default", "Home") + processed_content
    
    # Write the Jekyll index file
    File.write(@jekyll_index_path, jekyll_content)
    puts "    ✅ Generated #{@jekyll_index_path}"
  end

  def create_jekyll_frontmatter(layout, title, permalink = nil)
    frontmatter = "---\nlayout: #{layout}\ntitle: #{title}\n"
    frontmatter += "permalink: #{permalink}\n" if permalink
    frontmatter += "---\n\n"
    frontmatter
  end

  def process_content_for_jekyll(content)
    # Fix relative links to work with Jekyll
    content = content.gsub(/\]\(\.\/([^)]+)\)/, ']({{ site.baseurl }}/\1)')
    
    # Convert CONTRIBUTING.md links to GitHub page
    content = content.gsub(/\[([^\]]*)\]\(CONTRIBUTING\.md\)/, '[\1](https://github.com/jmervine/auto-spotify/blob/master/CONTRIBUTING.md)')
    content = content.gsub(/\[([^\]]*)\]\(\.\/CONTRIBUTING\.md\)/, '[\1](https://github.com/jmervine/auto-spotify/blob/master/CONTRIBUTING.md)')
    
    # Convert README.md links to home page
    content = content.gsub(/\[([^\]]*)\]\(README\.md\)/, '[\1]({{ site.baseurl }}/)')
    content = content.gsub(/\[([^\]]*)\]\(\.\/README\.md\)/, '[\1]({{ site.baseurl }}/)')
    
    # Keep GitHub releases links as-is (binaries are now served from GitHub Releases)
    # No conversion needed - links should point directly to GitHub releases
    
    # Convert other .md files to Jekyll-style links (remove .md extension)
    content = content.gsub(/\[([^\]]*)\]\(([^)]*?)\.md\)/) do |match|
      link_text = $1
      file_path = $2
      if file_path.start_with?('http')
        match # Keep external links as-is
      else
        "[#{link_text}]({{ site.baseurl }}/#{file_path.gsub(/^\.\//, '')}/)"
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