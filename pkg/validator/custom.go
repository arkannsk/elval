package validator

import (
	"fmt"
	"sync"
)

type CustomValidator func(value any, params string) error

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
func ValidateCustom(name string, value any, params string) error {
	customMu.RLock()
	defer customMu.RUnlock()

	fn, ok := customs[name]
	if !ok {
		return fmt.Errorf("кастомный валидатор %s не зарегистрирован", name)
	}
	return fn(value, params)
}
