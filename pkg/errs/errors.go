package errs

import "fmt"

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string
	Rule    string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("field '%s': %s", e.Field, e.Message)
	}
	return e.Message
}

// NewValidationError создает новую ошибку валидации.
// field: имя поля (может быть пустым, если будет установлено валидатором)
// rule: название правила (например, "required")
// message: сообщение об ошибке (поддерживает fmt.Sprintf директивы)
// args: аргументы для форматирования сообщения
func NewValidationError(field, rule, message string, args ...any) *ValidationError {
	msg := message
	if len(args) > 0 {
		msg = fmt.Sprintf(message, args...)
	}
	return &ValidationError{
		Field:   field,
		Rule:    rule,
		Message: msg,
	}
}

// Стандартные ошибки
var (
	ErrInvalidEmail = NewValidationError("", "email", "invalid email")
	ErrInvalidPhone = NewValidationError("", "phone", "invalid phone number")
	ErrInvalidUUID  = NewValidationError("", "uuid", "invalid UUID")
)
