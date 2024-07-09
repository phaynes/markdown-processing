package proofer

import (
	"fmt"

	"github.com/yourusername/markdown-proofing-tool/internal/config"
)

func ProofText(input string, prompt string, appConfig *config.AppConfig) (string, error) {
	switch appConfig.AIProvider {
	case "openai":
		return ProofTextOpenAI(input, prompt, appConfig.OpenAIKey)
	case "anthropic":
		return ProofTextAnthropic(input, prompt, appConfig.AnthropicKey)
	default:
		return "", fmt.Errorf("Invalid AI provider: %s", appConfig.AIProvider)
	}
}
