package proofing

import (
	"fmt"
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
	if s.appConfig.ProofingType == "command_line" {
		// For command line input, content is already in InputFile
		return s.appConfig.InputFile, nil
	}

	// Read from the input file
	content, err := os.ReadFile(s.appConfig.InputFile)
	if err != nil {
		return "", fmt.Errorf("error reading input file: %v", err)
	}
	return string(content), nil
}

func (s *StandardProofingStages) ExecuteProofing(content string, promptText string) (string, error) {
	return proofer.ProofText(content, promptText, s.appConfig)
}

func (s *StandardProofingStages) ExecuteReview(content string, promptText string) (string, error) {
	return proofer.ReviewText(content, promptText, s.appConfig)
}

func (s *StandardProofingStages) HandleOutput(proofedContent string) error {
	if proofedContent == "" {
		return nil // No changes to write
	}

	if s.appConfig.OutputFile != "" {
		return os.WriteFile(s.appConfig.OutputFile, []byte(proofedContent), 0644)
	}

	// If no output file is specified, print to console
	fmt.Println(proofedContent)
	return nil
}

func (s *StandardProofingStages) Finalize() error {
	// No finalization needed for standard proofing
	return nil
}
