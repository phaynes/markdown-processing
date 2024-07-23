package proofer

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func ProofTextOpenAI(input string, prompt string, apiKey string) (string, error) {
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
                    Content: "Please review and improve the following text:\n\n" + input,
                },
            },
        },
    )

    if err != nil {
        return "", fmt.Errorf("ChatCompletion error: %v", err)
    }

    return resp.Choices[0].Message.Content, nil
}
