package generator

import (
	"embed"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/arkannsk/elval/pkg/parser"
)

//go:embed templates
var templatesFS embed.FS

type Generator struct {
	outputDir       string
	tmpl            *template.Template
	generateOpenAPI bool
}

func NewGenerator(outputDir string, generateOpenAPI bool) (*Generator, error) {
	tmpl := template.New("").
		Funcs(template.FuncMap{
			"hasTime": func(structs []parser.Struct) bool {
				for _, s := range structs {
					for _, f := range s.Fields {
						if f.Type.Name == "time.Time" || f.Type.Name == "time.Duration" {
							return true
						}
					}
				}
				return false
			},
			"dict": func(values ...interface{}) (map[string]interface{}, error) {
				if len(values)%2 != 0 {
					return nil, fmt.Errorf("invalid dict call")
				}
				dict := make(map[string]interface{}, len(values)/2)
				for i := 0; i < len(values); i += 2 {
					key, ok := values[i].(string)
					if !ok {
						return nil, fmt.Errorf("dict keys must be strings")
					}
					dict[key] = values[i+1]
				}
				return dict, nil
			},
			"hasOptional": func(directives []parser.Directive) bool {
				for _, d := range directives {
					if d.Type == "optional" {
						return true
					}
				}
				return false
			},
			"itoa": func(i int) string {
				return strconv.Itoa(i)
			},
			"hasCustomDirective": func(directives []parser.Directive) bool {
				for _, d := range directives {
					if strings.HasPrefix(d.Type, "x-") {
						return true
					}
				}
				return false
			},
			"isCustomDirective": func(dirType string) bool {
				return strings.HasPrefix(dirType, "x-")

			},
			"hasHTTPContext": func(structs []parser.Struct) bool {
				for _, s := range structs {
					for _, f := range s.Fields {
						for _, d := range f.Decorators {
							if d.Type == "httpctx-get" {
								return true
							}
						}
					}
				}
				return false
			},
			"hasEnvGet": func(structs []parser.Struct) bool {
				for _, s := range structs {
					for _, f := range s.Fields {
						for _, d := range f.Decorators {
							if d.Type == "env-get" {
								return true
							}
						}
					}
				}
				return false
			},
			"hasUUIDGen": func(structs []parser.Struct) bool {
				for _, s := range structs {
					for _, f := range s.Fields {
						for _, d := range f.Decorators {
							if d.Type == "uuid-gen" {
								return true
							}
						}
					}
				}
				return false
			},
			"hasOaSchema": func(structs []parser.Struct) bool {
				// Проверяем, есть ли у структур OaAnnotations
				for _, s := range structs {
					for _, f := range s.Fields {
						if len(f.OaAnnotations) > 0 {
							return true
						}
					}
				}
				return false
			},
			"toLower": strings.ToLower,
		})

	// Рекурсивно обходим все файлы шаблонов
	err := fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".tmpl") {
			content, err := fs.ReadFile(templatesFS, path)
			if err != nil {
				return err
			}
			// Извлекаем короткое имя (без пути и расширения)
			// templates/header.tmpl -> header
			// templates/field/slice.tmpl -> slice_validation
			name := strings.TrimPrefix(path, "templates/")
			name = strings.TrimSuffix(name, ".tmpl")
			name = strings.ReplaceAll(name, "/", "_")

			// Для field файлов добавляем суффикс для ясности
			if strings.Contains(path, "/field/") {
				name = strings.TrimPrefix(name, "field_")
			}

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
	}, nil
}

func (g *Generator) Generate(parseResult *parser.ParseResult, sourceFile string) error {
	if len(parseResult.Structs) == 0 {
		// Нет структур с аннотациями - не генерируем файл
		return nil
	}

	data := struct {
		Package         string
		Structs         []parser.Struct
		SourceFile      string
		GenerateOpenAPI bool
	}{
		Package:         parseResult.Package,
		Structs:         parseResult.Structs,
		SourceFile:      filepath.Base(sourceFile),
		GenerateOpenAPI: g.generateOpenAPI,
	}

	var buf strings.Builder
	if err := g.tmpl.ExecuteTemplate(&buf, "validator", data); err != nil {
		return fmt.Errorf("ошибка выполнения шаблона: %w", err)
	}

	// Форматируем код
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		// Сохраняем неформатированный для отладки
		debugFile := strings.TrimSuffix(sourceFile, ".go") + ".debug.go"
		_ = os.WriteFile(debugFile, []byte(buf.String()), 0644)
		return fmt.Errorf("ошибка форматирования: %w", err)
	}

	// Сохраняем результат
	outputPath := filepath.Join(g.outputDir, strings.TrimSuffix(filepath.Base(sourceFile), ".go")+".gen.go")
	return os.WriteFile(outputPath, formatted, 0644)
}
