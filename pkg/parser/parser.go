package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

// Parser парсит Go-файлы и извлекает аннотации валидации
type Parser struct {
	fset             *token.FileSet
	verbose          bool
	typeParser       *TypeParser
	annotationParser *AnnotationParser
}

// NewParser создаёт новый парсер
func NewParser(verbose bool) *Parser {
	return &Parser{
		fset:             token.NewFileSet(),
		verbose:          verbose,
		typeParser:       NewTypeParser(make(map[string]string), "", verbose),
		annotationParser: NewAnnotationParser(verbose),
	}
}

func (p *Parser) ParseFile(filename string) (*ParseResult, error) {
	absFile, err := filepath.Abs(filename)
	if err != nil {
		absFile = filename
	}

	node, err := parser.ParseFile(p.fset, absFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга файла %s: %w", absFile, err)
	}

	// 1. Получаем информацию о модуле
	modInfo, err := resolveModuleInfo(absFile)
	if err != nil {
		return nil, err
	}

	// 2. Инициализируем под-парсеры с контекстом
	p.typeParser = NewTypeParser(make(map[string]string), modInfo.Package, p.verbose)
	p.annotationParser = NewAnnotationParser(p.verbose)

	result := &ParseResult{
		Package: node.Name.Name,
		Structs: make([]Struct, 0),
		Errors:  make([]error, 0),
	}

	if p.verbose {
		log.Printf("DEBUG: File=%s | Module=%s | Root=%s | PkgPath=%s",
			filename, modInfo.Module, modInfo.ModuleRoot, modInfo.PackagePath)
	}

	// 3. Первый проход: собираем все структуры и алиасы
	allStructs, typeAliases := p.parseStructsFirstPass(node, filename, modInfo)
	p.typeParser.typeAliases = typeAliases

	// 4. Второй проход: парсим поля структур и заполняем результат
	p.parseFieldsSecondPass(node, allStructs, result, filename, modInfo)

	if p.verbose {
		log.Printf("DEBUG: Parsed %s, found %d structs in result", filename, len(result.Structs))
	}

	return result, nil
}

// parseStructsFirstPass проходит по всем объявлениям типов и создает карту структур.
func (p *Parser) parseStructsFirstPass(node *ast.File, filename string, modInfo *ModuleInfo) (map[string]*Struct, map[string]string) {
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

			// Проверяем, является ли тип структурой
			if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {

				// Проверяем @oa:ignore на уровне структуры ДО добавления в карту
				structOaAnnotations := p.annotationParser.ParseStructOaAnnotations(genDecl, typeSpec)
				isIgnored := false
				for _, ann := range structOaAnnotations {
					if ann.Type == "ignore" {
						isIgnored = true
						break
					}
				}

				if isIgnored {
					if p.verbose {
						log.Printf("DEBUG: Skipping ignored struct %s in %s", name, filename)
					}
					continue // Пропускаем структуру целиком
				}

				allStructs[name] = &Struct{
					Name:        name,
					File:        filename,
					Fields:      []Field{},
					Package:     modInfo.Package,
					PackagePath: modInfo.PackagePath,
					Module:      modInfo.Module,
					IsIgnored:   false,
				}
			} else {
				// Обработка алиасов
				switch base := typeSpec.Type.(type) {
				case *ast.Ident:
					typeAliases[name] = base.Name
				case *ast.SelectorExpr:
					typeAliases[name] = exprToString(base)
				}
			}
		}
	}

	return allStructs, typeAliases
}

