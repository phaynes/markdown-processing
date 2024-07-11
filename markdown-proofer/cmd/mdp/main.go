package main

import (
	"log"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/prompt"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/proofing"
)

func main() {
	appConfig, err := config.Setup()
	if err != nil {
		log.Fatalf("Error setting up configuration: %v", err)
	}

	stages, err := proofing.NewProofingStages(appConfig)
	if err != nil {
		log.Fatalf("Error creating proofing stages: %v", err)
	}

	// Initialize
	if err := stages.Initialize(); err != nil {
		log.Fatalf("Error initializing proofing stages: %v", err)
	}

	// Prepare Content
	contentToProof, err := stages.PrepareContent()
	if err != nil {
		log.Fatalf("Error preparing content to proof: %v", err)
	}

	// Build proofing prompt
	promptText, err := prompt.BuildProofingPrompt(appConfig.ProofingPrompts, appConfig.ProofType)
	if err != nil {
		log.Fatalf("Error building proofing prompt: %v", err)
	}

	// Execute Proofing or Review
	var result string
	if appConfig.Mode == "proof" {
		result, err = stages.ExecuteProofing(contentToProof, promptText)
	} else if appConfig.Mode == "review" {
		result, err = stages.ExecuteReview(contentToProof, promptText)
	} else {
		log.Fatalf("Invalid mode: %s. Must be 'proof' or 'review'", appConfig.Mode)
	}

	if err != nil {
		log.Fatalf("Error executing %s: %v", appConfig.Mode, err)
	}

	// Handle Output
	if err := stages.HandleOutput(result); err != nil {
		log.Fatalf("Error handling output: %v", err)
	}

	// Finalize
	if err := stages.Finalize(); err != nil {
		log.Fatalf("Error finalizing proofing: %v", err)
	}

}
