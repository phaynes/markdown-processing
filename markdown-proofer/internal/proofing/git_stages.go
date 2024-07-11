package proofing

import (
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/git"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/proofer"
)

type GitProofingStages struct {
	appConfig *config.AppConfig
}

func NewGitProofingStages(appConfig *config.AppConfig) *GitProofingStages {
	return &GitProofingStages{appConfig: appConfig}
}

func (g *GitProofingStages) Initialize() error {
	return git.PrepareRepository(g.appConfig.InputFile)
}

func (g *GitProofingStages) PrepareContent() (string, error) {
	return git.GetContentToProof(g.appConfig.InputFile, g.appConfig.AIConfig.ProofGitDiff)
}

func (g *GitProofingStages) ExecuteProofing(content string, promptText string) (string, error) {
	return proofer.ProofText(content, promptText, g.appConfig)
}

func (g *GitProofingStages) ExecuteReview(content string, promptText string) (string, error) {
	// Implement review logic here
	return proofer.ReviewText(content, promptText, g.appConfig)
}

func (g *GitProofingStages) HandleOutput(proofedContent string) error {
	return git.WriteProofedContent(g.appConfig.InputFile, proofedContent)
}

func (g *GitProofingStages) Finalize() error {
	return git.CompleteWorkflow(g.appConfig.InputFile)
}
