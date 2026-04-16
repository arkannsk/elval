package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func validateDirective(dir Directive, ft FieldType) error {
	// 1. Кастомные директивы с префиксом x- всегда валидны
	if strings.HasPrefix(dir.Type, "x-") {
		return nil
	}

	// 2. Проверяем существование директивы
	info, ok := SupportedDirectives[DirectiveType(dir.Type)]
	if !ok {
		return fmt.Errorf("неизвестная директива. Доступные: %s", getAvailableDirectives())
	}

	// 3. Проверяем количество параметров
	if len(dir.Params) != info.ParamCount && info.ParamCount != -1 {
		if info.ParamCount == 0 {
			return fmt.Errorf("директива не должна иметь параметров")
		}
		return fmt.Errorf("требуется %d параметр(ов), получено %d", info.ParamCount, len(dir.Params))
	}

	// 4. Определяем базовый тип
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

	// 5. Проверяем тип поля
	typeSupported := false
	for _, allowedType := range info.AllowedTypes {
		if allowedType == baseType {
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
		return fmt.Errorf("не поддерживается для типа %s. Поддерживаемые типы: %s",
			actualType, strings.Join(info.AllowedTypes, ", "))
	}

	// 6. Специфичные проверки
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

	return nil
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

func validatePatternParam(dir Directive, info DirectiveInfo) error {
	if len(dir.Params) == 0 {
		return nil
	}

	pattern := dir.Params[0]

	// Предопределённые паттерны
	if _, ok := info.PredefinedPatterns[pattern]; ok {
		return nil
	}

	// Проверка regexp
	if _, err := regexp.Compile(pattern); err != nil {
		return fmt.Errorf("невалидное регулярное выражение: %s", err.Error())
	}

	return nil
}

func validateEnumParam(dir Directive) error {
	if len(dir.Params) == 0 {
		return fmt.Errorf("требуется хотя бы одно значение")
	}
	return nil
}

func validateDateParam(dir Directive) error {
	if len(dir.Params) == 0 {
		return nil
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
			return nil
		}
	}

	return fmt.Errorf("невалидный формат даты. Ожидается: YYYY-MM-DD или RFC3339, получено '%s'", date)
}

func validateMinMaxParam(dir Directive, baseType string) error {
	if len(dir.Params) == 0 {
		return nil
	}

	param := dir.Params[0]

	if baseType == "string" || baseType == "slice" {
		if _, err := strconv.Atoi(param); err != nil {
			return fmt.Errorf("параметр должен быть целым числом, получено '%s'", param)
		}
	} else if baseType == "time.Duration" {
		if _, err := time.ParseDuration(param); err != nil {
			return fmt.Errorf("параметр должен быть валидной длительностью (например, 1s, 1m, 1h), получено '%s'", param)
		}
	} else {
		if _, err := strconv.ParseFloat(param, 64); err != nil {
			return fmt.Errorf("параметр должен быть числом, получено '%s'", param)
		}
	}

	return nil
}

func validateLenParam(dir Directive, baseType string) error {
	if len(dir.Params) == 0 {
		return nil
	}

	if baseType != "string" && baseType != "slice" {
		return fmt.Errorf("len поддерживается только для string и slice")
	}

	param := dir.Params[0]
	if _, err := strconv.Atoi(param); err != nil {
		return fmt.Errorf("параметр должен быть целым числом, получено '%s'", param)
	}

	return nil
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
