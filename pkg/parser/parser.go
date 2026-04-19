package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"regexp"
	"strings"
)

// Parser парсит Go файлы и извлекает аннотации валидации
type Parser struct {
	fset        *token.FileSet
	verbose     bool
	typeAliases map[string]string
}

// NewParser создает новый парсер
func NewParser(verbose bool) *Parser {
	return &Parser{
		fset:    token.NewFileSet(),
		verbose: verbose,
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

	// Сначала собираем все структуры в файле (даже без аннотаций)
	allStructs := make(map[string]*Struct)
	typeAliases := make(map[string]string)

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

			name := typeSpec.Name.Name

			switch base := typeSpec.Type.(type) {
			case *ast.StructType:
				// Настоящая структура — добавляем в allStructs
				allStructs[name] = &Struct{Name: name, File: filename, Fields: []Field{}}

			case *ast.Ident:
				// Алиас на примитив: type Theme string
				typeAliases[name] = base.Name // "Theme" → "string"

			case *ast.SelectorExpr:
				// Алиас на внешний тип: type Timestamp time.Time
				typeAliases[name] = p.exprToString(base) // "Timestamp" → "time.Time"
			}
			// Остальные случаи (slice, pointer, etc.) можно игнорировать для алиасов
		}
	}

	p.typeAliases = typeAliases

	// Второй проход: парсим поля структур
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

			s := allStructs[typeSpec.Name.Name]
			if s == nil {
				continue
			}

			for _, field := range structType.Fields.List {
				fieldName := getFieldName(field)
				if fieldName == "" {
					continue
				}

				fieldType := p.getFieldTypeWithStructs(field, allStructs)
				directives := p.parseFieldDirectives(field)

				for _, dir := range directives {
					if err := ValidateDirective(dir, fieldType); err != nil {
						result.Errors = append(result.Errors, fmt.Errorf("%s:%d: поле %s: %w",
							filename, p.fset.Position(field.Pos()).Line, fieldName, err))
					}
					if p.verbose {
						log.Printf("directive: %v", dir)
					}
				}

				// Добавляем поле даже если нет директив (оно может быть вложенной структурой)
				s.Fields = append(s.Fields, Field{
					Name:          fieldName,
					Type:          fieldType,
					Directives:    directives,
					Decorators:    p.parseFieldDecorators(field),
					Line:          p.fset.Position(field.Pos()).Line,
					OaAnnotations: p.parseFieldOaAnnotations(field),
				})
			}
			// Добавляем ВСЕ структуры в результат
			// Фильтрация будет в генераторе
			result.Structs = append(result.Structs, *s)
		}
	}

	return result, nil
}

func (p *Parser) getFieldTypeWithStructs(field *ast.Field, allStructs map[string]*Struct) FieldType {
	return p.parseFieldTypeExpr(field.Type, allStructs)
}

