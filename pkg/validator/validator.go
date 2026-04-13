package validator

import (
	"fmt"
)

// ValidationRule функция проверки
type ValidationRule[T any] func(T) error

// FieldValidator валидатор поля
type FieldValidator[T any] struct {
	fieldName string
	rules     []ValidationRule[T]
}

// New создает новый валидатор
func New[T any](fieldName string) *FieldValidator[T] {
	return &FieldValidator[T]{
		fieldName: fieldName,
		rules:     make([]ValidationRule[T], 0),
	}
}

// AddRule добавляет правило валидации
func (fv *FieldValidator[T]) AddRule(rule ValidationRule[T]) *FieldValidator[T] {
	fv.rules = append(fv.rules, rule)
	return fv
}

// Validate применяет все правила к значению
func (fv *FieldValidator[T]) Validate(value T) error {
	for _, rule := range fv.rules {
		if err := rule(value); err != nil {
			return fmt.Errorf("поле %s: %w", fv.fieldName, err)
		}
	}
	return nil
}

// ValidateFunc удобная функция для быстрой валидации
func ValidateFunc[T any](fieldName string, value T, rules ...ValidationRule[T]) error {
	v := New[T](fieldName)
	for _, rule := range rules {
		v.AddRule(rule)
	}
	return v.Validate(value)
}
