package validator

import "github.com/arkannsk/elval/pkg/errs"

// Eq проверяет что значение равно expected
func Eq[T comparable](expected T) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if value != expected {
			return errs.NewValidationError("", "eq", "value must be equal to %v", expected)
		}
		return nil
	}
}

// Neq проверяет что значение не равно expected
func Neq[T comparable](expected T) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if value == expected {
			return errs.NewValidationError("", "neq", "value must not be equal to %v", expected)
		}
		return nil
	}
}

// Lt проверяет что значение меньше expected
func Lt[T Number](expected T) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if value >= expected {
			return errs.NewValidationError("", "lt", "value must be less than %v", expected)
		}
		return nil
	}
}

// Lte проверяет что значение меньше или равно expected
func Lte[T Number](expected T) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if value > expected {
			return errs.NewValidationError("", "lte", "value must be less than or equal to %v", expected)
		}
		return nil
	}
}

// Gt проверяет что значение больше expected
func Gt[T Number](expected T) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if value <= expected {
			return errs.NewValidationError("", "gt", "value must be greater than %v", expected)
		}
		return nil
	}
}

// Gte проверяет что значение больше или равно expected
func Gte[T Number](expected T) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if value < expected {
			return errs.NewValidationError("", "gte", "value must be greater than or equal to %v", expected)
		}
		return nil
	}
}
