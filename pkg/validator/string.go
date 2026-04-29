package validator

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/arkannsk/elval/pkg/errs"
)

// MinLen возвращает правило минимальной длины строки.
func MinLen(min int) ValidationRule[string] {
	return func(s string) *errs.ValidationError {
		if len(s) < min {
			return errs.NewValidationError("", "minlen", "length must be at least %d", min)
		}
		return nil
	}
}

// MaxLen возвращает правило максимальной длины строки.
func MaxLen(max int) ValidationRule[string] {
	return func(s string) *errs.ValidationError {
		if len(s) > max {
			return errs.NewValidationError("", "maxlen", "length must be at most %d", max)
		}
		return nil
	}
}

// LenRange возвращает правило диапазона длины строки.
func LenRange(min, max int) ValidationRule[string] {
	return func(s string) *errs.ValidationError {
		if len(s) < min || len(s) > max {
			return errs.NewValidationError("", "lenrange", "length must be between %d and %d", min, max)
		}
		return nil
	}
}

// MatchRegexp возвращает правило соответствия регулярному выражению.
func MatchRegexp(pattern string) ValidationRule[string] {
	re := regexp.MustCompile(pattern)
	return func(s string) *errs.ValidationError {
		if !re.MatchString(s) {
			return errs.NewValidationError("", "pattern", "value does not match pattern: %s", pattern)
		}
		return nil
	}
}

// Email возвращает правило проверки формата email.
func Email() ValidationRule[string] {
	// Простая проверка regex для email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return func(s string) *errs.ValidationError {
		if !emailRegex.MatchString(s) {
			return errs.ErrInvalidEmail
		}
		return nil
	}
}

