package main

import (
	"fmt"
	"strings"

	"github.com/arkannsk/elval/pkg/validator"
	"github.com/samber/mo"
)

func init() {
	// Валидатор для x-option-present - проверяет что Option не пустой
	validator.RegisterCustom("x-option-present", func(value any, params string) error {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return fmt.Errorf("ожидается mo.Option[string]")
		}
		if opt.IsAbsent() {
			return fmt.Errorf("значение обязательно")
		}
		return nil
	})

	// Валидатор для x-option-absent - проверяет что Option пустой
	validator.RegisterCustom("x-option-absent", func(value any, params string) error {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return fmt.Errorf("ожидается mo.Option[string]")
		}
		if opt.IsPresent() {
			return fmt.Errorf("значение должно быть пустым")
		}
		return nil
	})

	// Валидатор для x-option-value-min - проверяет минимальную длину значения внутри Option
	validator.RegisterCustom("x-option-value-min", func(value any, params string) error {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return fmt.Errorf("ожидается mo.Option[string]")
		}
		if opt.IsAbsent() {
			return nil // пустой Option пропускаем
		}

		val := opt.MustGet()
		if len(val) < 3 {
			return fmt.Errorf("значение должно содержать минимум 3 символа, текущая длина: %d", len(val))
		}
		return nil
	})

	// Валидатор для x-option-value-max - проверяет максимальную длину значения внутри Option
	validator.RegisterCustom("x-option-value-max", func(value any, params string) error {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return fmt.Errorf("ожидается mo.Option[string]")
		}
		if opt.IsAbsent() {
			return nil
		}

		val := opt.MustGet()
		if len(val) > 50 {
			return fmt.Errorf("значение должно содержать максимум 50 символов, текущая длина: %d", len(val))
		}
		return nil
	})

	// Валидатор для x-option-value-eq - проверяет равенство значения внутри Option
	validator.RegisterCustom("x-option-value-eq", func(value any, params string) error {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return fmt.Errorf("ожидается mo.Option[string]")
		}
		if opt.IsAbsent() {
			return nil
		}

		val := opt.MustGet()
		if val != params {
			return fmt.Errorf("значение должно быть равно '%s', текущее: '%s'", params, val)
		}
		return nil
	})

	// Валидатор для x-option-value-contains - проверяет содержит ли значение подстроку
	validator.RegisterCustom("x-option-value-contains", func(value any, params string) error {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return fmt.Errorf("ожидается mo.Option[string]")
		}
		if opt.IsAbsent() {
			return nil
		}

		val := opt.MustGet()
		if !strings.Contains(val, params) {
			return fmt.Errorf("значение должно содержать '%s', текущее: '%s'", params, val)
		}
		return nil
	})
}
