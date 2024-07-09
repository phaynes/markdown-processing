import json
import subprocess
import sys
import os
import re
import argparse

def run_command(command, description):
    print(f"Running {description} command")
    try:
        subprocess.run(command, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        print(f"{description} completed successfully.")
    except subprocess.CalledProcessError as e:
        print(f"Error running {description}")
        raise

def remove_files(base_name):
    extensions = ['.log', '.aux', '.bbl', '.bcf', '.blg', '.run.xml', '.out', '.xdv', '.tex']
    for ext in extensions:
        try:
            os.remove(f"{base_name}{ext}")
        except OSError:
            pass

def post_process_tex(tex_file):
    with open(tex_file, 'r') as file:
        content = file.read()

    # Modify the CSLReferences environment
    content = re.sub(r'(\\begin\{CSLReferences\})\{1\}\{0\}', r'\1{1}{}', content)

    # Remove any standalone '0' that might appear right after CSLReferences
    content = re.sub(r'(\\begin\{CSLReferences\}\{1\}\{\})\s*0', r'\1', content)

    with open(tex_file, 'w') as file:
        file.write(content)

def run_pandoc(config, output_format):
    input_files = config['input_files']
    output_file = config['output_file']

    # Update the output file extension based on the format
    output_file = os.path.splitext(output_file)[0] + f'.{output_format}'
    config['output_file'] = output_file

    if output_format == 'pdf':
        return run_pandoc_pdf(config)
    elif output_format == 'docx':
        return run_pandoc_docx(config)
    elif output_format == 'html':
        return run_pandoc_html(config)
    elif output_format == 'md':
        return run_pandoc_markdown(config)
    else:
        raise ValueError(f"Unsupported output format: {output_format}")

def run_pandoc_pdf(config):
    input_files = config['input_files']
    output_file = config['output_file']
    latex_file = output_file.replace('.pdf', '.tex')
    base_name = latex_file.replace('.tex', '')

    pandoc_command = [
        'pandoc',
        '--from=markdown+tex_math_single_backslash',
        '--to=latex',
        f'--output={latex_file}',
        f'--template=/usr/local/share/pandoc/data/apa7.latex',
        f'--bibliography={config["bibliography"]}',
        f'--csl=/usr/local/share/pandoc/data/apa.csl',
        '--citeproc',
        '--standalone',
        '--variable', 'documentclass=apa7',
        '--variable', 'classoption=man',
        '--variable', 'biblatexoptions=style=apa,sortcites=true,sorting=nyt,backend=biber',
        '--mathjax',

    ]

    add_metadata(pandoc_command, config)
    pandoc_command.extend(input_files)

    run_command(pandoc_command, "Pandoc")

    # Post-process the .tex file to modify the CSLReferences environment
    post_process_tex(latex_file)

    # Run xelatex
    xelatex_command = ['xelatex', '-interaction=nonstopmode', '-no-pdf', latex_file]
    run_command(xelatex_command, "XeLaTeX (first run)")

    # Run biber
    biber_command = ['biber', base_name]
    run_command(biber_command, "Biber")

    # Run xelatex twice more
    xelatex_command = ['xelatex', '-interaction=nonstopmode', latex_file]
    run_command(xelatex_command, "XeLaTeX (second run)")
    run_command(xelatex_command, "XeLaTeX (final run)")

    # Remove all unnecessary files, keeping only the PDF
    remove_files(base_name)

def run_pandoc_markdown(config):
    input_files = config['input_files']
    output_file = config['output_file']
    latex_file = output_file.replace('.md', '.tex')

    # Step 1: Convert to LaTeX (similar to PDF process)
    latex_command = [
        'pandoc',
        '--from=markdown+tex_math_single_backslash',
        '--to=latex',
        f'--output={latex_file}',
        f'--template=/usr/local/share/pandoc/data/apa7.latex',
        f'--bibliography={config["bibliography"]}',
        f'--csl=/usr/local/share/pandoc/data/apa.csl',
        '--citeproc',
        '--standalone',
        '--variable', 'documentclass=apa7',
        '--variable', 'classoption=man',
        '--variable', 'biblatexoptions=style=apa,sortcites=true,sorting=nyt,backend=biber',
    ]
    add_metadata(latex_command, config)
    latex_command.extend(input_files)
    run_command(latex_command, "Pandoc (Markdown to LaTeX)")

    # Post-process the .tex file to modify the CSLReferences environment
    post_process_tex(latex_file)

    # Step 2: Convert LaTeX back to Markdown
    md_command = [
        'pandoc',
        '--from=latex',
        '--to=markdown_strict-raw_html-fenced_divs-bracketed_spans',
        f'--output={output_file}',
        '--wrap=preserve',
        '--markdown-headings=atx',
        latex_file
    ]
    run_command(md_command, "Pandoc (LaTeX to Markdown)")


    # Post-process the markdown file
    with open(output_file, 'r') as file:
        content = file.read()

    # Remove any remaining div-like structures
    content = re.sub(r':::.*?\n', '', content, flags=re.DOTALL)
    content = re.sub(r'{#.*?}', '', content)


    # Ensure there's a blank line before the References section
    content = re.sub(r'(# References)', r'\n\1', content)

    # Remove any extra newlines at the end of the document
    content = content.rstrip() + '\n'

    with open(output_file, 'w') as file:
        file.write(content)

    print(f"Markdown file with processed citations and references created: {output_file}")

    # Optionally, remove the temporary LaTeX file
    os.remove(latex_file)

    return output_file

def add_metadata(pandoc_command, config):
    for key in ['title', 'author', 'affiliation', 'course', 'instructor', 'date', 'shorttitle', 'keywords', 'bibliography']:
        if key in config:
            value = config[key]
            if isinstance(value, list):
                for item in value:
                    pandoc_command.extend([f'--metadata={key}:{item}'])
            else:
                pandoc_command.extend([f'--metadata={key}:{value}'])

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Convert Markdown to APA7 format')
    parser.add_argument('config', help='Path to the config JSON file')
    parser.add_argument('--format', choices=['pdf', 'docx', 'html', 'md'], default='pdf',
                        help='Output format (default: pdf)')
    args = parser.parse_args()

    if not os.path.exists(args.config):
        print(f"Error: Config file '{args.config}' not found.")
        sys.exit(1)

    with open(args.config, 'r') as f:
        config = json.load(f)

    try:
        run_pandoc(config, args.format)
    except Exception as e:
        print(f"Error during conversion: {e}")
        sys.exit(1)

    print("Conversion completed successfully.")
