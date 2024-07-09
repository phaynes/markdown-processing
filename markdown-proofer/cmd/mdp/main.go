package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/phaynes/markdown-proofing-tool/internal/config"
	"github.com/phaynes/markdown-proofing-tool/internal/prompt"
	"github.com/phaynes/markdown-proofing-tool/internal/proofer"
)

func main() {
	appConfig, err := config.Setup()
	if err != nil {
		log.Fatalf("Error setting up configuration: %v", err)
	}

	input := strings.Join(flag.Args(), " ")
	promptText, err := prompt.BuildProofingPrompt(appConfig.ProofingPrompts, appConfig.ProofType)
	if err != nil {
		log.Fatalf("Error building proofing prompt: %v", err)
	}

	proofedText, err := proofer.ProofText(input, promptText, appConfig)
	if err != nil {
		log.Fatalf("Error proofing text: %v", err)
	}

	fmt.Println(proofedText)
}
