package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

type DirectiveType string

const (
	DirRequired DirectiveType = "required"
	DirMinMax   DirectiveType = "min-max"
	DirPattern  DirectiveType = "pattern"
)

type Validator[T any] interface {
	Validate(value T) error
	ParseParams(params []string) error
}

type RequiredValidator[T any] struct {
	ZeroValue T
}

func (r *RequiredValidator[T]) ParseParams(params []string) error {
	if len(params) > 0 {
		return fmt.Errorf("required не принимает параметры")
	}
	return nil
}

func (r *RequiredValidator[T]) Validate(value T) error {
	// Сравниваем с zero value
	if any(value) == any(r.ZeroValue) {
		return fmt.Errorf("поле обязательно")
	}
	return nil
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

type MinMaxValidator[T Number] struct {
	Min T
	Max T
}

func (m *MinMaxValidator[T]) ParseParams(params []string) error {
	if len(params) != 2 {
		return fmt.Errorf("min-max требует 2 параметра")
	}

	var min, max T
	var err error

	// Определяем тип и парсим
	switch any(min).(type) {
	case int, int8, int16, int32, int64:
		var iMin, iMax int64
		iMin, err = strconv.ParseInt(params[0], 10, 64)
		if err == nil {
			iMax, err = strconv.ParseInt(params[1], 10, 64)
			min = T(iMin)
			max = T(iMax)
		}
	case float32, float64:
		var fMin, fMax float64
		fMin, err = strconv.ParseFloat(params[0], 64)
		if err == nil {
			fMax, err = strconv.ParseFloat(params[1], 64)
			min = T(fMin)
			max = T(fMax)
		}
	}

	if err != nil {
		return fmt.Errorf("невалидные параметры: %w", err)
	}

	if min > max {
		return fmt.Errorf("min (%v) не может быть больше max (%v)", min, max)
	}

	m.Min = min
	m.Max = max
	return nil
}

func (m *MinMaxValidator[T]) Validate(value T) error {
	if value < m.Min || value > m.Max {
		return fmt.Errorf("значение должно быть между %v и %v", m.Min, m.Max)
	}
	return nil
}

// MinMaxLenValidator для строк
type MinMaxLenValidator struct {
	Min int
	Max int
}

func (m *MinMaxLenValidator) ParseParams(params []string) error {
	if len(params) != 2 {
		return fmt.Errorf("min-max требует 2 параметра")
	}

	min, err := strconv.Atoi(params[0])
	if err != nil || min < 0 {
		return fmt.Errorf("min должен быть неотрицательным целым числом")
	}

	max, err := strconv.Atoi(params[1])
	if err != nil || max < 0 {
		return fmt.Errorf("max должен быть неотрицательным целым числом")
	}

	if min > max {
		return fmt.Errorf("min (%d) не может быть больше max (%d)", min, max)
	}

	m.Min = min
	m.Max = max
	return nil
}

func (m *MinMaxLenValidator) Validate(value string) error {
	length := len(value)
	if length < m.Min || length > m.Max {
		return fmt.Errorf("длина должна быть между %d и %d", m.Min, m.Max)
	}
	return nil
}

type PatternValidator struct {
	Regex *regexp.Regexp
}

func (p *PatternValidator) ParseParams(params []string) error {
	if len(params) != 1 {
		return fmt.Errorf("pattern требует 1 параметр")
	}

	patterns := map[string]string{
		"email": `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		"phone": `^\+?[1-9][0-9]{7,14}$`,
		"uuid":  `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
	}

	pattern := params[0]
	if predefined, ok := patterns[pattern]; ok {
		pattern = predefined
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("невалидный regexp: %w", err)
	}

	p.Regex = regex
	return nil
}

func (p *PatternValidator) Validate(value string) error {
	if !p.Regex.MatchString(value) {
		return fmt.Errorf("значение не соответствует паттерну")
	}
	return nil
}