// parseFieldTypeExpr рекурсивно разбирает AST-выражение типа поля.
func (p *Parser) parseFieldTypeExpr(expr ast.Expr, allStructs map[string]*Struct) FieldType {
	var ft FieldType

	switch t := expr.(type) {
	case *ast.Ident:
		ft.Name = t.Name

		// 1. Проверяем, структура ли это
		if _, ok := allStructs[t.Name]; ok {
			ft.IsStruct = true
			return ft
		}

		if baseType, ok := p.typeAliases[t.Name]; ok {
			ft.BaseType = baseType

			// Рекурсивно парсим базовый тип, чтобы унаследовать IsStruct/IsCustom
			var baseExpr ast.Expr
			if idx := strings.Index(baseType, "."); idx > 0 {
				baseExpr = &ast.SelectorExpr{
					X:   &ast.Ident{Name: baseType[:idx]},
					Sel: &ast.Ident{Name: baseType[idx+1:]},
				}
			} else {
				baseExpr = &ast.Ident{Name: baseType}
			}
			baseFt := p.parseFieldTypeExpr(baseExpr, allStructs)

			ft.IsStruct = baseFt.IsStruct
			ft.IsCustom = baseFt.IsCustom
			return ft
		}

		// 3. Обычный тип
		ft.IsCustom = !isBuiltin(ft.Name)
		return ft
	case *ast.SelectorExpr:
		ft.Name = p.exprToString(t)
		ft.IsCustom = true

	case *ast.StarExpr: // *T
		inner := p.parseFieldTypeExpr(t.X, allStructs)
		ft.Name = "*" + inner.Name
		ft.IsPointer = true
		ft.IsStruct = inner.IsStruct
		ft.IsCustom = inner.IsCustom
		ft.IsGeneric = inner.IsGeneric
		ft.GenericBase = inner.GenericBase
		ft.GenericArgs = inner.GenericArgs

	case *ast.ArrayType: // []T
		if t.Len == nil {
			inner := p.parseFieldTypeExpr(t.Elt, allStructs)
			ft.Name = "[]" + inner.Name
			ft.IsSlice = true
			if inner.IsGeneric {
				ft.IsGeneric = inner.IsGeneric
				ft.GenericBase = inner.GenericBase
				ft.GenericArgs = inner.GenericArgs
			}
			ft.IsStruct = inner.IsStruct
			ft.IsCustom = inner.IsCustom
		} else {
			ft.Name = getTypeString(expr)
			ft.IsCustom = true
		}

	case *ast.IndexExpr: // Option[T]
		base := p.exprToString(t.X)
		inner := p.parseFieldTypeExpr(t.Index, allStructs)

		ft.Name = fmt.Sprintf("%s[%s]", base, inner.Name)
		ft.IsGeneric = true
		ft.GenericBase = base
		ft.GenericArgs = append(ft.GenericArgs, inner)
		ft.IsCustom = true

	case *ast.IndexListExpr: // T[K,V]
		base := p.exprToString(t.X)
		var args []string
		for _, idx := range t.Indices {
			arg := p.parseFieldTypeExpr(idx, allStructs)
			ft.GenericArgs = append(ft.GenericArgs, arg)
			args = append(args, arg.Name)
		}
		ft.Name = fmt.Sprintf("%s[%s]", base, strings.Join(args, ", "))
		ft.IsGeneric = true
		ft.GenericBase = base
		ft.IsCustom = true

	case *ast.MapType:
		ft.Name = fmt.Sprintf("map[%s]%s",
			p.exprToString(t.Key),
			p.exprToString(t.Value))
		ft.IsCustom = true

	default:
		ft.Name = getTypeString(expr)
		ft.IsCustom = true
	}
	if p.verbose {
		fmt.Printf("PARSED: %s → IsGeneric=%v, Base=%q, Args=%+v\n",
			getTypeString(expr), ft.IsGeneric, ft.GenericBase, ft.GenericArgs)
	}

	return ft
}

func (p *Parser) exprToString(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return p.exprToString(v.X) + "." + v.Sel.Name
	case *ast.StarExpr:
		return "*" + p.exprToString(v.X)
	case *ast.ArrayType:
		if v.Len == nil {
			return "[]" + p.exprToString(v.Elt)
		}
		return fmt.Sprintf("[%s]%s", p.exprToString(v.Len), p.exprToString(v.Elt))
	case *ast.IndexExpr:
		return p.exprToString(v.X) + "[" + p.exprToString(v.Index) + "]"
	case *ast.IndexListExpr:
		var parts []string
		for _, idx := range v.Indices {
			parts = append(parts, p.exprToString(idx))
		}
		return p.exprToString(v.X) + "[" + strings.Join(parts, ", ") + "]"
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", p.exprToString(v.Key), p.exprToString(v.Value))
	case *ast.ParenExpr:
		return "(" + p.exprToString(v.X) + ")"
	default:
		return "unknown"
	}
}

