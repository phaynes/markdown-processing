package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PrepareRepository prepares the git repository for proofing
func PrepareRepository(inputFile string) error {
	if err := checkGitRepository(); err != nil {
		return fmt.Errorf("git repository check failed: %v", err)
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

// GetContentToProof retrieves the content to be proofed
func GetContentToProof(inputFile string, isGitDiff bool) (string, error) {
	if isGitDiff {
		return getGitDiff(inputFile)
	}
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	return string(content), nil
}

// WriteProofedContent writes the proofed content back to the file
func WriteProofedContent(inputFile string, proofedContent string) error {
	return os.WriteFile(inputFile, []byte(proofedContent), 0644)
}

// CompleteWorkflow completes the git workflow after proofing
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not in a git repository: %v", err)
	}
	return nil
}

func addAndCommit(filepath string, message string) error {
	// Check if the file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filepath)
	}

	// Git add
	addCmd := exec.Command("git", "add", filepath)
	addOutput, err := addCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git add failed: %v\nOutput: %s", err, string(addOutput))
	}

	// Check if there are changes to commit
	statusCmd := exec.Command("git", "status", "--porcelain", filepath)
	statusOutput, err := statusCmd.Output()
	if err != nil {
		return fmt.Errorf("git status failed: %v", err)
	}

	if len(statusOutput) == 0 {
		// No changes to commit
		return nil
	}

	// Git commit
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitOutput, err := commitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit failed: %v\nOutput: %s", err, string(commitOutput))
	}

	return nil
}

func createAndCheckoutBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create and checkout branch: %v\nOutput: %s", err, string(output))
	}
	return nil
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

func getGitDiff(inputFile string) (string, error) {
	cmd := exec.Command("git", "diff", "HEAD", inputFile)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git diff: %v\nOutput: %s", err, string(output))
	}
	return string(output), nil
}

// GetCurrentBranch gets the current Git branch
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting current branch: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
