package generator

import (
	"embed"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/arkannsk/elval/pkg/parser"
)

//go:embed templates
var templatesFS embed.FS

type Generator struct {
	outputDir                string
	tmpl                     *template.Template
	generateOpenAPI, verbose bool
}

// Маппинг примитивов (включая ваши float64 и duration)
var primitives = map[string]bool{
	"string": true, "int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"float32": true, "float64": true,
	"bool":      true,
	"time.Time": true, "time.Duration": true,
}

func NewGenerator(outputDir string, generateOpenAPI, verbose bool) (*Generator, error) {
	tmpl := template.New("").Funcs(templateFucMap)

	// Загружаем все шаблоны рекурсивно
	err := fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".tmpl") {
			content, err := fs.ReadFile(templatesFS, path)
			if err != nil {
				return err
			}
			// Используем путь как имя шаблона (без расширения .tmpl)
			name := strings.TrimPrefix(path, "templates/")
			name = strings.TrimSuffix(name, ".tmpl")

			_, err = tmpl.New(name).Parse(string(content))
			if err != nil {
				return fmt.Errorf("ошибка парсинга %s: %w", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки шаблонов: %w", err)
	}

	return &Generator{
		outputDir:       outputDir,
		tmpl:            tmpl,
		generateOpenAPI: generateOpenAPI,
		verbose:         verbose,
	}, nil
}

func (g *Generator) Generate(parseResult *parser.ParseResult, sourceFile string) error {
	baseName := strings.TrimSuffix(filepath.Base(sourceFile), ".go")

	// Отбираем структуры для валидации (есть директивы)
	structsForValidation := make([]parser.Struct, 0)
	for _, s := range parseResult.Structs {
		if s.HasDirectives() {
			structsForValidation = append(structsForValidation, s)
		}
	}

	// Отбираем структуры для OpenAPI
	structsForOpenAPI := make([]parser.Struct, 0)
	for _, s := range parseResult.Structs {
		if s.ShouldGenerateOpenAPI(g.generateOpenAPI) {
			structsForOpenAPI = append(structsForOpenAPI, s)
		}
	}

	// 1. Генерируем файл валидации
	if len(structsForValidation) > 0 {
		data := struct {
			Package            string
			Structs            []parser.Struct
			SourceFile         string
			GenerateValidation bool
			Imports            map[string]string
		}{
			Package:            parseResult.Package,
			Structs:            structsForValidation,
			SourceFile:         filepath.Base(sourceFile),
			GenerateValidation: true,
			Imports:            parser.CollectValidationImports(parseResult.Structs),
		}

		if g.verbose {
			log.Printf("import list: %v in file: %s", parseResult.Imports, sourceFile)
		}

		var buf strings.Builder
		if err := g.tmpl.ExecuteTemplate(&buf, "validation/validation", data); err != nil {
			return fmt.Errorf("ошибка выполнения шаблона валидации: %w", err)
		}

		formatted, err := format.Source([]byte(buf.String()))
		if err != nil {
			debugFile := strings.TrimSuffix(sourceFile, ".go") + ".debug.go"
			_ = os.WriteFile(debugFile, []byte(buf.String()), 0644)
			return fmt.Errorf("ошибка форматирования: %w", err)
		}

		outputPath := filepath.Join(g.outputDir, baseName+".gen.go")
		if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
			return err
		}
	}

	// 2. Генерируем OpenAPI файл
	if g.generateOpenAPI && len(structsForOpenAPI) > 0 {
		data := struct {
			Package         string
			Structs         []parser.Struct
			SourceFile      string
			GenerateOpenAPI bool
			Imports         map[string]string
		}{
			Package:         parseResult.Package,
			Structs:         structsForOpenAPI,
			SourceFile:      filepath.Base(sourceFile),
			GenerateOpenAPI: true,
			Imports:         parser.CollectOpenAPIImports(parseResult.Structs),
		}

		var buf strings.Builder
		if err := g.tmpl.ExecuteTemplate(&buf, "openapi/openapi", data); err != nil {
			return fmt.Errorf("ошибка выполнения шаблона OpenAPI: %w", err)
		}

		formatted, err := format.Source([]byte(buf.String()))
		if err != nil {
			debugFile := strings.TrimSuffix(sourceFile, ".go") + ".openapi.debug.go"
			_ = os.WriteFile(debugFile, []byte(buf.String()), 0644)
			return fmt.Errorf("ошибка форматирования OpenAPI: %w", err)
		}

		outputPath := filepath.Join(g.outputDir, baseName+".oa.gen.go")
		if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
			return err
		}
	}

	return nil
}
