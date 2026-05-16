package validator

import (
	"github.com/arkannsk/elval/pkg/errs"
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

// SkipIfEmpty оборачивает валидатор, чтобы пропускать пустые слайсы
func (sv *SliceValidator[T]) SkipIfEmpty() *SliceValidator[T] {
	sv.required = false
	sv.notZero = false
	return sv
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
func (sv *SliceValidator[T]) Validate(value []T) *errs.ValidationError {
	if !sv.notZero && !sv.required && len(value) == 0 {
		return nil
	}

	if sv.required && len(value) == 0 {
		return errs.NewValidationError("", "required", "field %s: slice cant be nil", sv.fieldName)
	}

	if sv.notZero && (len(value) == 0) {
		return errs.NewValidationError("", "notzero", "field %s: slice cant be empty", sv.fieldName)
	}

	size := len(value)
	if sv.exactSize >= 0 && size != sv.exactSize {
		return errs.NewValidationError("", "exact_size", "field %s: expect size %d, received %d", sv.fieldName, sv.exactSize, size)
	}

	if sv.minSize >= 0 && size < sv.minSize {
		return errs.NewValidationError("", "min_size", "field %s: min len %d, received %d", sv.fieldName, sv.minSize, size)
	}

	if sv.maxSize >= 0 && size > sv.maxSize {
		return errs.NewValidationError("", "max_size", "field %s: max len %d, received %d", sv.fieldName, sv.maxSize, size)
	}

	if sv.elementValidator != nil {
		for i, elem := range value {
			if err := sv.elementValidator.Validate(elem); err != nil {
				return errs.NewValidationError("", "element_validation", "field %s[%d]: %s", sv.fieldName, i, err.Error())
			}
		}
	}

	return nil
}
