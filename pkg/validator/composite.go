package validator

import (
	"fmt"
	"strings"
)

// And комбинирует несколько правил - все должны выполниться
func And[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return func(value T) error {
		var errors []string
		for _, rule := range rules {
			if err := rule(value); err != nil {
				errors = append(errors, err.Error())
			}
		}
		if len(errors) > 0 {
			return fmt.Errorf("не выполнены условия: %s", strings.Join(errors, "; "))
		}
		return nil
	}
}

// Or комбинирует правила - достаточно выполнения хотя бы одного
func Or[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return func(value T) error {
		var errors []string
		for _, rule := range rules {
			if err := rule(value); err == nil {
				return nil
			} else {
				errors = append(errors, err.Error())
			}
		}
		return fmt.Errorf("ни одно из условий не выполнено: %s", strings.Join(errors, "; "))
	}
}

// IfThen условное правило
func IfThen[T any](condition func(T) bool, rule ValidationRule[T]) ValidationRule[T] {
	return func(value T) error {
		if condition(value) {
			return rule(value)
		}
		return nil
	}
}
