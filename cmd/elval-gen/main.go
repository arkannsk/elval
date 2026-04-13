package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/arkannsk/elval/pkg/generator"
	"github.com/arkannsk/elval/pkg/parser"
)

func main() {
	var inputDir string
	var outputDir string
	var verbose bool

	flag.StringVar(&inputDir, "input", ".", "директория с исходными .go файлами")
	flag.StringVar(&outputDir, "output", "", "директория для сгенерированных файлов (по умолчанию та же, что и input)")
	flag.BoolVar(&verbose, "v", false, "подробный вывод")
	flag.Parse()

	if outputDir == "" {
		outputDir = inputDir
	}

	p := parser.NewParser()
	gen, err := generator.NewGenerator(outputDir)
	if err != nil {
		log.Fatal(err)
	}

	// Ищем все .go файлы
	files, err := filepath.Glob(filepath.Join(inputDir, "*.go"))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file, ".gen.go") {
			continue
		}

		if verbose {
			fmt.Printf("Обработка: %s\n", file)
		}

		result, err := p.ParseFile(file)
		if err != nil {
			log.Printf("Ошибка парсинга %s: %v", file, err)
			continue
		}

		if len(result.Errors) > 0 {
			for _, err := range result.Errors {
				log.Printf("Ошибка в %s: %v", file, err)
			}
		}

		if err := gen.Generate(result, file); err != nil {
			log.Printf("Ошибка генерации для %s: %v", file, err)
			continue
		}

		if verbose {
			fmt.Printf("  ✅ Сгенерирован %s\n", strings.TrimSuffix(filepath.Base(file), ".go")+".gen.go")
		}
	}
}
