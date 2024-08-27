package pdf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type Metadata struct {
	Title   string `json:"title"`
	Authors string `json:"authors"`
	Year    int    `json:"year"`
}

func ExtractContent(path string) (string, Metadata, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", Metadata{}, fmt.Errorf("error opening PDF: %v", err)
	}
	defer file.Close()

	// Create a more lenient configuration
	config := &model.Configuration{
		ValidationMode:   model.ValidationRelaxed,
		Reader15:         true, // Enable reading of PDF 1.5+ files
		DecodeAllStreams: true,
	}

	ctx, err := pdfcpu.Read(file, config)
	if err != nil {
		return "", Metadata{}, fmt.Errorf("error reading PDF: %v", err)
	}

	var content bytes.Buffer
	for i := 1; i <= ctx.PageCount; i++ {
		pageContent, err := pdfcpu.ExtractPageContent(ctx, i)
		if err != nil {
			fmt.Printf("Warning: Error extracting content from page %d: %v\n", i, err)
			continue // Skip this page but continue with others
		}

		// Create a buffer to read the page content
		var pageBuffer bytes.Buffer
		_, err = io.Copy(&pageBuffer, pageContent)
		if err != nil {
			fmt.Printf("Warning: Error copying content from page %d: %v\n", i, err)
			continue // Skip this page but continue with others
		}

		if pageBuffer.Len() == 0 {
			fmt.Printf("Warning: No text extracted from page %d of file %s\n", i, path)
			continue // Skip this page but continue with others
		}

		// Copy page buffer to content buffer
		content.Write(pageBuffer.Bytes())
	}

	if content.Len() == 0 {
		return "", Metadata{}, fmt.Errorf("extracted text is empty for file %s", path)
	}

	metadata, err := extractMetadata(ctx)
	if err != nil {
		return "", Metadata{}, fmt.Errorf("error extracting metadata: %v", err)
	}

	return content.String(), metadata, nil
}

func extractMetadata(ctx *model.Context) (Metadata, error) {
	metadata := Metadata{
		Title:   "Unknown Title",
		Authors: "Unknown Author",
		Year:    0,
	}

	if ctx.XRefTable.Info != nil {
		info, err := ctx.XRefTable.DereferenceDict(*ctx.XRefTable.Info)
		if err != nil {
			return metadata, fmt.Errorf("error dereferencing info dictionary: %v", err)
		}

		if title := info.StringEntry("Title"); title != nil && *title != "" {
			metadata.Title = *title
		}

		if author := info.StringEntry("Author"); author != nil && *author != "" {
			metadata.Authors = *author
		}

		if creationDate := info.StringEntry("CreationDate"); creationDate != nil && *creationDate != "" {
			if t, err := time.Parse("20060102150405-07'00'", *creationDate); err == nil {
				metadata.Year = t.Year()
			}
		}
	}

	return metadata, nil
}

func ExtractText(path string) (string, error) {
	content, _, err := ExtractContent(path)
	return content, err
}
