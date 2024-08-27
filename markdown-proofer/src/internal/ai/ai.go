package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/pdf"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/proofer"
)

func EnhanceMetadata(text string, metadata pdf.Metadata, appConfig *config.AppConfig) (pdf.Metadata, error) {
	prompt := "Extract or enhance the following metadata from the text: title, authors (as a comma-separated list), and year. If any of this information is missing or seems incorrect in the provided metadata, please provide the correct information. Format the response as JSON."

	response, err := executeAIRequest(text, prompt, appConfig)
	if err != nil {
		return metadata, err
	}

	enhancedMetadata := parseAIResponse(response, metadata)
	return enhancedMetadata, nil
}

func SummarizeText(content string, criteria string, appConfig *config.AppConfig) (string, error) {
	prompt := fmt.Sprintf("Summarize the following text based on these criteria: %s", criteria)

	return executeAIRequest(content, prompt, appConfig)
}

func executeAIRequest(content string, prompt string, appConfig *config.AppConfig) (string, error) {
	// Combine the prompt and content
	fullPrompt := fmt.Sprintf("%s\n\nText to process:\n%s", prompt, content)

	// Use the existing proofer.ProofText function
	response, err := proofer.ProofText(fullPrompt, prompt, appConfig)
	if err != nil {
		return "", fmt.Errorf("error executing AI request: %v", err)
	}

	return response, nil
}

func parseAIResponse(response string, originalMetadata pdf.Metadata) pdf.Metadata {

	var enhancedMetadata pdf.Metadata
	cleanedResponse := cleanJSONResponse(response)

	err := json.Unmarshal([]byte(cleanedResponse), &enhancedMetadata)
	if err != nil {
		// If there's an error parsing the JSON, fall back to the original metadata
		fmt.Printf("Error parsing JSON: %v\n", err)
		return originalMetadata
	}

	// Convert authors string to slice if necessary
	if enhancedMetadata.Authors == "" {
		enhancedMetadata.Authors = originalMetadata.Authors
	}

	// If any field in the enhanced metadata is empty, use the original value
	if enhancedMetadata.Title == "" {
		enhancedMetadata.Title = originalMetadata.Title
	}
	if len(enhancedMetadata.Authors) == 0 {
		enhancedMetadata.Authors = originalMetadata.Authors
	}
	if enhancedMetadata.Year == 0 {
		enhancedMetadata.Year = originalMetadata.Year
	}

	return enhancedMetadata
}

// Helper function to clean up JSON response string
func cleanJSONResponse(response string) string {
	// Remove the starting ```json and ending ```
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)
	return response
}
