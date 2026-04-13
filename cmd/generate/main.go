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
	var inputFile string
	var outputDir string

	flag.StringVar(&inputFile, "file", "", "исходный .go файл")
	flag.StringVar(&outputDir, "out", ".", "директория для вывода")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("укажите файл для генерации: -file")
	}

	p := parser.NewParser()
	gen, err := generator.NewGenerator(outputDir)
	if err != nil {
		log.Fatal(err)
	}

	result, err := p.ParseFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := gen.Generate(result, inputFile); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Сгенерирован %s\n", strings.TrimSuffix(filepath.Base(inputFile), ".go")+".gen.go")
}
