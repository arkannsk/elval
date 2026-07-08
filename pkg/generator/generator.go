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
	basePath := strings.TrimSuffix(sourceFile, ".go")

	// Отбираем структуры для валидации
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

	if g.verbose {
		log.Printf("DEBUG: Generating validation for %s, Package=%q", sourceFile, parseResult.Package)
	}

	// 1. Генерация файла валидации
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
			SourceFile:         filepath.Base(sourceFile), // имя файла для шаблона
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
			debugFile := basePath + ".debug.go"
			_ = os.WriteFile(debugFile, []byte(buf.String()), 0644)
			return fmt.Errorf("ошибка форматирования: %w (debug: %s)", err, debugFile)
		}

		outputPath := basePath + ".gen.go"

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
			return err
		}
	}

	// 2. Генерация OpenAPI файла
	if g.generateOpenAPI && len(structsForOpenAPI) > 0 {
		// Заполняем DiscriminatorPropertyName для дочерних структур
		g.propagateDiscriminatorPropertyName(structsForOpenAPI)

		// Автоматически добавляем недостающие media types на основе struct тегов
		g.normalizeResponseMediaTypes(structsForOpenAPI)

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
			debugFile := basePath + ".openapi.debug.go"
			_ = os.WriteFile(debugFile, []byte(buf.String()), 0644)
			return fmt.Errorf("ошибка форматирования OpenAPI: %w (debug: %s)", err, debugFile)
		}

		// outputPath для OpenAPI
		outputPath := basePath + ".oa.gen.go"

		// Создаём директорию перед записью
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
			return err
		}
	}

	return nil
}

// propagateDiscriminatorPropertyName проходит по структурам, находит родителей с дискриминатором
// и устанавливает DiscriminatorPropertyName для дочерних структур, указанных в маппинге.
func (g *Generator) propagateDiscriminatorPropertyName(structs []parser.Struct) {
	// Строим индекс: имя структуры -> указатель на структуру
	index := make(map[string]*parser.Struct, len(structs))
	for i := range structs {
		index[structs[i].Name] = &structs[i]
	}

	// Для каждой структуры с дискриминатором заполняем дочерние
	for i := range structs {
		s := &structs[i]
		if s.Discriminator == nil || len(s.Discriminator.Mapping) == 0 {
			continue
		}
		propName := s.Discriminator.PropertyName
		for _, childName := range s.Discriminator.Mapping {
			child, ok := index[childName]
			if ok && child.DiscriminatorPropertyName == "" {
				child.DiscriminatorPropertyName = propName
			}
		}
	}
}

// normalizeResponseMediaTypes автоматически добавляет недостающие media types
// в @oa:response аннотации на основе наличия xml/yaml тегов в полях структуры.
// Если у структуры есть xml теги, но в @oa:response нет application/xml — добавляется.
// Если есть yaml теги, но нет application/x-yaml — добавляется.
func (g *Generator) normalizeResponseMediaTypes(structs []parser.Struct) {
	for i := range structs {
		s := &structs[i]
		if len(s.OaResponses) == 0 {
			continue
		}

		// Собираем все media types из всех response аннотаций
		existingMediaTypes := make(map[string]bool)
		for _, resp := range s.OaResponses {
			for _, mt := range resp.MediaTypes {
				existingMediaTypes[mt] = true
			}
		}

		// Если есть xml теги, но нет application/xml — добавляем
		if s.HasXmlTags && !existingMediaTypes["application/xml"] {
			for j := range s.OaResponses {
				s.OaResponses[j].MediaTypes = append(s.OaResponses[j].MediaTypes, "application/xml")
			}
		}

		// Если есть yaml теги, но нет text/yaml — добавляем (RFC 9512)
		if s.HasYamlTags && !existingMediaTypes["text/yaml"] {
			for j := range s.OaResponses {
				s.OaResponses[j].MediaTypes = append(s.OaResponses[j].MediaTypes, "text/yaml")
			}
		}
	}
}
