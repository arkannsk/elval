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
	DirRequired DirectiveType = "required"
	DirOptional DirectiveType = "optional"
	DirMin      DirectiveType = "min"
	DirMax      DirectiveType = "max"
	DirLen      DirectiveType = "len"
	DirMinMax   DirectiveType = "min-max"
	DirPattern  DirectiveType = "pattern"
	DirNotZero  DirectiveType = "not-zero"
	DirBefore   DirectiveType = "before" // для time.Time: должно быть до указанной даты
	DirAfter    DirectiveType = "after"  // для time.Time: должно быть после указанной даты

	DirEq  DirectiveType = "eq"  // равно
	DirNeq DirectiveType = "neq" // не равно
	DirLt  DirectiveType = "lt"  // меньше
	DirLte DirectiveType = "lte" // меньше или равно
	DirGt  DirectiveType = "gt"  // больше
	DirGte DirectiveType = "gte" // больше или равно
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
		ParamCount:   0,
		Example:      "@evl:validate required",
	},
	DirOptional: {
		Description:  "поле опционально",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool", "slice", "pointer", "time.Time", "time.Duration"},
		ParamCount:   0,
		Example:      "@evl:validate optional",
	},
	DirMin: {
		Description:  "минимальное значение (для чисел), минимальная длина (для строк), минимальный размер (для слайсов), минимальная длительность (для Duration)",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "slice", "time.Duration"},
		ParamCount:   1,
		Example:      "@evl:validate min:18",
	},
	DirMax: {
		Description:  "максимальное значение (для чисел), максимальная длина (для строк), максимальный размер (для слайсов), максимальная длительность (для Duration)",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "slice", "time.Duration"},
		ParamCount:   1,
		Example:      "@evl:validate max:99",
	},
	DirLen: {
		Description:  "точная длина строки или точный размер слайса",
		AllowedTypes: []string{"string", "slice"},
		ParamCount:   1,
		Example:      "@evl:validate len:10",
	},
	DirMinMax: {
		Description:  "диапазон значений - устаревшая, используйте min и max",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "time.Duration"},
		ParamCount:   2,
		Example:      "@evl:validate min-max:3,50",
		Deprecated:   true,
	},
	DirPattern: {
		Description:  "проверка по регулярному выражению или предопределенному паттерну",
		AllowedTypes: []string{"string"},
		ParamCount:   1,
		Example:      "@evl:validate pattern:email",
		PredefinedPatterns: map[string]string{
			"email": `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
			"phone": `^\+?[0-9]{8,15}$`,
			"uuid":  `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
		},
	},
	DirNotZero: {
		Description:  "значение не должно быть нулевым (для слайсов - не пустым, для Time - не zero time)",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64", "slice", "time.Time", "time.Duration"},
		ParamCount:   0,
		Example:      "@evl:validate not-zero",
	},
	DirBefore: {
		Description:  "для time.Time: значение должно быть до указанной даты",
		AllowedTypes: []string{"time.Time"},
		ParamCount:   1,
		Example:      "@evl:validate before:2024-01-01",
	},
	DirAfter: {
		Description:  "для time.Time: значение должно быть после указанной даты",
		AllowedTypes: []string{"time.Time"},
		ParamCount:   1,
		Example:      "@evl:validate after:2020-01-01",
	},
	DirEq: {
		Description:  "значение должно быть равно указанному",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool"},
		ParamCount:   1,
		Example:      "@evl:validate eq:10",
	},
	DirNeq: {
		Description:  "значение не должно быть равно указанному",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool"},
		ParamCount:   1,
		Example:      "@evl:validate neq:0",
	},
	DirLt: {
		Description:  "значение должно быть меньше указанного",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   1,
		Example:      "@evl:validate lt:100",
	},
	DirLte: {
		Description:  "значение должно быть меньше или равно указанному",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   1,
		Example:      "@evl:validate lte:100",
	},
	DirGt: {
		Description:  "значение должно быть больше указанного",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   1,
		Example:      "@evl:validate gt:0",
	},
	DirGte: {
		Description:  "значение должно быть больше или равно указанному",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   1,
		Example:      "@evl:validate gte:18",
	},
}

// ValidateDirective обновляем с учетом базового типа для указателей
func ValidateDirective(dir Directive, ft FieldType) error {
	// Определяем базовый тип (без указателя)
	baseType := ft.Name

	// Для указателей убираем звездочку, но запоминаем что это указатель
	isPointer := ft.IsPointer
	if isPointer {
		baseType = strings.TrimPrefix(baseType, "*")
	}

	// Для слайсов используем тип "slice"
	actualType := baseType
	if ft.IsSlice {
		actualType = "slice"
	}

	// Проверяем существует ли директива
	info, ok := SupportedDirectives[DirectiveType(dir.Type)]
	if !ok {
		return fmt.Errorf("неизвестная директива: %s", dir.Type)
	}

	// Проверяем поддерживаемый тип поля
	typeSupported := false
	for _, allowedType := range info.AllowedTypes {
		if allowedType == actualType {
			typeSupported = true
			break
		}
		// Для указателей проверяем базовый тип
		if isPointer && allowedType == baseType {
			typeSupported = true
			break
		}
		// Для time.Time и time.Duration проверяем точное совпадение имени
		if allowedType == "time.Time" && baseType == "time.Time" {
			typeSupported = true
			break
		}
		if allowedType == "time.Duration" && baseType == "time.Duration" {
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
			if baseType == "time.Duration" {
				if _, err := time.ParseDuration(dir.Params[0]); err != nil {
					return fmt.Errorf("min параметр %s должен быть валидной длительностью", dir.Params[0])
				}
				if _, err := time.ParseDuration(dir.Params[1]); err != nil {
					return fmt.Errorf("max параметр %s должен быть валидной длительностью", dir.Params[1])
				}
			} else {
				if dir.Params[0] > dir.Params[1] && !strings.Contains(dir.Params[0], ".") {
					return fmt.Errorf("min (%s) не может быть больше max (%s)", dir.Params[0], dir.Params[1])
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
