package validator

import (
	"github.com/arkannsk/elval/pkg/errs"
)

// And объединяет несколько правил, все они должны пройти.
func And[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		for _, rule := range rules {
			if err := rule(value); err != nil {
				return err
			}
		}
		return nil
	}
}

// Or объединяет несколько правил, хотя бы одно должно пройти.
func Or[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		var lastErr *errs.ValidationError
		for _, rule := range rules {
			if err := rule(value); err == nil {
				return nil
			} else {
				lastErr = err
			}
		}
		if lastErr != nil {
			return lastErr
		}
		return errs.NewValidationError("", "or", "all rules failed")
	}
}

// IfThen применяет правило только если условие истинно.
func IfThen[T any](condition func(T) bool, rule ValidationRule[T]) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if condition(value) {
			return rule(value)
		}
		return nil
	}
}
