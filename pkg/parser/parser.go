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

func (p *Parser) ParseFile(filename string) (*ParseResult, error) {
	absFile, err := filepath.Abs(filename)
	if err != nil {
		absFile = filename
	}

	node, err := parser.ParseFile(p.fset, absFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга файла %s: %w", absFile, err)
	}

	modInfo, err := resolveModuleInfo(absFile)
	if err != nil {
		return nil, err
	}

	p.typeParser = NewTypeParser(make(map[string]string), modInfo.Package, p.verbose)

	result := &ParseResult{
		Package: node.Name.Name,
		Structs: make([]Struct, 0),
		Errors:  make([]error, 0),
	}

	if p.verbose {
		log.Printf("DEBUG: File=%s | Module=%s | Root=%s | PkgPath=%s",
			filename, modInfo.Module, modInfo.ModuleRoot, modInfo.PackagePath)
	}

	allStructs, typeAliases := p.parseStructsFirstPass(node, filename, modInfo)
	p.typeParser.typeAliases = typeAliases

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

			if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
				// Парсим аннотации один раз и сохраняем в структуру для повторного использования
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
					continue
				}

				allStructs[name] = &Struct{
					Name:             name,
					File:             filename,
					Fields:           []Field{},
					Package:          modInfo.Package,
					PackagePath:      modInfo.PackagePath,
					Module:           modInfo.Module,
					IsIgnored:        false,
					RawOaAnnotations: structOaAnnotations, // Сохраняем сырые аннотации
				}
			} else {
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
	_ *ModuleInfo) {
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

			s := allStructs[typeSpec.Name.Name]
			if s == nil {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			for _, field := range structType.Fields.List {
				fieldName := getFieldName(field)

				fieldType := p.typeParser.ParseExpr(field.Type, allStructs)

				if p.verbose {
					log.Printf(`PARSED: fieldName: %s, exprType: %T, val: %s → IsGeneric=%v, 
						Base=%q, Args=%+v, IsSlice=%v`,
						fieldName, field.Type, exprToString(field.Type), fieldType.IsGeneric,
						fieldType.GenericBase, fieldType.GenericArgs, fieldType.IsSlice)
				}

				directives := p.annotationParser.ParseFieldDirectives(field)

				for _, dir := range directives {
					if err := ValidateDirective(dir, fieldType); err != nil {
						result.Errors = append(result.Errors, fmt.Errorf("%s:%d: поле %s: %w",
							filename, p.fset.Position(field.Pos()).Line, fieldName, err))
					}
				}

				rawOaAnnotations := p.annotationParser.ParseFieldOaAnnotations(field)
				fAnot := p.annotationParser.ProcessFieldAnnotations(rawOaAnnotations)

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
					IsIgnored:     false,
					OaIn:          fAnot.OaIn,
					OaParamName:   fAnot.OaParamName,
				})
			}

			// Используем сохраненные сырые аннотации структуры
			p.annotationParser.ExtractDiscriminator(s, s.RawOaAnnotations)

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
}
