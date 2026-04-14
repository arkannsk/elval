package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// MinLen минимальная длина строки (в символах)
func MinLen(min int) ValidationRule[string] {
	return func(value string) error {
		if utf8.RuneCountInString(value) < min {
			return NewValidationError("min_len", "минимальная длина %d символов", min)
		}
		return nil
	}
}

// MaxLen максимальная длина строки (в символах)
func MaxLen(max int) ValidationRule[string] {
	return func(value string) error {
		if utf8.RuneCountInString(value) > max {
			return NewValidationError("max_len", "максимальная длина %d символов", max)
		}
		return nil
	}
}

// LenRange диапазон длины строки
func LenRange(min, max int) ValidationRule[string] {
	return func(value string) error {
		length := utf8.RuneCountInString(value)
		if length < min || length > max {
			return NewValidationError("len_range", "длина должна быть между %d и %d символами", min, max)
		}
		return nil
	}
}

// MatchRegexp проверка по регулярному выражению
func MatchRegexp(pattern string) ValidationRule[string] {
	re := regexp.MustCompile(pattern)
	return func(value string) error {
		if !re.MatchString(value) {
			return NewValidationError("regexp", "значение не соответствует паттерну: %s", pattern)
		}
		return nil
	}
}

// Email проверка email
func Email() ValidationRule[string] {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return func(value string) error {
		if !emailRegex.MatchString(value) {
			return ErrInvalidEmail
		}
		return nil
	}
}

// Phone проверка телефона
func Phone() ValidationRule[string] {
	// Регулярка для телефона: опциональный +, затем от 8 до 15 цифр
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{8,15}$`)
	return func(value string) error {
		if !phoneRegex.MatchString(value) {
			return ErrInvalidPhone
		}
		return nil
	}
}

// UUID проверка UUID
func UUID() ValidationRule[string] {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return func(value string) error {
		if !uuidRegex.MatchString(value) {
			return ErrInvalidUUID
		}
		return nil
	}
}

// NotEmpty проверка что строка не пустая
func NotEmpty() ValidationRule[string] {
	return func(value string) error {
		if value == "" {
			return ErrRequired
		}
		return nil
	}
}

// Contains проверяет что строка содержит подстроку
func Contains(substr string) ValidationRule[string] {
	return func(value string) error {
		if !strings.Contains(value, substr) {
			return NewValidationError("contains", "строка должна содержать '%s'", substr)
		}
		return nil
	}
}

// StartsWith проверяет что строка начинается с префикса
func StartsWith(prefix string) ValidationRule[string] {
	return func(value string) error {
		if !strings.HasPrefix(value, prefix) {
			return NewValidationError("starts_with", "строка должна начинаться с '%s'", prefix)
		}
		return nil
	}
}

// EndsWith проверяет что строка заканчивается суффиксом
func EndsWith(suffix string) ValidationRule[string] {
	return func(value string) error {
		if !strings.HasSuffix(value, suffix) {
			return NewValidationError("ends_with", "строка должна заканчиваться на '%s'", suffix)
		}
		return nil
	}
}
