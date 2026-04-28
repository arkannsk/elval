package annotations

import (
	"regexp"
	"strings"
)

const (
	PrefixOA = "@oa"
)

// OaAnnotation представляет одну аннотацию OpenAPI
type OaAnnotation struct {
	Type  string
	Value string
}

// Directive представляет директиву валидации
type Directive struct {
	Type   string
	Raw    string
	Params []string
}

// ParseOaAnnotationsFromTexts парсит OA аннотации из списка текстов комментариев
func ParseOaAnnotationsFromTexts(texts []string) []OaAnnotation {
	if len(texts) == 0 {
		return nil
	}

	re := regexp.MustCompile(PrefixOA + `[-:]([a-zA-Z0-9_.-]+)\s*(.*)`)

	var descriptions []string
	var otherAnnotations []OaAnnotation

	for _, text := range texts {
		matches := re.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}

			annType := match[1]
			rawValue := strings.TrimSpace(match[2])

			// Убираем ведущие разделители (:, -, пробел) из значения
			val := strings.TrimLeft(rawValue, " :-")
			val = trimQuotes(val)

			if annType == "description" && val != "" {
				descriptions = append(descriptions, val)
			} else {
				otherAnnotations = append(otherAnnotations, OaAnnotation{
					Type:  annType,
					Value: val,
				})
			}
		}
	}

	var oas []OaAnnotation
	if len(descriptions) > 0 {
		oas = append(oas, OaAnnotation{
			Type:  "description",
			Value: strings.Join(descriptions, "\n"),
		})
	}
	oas = append(oas, otherAnnotations...)

	return oas
}

// ParseDirectivesFromTexts парсит директивы валидации из списка текстов комментариев
func ParseDirectivesFromTexts(texts []string) []Directive {
	re := regexp.MustCompile(`@evl:validate\s+([a-zA-Z_-]+)(?::([^\s@]+))?`)

	var directives []Directive

	for _, text := range texts {
		matches := re.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				dir := Directive{
					Type: match[1],
					Raw:  match[0],
				}
				if len(match) >= 3 && match[2] != "" {
					param := strings.TrimSpace(strings.Trim(match[2], ` "'`))
					if dir.Type == "pattern" {
						dir.Params = []string{param}
					} else {
						params := strings.Split(param, ",")
						for i, p := range params {
							params[i] = strings.TrimSpace(p)
						}
						dir.Params = params
					}
				}
				directives = append(directives, dir)
			}
		}
	}
	return directives
}

// CleanComment удаляет маркеры комментариев из строки
func CleanComment(text string) string {
	text = strings.TrimPrefix(text, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	text = strings.TrimSuffix(text, "/")
	return strings.TrimSpace(text)
}
