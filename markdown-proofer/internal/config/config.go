package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	InputFiles []string `json:"input_files"`
	OutputFile string   `json:"output_file"`
}

type AIConfig struct {
	DefaultAI       string `json:"default_ai"`
	OpenAIKey       string `json:"openAI_key"`
	AnthropicKey    string `json:"anthropic_key"`
	ProofingPrompts string `json:"proofing_prompts"`
	UseGit          bool   `json:"use_git"`
	ProofGitDiff    bool   `json:"proof_git_diff"`
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
	AIConfig        *AIConfig
	ProofingPrompts []ProofingPrompt
	ProofType       string
	AIProvider      string
	Mode            string // "proof" or "review"
	InputFile       string
	OutputFile      string
	ProofingType    string // "command_line", "git_full", "git_diff", or "standard"
	APIKey          string // We'll store the actual API key here
	LineRange       string // New field for line range

}

func Setup() (*AppConfig, error) {
	configFile := flag.String("config", "config.json", "Path to the configuration file")
	aiConfigFile := flag.String("ai-config", "ai_config.json", "Path to the AI configuration file")
	proofType := flag.String("type", "basic-proof", "Type of proofing to perform")
	aiProvider := flag.String("ai", "", "AI provider to use (openai or anthropic)")
	useGit := flag.Bool("use-git", false, "Use git features")
	proofGitDiff := flag.Bool("proof-git-diff", false, "Proof only git diff")
	mode := flag.String("mode", "proof", "Mode of operation: 'proof' or 'review'")
	inputFile := flag.String("input", "", "Input file to proof")
	outputFile := flag.String("output", "", "Output file for proofed content")
	lineRange := flag.String("n", "", "Line range to proof (e.g., '6-10' or '6')")
	flag.Parse()

	aiConfig, err := loadAIConfig(*aiConfigFile)
	if err != nil {
		return nil, fmt.Errorf("error loading AI configuration: %v", err)
	}

	proofingPrompts, err := loadProofingPrompts(aiConfig.ProofingPrompts)
	if err != nil {
		return nil, fmt.Errorf("error loading proofing prompts: %v", err)
	}

	// Determine AI provider
	provider := *aiProvider
	if provider == "" {
		provider = aiConfig.DefaultAI
	}

	// Read the appropriate API key
	var apiKey string
	switch provider {
	case "openai":
		apiKey, err = readAPIKey(aiConfig.OpenAIKey)
	case "anthropic":
		apiKey, err = readAPIKey(aiConfig.AnthropicKey)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", provider)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading API key for %s: %v", provider, err)
	}

	isCommandLineInput := flag.NArg() > 0

	// Determine proofing type
	var proofingType string
	if *lineRange != "" {
		proofingType = "line_range"
	} else if isCommandLineInput {
		proofingType = "command_line"
	} else if (*useGit || aiConfig.UseGit) && *lineRange == "" {
		if *proofGitDiff || aiConfig.ProofGitDiff {
			proofingType = "git_diff"
		} else {
			proofingType = "git_full"
		}
	} else {
		proofingType = "standard"
	}

	// Determine input file
	inputFilePath := *inputFile
	if inputFilePath == "" && !isCommandLineInput {
		config, err := loadConfig(*configFile)
		if err != nil {
			return nil, fmt.Errorf("error loading configuration: %v", err)
		}
		if len(config.InputFiles) > 0 {
			inputFilePath = config.InputFiles[0]
		} else {
			return nil, fmt.Errorf("no input file specified")
		}
	}

	// Determine output file or console output
	var outputFilePath string
	if *outputFile != "" {
		outputFilePath = *outputFile // Use specified output file if provided
	}
	// Note: If outputFilePath is empty, output will go to console

	return &AppConfig{
		AIConfig:        aiConfig,
		ProofingPrompts: proofingPrompts,
		ProofType:       *proofType,
		AIProvider:      provider,
		Mode:            *mode,
		InputFile:       inputFilePath,
		OutputFile:      outputFilePath,
		ProofingType:    proofingType,
		APIKey:          apiKey,
		LineRange:       *lineRange,
	}, nil
}
func readAPIKey(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
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
	if err := json.Unmarshal(data, &aiConfig); err != nil {
		return nil, err
	}

	return &aiConfig, nil
}

func loadProofingPrompts(filename string) ([]ProofingPrompt, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var prompts struct {
		ProofingPrompts []ProofingPrompt `json:"proofing_prompts"`
	}
	if err := json.Unmarshal(data, &prompts); err != nil {
		return nil, err
	}

	return prompts.ProofingPrompts, nil
}
