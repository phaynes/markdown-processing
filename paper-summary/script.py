import os
import xml.etree.ElementTree as ET
from openai import OpenAI
from typing import List, Dict, Optional
import logging
import bibtexparser
from dataclasses import dataclass

@dataclass
class Paper:
    title: Optional[str]
    authors: List[str]
    year: Optional[str]
    abstract: Optional[str]
    introduction: Optional[str]
    conclusion: Optional[str]
    body: Optional[str]
    filename: str

def load_api_key(file_path: str) -> str:
    with open(file_path, 'r') as file:
        return file.read().strip()

client = OpenAI(api_key=load_api_key("/usr/src/app/openai_api_key.txt"))

def safe_find_text(element: Optional[ET.Element], xpath: str, namespaces: Dict[str, str]) -> Optional[str]:
    if element is None:
        return None
    found = element.find(xpath, namespaces)
    return found.text if found is not None else None

def parse_xml(file_path: str) -> Optional[Paper]:
    try:
        tree = ET.parse(file_path)
        root = tree.getroot()
        ns = {'tei': 'http://www.tei-c.org/ns/1.0'}

        title = safe_find_text(root, ".//tei:titleStmt/tei:title", ns)
        authors = [safe_find_text(author, "tei:surname", ns) or "" for author in root.findall(".//tei:author", ns)]
        year = safe_find_text(root, ".//tei:date", ns)
        abstract = safe_find_text(root, ".//tei:abstract", ns)
        introduction = safe_find_text(root, ".//tei:body//tei:div[@type='introduction']", ns)
        conclusion = safe_find_text(root, ".//tei:body//tei:div[@type='conclusion']", ns)

        # Extract the full body text
        body_element = root.find(".//tei:body", ns)
        body = ET.tostring(body_element, encoding='unicode', method='text') if body_element is not None else ""

        return Paper(title, authors, year, abstract, introduction, conclusion, body, os.path.basename(file_path))
    except ET.ParseError as e:
        logging.error(f"XML parsing error in {file_path}: {e}")
        return None

def summarize_text(text: str, query: str) -> str:
    try:
        # Truncate the text if it's too long
        max_tokens = 4000  # Adjust this value based on your needs and model limits
        truncated_text = text[:max_tokens * 4]  # Approximate 1 token to 4 characters

        response = client.chat.completions.create(
            model="gpt-4",
            messages=[
                {"role": "system", "content": "You are a helpful assistant that summarizes academic papers."},
                {"role": "user", "content": f"Summarize the following text in relation to this query: {query}\n\n{truncated_text}"}
            ],
            max_tokens=500  # Increased for a more comprehensive summary
        )
        return response.choices[0].message.content.strip()
    except Exception as e:
        logging.error(f"Error in summarization: {e}")
        return f"Error in summarization: {str(e)}"

def generate_apa_reference(paper: Paper) -> str:
    authors = ", ".join(filter(None, paper.authors))
    return f"{authors} ({paper.year or 'n.d.'}). {paper.title or 'Untitled'}. Journal Title. Volume(Issue), pages."

def generate_bibtex_reference(paper: Paper) -> str:
    authors = " and ".join(filter(None, paper.authors))
    return f"""@article{{{paper.authors[0].lower() if paper.authors else 'unknown'}{paper.year or 'unknown'},
  title={{{paper.title or 'Untitled'}}},
  author={{{authors}}},
  year={{{paper.year or 'n.d.'}}},
  journal={{Journal Title}},
  volume={{Volume}},
  number={{Issue}},
  pages={{pages}}
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

                output_filename = f"{paper.year or 'unknown'}-({paper.authors[0] if paper.authors else 'unknown'})-{(paper.title or 'Untitled').replace(' ', '-')[:50]}-summary.md"
                output_path = os.path.join("/usr/src/app/output", output_filename)
                with open(output_path, 'w') as md_file:
                    md_file.write(markdown_content)

                summaries.append({"filename": output_filename, "content": markdown_content, "relevance": len(summary)})
                logging.info(f"Processed {filename}")
            except Exception as e:
                logging.error(f"Error processing {filename}: {e}")

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
