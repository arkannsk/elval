package annotations

import (
	"regexp"
	"strings"
)

// Decorator представляет декоратор поля
type Decorator struct {
	Type   string   // ctx-get, httpctx-get, env-get и т.д.
	Params []string // параметры декоратора
	Raw    string
}

// ParseDecoratorsFromTexts парсит декораторы @evl:decor из списка текстов комментариев
func ParseDecoratorsFromTexts(texts []string) []Decorator {
	// Измененный regex:
	// 1. @evl:decor\s+ - префикс
	// 2. ([a-zA-Z_-]+) - тип декоратора
	// 3. \s*(.*) - необязательные параметры (всё остальное после типа)
	re := regexp.MustCompile(`@evl:decor\s+([a-zA-Z_-]+)\s*(.*)`)

	var decorators []Decorator

	for _, text := range texts {
		matches := re.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}

			decType := match[1]
			rawParams := strings.TrimSpace(match[2])

			dec := Decorator{
				Type: decType,
				Raw:  match[0],
			}

			if rawParams != "" {
				// Разделяем по запятой, если есть несколько параметров
				params := strings.Split(rawParams, ",")
				for i, p := range params {
					params[i] = strings.TrimSpace(p)
				}
				dec.Params = params
			}

			decorators = append(decorators, dec)
		}
	}
	return decorators
}
