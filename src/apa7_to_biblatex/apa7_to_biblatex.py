import os
import sys
import openai

def read_api_key(file_path):
    with open(file_path, 'r') as file:
        return file.read().strip()

def convert_apa7_to_biblatex(apa7_reference):
    prompt = f"""
    Convert the following APA7 reference in markdown format to BibLaTeX format:

    {apa7_reference}

    Please provide only the BibLaTeX entry without any additional text, explanations, or markdown formatting.
    Ensure the entry type is a standard BibLaTeX type (e.g., article, book, inproceedings, etc.).
    Include all relevant fields, such as author, title, year, journal, volume, number, pages, doi, url, etc.    """

    response = openai.ChatCompletion.create(
        model="gpt-4o",
        messages=[
            {"role": "system", "content": "You are an expert in converting APA7 references in markdown format to BibLaTeX format."},
            {"role": "user", "content": prompt}
        ]
    )

    biblatex_entry = response.choices[0].message['content'].strip()

    # Remove potential markdown artifacts
    biblatex_entry = biblatex_entry.replace("bibtex", "").replace("bib", "").replace("latex", "").replace("```", "").strip()

    return biblatex_entry

def process_references(input_file, output_file):
    with open(input_file, 'r') as file:
        apa7_references = file.read().strip().split('\n\n')

    biblatex_entries = []

    total_refs = len(apa7_references)
    for i, ref in enumerate(apa7_references, 1):
        print(f"Processing reference {i}/{total_refs}")
        try:
            biblatex_entry = convert_apa7_to_biblatex(ref)
            biblatex_entries.append(biblatex_entry)
            print(f"Successfully converted reference {i}")
        except Exception as e:
            print(f"Error processing entry {i}: {str(e)}")
            print(f"Problematic entry: {ref}")

    with open(output_file, 'w') as bibfile:
        bibfile.write('\n\n'.join(biblatex_entries))

def main():
    input_file = '/data/input.md'
    output_file = '/data/output.bib'
    api_key_file = os.environ.get('OPENAI_API_KEY_FILE')

    if not api_key_file:
        print("Error: OPENAI_API_KEY_FILE environment variable is not set.")
        sys.exit(1)

    openai.api_key = read_api_key(api_key_file)

    if not openai.api_key:
        print("Error: Unable to read OpenAI API key from file.")
        sys.exit(1)

    print("Converting APA7 references to BibLaTeX format...")
    process_references(input_file, output_file)
    print(f"Conversion complete. BibLaTeX entries saved to {output_file}")

if __name__ == "__main__":
    main()
