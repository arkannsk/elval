package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DirectiveType тип директивы
type DirectiveType string

// Константы директив
const (
	DirRequired   DirectiveType = "required"
	DirOptional   DirectiveType = "optional"
	DirMin        DirectiveType = "min"
	DirMax        DirectiveType = "max"
	DirLen        DirectiveType = "len"
	DirMinMax     DirectiveType = "min-max"
	DirPattern    DirectiveType = "pattern"
	DirNotZero    DirectiveType = "not-zero"
	DirEnum       DirectiveType = "enum"   // проверка на вхождение в список
	DirBefore     DirectiveType = "before" // для time.Time: должно быть до указанной даты
	DirAfter      DirectiveType = "after"  // для time.Time: должно быть после указанной даты
	DirContains   DirectiveType = "contains"
	DirStartsWith DirectiveType = "starts_with"
	DirEndsWith   DirectiveType = "ends_with"
	DirURL        DirectiveType = "url"
	DirHTTPURL    DirectiveType = "http_url"
	DirDSN        DirectiveType = "dsn"
	DirDate       DirectiveType = "date"

	DirEq  DirectiveType = "eq"  // равно
	DirNeq DirectiveType = "neq" // не равно
	DirLt  DirectiveType = "lt"  // меньше
	DirLte DirectiveType = "lte" // меньше или равно
	DirGt  DirectiveType = "gt"  // больше
	DirGte DirectiveType = "gte" // больше или равно
)

const (
	paramsCntNone      = 0
	paramsCntOne       = 1
	paramsCntTwo       = 2
	paramsCntUnbounded = -1
)

// DirectiveInfo информация о директиве
type DirectiveInfo struct {
	Description        string
	AllowedTypes       []string
	ParamCount         int
	Example            string
	Deprecated         bool
	PredefinedPatterns map[string]string // только для DirPattern
	ValidateFunc       func(fieldType string, params []string) error
}

