package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

func AddAndCommitChanges(filepath string) error {
	addCmd := exec.Command("git", "add", filepath)
	err := addCmd.Run()
	if err != nil {
		return fmt.Errorf("error adding file to git: %v", err)
	}

	commitCmd := exec.Command("git", "commit", "-m", "Auto-commit before proofing")
	err = commitCmd.Run()
	if err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}

	return nil
}

func GetChangedParagraphs(filepath string) ([]string, error) {
	diffOutput, err := getGitDiff(filepath)
	if err != nil {
		return nil, err
	}

	changedLines := extractChangedLines(diffOutput)
	fileContent, err := readFile(filepath)
	if err != nil {
		return nil, err
	}

	return extractRelevantParagraphs(fileContent, changedLines), nil
}

func getGitDiff(filepath string) (string, error) {
	cmd := exec.Command("git", "diff", "HEAD~1", "--", filepath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git diff: %v", err)
	}
	return string(output), nil
}

func extractChangedLines(diff string) []int {
	var changedLines []int
	lineNum := 0
	diffLines := strings.Split(diff, "\n")
	for _, line := range diffLines {
		if strings.HasPrefix(line, "@@") {
			parts := strings.Split(line, " ")
			if len(parts) > 2 {
				lineInfo := strings.Split(parts[2], ",")
				if len(lineInfo) > 0 {
					lineNum, _ = strconv.Atoi(lineInfo[0][1:])
				}
			}
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			changedLines = append(changedLines, lineNum)
		}
		if !strings.HasPrefix(line, "-") {
			lineNum++
		}
	}
	return changedLines
}

func readFile(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	return string(content), nil
}

func extractRelevantParagraphs(content string, changedLines []int) []string {
	var relevantParagraphs []string
	paragraphs := splitIntoParagraphs(content)

	for _, paragraph := range paragraphs {
		if paragraphContainsChangedLines(paragraph, changedLines) {
			relevantParagraphs = append(relevantParagraphs, paragraph.text)
		}
	}

	return relevantParagraphs
}

type Paragraph struct {
	startLine int
	endLine   int
	text      string
}

func splitIntoParagraphs(content string) []Paragraph {
	var paragraphs []Paragraph
	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentParagraph strings.Builder
	lineNum := 1
	paragraphStart := 1

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" && currentParagraph.Len() > 0 {
			paragraphs = append(paragraphs, Paragraph{
				startLine: paragraphStart,
				endLine:   lineNum - 1,
				text:      strings.TrimSpace(currentParagraph.String()),
			})
			currentParagraph.Reset()
			paragraphStart = lineNum + 1
		} else {
			currentParagraph.WriteString(line)
			currentParagraph.WriteString("\n")
		}
		lineNum++
	}

	if currentParagraph.Len() > 0 {
		paragraphs = append(paragraphs, Paragraph{
			startLine: paragraphStart,
			endLine:   lineNum - 1,
			text:      strings.TrimSpace(currentParagraph.String()),
		})
	}

	return paragraphs
}

func paragraphContainsChangedLines(p Paragraph, changedLines []int) bool {
	for _, line := range changedLines {
		if line >= p.startLine && line <= p.endLine {
			return true
		}
	}
	return false
}
