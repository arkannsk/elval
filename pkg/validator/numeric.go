package validator

import (
	"github.com/arkannsk/elval/pkg/errs"
)

// Number интерфейс для числовых типов
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

// RequiredNum возвращает правило обязательности для числовых типов.
func RequiredNum[T Number]() ValidationRule[T] {
	var zero T
	return func(value T) *errs.ValidationError {
		if value == zero {
			return errs.NewValidationError("", "required", "field is required")
		}
		return nil
	}
}

// Min возвращает правило минимального значения для числа.
func Min[T Number](min T) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if value < min {
			return errs.NewValidationError("", "min", "value must be >= %v", min)
		}
		return nil
	}
}

// Max возвращает правило максимального значения для числа.
func Max[T Number](max T) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if value > max {
			return errs.NewValidationError("", "max", "value must be <= %v", max)
		}
		return nil
	}
}

// MinMax возвращает правило диапазона значений для числа.
func MinMax[T Number](min, max T) ValidationRule[T] {
	return func(value T) *errs.ValidationError {
		if value < min {
			return errs.NewValidationError("", "min", "value must be >= %v", min)
		}
		if value > max {
			return errs.NewValidationError("", "max", "value must be <= %v", max)
		}
		return nil
	}
}

// NotZero возвращает правило, запрещающее нулевое значение.
func NotZero[T Number]() ValidationRule[T] {
	var zero T
	return func(value T) *errs.ValidationError {
		if value == zero {
			return errs.NewValidationError("", "notzero", "value must not be zero")
		}
		return nil
	}
}

// Positive возвращает правило положительного значения.
func Positive[T Number]() ValidationRule[T] {
	var zero T
	return func(value T) *errs.ValidationError {
		if value <= zero {
			return errs.NewValidationError("", "positive", "value must be positive")
		}
		return nil
	}
}

// Negative возвращает правило отрицательного значения.
func Negative[T Number]() ValidationRule[T] {
	var zero T
	return func(value T) *errs.ValidationError {
		if value >= zero {
			return errs.NewValidationError("", "negative", "value must be negative")
		}
		return nil
	}
}
