package proofer

import (
	"fmt"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
)

func ProofText(content string, promptText string, appConfig *config.AppConfig) (string, error) {
	switch appConfig.AIProvider {
	case "openai":
		return ProofTextOpenAI(content, promptText, appConfig.APIKey)
	case "anthropic":
		return ProofTextAnthropic(content, promptText, appConfig.APIKey)
	default:
		return "", fmt.Errorf("unsupported AI provider: %s", appConfig.AIProvider)
	}
}

func ReviewText(content string, promptText string, appConfig *config.AppConfig) (string, error) {
	// For now, we'll use the same logic as ProofText
	return ProofText(content, promptText, appConfig)
}
