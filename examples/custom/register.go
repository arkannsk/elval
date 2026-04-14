package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arkannsk/elval/pkg/validator"
)

func init() {
	validator.RegisterCustom("x-color", func(value any, params string) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("ожидается строка")
		}
		validColors := map[string]bool{"red": true, "green": true, "blue": true}
		if !validColors[strings.ToLower(str)] {
			return fmt.Errorf("невалидный цвет: %s", str)
		}
		return nil
	})

	validator.RegisterCustom("x-even", func(value any, params string) error {
		num, ok := value.(int)
		if !ok {
			return fmt.Errorf("ожидается целое число")
		}
		if num%2 != 0 {
			return fmt.Errorf("число %d должно быть четным", num)
		}
		return nil
	})

	validator.RegisterCustom("x-between", func(value any, params string) error {
		num, ok := value.(int)
		if !ok {
			return fmt.Errorf("ожидается целое число")
		}
		// Парсим параметры: "10,90"
		parts := strings.Split(params, ",")
		if len(parts) != 2 {
			return fmt.Errorf("x-between требует 2 параметра: min,max")
		}
		min, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return fmt.Errorf("min должен быть числом")
		}
		max, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return fmt.Errorf("max должен быть числом")
		}
		if num < min || num > max {
			return fmt.Errorf("значение должно быть между %d и %d", min, max)
		}
		return nil
	})

	validator.RegisterCustom("x-contains", func(value any, params string) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("ожидается строка")
		}
		if params == "" {
			return fmt.Errorf("x-contains требует параметр: подстрока")
		}
		if !strings.Contains(str, params) {
			return fmt.Errorf("строка должна содержать '%s'", params)
		}
		return nil
	})
}
