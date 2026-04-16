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
	genFlags.StringVar(&inputDir, "input", ".", "")
	genFlags.StringVar(&inputDir, "i", ".", "")
	genFlags.StringVar(&outputDir, "output", "", "")
	genFlags.StringVar(&outputDir, "o", "", "")
	genFlags.BoolVar(&verbose, "v", false, "")
	genFlags.BoolVar(&generateOpenAPI, "openapi", false, "")

	genFlags.Parse(os.Args[2:])

	if outputDir == "" {
		outputDir = inputDir
	}

	if verbose {
		fmt.Printf("Generating code for %s\n", inputDir)
		if generateOpenAPI {
			fmt.Printf("OpenAPI schemas enabled\n")
		}
	}

	p := parser.NewParser()
	gen, err := generator.NewGenerator(outputDir, generateOpenAPI)
	if err != nil {
		log.Fatal(err)
	}

	files, err := filepath.Glob(filepath.Join(inputDir, "*.go"))
	if err != nil {
		log.Fatal(err)
	}

	generated := 0
	skipped := 0

	for _, file := range files {
		if strings.HasSuffix(file, ".gen.go") || strings.HasSuffix(file, ".oa.gen.go") {
			continue
		}

		if verbose {
			fmt.Printf("  Processing: %s\n", filepath.Base(file))
		}

		result, err := p.ParseFile(file)
		if err != nil {
			log.Printf("Parse error %s: %v", file, err)
			continue
		}

		if len(result.Structs) == 0 {
			if verbose {
				fmt.Printf("    Skip (no structs with annotations)\n")
			}
			skipped++
			continue
		}

		if err := gen.Generate(result, file); err != nil {
			log.Printf("Generation error for %s: %v", file, err)
			continue
		}

		generated++
		if verbose {
			fmt.Printf("    Generated %s\n", strings.TrimSuffix(filepath.Base(file), ".go")+".gen.go")
			if generateOpenAPI {
				fmt.Printf("    Generated %s\n", strings.TrimSuffix(filepath.Base(file), ".go")+".oa.gen.go")
			}
		}
	}

	if verbose {
		fmt.Printf("\nStatistics: generated %d, skipped %d\n", generated, skipped)
	}
}

// cmd/elval-gen/main.go

func lintCmd() {
	var inputDir string
	var verbose bool
	var exclude string

	lintFlags := flag.NewFlagSet("lint", flag.ExitOnError)
	lintFlags.StringVar(&inputDir, "input", ".", "")
	lintFlags.StringVar(&inputDir, "i", ".", "")
	lintFlags.BoolVar(&verbose, "v", false, "")
	lintFlags.StringVar(&exclude, "exclude", "", "comma-separated patterns to exclude (e.g. vendor,testdata)")

	lintFlags.Parse(os.Args[2:])

	p := parser.NewParser()

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

	for _, file := range files {
		totalFiles++

		result, err := p.ParseFile(file)
		if err != nil {
			log.Printf("Parse error %s: %v", file, err)
			hasErrors = true
			continue
		}

		errors := result.ValidateDirectives()
		if len(errors) > 0 {
			hasErrors = true
			errorCount += len(errors)
			for _, err := range errors {
				fmt.Println(err.Error())
			}
		} else if verbose {
			fmt.Printf("%s - all annotations valid\n", filepath.Base(file))
		}
	}

	if verbose {
		fmt.Printf("\nFiles checked: %d, errors found: %d\n", totalFiles, errorCount)
	}

	if hasErrors {
		os.Exit(1)
	}
}

func versionCmd() {
	fmt.Println("elval-gen version 0.1.0")
}
