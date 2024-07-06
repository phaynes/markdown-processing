import sys
import os
import openai
import bibtexparser

def read_api_key(file_path):
    with open(file_path, 'r') as file:
        return file.read().strip()

def check_apa7_compliance(markdown_file):
    with open(markdown_file, 'r') as file:
        content = file.read()

    prompt = f"""
    Please review the following list of references and check if they comply with APA 7th edition format.
    Provide a brief summary of any errors or inconsistencies found.

    References:
    {content}
    """

    response = openai.ChatCompletion.create(
        model="gpt-4o",
        messages=[
            {"role": "system", "content": "You are an expert in APA 7th edition formatting."},
            {"role": "user", "content": prompt}
        ]
    )

    return response.choices[0].message['content']

def check_bibtex_correctness(bibtex_file):
    with open(bibtex_file, 'r') as bibtex_file:
        bib_database = bibtexparser.load(bibtex_file)

    entries = bib_database.entries
    errors = []

    for entry in entries:
        if 'author' not in entry or 'year' not in entry or 'title' not in entry:
            errors.append(f"Entry {entry.get('ID', 'Unknown')} is missing required fields (author, year, or title)")

    if errors:
        return "\n".join(errors)
    else:
        return "All BibTeX entries appear to be correctly formatted."

def main():
    print(f"Current working directory: {os.getcwd()}")
    print(f"Contents of /data directory: {os.listdir('/data')}")

    markdown_file = '/data/input.md'
    bibtex_file = '/data/output.bib'
    api_key_file = os.environ.get('OPENAI_API_KEY_FILE')

    if not api_key_file:
        print("Error: OPENAI_API_KEY_FILE environment variable is not set.")
        sys.exit(1)

    openai.api_key = read_api_key(api_key_file)

    if not openai.api_key:
        print("Error: Unable to read OpenAI API key from file.")
        sys.exit(1)

    print("Checking APA7 compliance...")
    apa7_result = check_apa7_compliance(markdown_file)
    print(apa7_result)

    print("\nChecking BibTeX correctness...")
    bibtex_result = check_bibtex_correctness(bibtex_file)
    print(bibtex_result)

if __name__ == "__main__":
    main()
