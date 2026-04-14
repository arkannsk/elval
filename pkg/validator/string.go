package validator

import (
	"net/url"
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

// URL проверяет что строка является валидным URL (любая схема)
func URL() ValidationRule[string] {
	return func(value string) error {
		if value == "" {
			return nil
		}

		u, err := url.Parse(value)
		if err != nil {
			return NewValidationError("url", "невалидный URL: %s", err.Error())
		}

		// Проверяем наличие схемы
		if u.Scheme == "" {
			return NewValidationError("url", "URL должен содержать схему (например, http://, https://, postgres://)")
		}

		// Проверяем наличие хоста или пути
		if u.Host == "" && u.Path == "" {
			return NewValidationError("url", "невалидный URL: отсутствует хост или путь")
		}

		return nil
	}
}

// URLWithScheme проверяет что URL имеет определённую схему
func URLWithScheme(allowedSchemes ...string) ValidationRule[string] {
	return func(value string) error {
		if value == "" {
			return nil
		}

		u, err := url.Parse(value)
		if err != nil {
			return NewValidationError("url", "невалидный URL: %s", err.Error())
		}

		if u.Scheme == "" {
			return NewValidationError("url", "URL должен содержать схему")
		}

		// Проверяем схему
		schemeAllowed := false
		for _, scheme := range allowedSchemes {
			if u.Scheme == scheme {
				schemeAllowed = true
				break
			}
		}

		if !schemeAllowed {
			return NewValidationError("url", "неподдерживаемая схема: %s, разрешены: %v", u.Scheme, allowedSchemes)
		}

		if u.Host == "" {
			return NewValidationError("url", "URL должен содержать хост")
		}

		return nil
	}
}

// HTTPURL проверяет что URL имеет схему http или https
func HTTPURL() ValidationRule[string] {
	return URLWithScheme("http", "https")
}

// DSN проверяет строку подключения к БД (postgres, mysql, clickhouse и т.д.)
func DSN() ValidationRule[string] {
	return func(value string) error {
		if value == "" {
			return nil
		}

		// Проверяем что это похоже на DSN
		// postgres://user:pass@localhost:5432/db
		// mysql://user:pass@localhost:3306/db
		// clickhouse://user:pass@localhost:9000/db
		u, err := url.Parse(value)
		if err != nil {
			return NewValidationError("dsn", "невалидный DSN: %s", err.Error())
		}

		if u.Scheme == "" {
			return NewValidationError("dsn", "DSN должен содержать схему (postgres://, mysql://, clickhouse://)")
		}

		if u.Host == "" {
			return NewValidationError("dsn", "DSN должен содержать хост")
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
