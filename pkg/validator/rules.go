package validator

// Required проверяет что поле не является zero value
func Required[T any]() ValidationRule[T] {
	return func(value T) error {
		var zero T
		if any(value) == any(zero) {
			return ErrRequired
		}
		return nil
	}
}

// Optional помечает поле как опциональное (просто пропускает)
func Optional[T any]() ValidationRule[T] {
	return func(value T) error {
		return nil
	}
}

// Custom создает кастомное правило
func Custom[T any](fn func(T) error) ValidationRule[T] {
	return fn
}

// SkipIfZero пропускает валидацию если значение zero
func SkipIfZero[T comparable](rule ValidationRule[T]) ValidationRule[T] {
	var zero T
	return func(value T) error {
		if value == zero {
			return nil
		}
		return rule(value)
	}
}

// Enum проверяет что значение входит в список разрешенных
func Enum[T comparable](allowed ...T) ValidationRule[T] {
	allowedMap := make(map[T]struct{}, len(allowed))
	for _, v := range allowed {
		allowedMap[v] = struct{}{}
	}

	return func(value T) error {
		if _, ok := allowedMap[value]; !ok {
			return NewValidationError("enum", "значение должно быть одним из: %v", allowed)
		}
		return nil
	}
}
