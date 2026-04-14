package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

// Parser парсит Go файлы и извлекает аннотации валидации
type Parser struct {
	fset *token.FileSet
}

// NewParser создает новый парсер
func NewParser() *Parser {
	return &Parser{
		fset: token.NewFileSet(),
	}
}

// ParseFile парсит файл и возвращает структуры с аннотациями
func (p *Parser) ParseFile(filename string) (*ParseResult, error) {
	node, err := parser.ParseFile(p.fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга файла %s: %w", filename, err)
	}

	result := &ParseResult{
		Package: node.Name.Name,
		Structs: make([]Struct, 0),
		Errors:  make([]error, 0),
	}

	// Обходим все объявления в файле
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// Парсим структуру
			s := Struct{
				Name:   typeSpec.Name.Name,
				File:   filename,
				Fields: make([]Field, 0),
			}

			for _, field := range structType.Fields.List {
				// Получаем имя поля
				fieldName := getFieldName(field)
				if fieldName == "" {
					continue
				}

				fieldType := getFieldType(field)

				// Собираем аннотации из комментариев
				directives := p.parseFieldDirectives(field)

				if len(directives) > 0 {
					s.Fields = append(s.Fields, Field{
						Name:       fieldName,
						Type:       fieldType,
						Directives: directives,
						Line:       p.fset.Position(field.Pos()).Line,
					})
				}
			}

			if len(s.Fields) > 0 {
				result.Structs = append(result.Structs, s)
			}
		}
	}

	return result, nil
}

// parseFieldDirectives парсит директивы из комментариев поля
func (p *Parser) parseFieldDirectives(field *ast.Field) []Directive {
	var directives []Directive

	// Регулярка для поиска @evl:validate
	re := regexp.MustCompile(`@evl:validate\s+([a-zA-Z-]+)(?::([^\s@]+))?`)

	// Проверяем комментарии после поля
	if field.Comment != nil {
		for _, comment := range field.Comment.List {
			dirs := p.extractDirectives(comment.Text, re)
			directives = append(directives, dirs...)
		}
	}

	// Проверяем документацию перед полем
	if field.Doc != nil {
		for _, comment := range field.Doc.List {
			dirs := p.extractDirectives(comment.Text, re)
			directives = append(directives, dirs...)
		}
	}

	return directives
}

// extractDirectives извлекает директивы из текста комментария
func (p *Parser) extractDirectives(text string, re *regexp.Regexp) []Directive {
	var directives []Directive

	text = strings.TrimPrefix(text, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	text = strings.TrimSpace(text)

	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			dir := Directive{
				Type: match[1],
				Raw:  match[0],
			}
			if len(match) >= 3 && match[2] != "" {
				// Параметры могут быть разделены запятыми
				params := strings.Split(match[2], ",")
				for i, param := range params {
					// Удаляем кавычки, если они есть
					param = strings.TrimSpace(param)
					param = strings.Trim(param, `"'`)
					params[i] = param
				}
				dir.Params = params
			}
			directives = append(directives, dir)
		}
	}

	return directives
}

// ValidateDirectives проверяет все директивы в структурах
func (r *ParseResult) ValidateDirectives() []error {
	var errors []error

	for _, s := range r.Structs {
		for _, field := range s.Fields {
			for _, dir := range field.Directives {
				if err := ValidateDirective(dir, field.Type); err != nil {
					errors = append(errors, fmt.Errorf("%s:%d: поле %s: %w",
						s.File, field.Line, field.Name, err))
				}
			}
		}
	}

	return errors
}

// getFieldName возвращает имя поля из ast.Field
func getFieldName(field *ast.Field) string {
	if len(field.Names) > 0 {
		return field.Names[0].Name
	}
	// Анонимное поле (встраивание)
	if ident, ok := field.Type.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// getFieldType возвращает тип поля с поддержкой слайсов и указателей
func getFieldType(field *ast.Field) FieldType {
	ft := FieldType{
		IsSlice:   false,
		IsPointer: false,
	}

	switch t := field.Type.(type) {
	case *ast.Ident:
		ft.Name = t.Name

	case *ast.StarExpr:
		ft.IsPointer = true
		// Определяем базовый тип для указателя
		if ident, ok := t.X.(*ast.Ident); ok {
			ft.Name = ident.Name
		} else if sel, ok := t.X.(*ast.SelectorExpr); ok {
			// *time.Time, *time.Duration
			if pkg, ok := sel.X.(*ast.Ident); ok {
				ft.Name = pkg.Name + "." + sel.Sel.Name
			}
		} else if arr, ok := t.X.(*ast.ArrayType); ok {
			// *[]T - указатель на слайс
			ft.IsSlice = true
			if ident, ok := arr.Elt.(*ast.Ident); ok {
				ft.Name = ident.Name
			}
		}

	case *ast.ArrayType:
		ft.IsSlice = true
		if ident, ok := t.Elt.(*ast.Ident); ok {
			ft.Name = ident.Name
		} else if star, ok := t.Elt.(*ast.StarExpr); ok {
			// []*T - слайс указателей
			if ident, ok := star.X.(*ast.Ident); ok {
				ft.Name = "*" + ident.Name
			}
		} else if sel, ok := t.Elt.(*ast.SelectorExpr); ok {
			// []time.Time - слайс импортированных типов
			if pkg, ok := sel.X.(*ast.Ident); ok {
				ft.Name = pkg.Name + "." + sel.Sel.Name
			}
		}

	case *ast.SelectorExpr:
		// Импортированный тип: time.Time, time.Duration
		if ident, ok := t.X.(*ast.Ident); ok {
			ft.Name = ident.Name + "." + t.Sel.Name
		}
	}

	return ft
}
