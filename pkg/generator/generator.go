package generator

import (
	"embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/arkannsk/elval/pkg/parser"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

type Generator struct {
	outputDir string
	tmpl      *template.Template
}

func NewGenerator(outputDir string) (*Generator, error) {
	tmpl := template.New("").
		Funcs(template.FuncMap{
			"hasTime": func(structs []parser.Struct) bool {
				for _, s := range structs {
					for _, f := range s.Fields {
						name := f.Type.Name
						if name == "time.Time" || name == "time.Duration" {
							return true
						}
					}
				}
				return false
			},
			"hasPattern": func(structs []parser.Struct) bool {
				// regexp не нужен, так как MatchRegexp создает regexp внутри
				return false
			},
			"printf": fmt.Sprintf,
		})

	tmpl, err := tmpl.ParseFS(templatesFS, "templates/validator.tmpl")
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга шаблона: %w", err)
	}

	return &Generator{
		outputDir: outputDir,
		tmpl:      tmpl,
	}, nil
}

func (g *Generator) Generate(parseResult *parser.ParseResult, sourceFile string) error {
	if len(parseResult.Structs) == 0 {
		return nil
	}

	data := struct {
		Package    string
		Structs    []parser.Struct
		SourceFile string
	}{
		Package:    parseResult.Package,
		Structs:    parseResult.Structs,
		SourceFile: filepath.Base(sourceFile),
	}

	var buf strings.Builder
	if err := g.tmpl.ExecuteTemplate(&buf, "validator.tmpl", data); err != nil {
		return fmt.Errorf("ошибка выполнения шаблона: %w", err)
	}

	// Форматируем код
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		// Сохраняем неформатированный для отладки
		debugFile := strings.TrimSuffix(sourceFile, ".go") + ".debug.go"
		os.WriteFile(debugFile, []byte(buf.String()), 0644)
		return fmt.Errorf("ошибка форматирования: %w", err)
	}

	// Сохраняем результат
	outputPath := filepath.Join(g.outputDir, strings.TrimSuffix(filepath.Base(sourceFile), ".go")+".gen.go")
	return os.WriteFile(outputPath, formatted, 0644)
}
