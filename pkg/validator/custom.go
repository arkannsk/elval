package validator

import (
	"sync"

	"github.com/arkannsk/elval/pkg/errs"
)

type CustomValidator func(value any, params string) *errs.ValidationError

var (
	customMu sync.RWMutex
	customs  = make(map[string]CustomValidator)
)

// RegisterCustom регистрирует кастомный валидатор
func RegisterCustom(name string, fn CustomValidator) {
	customMu.Lock()
	defer customMu.Unlock()
	customs[name] = fn
}

// ValidateCustom вызывает кастомный валидатор с параметрами (как строка)
func ValidateCustom(name string, value any, params string) *errs.ValidationError {
	customMu.RLock()
	defer customMu.RUnlock()

	fn, ok := customs[name]
	if !ok {
		return errs.NewValidationError("", "custom_validator_not_found", "custom validator '%s' not found", name)
	}
	return fn(value, params)
}

func Custom[T any](fn func(T) *errs.ValidationError) ValidationRule[T] {
	return fn
}
