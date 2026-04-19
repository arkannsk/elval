package parser

import (
	"fmt"
	"strings"
)

// FieldType представляет тип поля с поддержкой слайсов и указателей
type FieldType struct {
	Name        string
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

type OaAnnotation struct {
	Type  string // title, description, example, format, etc.
	Value string
}

// Field представляет поле структуры с аннотациями
type Field struct {
	Name          string      // имя поля
	Type          FieldType   // тип поля
	Directives    []Directive // список директив валидации
	Decorators    []Decorator // декораторы
	Line          int         // номер строки в файле (для ошибок)
	OaAnnotations []OaAnnotation
	Description   string
}

// Struct представляет структуру с полями для валидации
type Struct struct {
	Name        string  // имя структуры
	Fields      []Field // поля с аннотациями
	File        string  // путь к файлу
	Description string
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
	Package string   // имя пакета
	Structs []Struct // найденные структуры
	Errors  []error  // ошибки парсинга

	Imports map[string]string
}
