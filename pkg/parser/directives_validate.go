package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// pkg/parser/directives.go

func validateDirective(dir Directive, ft FieldType) (Severity, error) {
	// 1. Кастомные директивы с префиксом x- всегда валидны (warning, так как требуют регистрации)
	if strings.HasPrefix(dir.Type, "x-") {
		return SeverityWarning, nil
	}

	// 2. Проверяем существование директивы
	info, ok := SupportedDirectives[DirectiveType(dir.Type)]
	if !ok {
		return SeverityError, fmt.Errorf("неизвестная директива. Доступные: %s", getAvailableDirectives())
	}

	// 3. Deprecated директивы - warning
	if info.Deprecated {
		return SeverityWarning, fmt.Errorf("директива устарела, используйте min и max")
	}

	// 4. Проверяем количество параметров
	if len(dir.Params) != info.ParamCount && info.ParamCount != -1 {
		if info.ParamCount == 0 {
			return SeverityError, fmt.Errorf("директива не должна иметь параметров")
		}
		return SeverityError, fmt.Errorf("требуется %d параметр(ов), получено %d", info.ParamCount, len(dir.Params))
	}

	// 5. Определяем базовый тип
	baseType := ft.Name
	if ft.IsPointer {
		baseType = strings.TrimPrefix(baseType, "*")
	}
	if ft.IsSlice {
		baseType = "slice"
	}
	if baseType == "time.Time" || baseType == "time.Duration" {
		baseType = baseType
	}

	// 6. Проверяем тип поля
	typeSupported := false
	for _, allowedType := range info.AllowedTypes {
		if allowedType == baseType {
			typeSupported = true
			break
		}
		if ft.IsCustom && allowedType == "any" {
			typeSupported = true
			break
		}
		if ft.IsSlice && allowedType == "slice" {
			typeSupported = true
			break
		}
		if ft.IsPointer && allowedType == "pointer" {
			typeSupported = true
			break
		}
		if baseType == "time.Time" && allowedType == "time.Time" {
			typeSupported = true
			break
		}
		if baseType == "time.Duration" && allowedType == "time.Duration" {
			typeSupported = true
			break
		}
	}

	if !typeSupported {
		actualType := getActualTypeString(ft)
		// Для указателей на примитивы - warning, для остальных - error
		if ft.IsPointer && !ft.IsStruct {
			return SeverityWarning, fmt.Errorf("не рекомендуется для типа %s. Поддерживаемые типы: %s",
				actualType, strings.Join(info.AllowedTypes, ", "))
		}
		return SeverityError, fmt.Errorf("не поддерживается для типа %s. Поддерживаемые типы: %s",
			actualType, strings.Join(info.AllowedTypes, ", "))
	}

	// 7. Специфичные проверки
	switch dir.Type {
	case "pattern":
		return validatePatternParam(dir, info)
	case "enum":
		return validateEnumParam(dir)
	case "after", "before":
		return validateDateParam(dir)
	case "min", "max":
		return validateMinMaxParam(dir, baseType)
	case "len":
		return validateLenParam(dir, baseType)
	}

	return SeverityError, nil
}

func validatePatternParam(dir Directive, info DirectiveInfo) (Severity, error) {
	if len(dir.Params) == 0 {
		return SeverityError, nil
	}

	pattern := dir.Params[0]

	// Предопределённые паттерны
	if _, ok := info.PredefinedPatterns[pattern]; ok {
		return SeverityError, nil
	}

	// Проверка regexp
	if _, err := regexp.Compile(pattern); err != nil {
		return SeverityError, fmt.Errorf("невалидное регулярное выражение: %s", err.Error())
	}

	return SeverityError, nil
}

func validateEnumParam(dir Directive) (Severity, error) {
	if len(dir.Params) == 0 {
		return SeverityError, fmt.Errorf("требуется хотя бы одно значение")
	}
	return SeverityError, nil
}

func validateDateParam(dir Directive) (Severity, error) {
	if len(dir.Params) == 0 {
		return SeverityError, nil
	}

	date := dir.Params[0]
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05+07:00",
		time.RFC3339,
	}

	for _, format := range formats {
		if _, err := time.Parse(format, date); err == nil {
			return SeverityError, nil
		}
	}

	return SeverityError, fmt.Errorf("невалидный формат даты. Ожидается: YYYY-MM-DD или RFC3339, получено '%s'", date)
}

func validateMinMaxParam(dir Directive, baseType string) (Severity, error) {
	if len(dir.Params) == 0 {
		return SeverityError, nil
	}

	param := dir.Params[0]

	if baseType == "string" || baseType == "slice" {
		if _, err := strconv.Atoi(param); err != nil {
			return SeverityError, fmt.Errorf("параметр должен быть целым числом, получено '%s'", param)
		}
	} else if baseType == "time.Duration" {
		if _, err := time.ParseDuration(param); err != nil {
			return SeverityError, fmt.Errorf("параметр должен быть валидной длительностью (например, 1s, 1m, 1h), получено '%s'", param)
		}
	} else {
		if _, err := strconv.ParseFloat(param, 64); err != nil {
			return SeverityError, fmt.Errorf("параметр должен быть числом, получено '%s'", param)
		}
	}

	return SeverityError, nil
}

func validateLenParam(dir Directive, baseType string) (Severity, error) {
	if len(dir.Params) == 0 {
		return SeverityError, nil
	}

	if baseType != "string" && baseType != "slice" {
		return SeverityError, fmt.Errorf("len поддерживается только для string и slice")
	}

	param := dir.Params[0]
	if _, err := strconv.Atoi(param); err != nil {
		return SeverityError, fmt.Errorf("параметр должен быть целым числом, получено '%s'", param)
	}

	return SeverityError, nil
}

func getActualTypeString(ft FieldType) string {
	if ft.IsSlice {
		return "[]" + ft.Name
	}
	if ft.IsPointer {
		return "*" + ft.Name
	}
	if ft.Name == "time.Time" || ft.Name == "time.Duration" {
		return ft.Name
	}
	return ft.Name
}

func getAvailableDirectives() string {
	var list []string
	for d := range SupportedDirectives {
		list = append(list, string(d))
	}
	// добавляем кастомные
	list = append(list, "x-*")
	return strings.Join(list, ", ")
}
