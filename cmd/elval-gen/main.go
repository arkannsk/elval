package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/arkannsk/elval/pkg/errs"
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
	var exclude string
	var generateOpenAPI bool
	var warningsAsErrors bool
	var noColor bool

	genFlags := flag.NewFlagSet("generate", flag.ExitOnError)
	genFlags.StringVar(&inputDir, "input", ".", "input directory")
	genFlags.StringVar(&inputDir, "i", ".", "alias for -input")
	genFlags.StringVar(&exclude, "exclude", "", "comma-separated patterns to exclude (e.g. vendor,testdata)")
	genFlags.BoolVar(&verbose, "v", false, "verbose output")
	genFlags.BoolVar(&generateOpenAPI, "openapi", false, "generate OpenAPI schemas")
	genFlags.BoolVar(&warningsAsErrors, "Werror", false, "treat warnings as errors")
	genFlags.BoolVar(&noColor, "no-color", false, "disable colored output")

	genFlags.Parse(os.Args[2:])

	colors := errs.DefaultColors()
	if noColor {
		colors = errs.Color{}
	}

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

	var files []string
	err = filepath.WalkDir(inputDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if exclude != "" {
				patterns := strings.Split(exclude, ",")
				for _, pattern := range patterns {
					if strings.Contains(path, strings.TrimSpace(pattern)) {
						return filepath.SkipDir
					}
				}
			}
			name := d.Name()
			if name == "vendor" || name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		// Берём только .go файлы, исключая сгенерированные
		if strings.HasSuffix(path, ".go") &&
			!strings.HasSuffix(path, "_test.go") &&
			!strings.HasSuffix(path, ".gen.go") &&
			!strings.HasSuffix(path, ".oa.gen.go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if verbose {
		fmt.Printf("Found %d Go files to process\n", len(files))
	}

	generated := 0
	skipped := 0
	totalErrors := 0
	totalWarnings := 0

	for _, file := range files {
		if verbose {
			// Показываем относительный путь для читаемости
			rel, _ := filepath.Rel(inputDir, file)
			fmt.Printf("  Processing: %s\n", rel)
		}

		result, err := p.ParseFile(file)
		if err != nil {
			log.Printf("Parse error %s: %v", file, err)
			continue
		}

		// Обработка диагностик
		fileHasErrors := false
		fileHasWarnings := false

		for _, d := range result.Diagnostics {
			fmt.Fprintln(os.Stderr, errs.FormatDiagnostic(d, colors))

			switch d.Severity {
			case errs.SeverityError:
				totalErrors++
				fileHasErrors = true
			case errs.SeverityWarning:
				totalWarnings++
				fileHasWarnings = true
			}
		}

		if warningsAsErrors && fileHasWarnings {
			fileHasErrors = true
		}

		if fileHasErrors {
			if verbose {
				fmt.Printf("    Skipped due to validation errors\n")
			}
			skipped++
			continue
		}

		if len(result.Structs) == 0 {
			if verbose {
				fmt.Printf("    Skip (no structs with annotations)\n")
			}
			skipped++
			continue
		}

		relPath, err := filepath.Rel(inputDir, file)
		if err != nil {
			relPath = filepath.Base(file) // fallback
		}

		outputFile := filepath.Join(outputDir, relPath)

		outputDirForFile := filepath.Dir(outputFile)
		if err := os.MkdirAll(outputDirForFile, 0755); err != nil {
			log.Printf("Failed to create output directory %s: %v", outputDirForFile, err)
			continue
		}

		if err := gen.Generate(result, outputFile); err != nil {
			log.Printf("Generation error for %s: %v", file, err)
			continue
		}
	}

	// Финальная статистика
	if verbose || totalErrors > 0 || totalWarnings > 0 {
		fmt.Fprintf(os.Stderr, "\n")
		if totalErrors > 0 {
			fmt.Fprintf(os.Stderr, "%s%d error(s)%s\n", colors.Error, totalErrors, colors.Reset)
		}
		if totalWarnings > 0 {
			fmt.Fprintf(os.Stderr, "%s%d warning(s)%s\n", colors.Warning, totalWarnings, colors.Reset)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "\nFiles: generated %d, skipped %d\n", generated, skipped)
		}
	}

	if totalErrors > 0 || (warningsAsErrors && totalWarnings > 0) {
		os.Exit(1)
	}
}

func versionCmd() {
	fmt.Println("elval-gen version 0.1.0")
}
