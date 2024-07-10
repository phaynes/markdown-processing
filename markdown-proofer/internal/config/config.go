package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	InputFiles   []string `json:"input_files"`
	OutputFile   string   `json:"output_file"`
	Bibliography string   `json:"bibliography"`
}

type AIConfig struct {
	DefaultAI       string `json:"default_ai"`
	OpenAIKey       string `json:"openAI_key"`
	AnthropicKey    string `json:"anthropic_key"`
	ProofingPrompts string `json:"proofing_prompts"`
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

type AppConfig struct {
	Config          *Config
	AIConfig        *AIConfig
	ProofingPrompts []ProofingPrompt
	OpenAIKey       string
	AnthropicKey    string
	ProofType       string
	AIProvider      string
}

func Setup() (*AppConfig, error) {
	configFile := flag.String("config", "config.json", "Path to the essay configuration file")
	aiConfigFile := flag.String("ai-config", "ai_config.json", "Path to the AI configuration file")
	proofType := flag.String("type", "basic-proof", "Type of proofing to perform")
	aiProvider := flag.String("ai", "", "AI provider to use (openai or anthropic)")
	flag.Parse()

	if flag.NArg() < 1 {
		return nil, fmt.Errorf("Usage: mdp [options] <text>")
	}

	config, err := loadConfig(*configFile)
	if err != nil {
		return nil, fmt.Errorf("Error loading essay configuration: %v", err)
	}

	aiConfig, err := loadAIConfig(*aiConfigFile)
	if err != nil {
		return nil, fmt.Errorf("Error loading AI configuration: %v", err)
	}

	proofingPrompts, err := loadProofingPrompts(aiConfig.ProofingPrompts)
	if err != nil {
		return nil, fmt.Errorf("Error loading proofing prompts: %v", err)
	}

	openAIKey, err := readAPIKey(aiConfig.OpenAIKey)
	if err != nil {
		return nil, fmt.Errorf("Error reading OpenAI API key: %v", err)
	}

	anthropicKey, err := readAPIKey(aiConfig.AnthropicKey)
	if err != nil {
		return nil, fmt.Errorf("Error reading Anthropic API key: %v", err)
	}

	if *aiProvider == "" {
		*aiProvider = aiConfig.DefaultAI
	}

	return &AppConfig{
		Config:          config,
		AIConfig:        aiConfig,
		ProofingPrompts: proofingPrompts,
		OpenAIKey:       openAIKey,
		AnthropicKey:    anthropicKey,
		ProofType:       *proofType,
		AIProvider:      *aiProvider,
	}, nil
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
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

func loadAIConfig(filename string) (*AIConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var aiConfig AIConfig
	err = json.Unmarshal(data, &aiConfig)
	if err != nil {
		return nil, err
	}

	return &aiConfig, nil
}

func readAPIKey(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func loadProofingPrompts(filename string) ([]ProofingPrompt, error) {
	data, err := os.ReadFile(filename)
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
