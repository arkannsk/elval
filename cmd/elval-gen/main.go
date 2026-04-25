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

	p := parser.NewParser(verbose)
	gen, err := generator.NewGenerator(outputDir, generateOpenAPI, verbose)
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

func versionCmd() {
	fmt.Println("elval-gen version 0.1.0")
}
