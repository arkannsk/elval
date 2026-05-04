// pkg/errs/format.go
package errs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/term"
)

// Color содержит ANSI-коды для цветного вывода в терминале
type Color struct {
	Error   string
	Warning string
	Info    string
	Loc     string
	Dim     string
	Reset   string
}

// DefaultColors возвращает цвета с авто-детектом терминала
// Если stdout не является терминалом (пайп, файл) — возвращает пустые строки
func DefaultColors() Color {
	if !isTerminal() {
		return Color{Reset: ""}
	}
	return Color{
		Error:   "\033[31m", // красный
		Warning: "\033[33m", // жёлтый
		Info:    "\033[36m", // голубой
		Loc:     "\033[35m", // фиолетовый
		Dim:     "\033[2m",  // тусклый
		Reset:   "\033[0m",  // сброс
	}
}

// isTerminal проверяет, является ли stdout интерактивным терминалом
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// FormatDiagnostic форматирует одно диагностическое сообщение с цветами
// Если colors.Reset пустой — вывод будет без цветов (безопасно для логов/пайпов)
func FormatDiagnostic(d Diagnostic, colors Color) string {
	var parts []string

	// Уровень с цветом
	level := map[Severity]string{
		SeverityError:   colors.Error + "error" + colors.Reset,
		SeverityWarning: colors.Warning + "warning" + colors.Reset,
		SeverityInfo:    colors.Info + "info" + colors.Reset,
	}[d.Severity]
	parts = append(parts, fmt.Sprintf("[%s]", level))

	// Локация с цветом
	loc := d.Loc.File
	if d.Loc.Line > 0 {
		loc = fmt.Sprintf("%s:%d", loc, d.Loc.Line)
		if d.Loc.Column > 0 {
			loc = fmt.Sprintf("%s:%d", loc, d.Loc.Column)
		}
	}
	if loc != "" {
		parts = append(parts, colors.Loc+loc+colors.Reset+":")
	}

	// Компонент в скобках
	if d.Component != "" {
		parts = append(parts, fmt.Sprintf("(%s)", d.Component))
	}

	// Директива
	if d.Directive != "" {
		parts = append(parts, fmt.Sprintf("@evl:%s", d.Directive))
	}

	// Контекст: Структура.Поле
	if d.StructName != "" {
		ctx := d.StructName
		if d.FieldName != "" {
			ctx += "." + colors.Dim + d.FieldName + colors.Reset
		}
		parts = append(parts, ctx+":")
	}

	parts = append(parts, d.Message)
	result := strings.Join(parts, " ")

	// Подсказка
	if d.Suggestion != "" {
		hintText := colors.Info + "[Hint]" + colors.Reset
		result += fmt.Sprintf("\n %s: %s%s%s", hintText, colors.Dim, d.Suggestion, colors.Reset)
	}

	return result
}

// PrintDiagnosticsGrouped выводит диагностики, сгруппированные по файлам
// с заголовками и разделителями для лучшей читаемости при обработке нескольких файлов.
// baseDir: корневая директория для расчёта относительных путей (можно пустую строку)
func PrintDiagnosticsGrouped(diags []Diagnostic, verbose bool, colors Color, baseDir string) {
	if len(diags) == 0 {
		return
	}

	// Группировка по файлам
	byFile := make(map[string][]Diagnostic)
	for _, d := range diags {
		file := d.Loc.File
		if file == "" {
			file = "unknown"
		}
		byFile[file] = append(byFile[file], d)
	}

	// Сортировка файлов для стабильного вывода
	files := make([]string, 0, len(byFile))
	for f := range byFile {
		files = append(files, f)
	}
	sort.Strings(files)

	totalErrors, totalWarnings := 0, 0

	// Вывод по файлам
	for _, file := range files {
		fileDiags := byFile[file]

		// Фильтрация по уровню детализации
		if !verbose {
			filtered := make([]Diagnostic, 0)
			for _, d := range fileDiags {
				if d.Severity != SeverityInfo {
					filtered = append(filtered, d)
				}
			}
			fileDiags = filtered
		}
		if len(fileDiags) == 0 {
			continue
		}

		// Заголовок файла
		relPath := file
		if baseDir != "" {
			if r, err := filepath.Rel(baseDir, file); err == nil {
				relPath = r
			}
		}

		// Печатаем заголовок с разделителем
		fmt.Fprintf(os.Stderr, "\n%s%s%s:\n", colors.Loc, relPath, colors.Reset)
		fmt.Fprintf(os.Stderr, "%s%s%s\n", colors.Dim, strings.Repeat("-", len(relPath)+1), colors.Reset)

		// Вывод диагностик файла
		for _, d := range fileDiags {
			fmt.Fprintln(os.Stderr, FormatDiagnostic(d, colors))
			if d.Severity == SeverityError {
				totalErrors++
			} else if d.Severity == SeverityWarning {
				totalWarnings++
			}
		}
	}

	// Итоговая статистика
	if totalErrors+totalWarnings > 0 {
		fmt.Fprintln(os.Stderr, "")
		if totalErrors > 0 {
			fmt.Fprintf(os.Stderr, "%s%d error(s)%s\n", colors.Error, totalErrors, colors.Reset)
		}
		if totalWarnings > 0 {
			fmt.Fprintf(os.Stderr, "%s%d warning(s)%s\n", colors.Warning, totalWarnings, colors.Reset)
		}
	}
}

// CountBySeverity считает количество сообщений каждого уровня
// Удобно для финальной статистики
func CountBySeverity(diags []Diagnostic) (errors, warnings, infos int) {
	for _, d := range diags {
		switch d.Severity {
		case SeverityError:
			errors++
		case SeverityWarning:
			warnings++
		case SeverityInfo:
			infos++
		}
	}
	return
}
