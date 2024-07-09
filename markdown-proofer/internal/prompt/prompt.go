package prompt

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/phaynes/markdown-proofing-tool/internal/config"
	"github.com/sashabaranov/go-openai"
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

	// Prelude
	builder.WriteString("You are an AI assistant specialised in proofreading and improving  text. ")
	builder.WriteString("Your task is to review and enhance the given text based on the following criteria and instructions. ")
	// builder.WriteString("Please return the improved text in markdown format.\n\n")

	// Build the main prompt
	builder.WriteString("Primary task: ")
	builder.WriteString(prompt.Description)
	builder.WriteString("\nPrimary criteria: ")
	builder.WriteString(strings.Join(prompt.Criteria, ", "))
	builder.WriteString("\n")

	// Apply additional prompts if specified
	for _, label := range prompt.ApplyToAll {
		for _, additionalPrompt := range allPrompts {
			if additionalPrompt.Label == label {
				builder.WriteString("\nAdditional task: ")
				builder.WriteString(additionalPrompt.Description)
				builder.WriteString("\nAdditional criteria: ")
				builder.WriteString(strings.Join(additionalPrompt.Criteria, ", "))
				builder.WriteString("\n")
			}
		}
	}

	// Include file content if specified
	if prompt.IncludeFile && prompt.FilePath != "" {
		data, err := ioutil.ReadFile(prompt.FilePath)
		if err == nil {
			builder.WriteString("\nAdditional context:\n")
			builder.Write(data)
			builder.WriteString("\n")
		}
	}

	// Request for additional information if specified
	if prompt.RequestAdditionalInfo {
		builder.WriteString("\nIf you need any additional information or clarification to complete this task effectively, please state your questions clearly.")
	}

	builder.WriteString("\nPlease review and improve the provided text based on these instructions.")

	return builder.String()
}

func getProofingPrompt(prompts []ProofingPrompt, proofType string) (string, error) {
	for _, prompt := range prompts {
		if prompt.Label == proofType {
			return buildPrompt(prompt, prompts), nil
		}
	}
	return "", fmt.Errorf("proofing type not found: %s", proofType)
}

func proofText(input string, prompt string, apiKey string) (string, error) {
	// fmt.Println("The prompt is:")
	// fmt.Println(prompt)

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Please review and improve the following  text:\n\n" + input,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}
