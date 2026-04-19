package generator

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/arkannsk/elval/pkg/parser"
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
	"hasOptional": func(directives []parser.Directive) bool {
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
	"hasCustomDirective": func(directives []parser.Directive) bool {
		for _, d := range directives {
			if strings.HasPrefix(d.Type, "x-") {
				return true
			}
		}
		return false
	},
	"isCustomDirective": func(dirType string) bool {
		return strings.HasPrefix(dirType, "x-")

	},
	"hasHTTPContext": func(structs []parser.Struct) bool {
		for _, s := range structs {
			for _, f := range s.Fields {
				for _, d := range f.Decorators {
					if d.Type == "httpctx-get" {
						return true
					}
				}
			}
		}
		return false
	},
	"hasEnvGet": func(structs []parser.Struct) bool {
		for _, s := range structs {
			for _, f := range s.Fields {
				for _, d := range f.Decorators {
					if d.Type == "env-get" {
						return true
					}
				}
			}
		}
		return false
	},
	"hasUUIDGen": func(structs []parser.Struct) bool {
		for _, s := range structs {
			for _, f := range s.Fields {
				for _, d := range f.Decorators {
					if d.Type == "uuid-gen" {
						return true
					}
				}
			}
		}
		return false
	},
	"hasOaSchema": func(structs []parser.Struct) bool {
		for _, s := range structs {
			for _, f := range s.Fields {
				if len(f.OaAnnotations) > 0 {
					return true
				}
			}
		}
		return false
	},
	"toLower":   strings.ToLower,
	"contains":  strings.Contains,
	"hasSuffix": strings.HasSuffix,
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
	"hasDirective": func(directives []parser.Directive, name string) bool {
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
}
