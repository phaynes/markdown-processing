package prompt

import (
	"fmt"
	"strings"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
)

func BuildProofingPrompt(prompts []config.ProofingPrompt, proofType string) (string, error) {
	for _, prompt := range prompts {
		if prompt.Label == proofType {
			return buildPrompt(prompt, prompts), nil
		}
	}
	return "", fmt.Errorf("proofing type not found: %s", proofType)
}

func buildPrompt(prompt config.ProofingPrompt, allPrompts []config.ProofingPrompt) string {
	var builder strings.Builder

	builder.WriteString("You are an AI assistant specialized in proofreading and improving text. ")
	builder.WriteString("Your task is to review and enhance the given text based on the following criteria and instructions:\n\n")

	builder.WriteString("Primary task: ")
	builder.WriteString(prompt.Description)
	builder.WriteString("\nPrimary criteria: ")
	builder.WriteString(strings.Join(prompt.Criteria, ", "))
	builder.WriteString("\n\n")

	for _, label := range prompt.ApplyToAll {
		for _, additionalPrompt := range allPrompts {
			if additionalPrompt.Label == label {
				builder.WriteString("Additional task: ")
				builder.WriteString(additionalPrompt.Description)
				builder.WriteString("\nAdditional criteria: ")
				builder.WriteString(strings.Join(additionalPrompt.Criteria, ", "))
				builder.WriteString("\n\n")
			}
		}
	}

	if prompt.RequestAdditionalInfo {
		builder.WriteString("If you need any additional information or clarification to complete this task effectively, please state your questions clearly.\n\n")
	}

	builder.WriteString("Please review and improve the provided text based on these instructions.")

	return builder.String()
}
