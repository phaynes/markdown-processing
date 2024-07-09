package proofer

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropic-ai/anthropic-sdk-go"
)

func ProofTextAnthropic(input string, prompt string, apiKey string) (string, error) {
	client := anthropic.NewClient(apiKey)
	resp, err := client.CompleteStream(context.Background(), &anthropic.CompletionRequest{
		Prompt:            prompt + "\n\nHuman: Please review and improve the following text:\n\n" + input + "\n\nAssistant:",
		Model:             anthropic.Claude3Opus,
		MaxTokensToSample: 1000,
		StopSequences:     []string{"Human:"},
	})
	if err != nil {
		return "", fmt.Errorf("Anthropic Completion error: %v", err)
	}
	defer resp.Close()

	var fullResponse strings.Builder
	for {
		completion, err := resp.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return "", fmt.Errorf("Error receiving completion: %v", err)
		}
		fullResponse.WriteString(completion.Completion)
	}

	return fullResponse.String(), nil
}
