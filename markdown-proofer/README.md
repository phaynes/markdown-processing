# The Markdown Proofer (MDP) User Guide

`mdp` is a command-line tool for proofing and reviewing markdown files using AI API queries. It can work with text from the command line, individual files, file managed using git, as wells as specific line ranges within files. The system supports a range of customisable proofing modes, achieved by configuring a database of AI queries.

This guide covers all the features and options available in `mdp`.

## Table of Contents

1. [Installation](#installation)
2. [Basic Usage](#basic-usage)
3. [Proofing Modes](#proofing-modes)
4. [Command-Line Options](#command-line-options)
5. [Configuration Files](#configuration-files)
6. [Git Integration](#git-integration)
7. [AI Providers](#ai-providers)
8. [Zed Integration](#zed-integration)
9. [Output Handling](#output-handling)
10. [Examples](#examples)
11. [Troubleshooting](#proofing-prompts)
12. [Troubleshooting](#troubleshooting)

## 1. Installation

`mdp` has been developed on a m-series Mac, although it is likely to work with a broader set of tools. Additionally, the system assumes both Git and Go are installed on the machine.

To install `mdp`, clone the repository and build the binary:

```bash
git clone https://github.com/phaynes/markdown-processing.git
cd markdown/markdown-proofer
./build.sh
```

Move the `mdp` binary to a directory in your PATH to use it from anywhere.

The system assumes you have access to [OpenAI](https://openai.com/index/openai-api/) and/or [Anthropic](https://www.anthropic.com/api) API keys.

## 2. Basic Usage

The basic syntax for using `mdp` is:

```bash
mdp [options] [input_file]
```

If no input file is specified, `mdp` will use the first file listed in the `config.json` file.

## 3. Proofing Modes

`mdp` supports several proofing modes:

- **Command Line Proofing**: Proofs text passed to it on the command line.
- **Standard Proofing**: Proofs the entire file.
- **Git Full Proofing**: Proofs the entire file with Git integration.
- **Git Diff Proofing**: Proofs only the changes detected by Git.
- **Line Range Proofing**: Proofs a specific range of lines in the file.

## 4. Command-Line Options

`mdp` supports the following command-line options:

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

`mdp` uses two configuration files:

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

`mdp` can integrate with Git repositories:

- Use `--use-git` to enable Git features.
- Use `--proof-git-diff` along with `--use-git` to proof only the changes detected by Git.

When using Git integration, MDP will create a new branch, make changes, and merge back to the original branch.

## 7. AI Providers

MDP supports two AI providers:

- OpenAI (default)
- Anthropic

Specify the provider using the `-ai` flag or in the `ai_config.json` file.

## 8. Zed Integration

[Zed](https://zed.dev/) is a high-performance development environment with integrated support for markdown as well as mechanisms for task automation. Below is an example integration that enables `mdp` to be used directly from the development environment.

### Task Integration

Below is an example definition for specifying `mdp` tasks within Zed:

```json
[
  {
    "label": "proof",
    "command": "$HOME/bin/mdp -ai-config $HOME/.mdp/ai_config.json -input $ZED_FILE --use-git --proof-git-diff",
    "use_new_terminal": false,
    "allow_concurrent_runs": false,
    "reveal": "never"
  },
  {
    "label": "lproof",
    "command": "$HOME/bin/mdp -ai-config $HOME/.mdp/ai_config.json -input $ZED_FILE -n $ZED_ROW",
    "use_new_terminal": false,
    "allow_concurrent_runs": false,
    "reveal": "never"
  },
  {
    "label": "anthropic_proof",
    "command": "$HOME/bin/mdp -ai-config $HOME/.mdp/ai_config.json -ai anthropic -input $ZED_FILE -n $ZED_ROW",
    "use_new_terminal": false,
    "allow_concurrent_runs": false,
    "reveal": "never"
  },
  {
    "label": "quick_commit",
    "command": "$HOME/bin/quick_commit $ZED_FILE",
    "use_new_terminal": false,
    "allow_concurrent_runs": false,
    "reveal": "never"
  }
]
```

### Keystroke Bindings

The following are example Zed keystroke bindings for the `mdp` functions:

```json
[
  {
    "context": "Workspace",
    "bindings": {
      "alt-p": ["task::Spawn", { "proof": "Proof the document" }],
      "alt-l": ["task::Spawn", { "lproof": "Proof the current line" }],
      "alt-a": [
        "task::Spawn",
        { "anthropic_proof": "Proof the current line using Claude.ai" }
      ],
      "alt-q": ["task::Spawn", { "quick_commit": "Commit the current file." }]
    }
  }
]
```

## 9. Output Handling

- If an output file is specified with `-output`, the proofed content will be written to that file.
- If no output file is specified, the proofed content will either be printed to the console, or for git or line range proofing, update the file itself.

## 10. Examples

1. Proofing text from the command line:
   ```bash
   mdp "Please proof ths text"
   ```
2. Standard proofing:
   ```bash
   mdp -input document.md
   ```

3. Git full proofing:
   ```bash
   mdp --use-git -input document.md
   ```

4. Git diff proofing:
   ```bash
   mdp --use-git --proof-git-diff -input document.md
   ```

5. Line range proofing:
   ```bash
   mdp -input document.md -n 10-20
   ```

6. Using a specific AI provider:
   ```bash
   mdp -input document.md -ai anthropic
   ```

## 11. Proofing Prompts

## 12. Troubleshooting

- Ensure your API keys are correctly set in the `ai_config.json` file.
- When using Git features, make sure you're in a Git repository.
- If you encounter any "no changes to proof" messages, ensure there are actually changes in your file or Git diff.
- For line range proofing, make sure the specified range is valid for your input file.

For more help or to report issues, please visit the markdown repository on GitHub.
