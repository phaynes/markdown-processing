package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type DiffChange struct {
	StartLine int
	Content   string
}

// GetGitDiff retrieves the diff between the current version of the file and the last committed version
func GetGitDiff(inputFile string) ([]DiffChange, error) {
	// First, check if there are any changes
	checkCmd := exec.Command("git", "diff", "--exit-code", inputFile)
	err := checkCmd.Run()
	if err == nil {
		// No changes
		return []DiffChange{}, nil
	}

	// If there are changes, get the diff
	cmd := exec.Command("git", "diff", "--unified=0", inputFile)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting git diff: %v", err)
	}

	changes, err := parseDiff(string(output))
	if err != nil {
		return nil, fmt.Errorf("error parsing git diff: %v", err)
	}

	if len(changes) == 0 {
		return nil, nil
	}

	return changes, nil
}
func PrepareRepository(inputFile string) error {
	if err := checkGitRepository(); err != nil {
		return fmt.Errorf("not in a git repository: %v", err)
	}

	if err := addAndCommit(inputFile, "Commit before proofing"); err != nil {
		return fmt.Errorf("error in git add and commit: %v", err)
	}

	branchName := "ai-proof-branch"
	if err := createAndCheckoutBranch(branchName); err != nil {
		return fmt.Errorf("error creating and checking out branch: %v", err)
	}

	return nil
}

func CompleteWorkflow(inputFile string) error {
	if err := addAndCommit(inputFile, "Proofing completed"); err != nil {
		return fmt.Errorf("error in git add and commit after proofing: %v", err)
	}

	if err := checkoutBranch("main"); err != nil {
		return fmt.Errorf("error checking out main branch: %v", err)
	}

	if err := mergeBranch("ai-proof-branch"); err != nil {
		return fmt.Errorf("error merging branch: %v", err)
	}

	if err := deleteBranch("ai-proof-branch"); err != nil {
		return fmt.Errorf("error deleting branch: %v", err)
	}

	return nil
}

func checkGitRepository() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	return cmd.Run()
}

func addAndCommit(filepath string, message string) error {
	addCmd := exec.Command("git", "add", filepath)
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %v", err)
	}

	commitCmd := exec.Command("git", "commit", "-m", message)
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %v", err)
	}

	return nil
}

func createAndCheckoutBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	return cmd.Run()
}

func checkoutBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	return cmd.Run()
}

func mergeBranch(branchName string) error {
	cmd := exec.Command("git", "merge", branchName)
	return cmd.Run()
}

func deleteBranch(branchName string) error {
	cmd := exec.Command("git", "branch", "-d", branchName)
	return cmd.Run()
}

func parseDiff(diff string) ([]DiffChange, error) {
	var changes []DiffChange
	lines := strings.Split(diff, "\n")
	var currentStartLine int
	var currentContent strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "@@ ") {
			// If we have accumulated content, add it to changes
			if currentContent.Len() > 0 {
				changes = append(changes, DiffChange{StartLine: currentStartLine, Content: currentContent.String()})
				currentContent.Reset()
			}

			// Parse the new hunk header
			parts := strings.Split(line, " ")
			if len(parts) < 3 {
				continue
			}
			lineInfo := strings.TrimPrefix(parts[2], "+")
			startLine, err := parseLineNumber(lineInfo)
			if err != nil {
				return nil, fmt.Errorf("error parsing line number from '%s': %v", lineInfo, err)
			}
			currentStartLine = startLine
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			// Accumulate added lines
			currentContent.WriteString(strings.TrimPrefix(line, "+"))
			currentContent.WriteString("\n")
		}
	}

	// Add any remaining content
	if currentContent.Len() > 0 {
		changes = append(changes, DiffChange{StartLine: currentStartLine, Content: currentContent.String()})
	}

	return changes, nil
}

func parseLineNumber(lineInfo string) (int, error) {
	parts := strings.Split(lineInfo, ",")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid line info format")
	}

	lineNumber, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid line number: %v", err)
	}

	return lineNumber, nil
}

// ApplyProofedChanges applies the proofed changes back to the file
func ApplyProofedChanges(inputFile string, changes []DiffChange) error {
	// Read the entire file
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Apply changes
	for _, change := range changes {
		startLine := change.StartLine - 1 // Convert to 0-based index
		if startLine < 0 {
			startLine = 0
		}
		changeLines := strings.Split(change.Content, "\n")
		endLine := startLine + len(changeLines)

		if endLine > len(lines) {
			endLine = len(lines)
		}

		// Replace lines
		if startLine < len(lines) {
			lines = append(lines[:startLine], append(changeLines, lines[endLine:]...)...)
		} else {
			// If the change is beyond the current file content, append it
			lines = append(lines, changeLines...)
		}
	}

	// Write the updated content back to the file
	output := strings.Join(lines, "\n")
	return os.WriteFile(inputFile, []byte(output), 0644)
}

// GetFullFileContent retrieves the full content of the file
func GetFullFileContent(inputFile string) (string, error) {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	return string(content), nil
}

// WriteFullFileContent writes the full content back to the file
func WriteFullFileContent(inputFile string, content string) error {
	return os.WriteFile(inputFile, []byte(content), 0644)
}
