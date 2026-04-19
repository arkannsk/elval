package main

import (
	"strings"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/pkg/validator"
	"github.com/samber/mo"
)

func init() {
	// Валидатор для x-option-present - проверяет что Option не пустой
	validator.RegisterCustom("x-option-present", func(value any, params string) *errs.ValidationError {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return errs.NewValidationError("", "ожидается mo.Option[string]")
		}
		if opt.IsAbsent() {
			return errs.NewValidationError("", "значение обязательно")
		}
		return nil
	})

	// Валидатор для x-option-absent - проверяет что Option пустой
	validator.RegisterCustom("x-option-absent", func(value any, params string) *errs.ValidationError {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return errs.NewValidationError("", "ожидается mo.Option[string]")
		}
		if opt.IsPresent() {
			return errs.NewValidationError("", "значение должно быть пустым")

		}
		return nil
	})

	// Валидатор для x-option-value-min - проверяет минимальную длину значения внутри Option
	validator.RegisterCustom("x-option-value-min", func(value any, params string) *errs.ValidationError {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return errs.NewValidationError("", "ожидается mo.Option[string]")
		}
		if opt.IsAbsent() {
			return nil // пустой Option пропускаем
		}

		val := opt.MustGet()
		if len(val) < 3 {
			return errs.NewValidationError("",
				"значение должно содержать минимум 3 символа, текущая длина: %d", len(val))
		}
		return nil
	})

	// Валидатор для x-option-value-max - проверяет максимальную длину значения внутри Option
	validator.RegisterCustom("x-option-value-max", func(value any, params string) *errs.ValidationError {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return errs.NewValidationError("", "ожидается mo.Option[string]")
		}
		if opt.IsAbsent() {
			return nil
		}

		val := opt.MustGet()
		if len(val) > 50 {
			return errs.NewValidationError("",
				"значение должно содержать максимум 50 символов, текущая длина: %d", len(val))
		}
		return nil
	})

	// Валидатор для x-option-value-eq - проверяет равенство значения внутри Option
	validator.RegisterCustom("x-option-value-eq", func(value any, params string) *errs.ValidationError {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return errs.NewValidationError("", "ожидается mo.Option[string]")
		}
		if opt.IsAbsent() {
			return nil
		}

		val := opt.MustGet()
		if val != params {
			return errs.NewValidationError("", "значение должно содержать '%s', текущее: '%s'", params, val)
		}
		return nil
	})

	// Валидатор для x-option-value-contains - проверяет содержит ли значение подстроку
	validator.RegisterCustom("x-option-value-contains", func(value any, params string) *errs.ValidationError {
		opt, ok := value.(mo.Option[string])
		if !ok {
			return errs.NewValidationError("", "ожидается mo.Option[string]")
		}
		if opt.IsAbsent() {
			return nil
		}

		val := opt.MustGet()
		if !strings.Contains(val, params) {
			return errs.NewValidationError("", "значение должно содержать '%s', текущее: '%s'", params, val)
		}
		return nil
	})
}