// parseFieldsSecondPass проходит по структурам, парсит их поля и добавляет в result.Structs
func (p *Parser) parseFieldsSecondPass(
	node *ast.File,
	allStructs map[string]*Struct,
	result *ParseResult,
	filename string,
	modInfo *ModuleInfo) {
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

			// Берем структуру из карты. Если её там нет (например, это алиас или игнорируемая структура), пропускаем
			s := allStructs[typeSpec.Name.Name]
			if s == nil {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// Парсим поля
			for _, field := range structType.Fields.List {
				fieldName := getFieldName(field)
				if fieldName == "" {
					continue
				}

				fieldType := p.typeParser.ParseExpr(field.Type, allStructs)

				if p.verbose {
					log.Printf(`PARSED: fieldName: %s, exprType: %T, val: %s → IsGeneric=%v, 
						Base=%q, Args=%+v, IsSlice=%v`,
						fieldName, field.Type, exprToString(field.Type), fieldType.IsGeneric,
						fieldType.GenericBase, fieldType.GenericArgs, fieldType.IsSlice)
				}

				directives := p.annotationParser.ParseFieldDirectives(field)

				// Валидация директив
				for _, dir := range directives {
					if err := ValidateDirective(dir, fieldType); err != nil {
						result.Errors = append(result.Errors, fmt.Errorf("%s:%d: поле %s: %w",
							filename, p.fset.Position(field.Pos()).Line, fieldName, err))
					}
				}

				// Парсим OA-аннотации
				oaAnnotations := p.annotationParser.ParseFieldOaAnnotations(field)

				// Извлекаем специфичные поля и фильтруем список аннотаций
				fAnot := p.processFieldAnnotations(oaAnnotations)

				// Если поле игнорируется, пропускаем его добавление
				if fAnot.IsIgnored {
					if p.verbose {
						log.Printf("DEBUG: Skipping ignored field %s in struct %s", fieldName, s.Name)
					}
					continue
				}

				s.Fields = append(s.Fields, Field{
					Name:          fieldName,
					Type:          fieldType,
					Directives:    directives,
					Decorators:    p.parseFieldDecorators(field),
					Line:          p.fset.Position(field.Pos()).Line,
					OaAnnotations: fAnot.Remaining,
					IsEmbedded:    len(field.Names) == 0,
					OaRewriteRef:  fAnot.RewriteRef,
					OaRewriteType: fAnot.RewriteType,
					IsIgnored:     false, // Поле уже отфильтровано, но для надежности ставим false
					OaIn:          fAnot.OaIn,
					OaParamName:   fAnot.OaParamName,
				})
			}

			// Парсим аннотации структуры (discriminator, oneOf и т.д.)
			// Примечание: мы уже проверили ignore в первом проходе, но здесь нам нужны остальные аннотации
			structOaAnnotations := p.annotationParser.ParseStructOaAnnotations(genDecl, typeSpec)
			p.annotationParser.ExtractDiscriminator(s, structOaAnnotations)

			// Валидация discriminator + oneOf
			if s.Discriminator != nil && len(s.OaOneOf) > 0 {
				for _, typeName := range s.OaOneOf {
					if _, ok := s.Discriminator.Mapping[typeName]; !ok {
						if p.verbose {
							log.Printf("WARNING: Struct %s: oneOf type %q not in discriminator mapping",
								s.Name, typeName)
						}
					}
				}
			}

			// Добавляем структуру в результат
			result.Structs = append(result.Structs, *s)
		}
	}
}

// processFieldAnnotations обрабатывает список аннотаций поля, извлекая специфичные значения
// и возвращая очищенный список аннотаций для дальнейшего использования в шаблонах.
func (p *Parser) processFieldAnnotations(annotations []OaAnnotation) parseFieldAnnotations {
	result := parseFieldAnnotations{}
	remaining := make([]OaAnnotation, 0, len(annotations))

	for _, ann := range annotations {
		switch ann.Type {
		case "rewrite.ref":
			result.RewriteRef = trimQuotes(ann.Value)
		case "rewrite.type":
			result.RewriteType = trimQuotes(ann.Value)
		case "ignore":
			result.IsIgnored = true
		case "in":
			// Формат: "path id" или "query fields"
			parts := strings.Fields(ann.Value)
			if len(parts) >= 1 {
				result.OaIn = parts[0] // "path", "query", etc.
			}
			if len(parts) >= 2 {
				result.OaParamName = parts[1] // "id", "fields", etc.
			}
		default:
			remaining = append(remaining, ann)
		}
	}
	result.Remaining = remaining

	return result
}

type parseFieldAnnotations struct {
	RewriteRef  string
	RewriteType string
	IsIgnored   bool
	OaIn        string
	OaParamName string
	Remaining   []OaAnnotation
}
