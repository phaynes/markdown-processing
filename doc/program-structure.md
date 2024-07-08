# Markdown to APA7 PDF Converter

## Introd
 - command in path.
  - being able to see markdown in browser.

## System Overview
This system is designed to convert Markdown documents into PDF format, adhering to APA 7th edition guidelines. It uses a combination of Pandoc for Markdown processing, LaTeX for PDF generation, and custom Python scripts for orchestration and bibliography conversion. The system is containerized using Docker for consistent performance across different environments.

Additionally, it also has:
1. APAv7 to bibtex converter.
2. APAv7 reference checker (written in Python)

These additional elements work together to convert APA7 style references to BibTeX format and then check the references for compliance and accuracy.

## Directory Structure
```
project/
├── Dockerfile
├── apa7.latex
├── apa.csl
├── convert.py
├── config.json
├── input/
│   └── your-markdown-file.md
└── apa7_to_bibtex/
      ├── docker-compose.yml
      ├── data/
      │   ├── input.md
      │   └── output.bib (generated)
      ├── rust/
      │   ├── Dockerfile
      │   └── src/
      │       └── main.rs
      └── python/
      ├── Dockerfile
      └── reference_checker.py
```

## Components
1. **Dockerfile**: Defines the container environment, including all necessary dependencies.
2. **LaTeX Template (apa7.latex)**: A custom LaTeX template that ensures APA 7 compliance.
3. **CSL File (apa.csl)**: Citation Style Language file for APA 7 citations.
4. **Python Script (convert.py)**: Orchestrates the main conversion process.
5. **Configuration File (config.json)**: Contains metadata and settings for the document.
6. **docker-compose.yml**: Orchestrates the two services (apa7_to_bibtex and reference_checker).
7. **apa7_to_bibtex (Rust)**: Converts APA7 style references to BibTeX format.
8. **reference_checker (Python)**: Checks the converted references for accuracy and compliance.

## Docker Compose Configuration

The `docker-compose.yml` file defines two services:

1. **apa7_to_bibtex**:
   - Built from the Rust Dockerfile in the `./rust` directory.
   - Mounts the `./data` directory to `/data` in the container.
   - Runs the `apa7_to_bibtex` command on `input.md` to produce `output.bib`.

2. **reference_checker**:
   - Built from the Python Dockerfile in the `./python` directory.
   - Mounts the `./data` directory and an OpenAI API key file.
   - Runs the `reference_checker.py` script.
   - Depends on the `apa7_to_bibtex` service to complete first.

## Build and Run Instructions

1. Ensure Docker and Docker Compose are installed on your system.
2. Download the apa.csl file by running the command `./download-csl.sh`.
3. Build the container by running `docker build -t apa7-converter .`.
4. Place your input Markdown file with APA7 references in the `data/` directory as `input.md`.
5. Build and run the services:
   ```
   docker-compose up --build
   ```

This command will:
- Build the Rust and Python Docker images
- Run the APA7 to BibTeX conversion
- Run the reference checker on the converted BibTeX file

## APA7 to BibTeX Converter (Rust)

The Rust program in `rust/src/main.rs` converts APA7 style references to BibTeX format. It reads from `input.md` and writes to `output.bib` in the `data/` directory.

## Reference Checker (Python)

The Python script `python/reference_checker.py` checks the converted references for accuracy and compliance. It uses the OpenAI API to assist in the checking process.

## Maintenance and Customisation

- To modify the APA7 to BibTeX conversion, update the Rust code in `rust/src/main.rs`.
- To change the reference checking process, modify `python/reference_checker.py`.
- To adjust the Docker environments, update the respective Dockerfiles in the `rust/` and `python/` directories.
- To change how the services interact, modify the `docker-compose.yml` file.

## Troubleshooting

- For issues with the APA7 to BibTeX conversion, check the Rust service logs in the Docker Compose output.
- For problems with the reference checker, examine the Python service logs.
- Ensure the OpenAI API key file is correctly located and mounted in the reference_checker service.
- Verify that the `data/` directory contains the expected input file before running the services.

This system provides an automated pipeline for converting APA7 references to BibTeX format and then checking those references for accuracy, all within a containerised environment for consistency and ease of use.


Roadmap:

1. An ability to run the conversion to  produce a  word file or html page, in addition the pdf.
