import os
import xml.etree.ElementTree as ET
from openai import OpenAI
from typing import List, Dict, Optional
import logging
import bibtexparser
import unicodedata
import re
from datetime import datetime

from dataclasses import dataclass

@dataclass
class Paper:
    title: Optional[str]
    authors: List[str]
    year: Optional[str]
    full_date: Optional[str]
    journal_title: Optional[str]
    volume: Optional[str]
    issue: Optional[str]
    pages: Optional[str]
    doi: Optional[str]
    abstract: Optional[str]
    introduction: Optional[str]
    conclusion: Optional[str]
    body: Optional[str]
    filename: str
    publisher: Optional[str]

def load_api_key(file_path: str) -> str:
    with open(file_path, 'r') as file:
        return file.read().strip()

client = OpenAI(api_key=load_api_key("/usr/src/app/openai_api_key.txt"))

def safe_find_text(element: Optional[ET.Element], xpath: str, namespaces: Dict[str, str]) -> Optional[str]:
    if element is None:
        return None
    found = element.find(xpath, namespaces)
    return found.text if found is not None else None

def remove_diacritics(text):
    """Remove diacritics from the given text."""
    return ''.join(c for c in unicodedata.normalize('NFKD', text)
                   if unicodedata.category(c) != 'Mn')

def safe_filename(text):
    """Convert text to a safe filename."""
    text = remove_diacritics(text)
    return re.sub(r'[^a-zA-Z0-9]', '', text).lower()

def extract_paper_authors(root: ET.Element, namespaces: Dict[str, str]) -> List[str]:
    """Extract all authors from the <analytic> section."""
    authors = []
    analytic_section = root.find(".//tei:analytic", namespaces)
    if analytic_section is not None:
        for author in analytic_section.findall("tei:author", namespaces):
            persName = author.find("tei:persName", namespaces)
            if persName is not None:
                surname = persName.find("tei:surname", namespaces)
                forename = persName.find("tei:forename", namespaces)
                if surname is not None and forename is not None:
                    authors.append(f"{surname.text}, {forename.text[0]}.")
    return authors

def parse_date(date_string):
    if not date_string:
        return None

    date_string = date_string.strip()

    # Check if the date format is "YYYY-MM-DD" or "YYYY-MM"
    if '-' in date_string:
        return date_string.split('-')[0]

    # Try to find a 4-digit year anywhere in the string
    match = re.search(r'\b(\d{4})\b', date_string)
    if match:
        return match.group(1)

    # If no year could be parsed, return None
    return None


def parse_xml(file_path: str) -> Optional[Paper]:
    try:
        tree = ET.parse(file_path)
        root = tree.getroot()
        ns = {'tei': 'http://www.tei-c.org/ns/1.0'}

        title = safe_find_text(root, ".//tei:titleStmt/tei:title", ns)
        authors = extract_paper_authors(root, ns)

        # Extract year and full date
        date = safe_find_text(root, ".//tei:publicationStmt/tei:date[@type='published']", ns)
        year = parse_date(date)
        full_date = date

        # Extract journal information
        journal_title = safe_find_text(root, ".//tei:monogr/tei:title[@level='j']", ns)
        volume = safe_find_text(root, ".//tei:monogr/tei:imprint/tei:biblScope[@unit='volume']", ns)
        issue = safe_find_text(root, ".//tei:monogr/tei:imprint/tei:biblScope[@unit='issue']", ns)

        # Extract page numbers
        page_elem = root.find(".//tei:monogr/tei:imprint/tei:biblScope[@unit='page']", ns)
        pages = f"{page_elem.get('from')}-{page_elem.get('to')}" if page_elem is not None else None

        # Extract DOI
        doi = safe_find_text(root, ".//tei:idno[@type='DOI']", ns)

        # Extract publisher
        publisher = safe_find_text(root, ".//tei:publicationStmt/tei:publisher", ns)

        abstract = safe_find_text(root, ".//tei:abstract", ns)
        introduction = safe_find_text(root, ".//tei:body//tei:div[@type='introduction']", ns)
        conclusion = safe_find_text(root, ".//tei:body//tei:div[@type='conclusion']", ns)

        # Extract the full body text
        body_element = root.find(".//tei:body", ns)
        body = ET.tostring(body_element, encoding='unicode', method='text') if body_element is not None else ""

        output_filename = os.path.basename(file_path)

        return Paper(title, authors, year, full_date, journal_title, volume, issue, pages, doi, abstract, introduction, conclusion, body, output_filename, publisher)

    except ET.ParseError as e:
        logging.error(f"XML parsing error in {file_path}: {e}")
        return None
    except Exception as e:
        logging.error(f"Error processing {file_path}: {str(e)}")
        return None

def summarize_text(text: str, query: str) -> str:
    try:
        # Truncate the text if it's too long
        max_tokens = 10000  # Adjust this value based on your needs and model limits
        truncated_text = text[:max_tokens * 4]  # Approximate 1 token to 4 characters

        response = client.chat.completions.create(
            model="gpt-4o-2024-08-06",
            messages=[
                {"role": "system", "content": "You are a helpful assistant that summarizes academic papers."},
                {"role": "user", "content": f"Summarize the following text in relation to this query: {query}\n\n{truncated_text}"}
            ],
            max_tokens=400  # Increased for a more comprehensive summary
        )
        return response.choices[0].message.content.strip()
    except Exception as e:
        logging.error(f"Error in summarization: {e}")
        return f"Error in summarization: {str(e)}"

