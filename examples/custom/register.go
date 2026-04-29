package main

import (
	"strconv"
	"strings"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/pkg/validator"
)

func init() {
	validator.RegisterCustom("x-color", func(value any, params string) *errs.ValidationError {
		str, ok := value.(string)
		if !ok {
			return errs.NewValidationError("", "x-color", "ожидается строка")
		}
		validColors := map[string]bool{"red": true, "green": true, "blue": true}
		if !validColors[strings.ToLower(str)] {
			return errs.NewValidationError("", "x-color", "невалидный цвет: %s", str)
		}
		return nil
	})

	validator.RegisterCustom("x-even", func(value any, params string) *errs.ValidationError {
		num, ok := value.(int)
		if !ok {
			return errs.NewValidationError("", "x-even", "ожидается целое число")
		}
		if num%2 != 0 {
			return errs.NewValidationError("", "x-even", "число %d должно быть четным", num)
		}
		return nil
	})

	validator.RegisterCustom("x-between", func(value any, params string) *errs.ValidationError {
		num, ok := value.(int)
		if !ok {
			return errs.NewValidationError("", "x-color", "ожидается целое число")
		}
		// Парсим параметры: "10,90"
		parts := strings.Split(params, ",")
		if len(parts) != 2 {
			return errs.NewValidationError("", "x-color", "x-between требует 2 параметра: min,max")
		}
		min, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return errs.NewValidationError("", "x-color", "min должен быть числом")
		}
		max, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return errs.NewValidationError("", "x-color", "max должен быть числом")
		}
		if num < min || num > max {
			return errs.NewValidationError("", "x-color", "значение должно быть между %d и %d", min, max)
		}
		return nil
	})

	validator.RegisterCustom("x-contains", func(value any, params string) *errs.ValidationError {
		str, ok := value.(string)
		if !ok {
			return errs.NewValidationError("", "x-contains", "ожидается строка")
		}
		if params == "" {
			return errs.NewValidationError("", "x-contains", "x-contains требует параметр: подстрока")
		}
		if !strings.Contains(str, params) {
			return errs.NewValidationError("", "x-contains", "строка должна содержать '%s'", params)
		}
		return nil
	})
}
