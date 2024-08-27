package pdf

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

type Metadata struct {
	Title   string `json:"title"`
	Authors string `json:"authors"` // Change to string to match JSON response
	Year    int    `json:"year"`
}

func ExtractContent(path string) (string, Metadata, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", Metadata{}, fmt.Errorf("error opening PDF: %v", err)
	}
	defer file.Close()

	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return "", Metadata{}, fmt.Errorf("error creating PDF reader: %v", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", Metadata{}, fmt.Errorf("error getting number of pages: %v", err)
	}

	var buf bytes.Buffer
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return "", Metadata{}, fmt.Errorf("error getting page %d: %v", i, err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return "", Metadata{}, fmt.Errorf("error creating extractor for page %d: %v", i, err)
		}

		text, err := ex.ExtractText()
		if err != nil {
			return "", Metadata{}, fmt.Errorf("error extracting text from page %d: %v", i, err)
		}

		buf.WriteString(text)
		buf.WriteString("\n")
	}

	metadata, err := extractMetadata(pdfReader)
	if err != nil {
		return "", Metadata{}, fmt.Errorf("error extracting metadata: %v", err)
	}

	return buf.String(), metadata, nil
}

func extractMetadata(r *model.PdfReader) (Metadata, error) {
	metadata := Metadata{
		Title:   "Unknown Title",
		Authors: "Unknown Authors",
		Year:    0,
	}

	trailer, err := r.GetTrailer()
	if err != nil {
		return metadata, fmt.Errorf("error getting PDF trailer: %v", err)
	}

	if trailer != nil {
		infoDict, ok := trailer.Get("Info").(*core.PdfIndirectObject)
		if !ok {
			return metadata, nil
		}

		dict, ok := infoDict.PdfObject.(*core.PdfObjectDictionary)
		if !ok {
			return metadata, nil
		}

		if title, ok := dict.Get("Title").(*core.PdfObjectString); ok {
			metadata.Title = title.String()
		}

		if author, ok := dict.Get("Author").(*core.PdfObjectString); ok {
			metadata.Authors = author.String()
		}

		if creationDate, ok := dict.Get("CreationDate").(*core.PdfObjectString); ok {
			if t, err := parseTime(creationDate.String()); err == nil {
				metadata.Year = t.Year()
			}
		} else if modDate, ok := dict.Get("ModDate").(*core.PdfObjectString); ok {
			if t, err := parseTime(modDate.String()); err == nil {
				metadata.Year = t.Year()
			}
		}
	}

	return metadata, nil
}

func parseTime(dateString string) (time.Time, error) {
	// PDF time format: (D:YYYYMMDDHHmmSSOHH'mm')
	dateString = strings.Trim(dateString, "(D:")
	dateString = strings.Split(dateString, "+")[0]
	return time.Parse("20060102150405", dateString)
}

func ExtractText(path string) (string, error) {
	content, _, err := ExtractContent(path)
	return content, err
}