def generate_apa_reference(paper: Paper) -> str:
    if len(paper.authors) > 2:
        citation_authors = f"{paper.authors[0]}, et al."
    else:
        citation_authors =", ".join(paper.authors) if paper.authors else "Unknown"

    year = paper.year or "n.d."
    title = paper.title or "Untitled"
    journal = paper.journal_title or "Unknown Journal"
    doi = paper.doi or ""

    apa_ref = f"{citation_authors} ({year}). {title}. *{journal}*"
    if paper.volume:
        apa_ref += f", {paper.volume}"
        if paper.issue:
            apa_ref += f"({paper.issue})"
    if paper.pages:
        apa_ref += f", {paper.pages}"
    if doi:
        apa_ref += f". https://doi.org/{doi}"
    apa_ref += "."

    return apa_ref

def generate_bibtex_reference(paper: Paper) -> str:
    authors = " and ".join(paper.authors) if paper.authors else "Unknown"
    year = paper.year or "unknown"
    title = paper.title or "Untitled"

    key_authors = remove_diacritics(f"{paper.authors[0].split(',')[0] if paper.authors else 'Unknown'}")

    key = f"{key_authors}{year}"
    journal = paper.journal_title or "Unknown Journal"
    volume = paper.volume or ""
    number = paper.issue or ""
    pages = paper.pages or ""
    doi = paper.doi or ""
    publisher = paper.publisher or ""

    return f"""@article{{{key},
  author = {{{authors}}},
  title = {{{title}}},
  journal = {{{journal}}},
  year = {{{year}}},
  volume = {{{volume}}},
  number = {{{number}}},
  pages = {{{pages}}},
  doi = {{{doi}}},
  url = {{https://doi.org/{doi}}},
  publisher = {{{publisher}}},
  issn = {{2195-7185}}
}}"""

def create_markdown_summary(paper: Paper, summary: str, query: str) -> str:
    apa_reference = generate_apa_reference(paper)
    bibtex_reference = generate_bibtex_reference(paper)

    markdown_content = f"""# {paper.title or 'Untitled'}

**APA Reference:** {apa_reference}

**BibTeX Reference:**
```bibtex
{bibtex_reference}
```

**Original PDF File Name:** {paper.filename}

**Summary (in relation to: "{query}"):**
{summary}
"""
    return markdown_content

def process_papers(directory_path: str, query: str) -> List[Dict[str, str]]:
    summaries = []
    for filename in os.listdir(directory_path):
        if filename.endswith(".xml"):
            file_path = os.path.join(directory_path, filename)
            try:
                paper = parse_xml(file_path)
                if paper is None:
                    continue

                # Combine all available text content
                text_to_summarize = f"Title: {paper.title or ''}\n\n"
                text_to_summarize += f"Abstract: {paper.abstract or ''}\n\n"
                text_to_summarize += f"Introduction: {paper.introduction or ''}\n\n"
                text_to_summarize += f"Body: {paper.body or ''}\n\n"
                text_to_summarize += f"Conclusion: {paper.conclusion or ''}"

                # Trim the text if it's too long (adjust the limit as needed)
                max_chars = 10000  # Adjust this value based on your needs and API limits
                if len(text_to_summarize) > max_chars:
                    text_to_summarize = text_to_summarize[:max_chars] + "..."

                summary = summarize_text(text_to_summarize, query)

                markdown_content = create_markdown_summary(paper, summary, query)

                # Create a safe filename
                safe_title = ''.join(c if c.isalnum() else '-' for c in (paper.title or 'Untitled'))
                safe_title = safe_title[:50]  # Limit the length of the title in the filename
                safe_author = safe_filename(paper.authors[0].split(',')[0] if paper.authors else 'unknown')

                output_filename = f"{paper.year or 'unknown'}-({safe_author})-{safe_title}-summary.md"
                output_path = os.path.join("/usr/src/app/output", output_filename)
                with open(output_path, 'w') as md_file:
                    md_file.write(markdown_content)

                summaries.append({"filename": output_filename, "content": markdown_content, "relevance": len(summary)})
                logging.info(f"Processed {filename}")
            except Exception as e:
                logging.error(f"Error processing {filename}: {str(e)}")

    return sorted(summaries, key=lambda x: x['relevance'], reverse=True)

def compile_summaries(summaries: List[Dict[str, str]]) -> str:
    compiled_content = "# Compiled Summaries\n\n"
    for summary in summaries:
        compiled_content += f"## Summary from {summary['filename']}\n\n"
        compiled_content += summary['content'] + "\n\n---\n\n"
    return compiled_content

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    directory_path = "/usr/src/app/input"  # Fixed input directory path
    summarization_query = os.environ.get('SUMMARIZATION_QUERY')

    if not summarization_query:
        print("Error: SUMMARIZATION_QUERY environment variable not set.")
        exit(1)

    summaries = process_papers(directory_path, summarization_query)
    compiled_summaries = compile_summaries(summaries)
    with open("/usr/src/app/output/compiled_summaries.md", 'w') as f:
        f.write(compiled_summaries)
    print("Processing complete. See compiled_summaries.md in the output directory for results.")