// parseFieldDirectives парсит директивы из комментариев поля
func (p *Parser) parseFieldDirectives(field *ast.Field) []Directive {
	var directives []Directive

	// Регулярка для поиска @evl:validate
	re := regexp.MustCompile(`@evl:validate\s+([a-zA-Z_-]+)(?::([^\s@]+))?`)

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

func (p *Parser) parseFieldOaAnnotations(field *ast.Field) []OaAnnotation {
	var oas []OaAnnotation
	re := regexp.MustCompile(`@oa:([a-zA-Z_-]+)\s+(.+)`)

	// Check comments after field
	if field.Comment != nil {
		for _, comment := range field.Comment.List {
			oas = append(oas, p.extractOaAnnotations(comment.Text, re)...)
		}
	}
	// Check doc before field
	if field.Doc != nil {
		for _, comment := range field.Doc.List {
			oas = append(oas, p.extractOaAnnotations(comment.Text, re)...)
		}
	}
	return oas
}

func (p *Parser) extractOaAnnotations(text string, re *regexp.Regexp) []OaAnnotation {
	var oas []OaAnnotation
	text = strings.TrimPrefix(text, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	text = strings.TrimSpace(text)

	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			oas = append(oas, OaAnnotation{
				Type:  match[1],
				Value: strings.TrimSpace(match[2]),
			})
		}
	}
	return oas
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
				param := match[2]
				param = strings.TrimSpace(param)
				param = strings.Trim(param, `"'`)

				// Для pattern не разбиваем по запятым
				if dir.Type == "pattern" {
					dir.Params = []string{param}
				} else {
					// Для остальных директив разбиваем по запятым
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

	return directives
}

func (r *ParseResult) ValidateDirectives() []DirectiveError {
	var errors []DirectiveError

	for _, s := range r.Structs {
		for _, field := range s.Fields {
			for _, dir := range field.Directives {
				severity, err := validateDirective(dir, field.Type)
				if err != nil {
					errors = append(errors, DirectiveError{
						File:      s.File,
						Line:      field.Line,
						Struct:    s.Name,
						Field:     field.Name,
						Directive: dir.Type,
						Message:   err.Error(),
						Severity:  severity,
					})
				}
			}
		}
	}

	return errors
}

// CollectValidationImports — импорты ТОЛЬКО для валидации/декорирования
func CollectValidationImports(structs []Struct) map[string]string {
	required := make(map[string]string)

	required["errs"] = "github.com/arkannsk/elval/pkg/errs"
	required["validator"] = "github.com/arkannsk/elval/pkg/validator"
	required["context"] = "context" // для Decorate

	needsElval := false

	for _, s := range structs {
		for _, field := range s.Fields {
			// --- Анализ типов для elval/time ---
			var checkType func(ft FieldType)
			checkType = func(ft FieldType) {
				if ft.IsGeneric && len(ft.GenericArgs) > 0 {
					needsElval = true
					for _, arg := range ft.GenericArgs {
						checkType(arg)
					}
					return
				}
				if ft.Name == "time.Time" || ft.Name == "time.Duration" {
					required["time"] = "time"
				}
				if ft.IsSlice || ft.IsPointer {
					base := strings.TrimPrefix(ft.Name, "[]")
					base = strings.TrimPrefix(base, "*")
					if base == "time.Time" || base == "time.Duration" {
						required["time"] = "time"
					}
				}
			}
			checkType(field.Type)

			// --- Директивы валидации → импорты ---
			for _, dir := range field.Directives {
				switch dir.Type {
				case "uuid":
					required["uuid"] = "github.com/google/uuid"
				}
			}

			// --- Декораторы → импорты ---
			for _, dec := range field.Decorators {
				switch dec.Type {
				case "uuid-gen":
					required["uuid"] = "github.com/google/uuid"
				case "env-get", "env_default":
					required["os"] = "os"
				case "time-now":
					required["time"] = "time"
				case "httpctx-get":
					required["net/http"] = "net/http"
				case "trim", "lower", "upper":
					required["strings"] = "strings"
				}
			}
		}
	}

	if needsElval {
		required["elval"] = "github.com/arkannsk/elval"
	}

	return required
}

// CollectOpenAPIImports — импорты ТОЛЬКО для OpenAPI схем
func CollectOpenAPIImports(structs []Struct) map[string]string {
	required := make(map[string]string)
	required["oa"] = "github.com/arkannsk/elval/pkg/oa"

	return required
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

func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name

	case *ast.StarExpr:
		return "*" + getTypeString(t.X)

	case *ast.ArrayType:
		return "[]" + getTypeString(t.Elt)

	case *ast.SelectorExpr:
		return getTypeString(t.X) + "." + t.Sel.Name

	case *ast.IndexExpr:
		// mo.Option[string]
		return getTypeString(t.X) + "[" + getTypeString(t.Index) + "]"

	case *ast.IndexListExpr:
		result := getTypeString(t.X) + "["
		for i, idx := range t.Indices {
			if i > 0 {
				result += ", "
			}
			result += getTypeString(idx)
		}
		result += "]"
		return result

	default:
		return "unknown"
	}
}

// isBuiltin проверяет, является ли тип встроенным в Go.
func isBuiltin(name string) bool {
	_, ok := map[string]bool{
		"bool": true, "int": true, "int8": true, "int16": true, "int32": true, "int64": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true, "uintptr": true,
		"float32": true, "float64": true, "complex64": true, "complex128": true,
		"string": true, "byte": true, "rune": true, "error": true, "any": true,
	}[name]
	return ok
}
