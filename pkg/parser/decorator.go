package parser

import (
	"go/ast"
	"regexp"
	"strings"
)

type Decorator struct {
	Type   string   // ctx-get, httpctx-get, env-get и т.д.
	Params []string // параметры декоратора
	Raw    string
}

// parseFieldDecorators парсит декораторы из комментариев поля
func (p *Parser) parseFieldDecorators(field *ast.Field) []Decorator {
	var decorators []Decorator

	re := regexp.MustCompile(`@evl:decor\s+([a-zA-Z_-]+)(?::([^\s@]+))?`)

	// Проверяем комментарии после поля
	if field.Comment != nil {
		for _, comment := range field.Comment.List {
			decs := p.extractDecorators(comment.Text, re)
			decorators = append(decorators, decs...)
		}
	}

	// Проверяем документацию перед полем
	if field.Doc != nil {
		for _, comment := range field.Doc.List {
			decs := p.extractDecorators(comment.Text, re)
			decorators = append(decorators, decs...)
		}
	}

	return decorators
}

// extractDecorators извлекает декораторы из текста комментария
func (p *Parser) extractDecorators(text string, re *regexp.Regexp) []Decorator {
	var decorators []Decorator

	text = strings.TrimPrefix(text, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	text = strings.TrimSpace(text)

	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			dec := Decorator{
				Type: match[1],
				Raw:  match[0],
			}
			if len(match) >= 3 && match[2] != "" {
				dec.Params = strings.Split(match[2], ",")
			}
			decorators = append(decorators, dec)
		}
	}

	return decorators
}
