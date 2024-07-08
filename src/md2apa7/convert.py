import json
import subprocess
import sys
import os
import re

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

def run_pandoc(config):
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
    ]

    for key in ['title', 'author', 'affiliation', 'course', 'instructor', 'date', 'shorttitle', 'keywords', 'bibliography']:
        if key in config:
            value = config[key]
            if isinstance(value, list):
                for item in value:
                    pandoc_command.extend([f'--metadata={key}:{item}'])
            else:
                pandoc_command.extend([f'--metadata={key}:{value}'])

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

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: python3 convert.py <config.json>")
        sys.exit(1)

    config_path = sys.argv[1]
    if not os.path.exists(config_path):
        print(f"Error: Config file '{config_path}' not found.")
        sys.exit(1)

    with open(config_path, 'r') as f:
        config = json.load(f)

    try:
        run_pandoc(config)
    except Exception as e:
        print(f"Error during conversion: {e}")
        sys.exit(1)

    print("Conversion completed successfully.")
