package validator

// Eq проверяет что значение равно expected
func Eq[T comparable](expected T) ValidationRule[T] {
	return func(value T) error {
		if value != expected {
			return NewValidationError("eq", "значение должно быть равно %v", expected)
		}
		return nil
	}
}

// Neq проверяет что значение не равно expected
func Neq[T comparable](expected T) ValidationRule[T] {
	return func(value T) error {
		if value == expected {
			return NewValidationError("neq", "значение не должно быть равно %v", expected)
		}
		return nil
	}
}

// Lt проверяет что значение меньше expected
func Lt[T Number](expected T) ValidationRule[T] {
	return func(value T) error {
		if value >= expected {
			return NewValidationError("lt", "значение должно быть меньше %v", expected)
		}
		return nil
	}
}

// Lte проверяет что значение меньше или равно expected
func Lte[T Number](expected T) ValidationRule[T] {
	return func(value T) error {
		if value > expected {
			return NewValidationError("lte", "значение должно быть меньше или равно %v", expected)
		}
		return nil
	}
}

// Gt проверяет что значение больше expected
func Gt[T Number](expected T) ValidationRule[T] {
	return func(value T) error {
		if value <= expected {
			return NewValidationError("gt", "значение должно быть больше %v", expected)
		}
		return nil
	}
}

// Gte проверяет что значение больше или равно expected
func Gte[T Number](expected T) ValidationRule[T] {
	return func(value T) error {
		if value < expected {
			return NewValidationError("gte", "значение должно быть больше или равно %v", expected)
		}
		return nil
	}
}
