package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
)

type PDFConfig struct {
	InputDirectory  string `json:"input_directory"`
	OutputDirectory string `json:"output_directory"`
	SummaryCriteria string `json:"summary_criteria"`
	*AppConfig
}

func SetupPDFConfig() (*PDFConfig, error) {
	err := license.SetMeteredKey("de9fa24d31a17f487034f7cc6f01c481a9b298366f7e72f8f14605afb1d38cc7")
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}
	pdfConfigFile := flag.String("pdf-config", "pdfconfig.json", "Path to the PDF configuration file")
	flag.Parse()

	// Read PDF-specific config
	file, err := os.Open(*pdfConfigFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pdfConfig PDFConfig
	err = json.NewDecoder(file).Decode(&pdfConfig)
	if err != nil {
		return nil, err
	}

	// Setup the main AppConfig
	appConfig, err := Setup()
	if err != nil {
		return nil, fmt.Errorf("error setting up main configuration: %v", err)
	}

	pdfConfig.AppConfig = appConfig

	return &pdfConfig, nil
}
