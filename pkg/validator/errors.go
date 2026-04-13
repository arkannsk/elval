package validator

import (
	"errors"
	"fmt"
)

// Стандартные ошибки
var (
	ErrRequired     = errors.New("поле обязательно")
	ErrInvalidEmail = errors.New("невалидный email адрес")
	ErrInvalidPhone = errors.New("невалидный номер телефона")
	ErrInvalidUUID  = errors.New("невалидный UUID")
)

// ValidationError кастомная ошибка валидации
type ValidationError struct {
	Rule    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError создает новую ошибку валидации
func NewValidationError(rule, format string, args ...interface{}) error {
	return &ValidationError{
		Rule:    rule,
		Message: fmt.Sprintf(format, args...),
	}
}
