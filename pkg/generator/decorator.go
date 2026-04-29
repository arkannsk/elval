package generator

import "fmt"

// GenerateDecoratorCode генерирует код Go для конкретного декоратора
func GenerateDecoratorCode(decType string, paramName string, fieldName string) string {
	switch decType {
	case "ctx-get":
		return fmt.Sprintf(`if val := ctx.Value(%q); val != nil {
	if str, ok := val.(string); ok {
		v.%s = str
	}
}`, paramName, fieldName)
	case "httpctx-get":
		return fmt.Sprintf(`if r := ctx.Value("http.request"); r != nil {
	if req, ok := r.(*http.Request); ok {
		v.%s = req.Header.Get(%q)
	}
}`, fieldName, paramName)
	case "env-get":
		return fmt.Sprintf(`v.%s = os.Getenv(%q)`, fieldName, paramName)
	case "time-now":
		return fmt.Sprintf(`v.%s = time.Now()`, fieldName)
	case "uuid-gen":
		return fmt.Sprintf(`v.%s = uuid.New().String()`, fieldName)
	case "trim":
		return fmt.Sprintf(`v.%s = strings.TrimSpace(v.%s)`, fieldName, fieldName)
	case "lower":
		return fmt.Sprintf(`v.%s = strings.ToLower(v.%s)`, fieldName, fieldName)
	case "upper":
		return fmt.Sprintf(`v.%s = strings.ToUpper(v.%s)`, fieldName, fieldName)
	default:
		return ""
	}
}
