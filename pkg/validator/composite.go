package validator

import (
	"strings"

	"github.com/arkannsk/elval/pkg/errs"
)

// And комбинирует несколько правил - все должны выполниться
func And[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		var errors []string
		for _, rule := range rules {
			if err := rule(value); err != nil {
				errors = append(errors, err.Error())
			}
		}
		if len(errors) > 0 {
			return errs.NewValidationError("", "invalid condition: %s", strings.Join(errors, "; "))
		}
		return nil
	}
}

// Or комбинирует правила - достаточно выполнения хотя бы одного
func Or[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		var errors []string
		for _, rule := range rules {
			if err := rule(value); err == nil {
				return nil
			} else {
				errors = append(errors, err.Error())
			}
		}
		return errs.NewValidationError("",
			"all condition failed: %s", strings.Join(errors, "; "))
	}
}

// IfThen условное правило
func IfThen[T any](condition func(T) bool, rule ValidationRule[T]) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if condition(value) {
			return rule(value)
		}
		return nil
	}
}