// SupportedDirectives обновляем с поддержкой time.Time и time.Duration
var SupportedDirectives = map[DirectiveType]DirectiveInfo{
	DirRequired: {
		Description:  "поле обязательно для заполнения",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool", "slice", "pointer", "time.Time", "time.Duration"},
		ParamCount:   paramsCntNone,
		Example:      "@evl:validate required",
	},
	DirOptional: {
		Description:  "поле опционально",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool", "slice", "pointer", "time.Time", "time.Duration"},
		ParamCount:   paramsCntNone,
		Example:      "@evl:validate optional",
	},
	DirMin: {
		Description:  "минимальное значение (для чисел), минимальная длина (для строк), минимальный размер (для слайсов), минимальная длительность (для Duration)",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "slice", "time.Duration"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate min:18",
	},
	DirMax: {
		Description:  "максимальное значение (для чисел), максимальная длина (для строк), максимальный размер (для слайсов), максимальная длительность (для Duration)",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "slice", "time.Duration"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate max:99",
	},
	DirLen: {
		Description:  "точная длина строки или точный размер слайса",
		AllowedTypes: []string{"string", "slice"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate len:10",
	},
	DirMinMax: {
		Description:  "диапазон значений - устаревшая, используйте min и max",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "time.Duration"},
		ParamCount:   paramsCntTwo,
		Example:      "@evl:validate min-max:3,50",
		Deprecated:   true,
	},
	DirPattern: {
		Description:  "проверка по регулярному выражению или предопределенному паттерну",
		AllowedTypes: []string{"string"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate pattern:email",
		PredefinedPatterns: map[string]string{
			"email": `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
			"phone": `^\+?[0-9]{8,15}$`,
			"uuid":  `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
		},
	},
	DirEnum: {
		Description:  "значение должно быть из указанного списка",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64"},
		ParamCount:   paramsCntUnbounded, // переменное количество параметров
		Example:      "@evl:validate enum:active,inactive,pending",
	},
	DirURL: {
		Description:  "валидный URL (любая схема)",
		AllowedTypes: []string{"string"},
		ParamCount:   paramsCntNone,
		Example:      "@evl:validate url",
	},
	DirHTTPURL: {
		Description:  "валидный HTTP/HTTPS URL",
		AllowedTypes: []string{"string"},
		ParamCount:   paramsCntNone,
		Example:      "@evl:validate http_url",
	},
	DirDSN: {
		Description:  "валидный DSN для подключения к БД",
		AllowedTypes: []string{"string"},
		ParamCount:   paramsCntNone,
		Example:      "@evl:validate dsn",
	},
	DirContains: {
		Description:  "строка должна содержать подстроку",
		AllowedTypes: []string{"string"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate contains:admin",
	},
	DirStartsWith: {
		Description:  "строка должна начинаться с префикса",
		AllowedTypes: []string{"string"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate starts_with:https://",
	},
	DirEndsWith: {
		Description:  "строка должна заканчиваться суффиксом",
		AllowedTypes: []string{"string"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate ends_with:.go",
	},
	DirNotZero: {
		Description:  "значение не должно быть нулевым (для слайсов - не пустым, для Time - не zero time)",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64", "slice", "time.Time", "time.Duration"},
		ParamCount:   paramsCntNone,
		Example:      "@evl:validate not-zero",
	},
	DirBefore: {
		Description:  "для time.Time: значение должно быть до указанной даты",
		AllowedTypes: []string{"time.Time"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate before:2024-01-01",
	},
	DirAfter: {
		Description:  "для time.Time: значение должно быть после указанной даты",
		AllowedTypes: []string{"time.Time"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate after:2020-01-01",
	},
	DirDate: {
		Description:  "проверка что строка является валидной датой в одном из указанных форматов",
		AllowedTypes: []string{"string"},
		ParamCount:   paramsCntUnbounded, // форматы через запятую: date:RFC3339,RFC3339Nano,2006-01-02
		Example:      "@evl:validate date:RFC3339,2006-01-02",
	},
	DirEq: {
		Description:  "значение должно быть равно указанному",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate eq:10",
	},
	DirNeq: {
		Description:  "значение не должно быть равно указанному",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate neq:0",
	},
	DirLt: {
		Description:  "значение должно быть меньше указанного",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate lt:100",
	},
	DirLte: {
		Description:  "значение должно быть меньше или равно указанному",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate lte:100",
	},
	DirGt: {
		Description:  "значение должно быть больше указанного",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate gt:0",
	},
	DirGte: {
		Description:  "значение должно быть больше или равно указанному",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   paramsCntOne,
		Example:      "@evl:validate gte:18",
	},
}

// ValidateDirective обновляем с учетом базового типа для указателей
func ValidateDirective(dir Directive, ft FieldType) error {
	baseType := ft.Name
	isPointer := ft.IsPointer
	if isPointer {
		baseType = strings.TrimPrefix(baseType, "*")
	}

	// Для указателей используем тип "pointer"
	actualType := baseType
	if ft.IsSlice {
		actualType = "slice"
	} else if ft.IsPointer {
		actualType = "pointer"
	}

	info, ok := SupportedDirectives[DirectiveType(dir.Type)]
	if !ok {
		return fmt.Errorf("неизвестная директива: %s", dir.Type)
	}

	// Проверяем поддерживаемый тип
	typeSupported := false
	for _, allowedType := range info.AllowedTypes {
		if allowedType == actualType {
			typeSupported = true
			break
		}
		// Для указателей на структуры
		if ft.IsPointer && ft.IsStruct && allowedType == "pointer" {
			typeSupported = true
			break
		}
		// Для конкретных типов
		if allowedType == baseType {
			typeSupported = true
			break
		}
	}

	if !typeSupported {
		return fmt.Errorf("директива %s не поддерживается для типа %s", dir.Type, ft.String())
	}

	// Дополнительные проверки для конкретных директив (без изменений)
	switch dir.Type {
	case string(DirMin), string(DirMax):
		if len(dir.Params) > 0 {
			param := dir.Params[0]
			if actualType == "string" || actualType == "slice" {
				if val, err := strconv.Atoi(param); err != nil || val < 0 {
					return fmt.Errorf("параметр %s должен быть неотрицательным целым числом", param)
				}
			} else if baseType == "time.Duration" {
				if _, err := time.ParseDuration(param); err != nil {
					return fmt.Errorf("параметр %s должен быть валидной длительностью (например, 1h, 30m, 500ms)", param)
				}
			} else {
				if _, err := strconv.ParseFloat(param, 64); err != nil {
					return fmt.Errorf("параметр %s должен быть числом", param)
				}
			}
		}

	case string(DirBefore), string(DirAfter):
		if len(dir.Params) > 0 {
			param := dir.Params[0]
			formats := []string{
				"2006-01-02",
				"2006-01-02T15:04:05Z",
				"2006-01-02T15:04:05+07:00",
				time.RFC3339,
			}
			valid := false
			for _, format := range formats {
				if _, err := time.Parse(format, param); err == nil {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("параметр %s должен быть валидной датой (например, 2024-01-01 или 2024-01-01T15:04:05Z)", param)
			}
		}

	case string(DirLen):
		if len(dir.Params) > 0 {
			param := dir.Params[0]
			if val, err := strconv.Atoi(param); err != nil || val < 0 {
				return fmt.Errorf("параметр %s должен быть неотрицательным целым числом", param)
			}
		}

	case string(DirMinMax):
		if len(dir.Params) == 2 {
			minVal := dir.Params[0]
			maxVal := dir.Params[1]

			// Для строк проверяем что min и max - числа и min <= max
			if ft.IsSlice || ft.Name == "string" {
				minInt, errMin := strconv.Atoi(minVal)
				maxInt, errMax := strconv.Atoi(maxVal)
				if errMin != nil || errMax != nil {
					return fmt.Errorf("min и max должны быть целыми числами")
				}
				if minInt > maxInt {
					return fmt.Errorf("min (%d) не может быть больше max (%d)", minInt, maxInt)
				}
			} else {
				// Для чисел проверяем что min <= max
				minFloat, errMin := strconv.ParseFloat(minVal, 64)
				maxFloat, errMax := strconv.ParseFloat(maxVal, 64)
				if errMin != nil || errMax != nil {
					return fmt.Errorf("min и max должны быть числами")
				}
				if minFloat > maxFloat {
					return fmt.Errorf("min (%v) не может быть больше max (%v)", minFloat, maxFloat)
				}
			}
		}

	case string(DirPattern):
		if len(dir.Params) > 0 && dir.Params[0] == "" {
			return fmt.Errorf("pattern не может быть пустым")
		}
	}

	return nil
}

// hasDirectives проверяет есть ли у структуры поля с аннотациями
func (s *Struct) hasDirectives() bool {
	for _, field := range s.Fields {
		if len(field.Directives) > 0 {
			return true
		}
	}
	return false
}

// isUsedAsNested проверяет используется ли структура как вложенная в других структурах
func (s *Struct) isUsedAsNested(allStructs map[string]*Struct) bool {
	for _, other := range allStructs {
		for _, field := range other.Fields {
			if field.Type.Name == s.Name && field.Type.IsStruct {
				return true
			}
			// Проверяем слайсы структур
			if field.Type.IsSlice && field.Type.Name == s.Name {
				return true
			}
			// Проверяем указатели на структуры
			if field.Type.IsPointer && field.Type.Name == s.Name {
				return true
			}
		}
	}
	return false
}

const customPrefix = "x-"

// AddCustomDirective добавляет кастомную директиву (только с префиксом x-)
func AddCustomDirective(name string, allowedTypes []string, paramCount int, description string) error {
	if !strings.HasPrefix(name, customPrefix) {
		return fmt.Errorf("кастомная директива должна начинаться с '%s'", customPrefix)
	}

	if description == "" {
		description = fmt.Sprintf("кастомная директива %s", name)
	}

	SupportedDirectives[DirectiveType(name)] = DirectiveInfo{
		Description:  description,
		AllowedTypes: allowedTypes,
		ParamCount:   paramCount,
		Example:      fmt.Sprintf("@evl:validate %s", name),
	}

	return nil
}
