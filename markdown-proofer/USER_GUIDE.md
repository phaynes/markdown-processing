# MDP (Markdown Proofer) User Guide

MDP is a command-line tool for proofing and reviewing markdown files. It can work with individual files, Git repositories, and even specific line ranges within files. This guide covers all the features and options available in MDP.

## Table of Contents

1. [Installation](#installation)
2. [Basic Usage](#basic-usage)
3. [Proofing Modes](#proofing-modes)
4. [Command-Line Options](#command-line-options)
5. [Configuration Files](#configuration-files)
6. [Git Integration](#git-integration)
7. [AI Providers](#ai-providers)
8. [Output Handling](#output-handling)
9. [Examples](#examples)
10. [Troubleshooting](#troubleshooting)

## 1. Installation

To install MDP, clone the repository and build the binary:

```bash
git clone https://github.com/your-repo/mdp.git
cd mdp
go build -o mdp ./cmd/mdp
```

Move the `mdp` binary to a directory in your PATH to use it from anywhere.

## 2. Basic Usage

The basic syntax for using MDP is:

```bash
mdp [options] [input_file]
```

If no input file is specified, MDP will use the first file listed in the `config.json` file.

## 3. Proofing Modes

MDP supports several proofing modes:

- **Standard Proofing**: Proofs the entire file.
- **Git Full Proofing**: Proofs the entire file with Git integration.
- **Git Diff Proofing**: Proofs only the changes detected by Git.
- **Line Range Proofing**: Proofs a specific range of lines in the file.

## 4. Command-Line Options

MDP supports the following command-line options:

- `-config`: Path to the configuration file (default: `config.json`)
- `-ai-config`: Path to the AI configuration file (default: `ai_config.json`)
- `-type`: Type of proofing to perform (default: `basic-proof`)
- `-ai`: AI provider to use (`openai` or `anthropic`)
- `--use-git`: Use Git features
- `--proof-git-diff`: Proof only Git diff (requires `--use-git`)
- `-mode`: Mode of operation (`proof` or `review`)
- `-input`: Input file to proof
- `-output`: Output file for proofed content
- `-n`: Line range to proof (e.g., '6-10' or '6')

## 5. Configuration Files

MDP uses two configuration files:

### config.json

Contains general configuration:

```json
{
  "input_files": ["default_input.md"],
  "bibliography": "references.bib"
}
```

### ai_config.json

Contains AI-specific configuration:

```json
{
  "default_ai": "openai",
  "openAI_key": "path/to/openai_key.txt",
  "anthropic_key": "path/to/anthropic_key.txt",
  "proofing_prompts": "path/to/proofing_prompts.json",
  "use_git": false,
  "proof_git_diff": false
}
```

## 6. Git Integration

mdp can integrate with Git repositories:

- Use `--use-git` to enable Git features.
- Use `--proof-git-diff` along with `--use-git` to proof only the changes detected by Git.

When using Git integration, MDP will create a new branch, make changes, and merge back to the original branch.

## 7. AI Providers

MDP supports two AI providers:

- OpenAI (default)
- Anthropic

Specify the provider using the `-ai` flag or in the `ai_config.json` file.

## 8. Output Handling

- If an output file is specified with `-output`, the proofed content will be written to that file.
- If no output file is specified, the proofed content will be printed to the console.

## 9. Examples

1. Standard proofing:
   ```bash
   mdp -input document.md
   ```

2. Git full proofing:
   ```bash
   mdp --use-git -input document.md
   ```

3. Git diff proofing:
   ```bash
   mdp --use-git --proof-git-diff -input document.md
   ```

4. Line range proofing:
   ```bash
   mdp -input document.md -n 10-20
   ```

5. Using a specific AI provider:
   ```bash
   mdp -input document.md -ai anthropic
   ```

## 10. Troubleshooting

- Ensure your API keys are correctly set in the `ai_config.json` file.
- When using Git features, make sure you're in a Git repository.
- If you encounter any "no changes to proof" messages, ensure there are actually changes in your file or Git diff.
- For line range proofing, make sure the specified range is valid for your input file.

For more help or to report issues, please visit the markdown repository on GitHub.
