package elval

import (
	"time"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/pkg/validator"
)

// Validate валидирует значение любого типа
func Validate[T any](value T, rules ...ValidationRule[T]) error {
	return validator.ValidateFunc("value", value, rules...)
}

// ValidateField валидирует поле с именем (для структур)
func ValidateField[T any](fieldName string, value T, rules ...ValidationRule[T]) error {
	return validator.ValidateFunc(fieldName, value, rules...)
}

// Алиасы
type (
	ValidationRule[T any] = validator.ValidationRule[T]
	ValidationError       = errs.ValidationError
)

// Правила для строк
func Required() ValidationRule[string] {
	return validator.Required[string]()
}

func MinLen(min int) ValidationRule[string] {
	return validator.MinLen(min)
}

func MaxLen(max int) ValidationRule[string] {
	return validator.MaxLen(max)
}

func LenRange(min, max int) ValidationRule[string] {
	return validator.LenRange(min, max)
}

func MatchRegexp(pattern string) ValidationRule[string] {
	return validator.MatchRegexp(pattern)
}

func Email() ValidationRule[string] {
	return validator.Email()
}

func Phone() ValidationRule[string] {
	return validator.Phone()
}

func UUID() ValidationRule[string] {
	return validator.UUID()
}

func Enum[T comparable]() ValidationRule[T] {
	return validator.Enum[T]()
}

func NotEmpty() ValidationRule[string] {
	return validator.NotEmpty()
}

func EqString(expected string) ValidationRule[string] {
	return validator.Eq(expected)
}

func NeqString(expected string) ValidationRule[string] {
	return validator.Neq(expected)
}

// Правила для чисел
func RequiredNum[T validator.Number]() ValidationRule[T] {
	return validator.Required[T]()
}

func Min[T validator.Number](min T) ValidationRule[T] {
	return validator.Min(min)
}

func Max[T validator.Number](max T) ValidationRule[T] {
	return validator.Max(max)
}

func MinMax[T validator.Number](min, max T) ValidationRule[T] {
	return validator.MinMax(min, max)
}

func NotZero[T validator.Number]() ValidationRule[T] {
	return validator.NotZero[T]()
}

func Positive[T validator.Number]() ValidationRule[T] {
	return validator.Positive[T]()
}

func Negative[T validator.Number]() ValidationRule[T] {
	return validator.Negative[T]()
}

func Eq[T validator.Number](expected T) ValidationRule[T] {
	return validator.Eq(expected)
}

func Neq[T validator.Number](expected T) ValidationRule[T] {
	return validator.Neq(expected)
}

func Lt[T validator.Number](expected T) ValidationRule[T] {
	return validator.Lt(expected)
}

func Lte[T validator.Number](expected T) ValidationRule[T] {
	return validator.Lte(expected)
}

func Gt[T validator.Number](expected T) ValidationRule[T] {
	return validator.Gt(expected)
}

func Gte[T validator.Number](expected T) ValidationRule[T] {
	return validator.Gte(expected)
}

// Правила для time.Time
func RequiredTime() ValidationRule[time.Time] {
	return validator.Required[time.Time]()
}

func TimeNotZero() ValidationRule[time.Time] {
	return validator.TimeNotZero()
}

func After(layout, date string) ValidationRule[time.Time] {
	return validator.After(layout, date)
}

func Before(layout, date string) ValidationRule[time.Time] {
	return validator.Before(layout, date)
}

// Правила для time.Duration
func RequiredDuration() ValidationRule[time.Duration] {
	return validator.Required[time.Duration]()
}

func DurationNotZero() ValidationRule[time.Duration] {
	return validator.DurationNotZero()
}

func DurationMin(min string) ValidationRule[time.Duration] {
	return validator.DurationMin(min)
}

func DurationMax(max string) ValidationRule[time.Duration] {
	return validator.DurationMax(max)
}

func DurationRange(min, max time.Duration) ValidationRule[time.Duration] {
	return validator.DurationRange(min, max)
}

// Композиторы
func And[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return validator.And(rules...)
}

func Or[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return validator.Or(rules...)
}

func IfThen[T any](condition func(T) bool, rule ValidationRule[T]) ValidationRule[T] {
	return validator.IfThen(condition, rule)
}

// Кастомное правило
func Custom[T any](fn func(T) *errs.ValidationError) ValidationRule[T] {
	return validator.Custom(fn)
}

// NewSliceValidator создает валидатор для слайсов
func NewSliceValidator[T any](fieldName string) *validator.SliceValidator[T] {
	return validator.NewSliceValidator[T](fieldName)
}
