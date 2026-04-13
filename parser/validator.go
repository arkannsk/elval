package parser

import (
	"fmt"
)

type FieldValidator[T any] struct {
	fieldName   string
	validators  []any // храним любые валидаторы
	parseErrors []error
}

// NewFieldValidator создает типизированный валидатор
func NewFieldValidator[T any](fieldName string) *FieldValidator[T] {
	return &FieldValidator[T]{
		fieldName:   fieldName,
		validators:  make([]any, 0),
		parseErrors: make([]error, 0),
	}
}

func (fv *FieldValidator[T]) AddRequired() {
	var zero T
	v := &RequiredValidator[T]{ZeroValue: zero}

	if err := v.ParseParams([]string{}); err != nil {
		fv.parseErrors = append(fv.parseErrors, err)
		return
	}
	fv.validators = append(fv.validators, v)
}

func (fv *FieldValidator[T]) AddMinMax(params []string) {
	// Создаем валидатор в зависимости от типа T
	var validator any

	switch any(*new(T)).(type) {
	case int, int8, int16, int32, int64:
		v := &MinMaxValidator[int64]{}
		if err := v.ParseParams(params); err != nil {
			fv.parseErrors = append(fv.parseErrors, err)
			return
		}
		validator = v

	case float32, float64:
		v := &MinMaxValidator[float64]{}
		if err := v.ParseParams(params); err != nil {
			fv.parseErrors = append(fv.parseErrors, err)
			return
		}
		validator = v

	case string:
		v := &MinMaxLenValidator{}
		if err := v.ParseParams(params); err != nil {
			fv.parseErrors = append(fv.parseErrors, err)
			return
		}
		validator = v

	default:
		fv.parseErrors = append(fv.parseErrors, fmt.Errorf("min-max не поддерживается для типа %T", *new(T)))
		return
	}

	fv.validators = append(fv.validators, validator)
}

func (fv *FieldValidator[T]) AddPattern(params []string) {
	// Проверяем что T - строка
	if _, ok := any(*new(T)).(string); !ok {
		fv.parseErrors = append(fv.parseErrors, fmt.Errorf("pattern поддерживается только для string, текущий тип: %T", *new(T)))
		return
	}

	v := &PatternValidator{}
	if err := v.ParseParams(params); err != nil {
		fv.parseErrors = append(fv.parseErrors, err)
		return
	}
	fv.validators = append(fv.validators, v)
}

func (fv *FieldValidator[T]) Validate(value T) error {
	if len(fv.parseErrors) > 0 {
		return fmt.Errorf("ошибки конфигурации: %v", fv.parseErrors)
	}

	for _, validator := range fv.validators {
		var err error

		switch v := validator.(type) {
		case *RequiredValidator[T]:
			err = v.Validate(value)

		case *MinMaxValidator[int64]:
			// Конвертируем T в int64
			var intVal int64
			switch val := any(value).(type) {
			case int:
				intVal = int64(val)
			case int8:
				intVal = int64(val)
			case int16:
				intVal = int64(val)
			case int32:
				intVal = int64(val)
			case int64:
				intVal = val
			default:
				return fmt.Errorf("не удалось конвертировать %T в int64", value)
			}
			err = v.Validate(intVal)

		case *MinMaxValidator[float64]:
			// Конвертируем T в float64
			var floatVal float64
			switch val := any(value).(type) {
			case float32:
				floatVal = float64(val)
			case float64:
				floatVal = val
			default:
				return fmt.Errorf("не удалось конвертировать %T в float64", value)
			}
			err = v.Validate(floatVal)

		case *MinMaxLenValidator:
			if strVal, ok := any(value).(string); ok {
				err = v.Validate(strVal)
			} else {
				return fmt.Errorf("MinMaxLenValidator ожидает string, получен %T", value)
			}

		case *PatternValidator:
			if strVal, ok := any(value).(string); ok {
				err = v.Validate(strVal)
			} else {
				return fmt.Errorf("PatternValidator ожидает string, получен %T", value)
			}
		}

		if err != nil {
			return fmt.Errorf("поле %s: %w", fv.fieldName, err)
		}
	}

	return nil
}
