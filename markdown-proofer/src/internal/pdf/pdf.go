package pdf

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

type Metadata struct {
	Title   string `json:"title"`
	Authors string `json:"authors"`
	Year    int    `json:"year"`
}

func ExtractContent(path string) (string, Metadata, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", Metadata{}, fmt.Errorf("error opening PDF: %v", err)
	}
	defer f.Close()

	totalPage := r.NumPage()

	var content bytes.Buffer
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		text, err := p.GetPlainText(nil)
		if err != nil {
			fmt.Printf("Warning: Error extracting text from page %d: %v\n", pageIndex, err)
			continue
		}
		content.WriteString(text)
		content.WriteString("\n")
	}

	if content.Len() == 0 {
		return "", Metadata{}, fmt.Errorf("extracted text is empty for file %s", path)
	}

	metadata, err := extractMetadata(r)
	if err != nil {
		return "", Metadata{}, fmt.Errorf("error extracting metadata: %v", err)
	}

	return content.String(), metadata, nil
}

func extractMetadata(r *pdf.Reader) (Metadata, error) {
	metadata := Metadata{
		Title:   "Unknown Title",
		Authors: "Unknown Author",
		Year:    0,
	}

	info := r.Trailer().Key("Info")

	if title := info.Key("Title"); title.Kind() == pdf.String {
		metadata.Title = title.String()
	}

	if author := info.Key("Author"); author.Kind() == pdf.String {
		metadata.Authors = author.String()
	}

	if creationDate := info.Key("CreationDate"); creationDate.Kind() == pdf.String {
		dateStr := strings.Trim(creationDate.String(), "(D:")
		dateStr = strings.Split(dateStr, "+")[0] // Remove timezone part
		if t, err := time.Parse("20060102150405", dateStr); err == nil {
			metadata.Year = t.Year()
		}
	}

	return metadata, nil
}

func ExtractText(path string) (string, error) {
	content, _, err := ExtractContent(path)
	return content, err
}
