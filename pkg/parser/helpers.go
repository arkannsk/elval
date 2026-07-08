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

// extractSerializedName извлекает имя поля из struct тегов (json, yaml, xml).
// Приоритет: json > yaml > xml. Если тег не найден, возвращается lowercase Go-имя поля.
// Тег "-" означает игнорирование поля.
func extractSerializedName(field *ast.Field, goFieldName string) string {
	if field.Tag == nil {
		return strings.ToLower(goFieldName)
	}

	tagStr := field.Tag.Value
	// Убираем внешние обратные кавычки
	tagStr = strings.Trim(tagStr, "`")

	// Парсим теги: json, yaml, xml
	tagMap := parseStructTags(tagStr)

	// Приоритет: json > yaml > xml
	for _, tagName := range []string{"json", "yaml", "xml"} {
		if val, ok := tagMap[tagName]; ok {
			// Тег "-" означает игнорирование
			if val == "-" {
				return ""
			}
			// Убираем опции (omitempty и т.д.)
			parts := strings.Split(val, ",")
			if parts[0] != "" && parts[0] != "-" {
				return parts[0]
			}
		}
	}

	return strings.ToLower(goFieldName)
}

// parseStructTags парсит строку struct тегов в map[tag_name]tag_value
func parseStructTags(tagStr string) map[string]string {
	result := make(map[string]string)

	fields := strings.FieldsFunc(tagStr, func(r rune) bool {
		return r == ' '
	})

	for _, field := range fields {
		colonIdx := strings.Index(field, ":")
		if colonIdx <= 0 {
			continue
		}
		tagName := field[:colonIdx]
		tagValue := field[colonIdx+1:]
		// Убираем обратные кавычки и двойные кавычки
		tagValue = strings.Trim(tagValue, "`\"")
		result[tagName] = tagValue
	}

	return result
}
