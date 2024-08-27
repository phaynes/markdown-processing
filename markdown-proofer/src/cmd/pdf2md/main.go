package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/ai"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/pdf"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/reference"
)

func main() {
	log.SetFlags(0) // Remove timestamp from log output

	// Setup configuration
	pdfConfig, err := config.SetupPDFConfig()
	if err != nil {
		log.Fatalf("Error setting up configuration: %v", err)
	}

	// Process each PDF in the directory
	err = filepath.Walk(pdfConfig.InputDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return nil
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".pdf" {
			if err := processPDF(path, pdfConfig); err != nil {
				log.Printf("Error processing %s: %v", path, err)
				// Continue processing other files
				return nil
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through directory: %v", err)
	}
}

func processPDF(path string, pdfConfig *config.PDFConfig) error {
	// Extract text and metadata from PDF
	text, metadata, err := pdf.ExtractContent(path)
	if err != nil {
		return fmt.Errorf("error extracting content from PDF: %v", err)
	}

	if text == "" {
		return fmt.Errorf("extracted text is empty, skipping file")
	}

	// Use AI to enhance metadata if needed
	enhancedMetadata, err := ai.EnhanceMetadata(text, metadata, pdfConfig.AppConfig)
	if err != nil {
		log.Printf("Warning: error enhancing metadata: %v. Using original metadata.", err)
		enhancedMetadata = metadata
	}

	// Generate APA v7 reference
	apaReference, err := reference.GenerateAPAReference(enhancedMetadata)
	if err != nil {
		log.Printf("Warning: error generating APA reference: %v. Skipping reference.", err)
		apaReference = "Reference generation failed"
	}

	// Generate new filename
	newFilename := generateFilename(enhancedMetadata)

	// Combine reference, metadata, and content
	markdownContent := fmt.Sprintf("---\ntitle: %s\nauthors: %s\nyear: %d\nreference: %s\n---\n\n%s",
		enhancedMetadata.Title, enhancedMetadata.Authors, enhancedMetadata.Year, apaReference, text)

	// Write markdown content to file
	outputPath := filepath.Join(pdfConfig.OutputDirectory, newFilename)
	if err := os.WriteFile(outputPath, []byte(markdownContent), 0644); err != nil {
		return fmt.Errorf("error writing markdown file: %v", err)
	}

	fmt.Printf("Processed: %s -> %s\n", path, outputPath)
	return nil
}

func generateFilename(metadata pdf.Metadata) string {
	safeTitle := formatTitleForFilename(metadata.Title, 60)
	safeAuthor := formatSurnameForFilename(metadata.Authors)
	return fmt.Sprintf("%d-(%s)-%s.md", metadata.Year, safeAuthor, safeTitle)
}

// Function to limit the length of a title to 50 characters and clean it for use in filenames
func formatTitleForFilename(title string, maxLength int) string {
	// Remove any invalid filename characters
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	cleanTitle := re.ReplaceAllString(title, "")

	// Trim whitespace from both ends
	cleanTitle = strings.TrimSpace(cleanTitle)

	// Truncate the title to the max length if necessary
	if len(cleanTitle) > maxLength {
		cleanTitle = cleanTitle[:maxLength]
	}

	// Remove any trailing whitespace after truncation
	cleanTitle = strings.TrimSpace(cleanTitle)
	cleanTitle = strings.ReplaceAll(cleanTitle, " ", "-")

	return cleanTitle
}
func formatSurnameForFilename(authors string) string {
	// Split authors by commas and trim spaces
	authorList := strings.Split(authors, ",")
	if len(authorList) == 0 {
		return ""
	}

	// Get the first author and trim spaces
	firstAuthor := strings.TrimSpace(authorList[0])

	// Extract the last name (surname) from the first author
	parts := strings.Split(firstAuthor, " ")
	lastName := parts[len(parts)-1]

	// Remove any invalid filename characters
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	formattedSurname := re.ReplaceAllString(lastName, "")

	return formattedSurname
}
