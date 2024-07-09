package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Config struct {
	InputFiles      []string `json:"input_files"`
	OutputFile      string   `json:"output_file"`
	Bibliography    string   `json:"bibliography"`
	DefaultAI       string   `json:"default_ai"`
	OpenAIKey       string   `json:"openAI_key"`
	AnthropicKey    string   `json:"anthropic_key"`
	ProofingPrompts string   `json:"proofing_prompts"`
}

type ProofingPrompt struct {
	Label                 string   `json:"label"`
	Description           string   `json:"description"`
	Criteria              []string `json:"criteria"`
	ApplyToAll            []string `json:"apply_to_all"`
	IncludeFile           bool     `json:"include_file"`
	FilePath              string   `json:"file_path"`
	RequestAdditionalInfo bool     `json:"request_additional_info"`
}

func main() {
	configFile := flag.String("config", "config.json", "Path to the configuration file")
	proofType := flag.String("type", "basic-proof", "Type of proofing to perform")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: mdp [options] <text>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	openAIKey, err := readAPIKey(config.OpenAIKey)
	if err != nil {
		log.Fatalf("Error reading OpenAI API key: %v", err)
	}

	proofingPrompts, err := loadProofingPrompts(config.ProofingPrompts)
	if err != nil {
		log.Fatalf("Error loading proofing prompts: %v", err)
	}

	prompt, err := getProofingPrompt(proofingPrompts, *proofType)
	if err != nil {
		log.Fatalf("Error getting proofing prompt: %v", err)
	}

	input := strings.Join(flag.Args(), " ")
	proofedText, err := proofText(input, prompt, openAIKey)
	if err != nil {
		log.Fatalf("Error proofing text: %v", err)
	}

	fmt.Println(proofedText)
}

func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func readAPIKey(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func loadProofingPrompts(filename string) ([]ProofingPrompt, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var prompts struct {
		ProofingPrompts []ProofingPrompt `json:"proofing_prompts"`
	}
	err = json.Unmarshal(data, &prompts)
	if err != nil {
		return nil, err
	}

	return prompts.ProofingPrompts, nil
}

func buildPrompt(prompt ProofingPrompt, allPrompts []ProofingPrompt) string {
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
