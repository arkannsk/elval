package annotations

import "strings"

// FieldAnnotationResult содержит обработанные аннотации поля для использования в генераторе
type FieldAnnotationResult struct {
	RewriteRef  string
	RewriteType string
	IsIgnored   bool
	OaIn        string         // "path", "query", "header", "cookie"
	OaParamName string         // Имя параметра, если отличается от имени поля
	Remaining   []OaAnnotation // Остальные аннотации (description, format, enum и т.д.)
}

// ProcessFieldAnnotations обрабатывает сырые OA аннотации поля
func ProcessFieldAnnotations(annotations []OaAnnotation) FieldAnnotationResult {
	result := FieldAnnotationResult{
		Remaining: make([]OaAnnotation, 0, len(annotations)),
	}

	for _, ann := range annotations {
		switch ann.Type {
		case "rewrite.ref":
			result.RewriteRef = trimQuotes(ann.Value)
		case "rewrite.type":
			result.RewriteType = RewriteTypeToOa(ann.Value)
		case "ignore":
			result.IsIgnored = true
		case "in":
			parts := strings.Fields(ann.Value)
			if len(parts) >= 1 {
				result.OaIn = parts[0]
			}
			if len(parts) >= 2 {
				result.OaParamName = parts[1]
			}
		default:
			result.Remaining = append(result.Remaining, ann)
		}
	}
	return result
}

// RewriteTypeToOa конвертирует строковое представление типа в OpenAPI тип
func RewriteTypeToOa(val string) string {
	val = strings.ToLower(trimQuotes(val))
	switch val {
	case "string", "text", "strings":
		return "string"
	case "boolean", "bool":
		return "boolean"
	case "integer", "int":
		return "integer"
	case "number", "float":
		return "number"
	case "array", "[]":
		return "array"
	default:
		return "object"
	}
}
