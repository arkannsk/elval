package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/arkannsk/elval/pkg/generator"
	"github.com/arkannsk/elval/pkg/parser"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate", "gen":
		generateCmd()
	case "lint":
		lintCmd()
	case "version", "ver":
		versionCmd()
	case "help", "h", "":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`elval-gen - Code generator for elval validator

Usage:
  elval-gen generate [flags]   Generate validation code
  elval-gen lint [flags]       Validate annotations without generation
  elval-gen version            Show version

Flags for generate and lint:
  -input, -i string      Input directory with .go files (default ".")
  -output, -o string     Output directory (default same as input)
  -v                     Verbose output

Additional flags for generate:
  -openapi               Generate OpenAPI schemas

Examples:
  elval-gen generate -i ./user
  elval-gen gen -i ./user -openapi -v
  elval-gen lint -i ./user
  elval-gen lint -i ./user -v`)
}

func generateCmd() {
	var inputDir string
	var outputDir string
	var verbose bool
	var generateOpenAPI bool

	genFlags := flag.NewFlagSet("generate", flag.ExitOnError)
	// Используем разные переменные для кратких и полных флагов, чтобы не было конфликтов при парсинге
	genFlags.StringVar(&inputDir, "input", ".", "Input directory")
	genFlags.StringVar(&inputDir, "i", ".", "Short for -input")

	genFlags.StringVar(&outputDir, "output", "", "Output directory (default: same as input)")
	genFlags.StringVar(&outputDir, "o", "", "Short for -output")

	genFlags.BoolVar(&verbose, "v", false, "Verbose output")
	genFlags.BoolVar(&verbose, "verbose", false, "Verbose output") // Алиас

	genFlags.BoolVar(&generateOpenAPI, "openapi", false, "Generate OpenAPI schemas")

	genFlags.Parse(os.Args[2:])

	if outputDir == "" {
		outputDir = inputDir
	}

	// Нормализуем пути
	inputDir, _ = filepath.Abs(inputDir)
	outputDir, _ = filepath.Abs(outputDir)

	if verbose {
		fmt.Printf("Generating code for %s\n", inputDir)
		if generateOpenAPI {
			fmt.Printf("OpenAPI schemas enabled\n")
		}
	}

	p := parser.NewParser(verbose)
	gen, err := generator.NewGenerator(outputDir, generateOpenAPI, verbose)
	if err != nil {
		log.Fatal(err)
	}

	generated := 0
	skipped := 0
	errors := 0

	// Рекурсивный обход директорий
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории
		if info.IsDir() {
			return nil
		}

		name := info.Name()

		// ИГНОРИРУЕМ сгенерированные файлы, чтобы избежать person.gen.gen.go
		if strings.HasSuffix(name, ".gen.go") ||
			strings.HasSuffix(name, ".oa.gen.go") ||
			strings.HasSuffix(name, "_test.go") {
			return nil
		}

		// Обрабатываем только исходные .go файлы
		if !strings.HasSuffix(name, ".go") {
			return nil
		}

		if verbose {
			fmt.Printf("  Processing: %s\n", path)
		}

		result, err := p.ParseFile(path)
		if err != nil {
			log.Printf("Parse error %s: %v", path, err)
			errors++
			return nil // Продолжаем обработку других файлов
		}

		if len(result.Structs) == 0 {
			if verbose {
				fmt.Printf("    Skip (no structs with annotations)\n")
			}
			skipped++
			return nil
		}

		// Определяем относительный путь для сохранения структуры папок в outputDir
		relPath, err := filepath.Rel(inputDir, filepath.Dir(path))
		if err != nil {
			relPath = "."
		}

		targetDir := filepath.Join(outputDir, relPath)

		// Создаем целевую директорию, если её нет
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			log.Printf("Error creating dir %s: %v", targetDir, err)
			return nil
		}

		// Генерируем файл
		baseName := strings.TrimSuffix(filepath.Base(path), ".go")
		outFilePath := filepath.Join(targetDir, baseName+".gen.go")

		if err := gen.Generate(result, outFilePath); err != nil {
			log.Printf("Generation error for %s: %v", path, err)
			errors++
			return nil
		}

		generated++
		if verbose {
			fmt.Printf("    Generated %s\n", outFilePath)
			if generateOpenAPI {
				oaOutPath := filepath.Join(targetDir, baseName+".oa.gen.go")
				fmt.Printf("    Generated %s\n", oaOutPath)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Walk error: %v", err)
	}

	if verbose {
		fmt.Printf("\nStatistics: generated %d, skipped %d, errors %d\n", generated, skipped, errors)
	}
}

func versionCmd() {
	fmt.Println("elval-gen version 0.1.0")
}
