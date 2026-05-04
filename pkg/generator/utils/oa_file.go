package utils

import (
	"strings"

	"github.com/arkannsk/elval/pkg/parser"
)

// IsFileOrStreamType возвращает true, если тип должен маппиться
// в OpenAPI как бинарные данные (соответствует spec 4.7.14.3 / 4.14.7).
func IsFileOrStreamType(ft parser.FieldType) bool {
	// Раскрываем дженерики рекурсивно
	if ft.IsGeneric && len(ft.GenericArgs) > 0 {
		return IsFileOrStreamType(ft.GenericArgs[0])
	}

	name := ft.Name
	if ft.IsPointer && len(name) > 0 && name[0] == '*' {
		name = name[1:]
	}
	// os.File -> File
	if parts := strings.Split(name, "."); len(parts) > 1 {
		name = parts[1]
	}

	switch name {
	case "File", "Reader", "ReadCloser", "ReadSeeker":
		return true
	default:
		return strings.HasPrefix(name, "[]byte")
	}
}

// FileSchema OAS 3.0/3.1/3.2.
func FileSchema(ft parser.FieldType, description string) map[string]any {
	schema := map[string]any{"type": "string"}

	name := ft.Name
	if ft.IsPointer && len(name) > 0 && name[0] == '*' {
		name = name[1:]
	}

	// []byte -> byte (base64), other -> binary (raw)
	if strings.HasPrefix(name, "[]byte") {
		schema["format"] = "byte"
	} else {
		schema["format"] = "binary"
	}

	if description != "" {
		schema["description"] = description
	}
	return schema
}
