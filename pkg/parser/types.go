package parser

import (
	"fmt"
	"strings"

	"github.com/arkannsk/elval/pkg/errs"
	ann "github.com/arkannsk/elval/pkg/parser/annotations"
)

// FieldType представляет тип поля с поддержкой слайсов и указателей
type FieldType struct {
	Name        string
	Package     string
	PackagePath string
	Module      string
	IsSlice     bool
	IsPointer   bool
	IsStruct    bool
	IsNamedType bool // true для type X string/int/etc.
	IsCustom    bool // кастомный тип (из другого пакета или с дженериками)
	IsGeneric   bool
	BaseType    string      // for aliases
	GenericBase string      // "Option", "Result" и т.д.
	GenericArgs []FieldType // аргументы типа: [string], [int], [User]
}

// String возвращает строковое представление типа
func (ft FieldType) String() string {
	if ft.IsGeneric && len(ft.GenericArgs) > 0 {
		args := make([]string, len(ft.GenericArgs))
		for i, arg := range ft.GenericArgs {
			args[i] = arg.String()
		}
		return fmt.Sprintf("%s[%s]", ft.GenericBase, strings.Join(args, ", "))
	}
	if ft.IsSlice {
		return "[]" + ft.Name
	}
	if ft.IsPointer {
		return "*" + ft.Name
	}
	return ft.Name
}

// Directive представляет одну аннотацию валидации
type Directive struct {
	Type   string   // required, min, max, len, pattern, not-zero, optional
	Params []string // параметры директивы
	Raw    string   // исходный текст
}

// Field представляет поле структуры с аннотациями
type Field struct {
	Name          string // имя поля
	Package       string
	PackagePath   string
	Module        string    // модуль типа: "github.com/myorg/api"
	Type          FieldType // тип поля
	IsEmbedded    bool
	Directives    []ann.Directive // список директив валидации
	Decorators    []Decorator     // декораторы
	Line          int             // номер строки в файле (для ошибок)
	OaAnnotations []ann.OaAnnotation
	Description   string

	IsIgnored     bool
	Discriminator *ann.OaDiscriminator
	OaOneOf       []string // для полей с union-типами
	OaOneOfRefs   []string
	OaAnyOf       []string
	OaAnyOfRefs   []string

	OaRewriteRef  string
	OaRewriteType string

	OaIn        string // "path", "query", "header", "cookie" или "" (body)
	OaParamName string // Имя параметра в URL/Query/Header (если отличается от имени поля)
}

// Struct представляет структуру с полями для валидации
type Struct struct {
	Name        string // имя структуры
	Package     string
	PackagePath string
	Module      string  // путь модуля: "github.com/myorg/api"
	Fields      []Field // поля с аннотациями
	File        string  // путь к файлу
	IsIgnored   bool
	Description string

	RawOaAnnotations []ann.OaAnnotation

	Discriminator *ann.OaDiscriminator `json:"-"` // не сериализуем, только для генерации
	OaOneOf       []string             // имена типов: ["Cat", "Dog"]
	OaOneOfRefs   []string             // прямые рефы: ["#/components/schemas/Cat"]
	OaAnyOf       []string
	OaAnyOfRefs   []string
}

// HasDirectives проверяет есть ли у структуры поля с аннотациями валидации
func (s *Struct) HasDirectives() bool {
	for _, field := range s.Fields {
		if len(field.Directives) > 0 {
			return true
		}
	}
	return false
}

// HasOaAnnotations проверяет есть ли у структуры OpenAPI аннотации
func (s *Struct) HasOaAnnotations() bool {
	for _, field := range s.Fields {
		if len(field.OaAnnotations) > 0 {
			return true
		}
	}
	return false
}

// HasDecorators проверяет есть ли у структуры декораторы
func (s *Struct) HasDecorators() bool {
	for _, field := range s.Fields {
		if len(field.Decorators) > 0 {
			return true
		}
	}
	return false
}

// HasDescription проверяет есть ли описание у структуры или её полей
func (s *Struct) HasDescription() bool {
	if s.Description != "" {
		return true
	}
	for _, field := range s.Fields {
		if field.Description != "" {
			return true
		}
	}
	return false
}

// ShouldGenerateOpenAPI проверяет нужно ли генерировать OpenAPI для структуры
func (s *Struct) ShouldGenerateOpenAPI(generateOpenAPIFlag bool) bool {
	if generateOpenAPIFlag {
		return true
	}
	if s.HasOaAnnotations() || s.HasDescription() {
		return true
	}
	return false
}

// ParseResult результат парсинга файла
type ParseResult struct {
	Package     string
	Structs     []Struct
	Errors      []error // Critical Generation Errors
	Diagnostics []errs.Diagnostic

	Imports map[string]string
}

func (r *ParseResult) HasErrors() bool {
	for _, d := range r.Diagnostics {
		if d.Severity == errs.SeverityError {
			return true
		}
	}
	return false
}

func (r *ParseResult) HasWarnings() bool {
	for _, d := range r.Diagnostics {
		if d.Severity == errs.SeverityWarning {
			return true
		}
	}
	return false
}