// Phone возвращает правило проверки формата телефона.
func Phone() ValidationRule[string] {
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{8,15}$`)
	return func(s string) *errs.ValidationError {
		if !phoneRegex.MatchString(s) {
			return errs.ErrInvalidPhone
		}
		return nil
	}
}

// UUID возвращает правило проверки формата UUID v4.
func UUID() ValidationRule[string] {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return func(s string) *errs.ValidationError {
		if !uuidRegex.MatchString(s) {
			return errs.ErrInvalidUUID
		}
		return nil
	}
}

// URL проверяет что строка является валидным URL (любая схема)
func URL() ValidationRule[string] {
	return func(value string) *errs.ValidationError {
		if value == "" {
			return nil
		}

		u, err := url.Parse(value)
		if err != nil {
			return errs.NewValidationError("", "url", "невалидный URL: %s", err.Error())
		}

		// Проверяем наличие схемы
		if u.Scheme == "" {
			return errs.NewValidationError("", "url", "URL должен содержать схему (например, http://, https://, postgres://)")
		}

		// Проверяем наличие хоста или пути
		if u.Host == "" && u.Path == "" {
			return errs.NewValidationError("", "url", "невалидный URL: отсутствует хост или путь")
		}

		return nil
	}
}

// URLWithScheme проверяет что URL имеет определённую схему
func URLWithScheme(allowedSchemes ...string) ValidationRule[string] {
	return func(value string) *errs.ValidationError {
		if value == "" {
			return nil
		}

		u, err := url.Parse(value)
		if err != nil {
			return errs.NewValidationError("", "url", "невалидный URL: %s", err.Error())
		}

		if u.Scheme == "" {
			return errs.NewValidationError("", "url", "URL должен содержать схему")
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
			return errs.NewValidationError("", "url", "неподдерживаемая схема: %s, разрешены: %v", u.Scheme, allowedSchemes)
		}

		if u.Host == "" {
			return errs.NewValidationError("", "url", "URL должен содержать хост")
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
	return func(value string) *errs.ValidationError {
		if value == "" {
			return nil
		}

		// Проверяем что это похоже на DSN
		// postgres://user:pass@localhost:5432/db
		// mysql://user:pass@localhost:3306/db
		// clickhouse://user:pass@localhost:9000/db
		u, err := url.Parse(value)
		if err != nil {
			return errs.NewValidationError("", "dsn", "невалидный DSN: %s", err.Error())
		}

		if u.Scheme == "" {
			return errs.NewValidationError("", "dsn", "DSN должен содержать схему (postgres://, mysql://, clickhouse://)")
		}

		if u.Host == "" {
			return errs.NewValidationError("", "dsn", "DSN должен содержать хост")
		}

		return nil
	}
}

// NotEmpty возвращает правило, запрещающее пустые строки.
func NotEmpty() ValidationRule[string] {
	return func(s string) *errs.ValidationError {
		if strings.TrimSpace(s) == "" {
			return errs.NewValidationError("", "notempty", "field cannot be empty or whitespace")
		}
		return nil
	}
}

// Contains проверяет, содержит ли строка подстроку substr
func Contains(substr string) ValidationRule[string] {
	return func(s string) *errs.ValidationError {
		if !strings.Contains(s, substr) {
			return errs.NewValidationError("", "contains", "string must contain '%s'", substr)
		}
		return nil
	}
}

// StartsWith проверяет, начинается ли строка с префикса prefix
func StartsWith(prefix string) ValidationRule[string] {
	return func(s string) *errs.ValidationError {
		if !strings.HasPrefix(s, prefix) {
			return errs.NewValidationError("", "starts_with", "string must start with '%s'", prefix)
		}
		return nil
	}
}

// EndsWith проверяет, заканчивается ли строка суффиксом suffix
func EndsWith(suffix string) ValidationRule[string] {
	return func(s string) *errs.ValidationError {
		if !strings.HasSuffix(s, suffix) {
			return errs.NewValidationError("", "ends_with", "string must end with '%s'", suffix)
		}
		return nil
	}
}

// Date проверяет что строка является валидной датой в одном из форматов
func Date(formats ...string) ValidationRule[string] {
	// Предопределённые форматы
	predefined := map[string]string{
		"RFC3339":     time.RFC3339,
		"RFC3339Nano": time.RFC3339Nano,
		"RFC1123":     time.RFC1123,
		"RFC1123Z":    time.RFC1123Z,
		"RFC822":      time.RFC822,
		"RFC822Z":     time.RFC822Z,
		"RFC850":      time.RFC850,
		"Kitchen":     time.Kitchen,
		"Stamp":       time.Stamp,
		"StampMilli":  time.StampMilli,
		"StampMicro":  time.StampMicro,
		"StampNano":   time.StampNano,
		"ANSIC":       time.ANSIC,
		"UnixDate":    time.UnixDate,
		"RubyDate":    time.RubyDate,
	}

	// Разворачиваем форматы
	expandedFormats := make([]string, 0, len(formats))
	for _, f := range formats {
		if predefinedFormat, ok := predefined[f]; ok {
			expandedFormats = append(expandedFormats, predefinedFormat)
		} else {
			expandedFormats = append(expandedFormats, f)
		}
	}

	return func(value string) *errs.ValidationError {
		if value == "" {
			return nil
		}

		for _, format := range expandedFormats {
			if _, err := time.Parse(format, value); err == nil {
				return nil
			}
		}

		return errs.NewValidationError("", "date", "невалидная дата. Ожидаемые форматы: %v", formats)
	}
}

// DateRFC3339 проверяет RFC3339 формат
func DateRFC3339() ValidationRule[string] {
	return Date("RFC3339")
}

// DateRFC3339Nano проверяет RFC3339Nano формат
func DateRFC3339Nano() ValidationRule[string] {
	return Date("RFC3339Nano")
}

// DateISO проверяет ISO формат (2006-01-02)
func DateISO() ValidationRule[string] {
	return Date("2006-01-02")
}

// DateTimeISO проверяет ISO формат с временем (2006-01-02T15:04:05)
func DateTimeISO() ValidationRule[string] {
	return Date("2006-01-02T15:04:05")
}
