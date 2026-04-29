package generator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/arkannsk/elval/pkg/parser"
	ann "github.com/arkannsk/elval/pkg/parser/annotations"
)

var templateFucMap = template.FuncMap{
	"hasTime": func(structs []parser.Struct) bool {
		for _, s := range structs {
			for _, f := range s.Fields {
				if f.Type.Name == "time.Time" || f.Type.Name == "time.Duration" {
					return true
				}
			}
		}
		return false
	},
	"dict": func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, fmt.Errorf("invalid dict call")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, fmt.Errorf("dict keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	},
	"hasOptional": func(directives []ann.Directive) bool {
		for _, d := range directives {
			if d.Type == "optional" {
				return true
			}
		}
		return false
	},
	"itoa": func(i int) string {
		return strconv.Itoa(i)
	},
	"isCustomDirective": func(dirType string) bool {
		return strings.HasPrefix(dirType, "x-")

	},
	"toLower":    strings.ToLower,
	"contains":   strings.Contains,
	"hasSuffix":  strings.HasSuffix,
	"hasPrefix":  strings.HasPrefix,
	"regexMatch": regexp.MatchString,
	"split":      strings.Split,
	"title": func(s string) string {
		if len(s) == 0 {
			return s
		}
		return strings.ToUpper(s[:1]) + s[1:]
	},
	"trim":       strings.TrimSpace,
	"trimPrefix": strings.TrimPrefix,
	"trimSuffix": strings.TrimSuffix,
	"trimQuotes": func(s string) string {
		return trimQuotes(s)
	},
	"oaString": func(val string) string {
		// 1. Убираем внешние кавычки, если они есть
		val = trimQuotes(val)
		// 2. Оборачиваем в Go-строковый литерал
		return fmt.Sprintf("%q", val)
	},
	"trimStar": func(s string) string {
		return strings.TrimPrefix(s, "*")
	},
	"firstArg": func(args []parser.FieldType) parser.FieldType {
		if len(args) > 0 {
			return args[0]
		}
		return parser.FieldType{}
	},
	"isPrimitive": func(name string) bool {
		_, ok := primitives[name]
		return ok
	},
	"hasDirective": func(directives []ann.Directive, name string) bool {
		for _, d := range directives {
			if d.Type == name {
				return true
			}
		}
		return false
	},
	"trimBrackets": func(s string) string {
		return strings.TrimPrefix(strings.TrimSuffix(s, "]"), "[")
	},
	"baseType": func(ft parser.FieldType) string {
		name := ft.Name
		// Убираем указатели
		name = strings.TrimPrefix(name, "*")
		// Убираем слайсы (для []T → T)
		if strings.HasPrefix(name, "[") && strings.HasSuffix(name, "]") {
			name = strings.TrimPrefix(strings.TrimSuffix(name, "]"), "[")
		}
		// Если есть алиас — используем его базовый тип
		if ft.BaseType != "" {
			return ft.BaseType
		}
		return name
	},
	"uniqueRef": func(typeName, typePackage, typeModule, structPackage, structModule string) string {
		// Если тип из другого модуля — полный путь
		if typeModule != "" && typeModule != structModule {
			return fmt.Sprintf("%s/%s.%s", typeModule, typePackage, typeName)
		}
		// Если тип из другого пакета в том же модуле
		if typePackage != "" && typePackage != structPackage {
			return fmt.Sprintf("%s.%s", typePackage, typeName)
		}
		// Локальный тип
		return typeName
	},
	"safeExample": func(val string) string {
		val = strings.TrimSpace(val)
		// Если значение выглядит как JSON-объект или массив → оборачиваем в backticks
		if len(val) > 0 && (val[0] == '{' || val[0] == '[') {
			return fmt.Sprintf("`%s`", val)
		}
		// Для простых значений (строки, числа, bool) оставляем как есть
		return val
	},
	"globalRefFor": func(typeName, typePkgPath, typeMod, structPkgPath, structMod string) string {
		if typeMod != "" && typeMod != structMod {
			return fmt.Sprintf("%s/%s.%s", typeMod, typePkgPath, typeName)
		}
		if typePkgPath != "" && typePkgPath != structPkgPath {
			return fmt.Sprintf("%s.%s", typePkgPath, typeName)
		}
		return typeName
	},
	"buildGlobalRef": func(typeName, pkgPath, module string) string {
		// 1. Если typeName уже содержит слэш → это полный путь типа (internal/models.User)
		if strings.Contains(typeName, "/") {
			return fmt.Sprintf("%s/%s", module, typeName)
		}

		// 2. Если typeName содержит точку → это пакет.тип (user.User)
		//    Формируем: module/pkgPath/user.User (слэш перед pkg.Type)
		if strings.Contains(typeName, ".") {
			return fmt.Sprintf("%s/%s/%s", module, pkgPath, typeName)
		}

		// 3. Локальный тип (User) → module/pkgPath.Type (точка перед типом)
		return fmt.Sprintf("%s/%s.%s", module, pkgPath, typeName)
	},
	// Для разделения строки enum
	"splitList": func(s string, sep string) []string {
		parts := strings.Split(s, sep)
		var result []string
		for _, p := range parts {
			result = append(result, strings.TrimSpace(p))
		}
		return result
	},
	"generateSchemaTypeCode": func(typeName string, prefix string) string {
		switch typeName {
		case "string":
			return fmt.Sprintf("%s.Type = \"string\"", prefix)
		case "time.Time":
			return fmt.Sprintf(`%s.Type = "string" %s.Format = "date-time"`, prefix, prefix)
		case "bool":
			return fmt.Sprintf("%s.Type = \"boolean\"", prefix)
		// Integers
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64":
			return fmt.Sprintf("%s.Type = \"integer\"", prefix)
		// Numbers
		case "float32", "float64":
			return fmt.Sprintf("%s.Type = \"number\"", prefix)

		default:
			// Если тип не распознан как примитив, возвращаем пустую строку или ошибку
			return ""
		}
	},
	"GenerateDecoratorCode": GenerateDecoratorCode,
}

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
