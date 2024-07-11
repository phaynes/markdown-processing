package git

import (
	"fmt"
	"os"
	"os/exec"
)

func PrepareRepository(inputFile string) error {
	if err := addAndCommit(inputFile, "Commit before proofing"); err != nil {
		return fmt.Errorf("error in git add and commit: %v", err)
	}

	branchName := "ai-proof-branch"
	if err := createAndCheckoutBranch(branchName); err != nil {
		return fmt.Errorf("error creating and checking out branch: %v", err)
	}

	return nil
}

func GetContentToProof(inputFile string, proofGitDiff bool) (string, error) {
	if proofGitDiff {
		return getGitDiff(inputFile)
	}
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	return string(content), nil
}

func WriteProofedContent(inputFile string, proofedContent string) error {
	return os.WriteFile(inputFile, []byte(proofedContent), 0644)
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

func addAndCommit(filepath string, message string) error {
	cmd := exec.Command("git", "add", filepath)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "commit", "-m", message)
	return cmd.Run()
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

func getGitDiff(inputFile string) (string, error) {
	cmd := exec.Command("git", "diff", "HEAD", inputFile)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git diff: %v", err)
	}
	return string(output), nil
}
