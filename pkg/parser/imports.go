package parser

import "strings"

// CollectValidationImports собирает импорты, необходимые для валидации/декорирования
func CollectValidationImports(structs []Struct) map[string]string {
	required := make(map[string]string)
	required["errs"] = "github.com/arkannsk/elval/pkg/errs"
	required["validator"] = "github.com/arkannsk/elval/pkg/validator"
	required["context"] = "context"

	needsElval := false

	for _, s := range structs {
		for _, field := range s.Fields {
			checkTypeForImports(field.Type, required, &needsElval)

			// Директивы → импорты
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

	return required
}

// CollectOpenAPIImports собирает импорты для OpenAPI-схем
func CollectOpenAPIImports(structs []Struct) map[string]string {
	return map[string]string{
		"oa": "github.com/arkannsk/elval/pkg/oa",
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
