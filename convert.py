import json
import subprocess
import sys
import os

def run_pandoc(config):
    input_files = config['input_files']
    output_file = config['output_file']

    command = [
        'pandoc',
        '--from=markdown+tex_math_single_backslash',
        '--to=pdf',
        f'--output={output_file}',
        '--pdf-engine=xelatex',
        '--template=apa7.latex',
        f'--bibliography={config["bibliography"]}',
        '--csl=apa.csl',
        '--citeproc'
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

    command.extend(input_files)

    print("Running command:", ' '.join(command))

    try:
        result = subprocess.run(command, check=True, capture_output=True, text=True)
        print("Pandoc output:")
        print(result.stdout)
    except subprocess.CalledProcessError as e:
        print("Error running Pandoc:")
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
