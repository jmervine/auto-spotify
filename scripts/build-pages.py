#!/usr/bin/env python3
"""
Build GitHub Pages site from README and binaries using HTML templates
"""
import re
import os
import shutil

def load_template(template_path):
    """Load HTML template from file"""
    with open(template_path, 'r') as f:
        return f.read()

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
    
    # Load templates
    main_template = load_template('templates/index.html')
    download_section = load_template('templates/download-section.html')
    
    # Insert download section after installation header
    html_content = html_content.replace('<h2>📦 Installation</h2>', '<h2>📦 Installation</h2>' + download_section)
    
    # Replace template placeholders
    final_html = main_template.replace('{{CONTENT}}', html_content)
    final_html = final_html.replace('{{DOWNLOAD_SECTION}}', download_section)
    
    # Create docs-site directory
    os.makedirs('docs-site', exist_ok=True)
    
    # Write HTML file
    with open('docs-site/index.html', 'w') as f:
        f.write(final_html)
    
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
