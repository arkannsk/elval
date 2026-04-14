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
	"gopkg.in/yaml.v3"
)

type CustomDirective struct {
	Name        string   `yaml:"name"`
	Types       []string `yaml:"types"`
	ParamCount  int      `yaml:"param_count"`
	Description string   `yaml:"description"`
}

type Config struct {
	CustomDirectives []CustomDirective `yaml:"custom_directives"`
}

func main() {
	var inputDir string
	var configFile string
	var verbose bool

	flag.StringVar(&inputDir, "input", ".", "директория с исходными .go файлами")
	flag.StringVar(&configFile, "config", ".elval.yaml", "файл конфигурации с кастомными директивами")
	flag.BoolVar(&verbose, "v", false, "подробный вывод")
	flag.Parse()

	// Загружаем конфигурацию
	if _, err := os.Stat(configFile); err == nil {
		data, err := os.ReadFile(configFile)
		if err != nil {
			log.Printf("⚠️  Ошибка чтения конфига: %v", err)
		} else {
			var config Config
			if err := yaml.Unmarshal(data, &config); err != nil {
				log.Printf("⚠️  Ошибка парсинга конфига: %v", err)
			} else {
				for _, d := range config.CustomDirectives {
					// Проверяем префикс
					if !strings.HasPrefix(d.Name, "x-") {
						log.Printf("⚠️  Кастомная директива %s должна начинаться с 'x-', пропускаем", d.Name)
						continue
					}

					err := parser.AddCustomDirective(d.Name, d.Types, d.ParamCount, d.Description)
					if err != nil {
						log.Printf("⚠️  Ошибка добавления директивы: %v", err)
						continue
					}
					if verbose {
						fmt.Printf("  ✅ Зарегистрирована кастомная директива: %s (типы: %v, параметров: %d)\n",
							d.Name, d.Types, d.ParamCount)
					}
				}
			}
		}
	} else if verbose {
		fmt.Printf("  ℹ️  Файл конфигурации %s не найден, используем только стандартные директивы\n", configFile)
	}

	p := parser.NewParser()
	gen, err := generator.NewGenerator(inputDir)
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
		if strings.HasSuffix(file, ".gen.go") {
			continue
		}

		if verbose {
			fmt.Printf("Обработка: %s\n", filepath.Base(file))
		}

		result, err := p.ParseFile(file)
		if err != nil {
			log.Printf("Ошибка парсинга %s: %v", file, err)
			continue
		}

		if len(result.Structs) == 0 {
			if verbose {
				fmt.Printf("  ⏭️  Пропуск (нет аннотаций)\n")
			}
			skipped++
			continue
		}

		if err := gen.Generate(result, file); err != nil {
			log.Printf("Ошибка генерации для %s: %v", file, err)
			continue
		}

		generated++
		if verbose {
			outputFile := strings.TrimSuffix(filepath.Base(file), ".go") + ".gen.go"
			fmt.Printf("  ✅ Сгенерирован %s\n", outputFile)
		}
	}
}
