package validator

import (
	"fmt"
)

type SliceValidator[T any] struct {
	fieldName        string
	minSize          int
	maxSize          int
	exactSize        int
	notZero          bool
	required         bool
	elementValidator *FieldValidator[T]
}

func NewSliceValidator[T any](fieldName string) *SliceValidator[T] {
	return &SliceValidator[T]{
		fieldName: fieldName,
		minSize:   -1,
		maxSize:   -1,
		exactSize: -1,
		notZero:   false,
		required:  false,
	}
}

// Min устанавливает минимальный размер слайса
func (sv *SliceValidator[T]) Min(size int) *SliceValidator[T] {
	sv.minSize = size
	return sv
}

// Max устанавливает максимальный размер слайса
func (sv *SliceValidator[T]) Max(size int) *SliceValidator[T] {
	sv.maxSize = size
	return sv
}

// Len устанавливает точный размер слайса
func (sv *SliceValidator[T]) Len(size int) *SliceValidator[T] {
	sv.exactSize = size
	return sv
}

// NotZero требует чтобы слайс не был пустым
func (sv *SliceValidator[T]) NotZero() *SliceValidator[T] {
	sv.notZero = true
	return sv
}

// Required требует чтобы слайс не был nil
func (sv *SliceValidator[T]) Required() *SliceValidator[T] {
	sv.required = true
	return sv
}

// Each валидирует каждый элемент слайса
func (sv *SliceValidator[T]) Each(validator *FieldValidator[T]) *SliceValidator[T] {
	sv.elementValidator = validator
	return sv
}

// Validate проверяет слайс
func (sv *SliceValidator[T]) Validate(value []T) error {
	if sv.required && len(value) == 0 {
		return fmt.Errorf("поле %s: слайс не может быть nil", sv.fieldName)
	}

	// Проверка на пустоту
	if sv.notZero && (len(value) == 0) {
		return fmt.Errorf("поле %s: слайс не может быть пустым", sv.fieldName)
	}

	size := len(value)

	// Проверка точного размера
	if sv.exactSize >= 0 && size != sv.exactSize {
		return fmt.Errorf("поле %s: ожидался размер %d, получен %d", sv.fieldName, sv.exactSize, size)
	}

	// Проверка минимального размера
	if sv.minSize >= 0 && size < sv.minSize {
		return fmt.Errorf("поле %s: минимальный размер %d, получен %d", sv.fieldName, sv.minSize, size)
	}

	// Проверка максимального размера
	if sv.maxSize >= 0 && size > sv.maxSize {
		return fmt.Errorf("поле %s: максимальный размер %d, получен %d", sv.fieldName, sv.maxSize, size)
	}

	// Валидация каждого элемента
	if sv.elementValidator != nil {
		for i, elem := range value {
			if err := sv.elementValidator.Validate(elem); err != nil {
				return fmt.Errorf("поле %s[%d]: %w", sv.fieldName, i, err)
			}
		}
	}

	return nil
}
