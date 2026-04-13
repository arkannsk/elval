package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ParseComments читает файл и находит аннотации в комментариях
func ParseComments(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %s: %w", filename, err)
	}
	defer file.Close()

	var annotations []string
	re := regexp.MustCompile(`@evl:validate\s+[a-zA-Z-]+(?::[^\s]+)?`)

	scanner := bufio.NewScanner(file)
	inMultilineComment := false

	for scanner.Scan() {
		line := scanner.Text()

		// Обработка многострочных комментариев
		if strings.Contains(line, "/*") {
			inMultilineComment = true
		}

		var textToSearch string
		if inMultilineComment {
			textToSearch = line
		} else {
			if idx := strings.Index(line, "//"); idx != -1 {
				textToSearch = line[idx+2:]
			} else {
				continue
			}
		}

		// Ищем все аннотации в строке
		matches := re.FindAllString(textToSearch, -1)

		for _, match := range matches {
			match = strings.TrimSpace(match)
			match = strings.TrimSuffix(match, "*/")
			match = strings.ReplaceAll(match, "\\\\", "\\")

			if match != "" {
				annotations = append(annotations, match)
			}
		}

		if strings.Contains(line, "*/") {
			inMultilineComment = false
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	return annotations, nil
}
