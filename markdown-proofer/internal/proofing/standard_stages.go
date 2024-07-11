package proofing

import (
	"os"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/proofer"
)

type StandardProofingStages struct {
	appConfig *config.AppConfig
}

func NewStandardProofingStages(appConfig *config.AppConfig) *StandardProofingStages {
	return &StandardProofingStages{appConfig: appConfig}
}

func (s *StandardProofingStages) Initialize() error {
	// No initialization needed for standard proofing
	return nil
}

func (s *StandardProofingStages) PrepareContent() (string, error) {
	content, err := os.ReadFile(s.appConfig.InputFile)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s *StandardProofingStages) ExecuteProofing(content string, promptText string) (string, error) {
	return proofer.ProofText(content, promptText, s.appConfig)
}

func (s *StandardProofingStages) ExecuteReview(content string, promptText string) (string, error) {
	// Implement review logic here
	return proofer.ReviewText(content, promptText, s.appConfig)
}

func (s *StandardProofingStages) HandleOutput(proofedContent string) error {
	return os.WriteFile(s.appConfig.OutputFile, []byte(proofedContent), 0644)
}

func (s *StandardProofingStages) Finalize() error {
	// No finalization needed for standard proofing
	return nil
}
