package proofing

import (
	"fmt"
	"os"
	"strings"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/git"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/proofer"
)

// NoChangesError is a custom error type for when there are no changes to proof
type NoChangesError struct {
	FileName string
}

func (e *NoChangesError) Error() string {
	return fmt.Sprintf("The file '%s' is up to date. No changes need to be proofed.", e.FileName)
}

type GitProofingStages struct {
	appConfig   *config.AppConfig
	changes     []git.DiffChange
	fullContent string
}

func NewGitProofingStages(appConfig *config.AppConfig) *GitProofingStages {
	return &GitProofingStages{appConfig: appConfig}
}

func (g *GitProofingStages) Initialize() error {
	return git.PrepareRepository(g.appConfig.InputFile)
}

func (g *GitProofingStages) PrepareContent() (string, error) {
	if g.appConfig.AIConfig.ProofGitDiff {
		return g.prepareGitDiffContent()
	}
	return g.prepareFullFileContent()
}

func (g *GitProofingStages) prepareGitDiffContent() (string, error) {
	changes, err := git.GetGitDiff(g.appConfig.InputFile)
	if err != nil {
		return "", fmt.Errorf("error getting git diff: %v", err)
	}

	if changes == nil || len(changes) == 0 {
		return "", &NoChangesError{FileName: g.appConfig.InputFile}
	}

	g.changes = changes

	var contentToProof strings.Builder
	for _, change := range g.changes {
		contentToProof.WriteString(change.Content)
	}

	return contentToProof.String(), nil
}

func (g *GitProofingStages) prepareFullFileContent() (string, error) {
	content, err := os.ReadFile(g.appConfig.InputFile)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	g.fullContent = string(content)
	return g.fullContent, nil
}

func (g *GitProofingStages) ExecuteProofing(content string, promptText string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("no content to proof")
	}
	return proofer.ProofText(content, promptText, g.appConfig)
}

func (g *GitProofingStages) ExecuteReview(content string, promptText string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("no content to review")
	}
	return proofer.ReviewText(content, promptText, g.appConfig)
}

func (g *GitProofingStages) HandleOutput(proofedContent string) error {
	if proofedContent == "" {
		return nil // No changes to write
	}

	if g.appConfig.AIConfig.ProofGitDiff {
		return g.handleGitDiffOutput(proofedContent)
	}
	return g.handleFullFileOutput(proofedContent)
}

func (g *GitProofingStages) handleGitDiffOutput(proofedContent string) error {
	proofedLines := strings.Split(proofedContent, "\n")
	var proofedIndex int

	for i, change := range g.changes {
		changeLines := strings.Count(change.Content, "\n") + 1
		endIndex := proofedIndex + changeLines

		if endIndex > len(proofedLines) {
			endIndex = len(proofedLines)
		}

		if proofedIndex >= len(proofedLines) {
			// We've run out of proofed lines, so keep the original content
			continue
		}

		g.changes[i].Content = strings.Join(proofedLines[proofedIndex:endIndex], "\n")
		proofedIndex = endIndex

		if proofedIndex >= len(proofedLines) {
			break
		}
	}

	return git.ApplyProofedChanges(g.appConfig.InputFile, g.changes)
}

func (g *GitProofingStages) handleFullFileOutput(proofedContent string) error {
	return os.WriteFile(g.appConfig.InputFile, []byte(proofedContent), 0644)
}

func (g *GitProofingStages) Finalize() error {
	return git.CompleteWorkflow(g.appConfig.InputFile)
}
