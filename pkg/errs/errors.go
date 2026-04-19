package errs

import (
	"fmt"
)

// Стандартные ошибки
var (
	ErrRequired     = NewValidationError("", "field is required")
	ErrInvalidEmail = NewValidationError("", "invalid email")
	ErrInvalidPhone = NewValidationError("", "invalid phone number")
	ErrInvalidUUID  = NewValidationError("", "invalid UUID")
)

type ValidationError struct {
	Field   string
	Rule    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError создает новую ошибку валидации
func NewValidationError(rule, format string, args ...interface{}) *ValidationError {
	return &ValidationError{
		Rule:    rule,
		Message: fmt.Sprintf(format, args...),
	}
}
