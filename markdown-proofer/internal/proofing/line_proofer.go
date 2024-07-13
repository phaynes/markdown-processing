package proofing

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/phaynes/markdown-processing/markdown-proofer/internal/config"
	"github.com/phaynes/markdown-processing/markdown-proofer/internal/proofer"
)

type LineProofingStages struct {
	appConfig *config.AppConfig
}

func NewLineProofingStages(appConfig *config.AppConfig) *LineProofingStages {
	return &LineProofingStages{appConfig: appConfig}
}

func (l *LineProofingStages) Initialize() error {
	// No initialization needed for line proofing
	return nil
}

func (l *LineProofingStages) PrepareContent() (string, error) {
	start, end, err := parseLineRange(l.appConfig.LineRange)
	if err != nil {
		return "", err
	}

	lines, err := readLines(l.appConfig.InputFile)
	if err != nil {
		return "", err
	}

	if end > len(lines) {
		end = len(lines)
	}

	return strings.Join(lines[start-1:end], "\n"), nil
}

func (l *LineProofingStages) ExecuteProofing(content string, promptText string) (string, error) {
	return proofer.ProofText(content, promptText, l.appConfig)
}

func (l *LineProofingStages) ExecuteReview(content string, promptText string) (string, error) {
	return proofer.ReviewText(content, promptText, l.appConfig)
}

func (l *LineProofingStages) HandleOutput(proofedContent string) error {
	start, end, err := parseLineRange(l.appConfig.LineRange)
	if err != nil {
		return err
	}

	lines, err := readLines(l.appConfig.InputFile)
	if err != nil {
		return err
	}

	proofedLines := strings.Split(proofedContent, "\n")
	
	if end > len(lines) {
		end = len(lines)
	}

	// Replace the specified lines with the proofed content
	for i := start - 1; i < end && i-start+1 < len(proofedLines); i++ {
		lines[i] = proofedLines[i-start+1]
	}

	// Write the updated content back to the file
	return writeLines(l.appConfig.InputFile, lines)
}

func (l *LineProofingStages) Finalize() error {
	// No finalization needed for line proofing
	return nil
}

func parseLineRange(lineRange string) (int, int, error) {
	parts := strings.Split(lineRange, "-")
	if len(parts) == 1 {
		line, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid line number: %s", parts[0])
		}
		return line, line, nil
	} else if len(parts) == 2 {
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid start line: %s", parts[0])
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid end line: %s", parts[1])
		}
		if start > end {
			return 0, 0, fmt.Errorf("start line must be less than or equal to end line")
		}
		return start, end, nil
	}
	return 0, 0, fmt.Errorf("invalid line range format: %s", lineRange)
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}
