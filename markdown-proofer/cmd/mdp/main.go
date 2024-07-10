package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/prompt"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/proofer"
)

func main() {
	appConfig, err := config.Setup()
	if err != nil {
		log.Fatalf("Error setting up configuration: %v", err)
	}

	var input string
	if flag.NArg() > 0 {
		// Use command-line input if provided
		input = strings.Join(flag.Args(), " ")
	} else if appConfig.InputFile != "" {
		// Use input file from config if no command-line input
		input, err = readInputFile(appConfig.InputFile)
		if err != nil {
			log.Fatalf("Error reading input file %s: %v", appConfig.InputFile, err)
		}
	} else {
		log.Fatalf("No input provided. Please provide text as a command-line argument or specify an input file in the config.")
	}

	promptText, err := prompt.BuildProofingPrompt(appConfig.ProofingPrompts, appConfig.ProofType)
	if err != nil {
		log.Fatalf("Error building proofing prompt: %v", err)
	}

	// Handle additional information if required
	for _, p := range appConfig.ProofingPrompts {
		if p.Label == appConfig.ProofType && p.RequestAdditionalInfo {
			if appConfig.AdditionalInfo == "" {
				appConfig.AdditionalInfo = getAdditionalInfoFromUser()
			}
			promptText += "\nAdditional Information: " + appConfig.AdditionalInfo
			break
		}
	}

	proofedText, err := proofer.ProofText(input, promptText, appConfig)
	if err != nil {
		log.Fatalf("Error proofing text: %v", err)
	}

	// Only write to a file if explicitly specified, otherwise print to console
	if appConfig.OutputFile != "" {
		err = writeOutputFile(appConfig.OutputFile, proofedText)
		if err != nil {
			log.Fatalf("Error writing output file: %v", err)
		}
		fmt.Printf("Proofed content written to %s\n", appConfig.OutputFile)
	} else {
		fmt.Println(proofedText)
	}
}

func readInputFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func writeOutputFile(filename string, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func getAdditionalInfoFromUser() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter additional information for proofing: ")
	info, _ := reader.ReadString('\n')
	return strings.TrimSpace(info)
}
