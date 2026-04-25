package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/arkannsk/elval/pkg/parser"
)

func lintCmd() {
	var inputDir string
	var verbose bool
	var exclude string
	var warningsAsErrors bool

	lintFlags := flag.NewFlagSet("lint", flag.ExitOnError)
	// Исправлено: уникальные имена флагов
	lintFlags.StringVar(&inputDir, "input", ".", "directory to lint")
	lintFlags.StringVar(&inputDir, "i", ".", "alias for -input")
	lintFlags.BoolVar(&verbose, "v", false, "verbose output")
	lintFlags.StringVar(&exclude, "exclude", "", "comma-separated patterns to exclude (e.g. vendor,testdata)")
	lintFlags.BoolVar(&warningsAsErrors, "Werror", false, "treat warnings as errors")

	lintFlags.Parse(os.Args[2:])

	p := parser.NewParser(verbose)

	// Рекурсивно обходим все поддиректории
	var files []string
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Проверяем исключения
			if exclude != "" {
				patterns := strings.Split(exclude, ",")
				for _, pattern := range patterns {
					if strings.Contains(path, strings.TrimSpace(pattern)) {
						return filepath.SkipDir
					}
				}
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") &&
			!strings.HasSuffix(path, ".gen.go") &&
			!strings.HasSuffix(path, ".oa.gen.go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	hasErrors := false
	totalFiles := 0
	errorCount := 0
	warningCount := 0

	allStructsMap := make(map[string]*parser.Struct) // GlobalRef -> Struct

	var allResults []*parser.ParseResult

	for _, file := range files {
		result, err := p.ParseFile(file)
		if err != nil {
			log.Printf("Parse error %s: %v", file, err)
			hasErrors = true
			continue
		}
		allResults = append(allResults, result)

		// Индексируем структуры
		for i := range result.Structs {
			s := &result.Structs[i]
			ref := buildGlobalRefForLint(s)
			allStructsMap[ref] = s
		}
	}

	// Шаг 2: Линтинг
	for _, result := range allResults {
		totalFiles++
		filePath := "unknown_file"
		if len(result.Structs) > 0 {
			filePath = result.Structs[0].File
		} else if len(result.Errors) > 0 {
			log.Printf("parser result errors: %v", result.Errors)
		}

		// 1. Валидация директив (@evl:...)
		directiveErrors := result.ValidateDirectives()
		for _, err := range directiveErrors {
			if err.Severity == parser.SeverityError {
				errorCount++
				hasErrors = true
				fmt.Println(err.Error())
			} else {
				warningCount++
				if verbose {
					fmt.Println(err.Error())
				}
				if warningsAsErrors {
					hasErrors = true
				}
			}
		}

		// 2. Валидация OpenAPI аннотаций (@oa:...)
		oaErrors := validateOpenAPIAnnotations(result, allStructsMap, filePath)
		for _, err := range oaErrors {
			if err.Severity == parser.SeverityError {
				errorCount++
				hasErrors = true
				fmt.Println(err.Error())
			} else {
				warningCount++
				if verbose {
					fmt.Println(err.Error())
				}
				if warningsAsErrors {
					hasErrors = true
				}
			}
		}
	}

	if verbose {
		fmt.Printf("\nFiles checked: %d, errors found: %d, warnings: %d\n", totalFiles, errorCount, warningCount)
	}

	if hasErrors {
		os.Exit(1)
	}
}

func validateOpenAPIAnnotations(result *parser.ParseResult, structsMap map[string]*parser.Struct, filePath string) []parser.DirectiveError {
	var errors []parser.DirectiveError

	// Допустимые типы для rewrite.type
	validOaTypes := map[string]bool{
		"string": true, "integer": true, "number": true, "boolean": true, "array": true,
	}

	for _, s := range result.Structs {
		// Проверка уровня структуры (Discriminator)
		if s.Discriminator != nil {
			// Проверяем, что все типы из oneOf/anyOf есть в маппинге
			allMappedTypes := make(map[string]bool)
			for t := range s.Discriminator.Mapping {
				allMappedTypes[t] = true
			}

			checkOneOfRefs := func(refs []string, typeName string) {
				for _, refName := range refs {
					// Если refName не найден в маппинге, это ошибка
					if !allMappedTypes[refName] {
						errors = append(errors, parser.DirectiveError{
							File:      filePath,
							Line:      0, // Линия структуры неизвестна без AST позиции, можно улучшить
							Struct:    s.Name,
							Field:     "",
							Directive: "discriminator",
							Message:   fmt.Sprintf("type %q from %s is not in discriminator mapping", refName, typeName),
							Severity:  parser.SeverityError,
						})
					}
				}
			}

			checkOneOfRefs(s.OaOneOf, "oneOf")
			checkOneOfRefs(s.OaAnyOf, "anyOf")
		}

		// Проверка полей
		for _, f := range s.Fields {
			// 1. Проверка @oa:rewrite.ref
			if f.OaRewriteRef != "" {
				// Проверяем, существует ли такая структура в нашем проекте
				if _, exists := structsMap[f.OaRewriteRef]; !exists {
					errors = append(errors, parser.DirectiveError{
						File:      filePath,
						Line:      f.Line,
						Struct:    s.Name,
						Field:     f.Name,
						Directive: "rewrite.ref",
						Message:   fmt.Sprintf("target schema %q does not exist in project", f.OaRewriteRef),
						Severity:  parser.SeverityError,
					})
				}
			}

			// 2. Проверка @oa:rewrite.type
			if f.OaRewriteType != "" {
				if !validOaTypes[f.OaRewriteType] {
					errors = append(errors, parser.DirectiveError{
						File:      filePath,
						Line:      f.Line,
						Struct:    s.Name,
						Field:     f.Name,
						Directive: "rewrite.type",
						Message:   fmt.Sprintf("invalid OpenAPI type %q. Allowed: string, integer, number, boolean, array", f.OaRewriteType),
						Severity:  parser.SeverityError,
					})
				}
			}

			// 3. Проверка @oa:oneOf-ref и @oa:anyOf-ref на уровне поля
			checkFieldRefs := func(refs []string, directiveName string) {
				for _, ref := range refs {
					// Убираем префикс #/components/schemas/ если есть
					cleanRef := strings.TrimPrefix(ref, "#/components/schemas/")

					// Пытаемся найти структуру по cleanRef
					found := false
					for structRef := range structsMap {
						if structRef == cleanRef || strings.HasSuffix(structRef, "."+cleanRef) {
							found = true
							break
						}
					}

					if !found {
						errors = append(errors, parser.DirectiveError{
							File:      filePath,
							Line:      f.Line,
							Struct:    s.Name,
							Field:     f.Name,
							Directive: directiveName,
							Message:   fmt.Sprintf("referenced schema %q does not exist", ref),
							Severity:  parser.SeverityWarning, // Warning, так как может быть внешний реф
						})
					}
				}
			}

			checkFieldRefs(f.OaOneOfRefs, "oneOf-ref")
			checkFieldRefs(f.OaAnyOfRefs, "anyOf-ref")
		}
	}

	return errors
}

func buildGlobalRefForLint(s *parser.Struct) string {
	ref := s.Name
	if s.PackagePath != "" {
		ref = s.PackagePath + "." + s.Name
		if s.Module != "" {
			ref = s.Module + "/" + ref
		}
	} else if s.Package != "" {
		ref = s.Package + "." + s.Name
	}
	return ref
}
