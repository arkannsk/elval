package parser

import (
	"go/ast"
	"strings"
)

// getFieldName возвращает имя поля из ast.Field
func getFieldName(field *ast.Field) string {
	if len(field.Names) > 0 {
		return field.Names[0].Name
	}
	// Анонимное поле (встраивание)
	switch t := field.Type.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	}
	return ""
}

// isBuiltin проверяет, является ли тип встроенным в Go
func isBuiltin(name string) bool {
	_, ok := map[string]bool{
		"bool": true, "int": true, "int8": true, "int16": true, "int32": true, "int64": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true, "uintptr": true,
		"float32": true, "float64": true, "complex64": true, "complex128": true,
		"string": true, "byte": true, "rune": true, "error": true, "any": true,
	}[name]
	return ok
}

// trimQuotes убирает внешние кавычки из строки
func trimQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
