package validator

// Number ограничение на числовые типы
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Min минимальное значение
func Min[T Number](min T) ValidationRule[T] {
	return func(value T) error {
		if value < min {
			return NewValidationError("min", "значение должно быть не менее %v", min)
		}
		return nil
	}
}

// Max максимальное значение
func Max[T Number](max T) ValidationRule[T] {
	return func(value T) error {
		if value > max {
			return NewValidationError("max", "значение должно быть не более %v", max)
		}
		return nil
	}
}

// MinMax диапазон значений
func MinMax[T Number](min, max T) ValidationRule[T] {
	return func(value T) error {
		if value < min || value > max {
			return NewValidationError("min_max", "значение должно быть между %v и %v", min, max)
		}
		return nil
	}
}

// Positive положительное число
func Positive[T Number]() ValidationRule[T] {
	return func(value T) error {
		var zero T
		if value <= zero {
			return NewValidationError("positive", "значение должно быть положительным")
		}
		return nil
	}
}

// Negative отрицательное число
func Negative[T Number]() ValidationRule[T] {
	return func(value T) error {
		var zero T
		if value >= zero {
			return NewValidationError("negative", "значение должно быть отрицательным")
		}
		return nil
	}
}

// NotZero не нулевое значение
func NotZero[T Number]() ValidationRule[T] {
	return func(value T) error {
		var zero T
		if value == zero {
			return NewValidationError("not_zero", "значение не может быть нулевым")
		}
		return nil
	}
}
