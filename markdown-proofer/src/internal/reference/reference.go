package reference

import (
	"fmt"
	"strings"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/pdf"
)

func GenerateAPAReference(metadata pdf.Metadata) (string, error) {
	// This is a simplified APA v7 reference generator
	// You may need to adjust this based on the specific requirements and edge cases

	// Format authors
	authors := formatAuthors(metadata.Authors)

	// Format title (assuming it's an article title)
	title := formatTitle(metadata.Title)

	// Since we don't have journal information, we'll assume it's a report or paper
	// You might want to extract more information from the PDF to make this more accurate
	reference := fmt.Sprintf("%s (%d). %s [Report].", authors, metadata.Year, title)

	return reference, nil
}

func formatAuthors(authors string) string {
	var author_list []string
	author_list = strings.Split(authors, ",")

	if len(author_list) == 0 {
		return ""
	}
	if len(author_list) == 1 {
		return author_list[0]
	}
	if len(author_list) == 2 {
		return author_list[0] + " & " + author_list[1]
	}
	return strings.Join(author_list[:len(author_list)-1], ", ") + ", & " + author_list[len(author_list)-1]
}

func formatTitle(title string) string {
	// Capitalize the first letter of the title and any letter following a colon
	words := strings.Split(title, " ")
	for i, word := range words {
		if i == 0 || (i > 0 && words[i-1][len(words[i-1])-1] == ':') {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
