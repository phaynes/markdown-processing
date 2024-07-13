package proofing

import (
	"fmt"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
)

type ProofingStages interface {
	Initialize() error
	PrepareContent() (string, error)
	ExecuteProofing(content string, promptText string) (string, error)
	ExecuteReview(content string, promptText string) (string, error)
	HandleOutput(result string) error
	Finalize() error
}

func NewProofingStages(appConfig *config.AppConfig) (ProofingStages, error) {
	if appConfig.LineRange != "" {
		return NewLineProofingStages(appConfig), nil
	}

	switch appConfig.ProofingType {
	case "command_line":
		return NewCommandLineProofingStages(appConfig), nil
	case "git":
		return NewGitProofingStages(appConfig), nil
	case "standard":
		return NewStandardProofingStages(appConfig), nil
	default:
		return nil, fmt.Errorf("unknown proofing type: %s", appConfig.ProofingType)
	}
}
