package proofing

import (
	"flag"
	"fmt" // Add this import
	"strings"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/proofer"
)

type CommandLineProofingStages struct {
    appConfig *config.AppConfig
}

func NewCommandLineProofingStages(appConfig *config.AppConfig) *CommandLineProofingStages {
    return &CommandLineProofingStages{appConfig: appConfig}
}

func (c *CommandLineProofingStages) Initialize() error {
    // No initialization needed for command line input
    return nil
}

func (c *CommandLineProofingStages) PrepareContent() (string, error) {
    return strings.Join(flag.Args(), " "), nil
}

func (c *CommandLineProofingStages) ExecuteProofing(content string, promptText string) (string, error) {
    return proofer.ProofText(content, promptText, c.appConfig)
}

func (c *CommandLineProofingStages) ExecuteReview(content string, promptText string) (string, error) {
    return proofer.ReviewText(content, promptText, c.appConfig)
}

func (c *CommandLineProofingStages) HandleOutput(proofedContent string) error {
    // For command line input, we'll just print to stdout
    fmt.Println(proofedContent)
    return nil
}

func (c *CommandLineProofingStages) Finalize() error {
    // No finalization needed for command line input
    return nil
}
