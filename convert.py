import json
import subprocess
import sys
import os

def run_pandoc(config):
    input_files = config['input_files']
    output_file = config['output_file']
    latex_file = output_file.replace('.pdf', '.tex')

    command = [
        'pandoc',
        '--from=markdown+tex_math_single_backslash',
        '--to=latex',
        f'--output={latex_file}',  # Output to a .tex file
        '--verbose',
        '--template=apa7.latex',
        f'--bibliography={config["bibliography"]}',
        '--csl=apa.csl',
        '--citeproc',
        '--standalone',
        '--variable', 'documentclass=apa7',
        '--variable', 'classoption=man',
        '--variable', 'biblatexoptions=style=apa,sortcites=true,sorting=nyt,backend=biber',
    ]

    # Add metadata
    for key in ['title', 'author', 'affiliation', 'course', 'instructor', 'date', 'shorttitle', 'keywords', 'bibliography']:
        if key in config:
            value = config[key]
            if isinstance(value, list):
                for item in value:
                    command.extend([f'--metadata={key}:{item}'])
            else:
                command.extend([f'--metadata={key}:{value}'])

    # Ensure shorttitle is set
    if 'shorttitle' not in config:
        command.extend(['--metadata=shorttitle:' + config.get('title', '').split(':')[0]])

    command.extend(input_files)

    print("Running Pandoc command:", ' '.join(command))

    try:
        result = subprocess.run(command, check=True, capture_output=True, text=True)
        print("Pandoc output:")
        print(result.stdout)
        print("Pandoc error output:")
        print(result.stderr)

        # Now run xelatex on the generated LaTeX file
        xelatex_command = ['xelatex', '-interaction=nonstopmode', latex_file]
        print("Running XeLaTeX command:", ' '.join(xelatex_command))
        xelatex_result = subprocess.run(xelatex_command, check=True, capture_output=True, text=True)
        print("XeLaTeX output:")
        print(xelatex_result.stdout)
        print("XeLaTeX error output:")
        print(xelatex_result.stderr)

    except subprocess.CalledProcessError as e:
        print("Error running Pandoc or XeLaTeX:")
        print(e.stdout)
        print("Error output:")
        print(e.stderr)
        raise

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
