package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
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

// ParseFile парсит файл и возвращает структуры с аннотациями
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
				allStructs[name] = &Struct{
					Name:        name,
					File:        filename,
					Fields:      []Field{},
					Package:     modInfo.Package,
					PackagePath: modInfo.PackagePath,
					Module:      modInfo.Module,
				}
			case *ast.Ident:
				typeAliases[name] = base.Name
			case *ast.SelectorExpr:
				typeAliases[name] = exprToString(base)
			}
		}
	}

	p.typeParser.typeAliases = typeAliases

	// 4. Второй проход: парсим поля структур
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

			// Парсим поля
			for _, field := range structType.Fields.List {
				fieldName := getFieldName(field)
				if fieldName == "" {
					continue
				}

				fieldType := p.typeParser.ParseExpr(field.Type, allStructs)
				directives := p.annotationParser.ParseFieldDirectives(field)

				// Валидация директив
				for _, dir := range directives {
					if err := ValidateDirective(dir, fieldType); err != nil {
						result.Errors = append(result.Errors, fmt.Errorf("%s:%d: поле %s: %w",
							filename, p.fset.Position(field.Pos()).Line, fieldName, err))
					}
				}

				// Парсим OA-аннотации и обрабатываем @oa-rewrite: ref
				oaAnnotations := p.annotationParser.ParseFieldOaAnnotations(field)
				var rewriteRef string
				remainingOa := make([]OaAnnotation, 0, len(oaAnnotations))
				for _, ann := range oaAnnotations {
					if ann.Type == "rewrite.ref" {
						rewriteRef = trimQuotes(ann.Value)
					} else {
						remainingOa = append(remainingOa, ann)
					}
				}

				s.Fields = append(s.Fields, Field{
					Name:          fieldName,
					Type:          fieldType,
					Directives:    directives,
					Decorators:    p.parseFieldDecorators(field), // можно тоже вынести
					Line:          p.fset.Position(field.Pos()).Line,
					OaAnnotations: remainingOa,
					IsEmbedded:    len(field.Names) == 0,
					OaRewriteRef:  rewriteRef,
				})
			}

			// Парсим аннотации структуры
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

			result.Structs = append(result.Structs, *s)
		}
	}

	if p.verbose {
		log.Printf("DEBUG: Parsed %s, node.Name.Name=%q", filename, node.Name.Name)
	}

	return result, nil
}
