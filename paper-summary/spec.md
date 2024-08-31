### **Summarizer Specification for Document Summarization Program**

**Objective**: Develop a program that summarizes a set of academic papers, extracted in Grobid XML format, based on an academic question. The program should generate markdown files with specific metadata and content structure, and then compile all summaries into a single grouped file.

#### **Input**:

1. **Academic Papers**: A directory containing academic papers in Grobid XML format. Each XML file represents an extracted academic paper.
2. **Academic Question**: A guiding question or topic to focus the summaries on.

#### **Output**:

For each academic paper:

1. **Markdown Summary File**:
   - The program should generate a markdown file for each academic paper.
   - **Filename Format**: `{paper-year}-({first author surname})-{paper-name-summary}.md`
     - Example: `2024-(Smith)-Neural-Networks-and-Learning-summary.md`

2. **Content of Each Markdown File**:
   - **Header Information**:
     1. **Paper Name**: The full title of the paper.
     2. **APA 7th Edition Reference**: A citation formatted in APA 7th edition style.
     3. **BibTeX Reference**: The citation in `.bib` format for reference management.
     4. **Original PDF File Name**: The filename of the original PDF from which the paper was extracted.

   - **Summary Content**:
     - A concise summary of the paper, tailored to the provided academic question. The summary should focus on the paper's relevance to the question, key findings, methods, and conclusions.

3. **Compilation of Summaries**:
   - After all markdown summaries are created, the program should compile these summaries into a single markdown file.
   - **Grouping Criteria**:
     - Group summaries based on their relevance to the academic question and the quality of the paper. High-quality papers and those most relevant to the academic question should be prioritized.

#### **Program Requirements**:

- **Read and Parse XML**: The program must read and parse Grobid XML files to extract necessary information such as the paper title, authors, abstract, introduction, conclusion, and publication year.
- **Generate References**:
  - **APA Reference**: Automatically generate an APA 7th edition citation from the extracted metadata.
  - **BibTeX Reference**: Automatically generate a BibTeX citation entry from the extracted metadata.
- **Summarization**:
  - Utilize a language model or summarization algorithm to generate a summary of each paper guided by the academic question.
- **Markdown Formatting**:
  - Create markdown files with the specified filename and structure.
- **Output Compilation**:
  - Compile individual summaries into a single markdown file, organized by relevance and paper quality.

#### **Additional Considerations**:

- The program should handle errors gracefully, such as missing data in XML files or incorrect formatting.
- Ensure the output is formatted consistently and is easy to navigate.
- Provide logs or a report indicating the processing status of each file and any issues encountered.
