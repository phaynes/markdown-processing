# AI-Assisted Markdown Proofing Tool Development
Markdown provides an efficient method for preparing text as it can be transformed into various formats (PDF / HTML / DOC) and is able to be strongly version controlled using tools such as Git.

However, written work needs to be proofed and reviewed, and the types of tooling support available in word processors are not directly accessible. Instead, proofing the markdown text generally involves copying text segments into an AI chat engine, receiving corrections, and copying them back into the document. This approach is labour-intensive and typically involves respecifying common instructions such as "please proof this text using British English spelling with an active voice.

## Overview

This project aims to develop a command-line tool that automates the process of using AI to proof, review, and comment on markdown text. The tool will use AI technologies (primarily OpenAI, optionally Anthropic) and integrate with Git for version control.

## Key Elements

1. **Executable Applications**:
   - Applications will be developed using Go.
   - They can run as standalone executables or within a Docker container.

2. **Docker Integration**:
   - The application build environment will run out of Docker to support CI/CD.

3. **Version Control**:
   - Git is assumed locally, with either GitHub or Gitea as the remote.
   - Applications should support fully local operations where possible.
   - Advanced features may leverage remote tooling.

4. **AI Technology**:
   - Primarily OpenAI, with optional support for Anthropic.

5. **Tool Binaries**:
   - **mdp** (markdown proof): Fully updates the markdown document without user review.
   - **mdr** (markdown review): Prepares proposed changes for user acceptance or rejection.
   - **mdc** (markdown comment): Provides comments without making specific changes.
   - Tools should support both file inputs and command-line text, with output to the terminal or file.

6. **Configuration**:
   - The tool starts by reading a local config file, providing API keys and default file paths.
   - Example config file:

   ```json
   {
     "input_files": ["2024-05-11-assignment2.md"],
     "output_file": "assignment2.pdf",
     "bibliography": "assignment2-refs.bib",
     "default_ai": "openai",
     "openAI_key": "path\\to\\openaikey.txt",
     "anthropic_key": "path\\to\\AnthropicKey.txt",
     "proofing_prompts": "path\\proofing-prompts.json"
   }
   ```

7. **Proofing Prompts**:
   - A JSON file contains review instructions with labels and flags for criteria to be run on all or specific reviews.
   - An example is provided below.

## Workflow

1. **Simple Changes**:
   - Corrected text is directly updated in the document.

2. **Complex Changes**:
   - User reviews proposed changes, accepting or rejecting them.

3. **Version Control Integration**:
   - Check in the document to the current branch.
   - Create and switch to a new branch for AI proofing.
   - Run AI proofing and update the document.
   - Merge changes back to the main branch if accepted.

4. **Code Logic**:
   ```bash
   git add {document being edited}
   git commit -m "prior to proofing"
   git checkout -b ai-proof-branch
   Run AI proofing engine & update document.
   git add {document being edited}
   git commit -m "Type of proofing run"
   git checkout main
   git merge ai-proof-branch
   git branch -d ai-proof-branch
   ```

## Development Stages

1. **Initial Proof of Concept**:
   - Proof text from the command line, outputting to the console.

2. **Feature Selection**:
   - Allow selection of different proof types.

3. **API Integration**:
   - Support for both OpenAI and Anthropic APIs.

4. **File Handling**:
   - Read a markdown file for proofing, output changes to the terminal.

5. **Basic Git Functions**:
   - Add and check in the file before updates.

6. **Review and Diff**:
   - Perform reviews using both a diff and full version of the input file.
   - Support branch deletion after review.

7. **Full Automation**:
   - Achieve full automated proofing after verifying previous functions.

## Example Command Line Prompts

1. mdp
```bash
mdp input.md
mdp input.md -o output.md
mdp -c custom_config.json input.md
mdp --anthropic input.md
mdp "This is some text to proof directly from the command line."
```

```bash
mdr input.md
mdr input.md -o review.md
mdr -c custom_config.json input.md
mdr --anthropic input.md
mdr "Please review this text for clarity and style."
```

```bash
mdc input.md
mdc input.md -o comments.md
mdc -c custom_config.json input.md
mdc --anthropic input.md
mdc "Add comments to improve this paragraph's structure."
```
Explanation of the updated arguments:

The first argument is either the input file or the text to be processed.
-o or --output: Specifies the output file (optional, only for file inputs).
-c or --config: Specifies a custom configuration file (optional).
--anthropic: Overrides the default AI to use Anthropic instead of OpenAI.

## Proof File Example
Here's an updated version of the JSON file to reflect your requirements:

```json
{
  "proofing_prompts": [
    {
      "label": "basic-proof",
      "description": "Review the document for spelling and grammar, using British English with an active voice and academic writing style.",
      "criteria": [
        "spelling",
        "grammar",
        "British English",
        "active voice",
        "academic style"
      ],
      "apply_to_all": [],
      "include_file": false,
      "file_path": "",
      "request_additional_info": false
    },
    {
      "label": "technical-proof",
      "description": "Check for technical accuracy and consistency in terminology specific to the domain.",
      "criteria": [
        "technical accuracy",
        "terminology consistency"
      ],
      "apply_to_all": ["basic-proof"],
      "include_file": true,
      "file_path": "path\\to\\technical-file.md",
      "request_additional_info": true
    },
    {
      "label": "clarity-review",
      "description": "Ensure the text is clear and easy to understand for a general audience.",
      "criteria": [
        "clarity",
        "readability"
      ],
      "apply_to_all": ["basic-proof"],
      "include_file": false,
      "file_path": "",
      "request_additional_info": false
    },
    {
      "label": "style-check",
      "description": "Verify that the writing style adheres to the provided guidelines or style manual.",
      "criteria": [
        "style guide adherence"
      ],
      "apply_to_all": ["basic-proof"],
      "include_file": false,
      "file_path": "",
      "request_additional_info": false
    }
  ]
}
```

### Explanation:
- **proofing_prompts**: An array of objects, each representing a different type of proofing prompt.
  - **label**: A short identifier for the prompt.
  - **description**: A detailed description of what the prompt entails.
  - **criteria**: A list of specific criteria that the proofing should focus on.
  - **apply_to_all**: A set of labels referring to other prompts to be run in conjunction with this one.
  - **include_file**: A boolean to indicate whether to include an additional file for this specific prompt.
  - **file_path**: The path to the additional file to be included if `include_file` is true.
  - **request_additional_info**: A boolean to specify if additional information should be requested for this specific prompt.

This structure allows for greater flexibility by specifying optional file inclusion and user interaction for each proofing prompt and defining dependencies between prompts through the `apply_to_all` field.
