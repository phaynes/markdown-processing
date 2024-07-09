# AI-Assisted Markdown Proofing Tool Development

Markdown provides an efficient method for preparing text as it can be transformed into various formats (PDF / HTML / DOC) and is able to be strongly version controlled using tools such as Git.

However, written work needs to be proofed and reviewed, and the types of tooling support available in word processors are not directly accessible. Instead, proofing the markdown text generally involves copying text segments into an AI chat engine, receiving corrections, and copying them back into the document. This approach is labour-intensive and typically involves respecifying common instructions such as "please proof this text using British English spelling with an active voice."

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
   - Example proofing prompts file:

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

### Proofing
```bash
mdp input.md
mdp input.md -o output.md
mdp -c custom_config.json input.md
mdp --anthropic input.md
mdp "This is some text to proof directly from the command line."
```

### Reviewing
```bash
mdr input.md
mdr input.md -o review.md
mdr -c custom_config.json input.md
mdr --anthropic input.md
mdr "Please review this text for clarity and style."
```

### Commenting
```bash
mdc input.md
mdc input.md -o comments.md
mdc -c custom_config.json input.md
mdc --anthropic input.md
mdc "Add comments to improve this paragraph's structure."
```

**Explanation of the Updated Arguments**:
- The first argument is either the input file or the text to be processed.
- `-o` or `--output`: Specifies the output file (optional, only for file inputs).
- `-c` or `--config`: Specifies a custom configuration file (optional).
- `--anthropic`: Overrides the default AI to use Anthropic instead of OpenAI.

## Expected Outputs

### Simple Proofing Output
**Command**:
```bash
mdp --prompt basic-proof
```
**Expected Output**:
```plaintext
Using configuration file: config.json
Proofing file: 2024-05-11-assignment2.md
Output file: assignment2-proofed.md
```

### Proofing with Version Control
**Command**:
```bash
mdp --prompt basic-proof --use-git
```
**Expected Output**:
```plaintext
Using configuration file: config.json
Proofing file: 2024-05-11-assignment2.md
Creating Git branch: ai-proof-branch
Running AI proofing engine...
Proofed file: 2024-05-11-assignment2.md
Output file: assignment2-proofed.md
Committing changes to ai-proof-branch…
Merging changes into main branch…
Deleting branch ai-proof-branch…
```

### Reviewing with Version Control
**Command**:
```bash
mdr --prompt clarity-review --use-git
```
**Expected Output**:
```plaintext
Using configuration file: config.json
Reviewing file: 2024-05-11-assignment2.md
Creating Git branch: ai-review-branch
Running AI review engine...
Review file: 2024-05-11-assignment2-review.md
Proposed changes prepared.
Committing changes to ai-review-branch...
Please review the proposed changes in 2024-05-11-assignment2-review.md
```

### Commenting with Version Control
**Command**:
```bash
mdc --prompt style-check --use-git
```
**Expected Output**:
```plaintext
Using configuration file: config.json
Commenting on file: 2024-05-11-assignment2.md
Creating Git branch: ai-comment-branch
Running AI commenting engine...
Comments file: 2024-05-11-assignment2-comments.md
Committing comments to ai-comment-branch...
Please review the comments in 2024-05-11-assignment2-comments.md
```

### Error Handling and Logging
- **Logs**: All steps and any issues should be logged, with logs stored in a specified location.
- **Error Messages**: Clear and actionable error messages for issues such as missing files or invalid API keys.

### Security Considerations
- **API Key Management**: API keys should be securely stored and managed.
- **Data Security**: Measures to ensure the security and privacy of the data being processed.

### Performance Considerations
- **Efficiency**: The tool should handle large markdown files or high volumes of text efficiently.
- **Scalability**: Ability to scale the tool for larger projects or multiple concurrent users.

### User Documentation and Help
- **Documentation**: Comprehensive user documentation should be provided.
- **Help Command**: A built-in help command (`mdp --help`) should provide detailed usage instructions and examples.

### Customization and Extensibility
- **Adding Prompts**: Users should be able to add new proofing prompts or extend the tool’s functionality by modifying the `proofing_prompts.json` file.

### Integration Tests
- **Test Cases**: Include unit tests and integration tests to ensure the tool's reliability.
- **Example Tests**: Provide examples of test cases to cover different scenarios.
## Additional Q& A
After reviewing the comprehensive specification for the AI-Assisted Markdown Proofing Tool, I believe it's well-structured and covers most aspects of the project in detail. However, there are a few areas where additional clarity or consideration could be beneficial:

1. Error Handling and Validation:
   Q. While error handling is mentioned briefly, it might be helpful to specify more detailed requirements for input validation, API error handling, and how the tool should behave in case of network issues or AI service downtime.

   A. TBD.

2. Rate Limiting and Quota Management:
   Q. Given that the tool uses external AI services, it would be prudent to include guidelines on how to handle rate limiting and API quota management to prevent excessive costs or service disruptions.

   A. API quota management is out of scope.

3. Markdown Parsing:
   Q. The specification doesn't mention how the tool will parse markdown. It might be worth specifying whether a particular markdown parser will be used and how it will handle different markdown flavors or extensions.

   A. TBD during technical work.

4. Diff Format:
   q. For the review functionality (mdr), it might be helpful to specify the exact format of the diff output. Will it use a standard diff format, or a custom format designed for easier user review?

   a. TBD.

5. Internationalization:
   This is to be handled by the AI tool, and is out of scope

6. Continuous Integration/Continuous Deployment (CI/CD):
   Comment: While Docker is mentioned for the build environment, more specific CI/CD requirements or expectations could be beneficial, especially given the Git integration.

   Yes: In the first instance, just a working version is required.

7. Performance Metrics:
The predominant slow step will be access to API's. The choice of model to accelerate different functions will be a subsequent update.

8. Version Compatibility:
   A. The versions of Go, Git, and the AI APIs the tool should be compatible with could prevent future compatibility issues.

   At this stage, the latest versions of software will be used.

9. User Interface for Review:
   The UI for review will be provided through an external tool such as github desktop or gitea.

10. Backup and Recovery:
    Any change to a file must be through git.
