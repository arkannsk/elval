package parser

import (
	"strings"
)

// CollectValidationImports собирает импорты, необходимые для валидации/декорирования
func CollectValidationImports(structs []Struct) map[string]string {
	required := make(map[string]string)
	required["validator"] = "github.com/arkannsk/elval/pkg/validator"

	var (
		needsElval, hasCustomValidator, needsErrs, hasDecor bool
	)

	for _, s := range structs {
		if s.HasCustomValidator {
			hasCustomValidator = true
		}
		if s.HasDecorators() {
			hasDecor = true
		}

		for _, field := range s.Fields {
			// Проверяем, есть ли у поля директивы валидации
			hasDirectives := len(field.Directives) > 0

			// Если есть директивы, проверяем тип для импортов
			if hasDirectives {
				checkTypeForImports(field.Type, required, &needsElval)
				hasRealRules := false
				for _, d := range field.Directives {
					if d.Type != "required" && d.Type != "optional" {
						hasRealRules = true
					}
				}
				if hasRealRules {
					needsErrs = true
				}
			}

			// Директивы → импорты (uuid и т.д.)
			for _, dir := range field.Directives {
				if dir.Type == "uuid" {
					required["uuid"] = "github.com/google/uuid"
				}
			}

			// Декораторы → импорты
			for _, dec := range field.Decorators {
				switch dec.Type {
				case "uuid-gen":
					required["uuid"] = "github.com/google/uuid"
				case "env-get", "env_default":
					required["os"] = "os"
				case "time-now":
					required["time"] = "time"
				case "httpctx-get":
					required["net/http"] = "net/http"
				case "trim", "lower", "upper":
					required["strings"] = "strings"
				}
			}
		}
	}

	if needsElval {
		required["elval"] = "github.com/arkannsk/elval"
	}
	if hasCustomValidator || needsErrs {
		required["errs"] = "github.com/arkannsk/elval/pkg/errs"
	}
	if hasDecor {
		required["context"] = "context"
	}

	return required
}

// CollectOpenAPIImports собирает импорты для OpenAPI-схем и HTTP парсинга
func CollectOpenAPIImports(structs []Struct) map[string]string {
	required := map[string]string{
		"oa": "github.com/arkannsk/elval/pkg/openapi",
	}

	hasParsingErrors := false

	for _, s := range structs {
		for _, field := range s.Fields {
			if field.OaIn != "" {
				// Добавляем net/http для Parse(r *http.Request)
				required["net/http"] = "net/http"
				if field.Type.Name == "bool" {
					required["strconv"] = "strconv"
				}

				// Проверяем тип поля для strconv/time и потенциальных ошибок
				if needsParsingError(field.Type) {
					hasParsingErrors = true
				}

				checkTypeForHTTPImports(field.Type, required)
			}
		}
	}

	// Добавляем errs только если есть поля, требующие парсинга с проверкой ошибок
	if hasParsingErrors {
		required["errs"] = "github.com/arkannsk/elval/pkg/errs"
	}

	return required
}

// needsParsingError проверяет, требует ли тип парсинга, который может вернуть ошибку
func needsParsingError(ft FieldType) bool {
	baseName := ft.Name
	if ft.IsSlice {
		baseName = ft.GenericBase
	} else if ft.IsPointer {
		baseName = ft.GenericBase
	}

	switch baseName {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"bool",
		"time.Time", "time.Duration":
		return true
	default:
		return false
	}
}

// checkTypeForHTTPImports проверяет тип поля и добавляет strconv/time если нужно
func checkTypeForHTTPImports(ft FieldType, required map[string]string) {
	// Определяем базовый тип
	baseName := ft.Name
	if ft.IsSlice {
		baseName = ft.GenericBase // Для []int берем "int"
	} else if ft.IsPointer {
		baseName = ft.GenericBase // Для *int берем "int"
	}

	// Числовые типы требуют strconv
	if isNumericType(baseName) {
		required["strconv"] = "strconv"
	}

	// Time требует time
	if baseName == "time.Time" || baseName == "time.Duration" {
		required["time"] = "time"
	}

	// Рекурсивная проверка для сложных типов (если нужно)
	if ft.IsGeneric && len(ft.GenericArgs) > 0 {
		for _, arg := range ft.GenericArgs {
			checkTypeForHTTPImports(arg, required)
		}
	}
}

// checkTypeForImports рекурсивно проверяет тип на наличие time.Time/Duration
func checkTypeForImports(ft FieldType, required map[string]string, needsElval *bool) {
	if ft.IsGeneric && len(ft.GenericArgs) > 0 {
		*needsElval = true
		for _, arg := range ft.GenericArgs {
			checkTypeForImports(arg, required, needsElval)
		}
		return
	}
	if ft.Name == "time.Time" || ft.Name == "time.Duration" {
		required["time"] = "time"
	}
	if ft.IsSlice || ft.IsPointer {
		base := strings.TrimPrefix(ft.Name, "[]")
		base = strings.TrimPrefix(base, "*")
		if base == "time.Time" || base == "time.Duration" {
			required["time"] = "time"
		}
	}
}

// isNumericType проверяет, является ли строка именем числового типа Go
func isNumericType(name string) bool {
	switch name {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return true
	}
	return false
}
