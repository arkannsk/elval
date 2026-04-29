// Package elval provides a lightweight validation library for Go structs and values.
// It supports two modes of operation:
//  1. Static Code Generation (Recommended): Use the 'elval-gen' tool to generate optimized
//     Validate() methods at compile time. This avoids runtime reflection overhead.
//  2. Runtime Validation: Use the functions in this package directly for dynamic validation,
//     testing, or when code generation is not feasible.
package elval

import (
	"time"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/pkg/validator"
)

// Validate validates a single value against a list of rules.
// It returns an error if any rule fails.
//
// Example:
//
//	err := elval.Validate("user@example.com", elval.Email())
func Validate[T any](value T, rules ...ValidationRule[T]) error {
	return validator.ValidateFunc("value", value, rules...)
}

// ValidateField validates a field with a specific name.
// The fieldName is used in error messages to identify which field failed.
//
// Example:
//
//	err := elval.ValidateField("Email", user.Email, elval.Email())
func ValidateField[T any](fieldName string, value T, rules ...ValidationRule[T]) error {
	return validator.ValidateFunc(fieldName, value, rules...)
}

// ValidationError represents a validation failure.
type ValidationError = errs.ValidationError

// ValidationRule defines a generic rule for validating a value of type T.
type ValidationRule[T any] = validator.ValidationRule[T]

// --- String Rules ---

// Required returns a rule that ensures the string is not empty.
func Required() ValidationRule[string] {
	return validator.Required[string]()
}

// MinLen returns a rule that ensures the string length is at least min.
func MinLen(min int) ValidationRule[string] {
	return validator.MinLen(min)
}

// MaxLen returns a rule that ensures the string length is at most max.
func MaxLen(max int) ValidationRule[string] {
	return validator.MaxLen(max)
}

// LenRange returns a rule that ensures the string length is between min and max.
func LenRange(min, max int) ValidationRule[string] {
	return validator.LenRange(min, max)
}

// MatchRegexp returns a rule that ensures the string matches the given regex pattern.
func MatchRegexp(pattern string) ValidationRule[string] {
	return validator.MatchRegexp(pattern)
}

// Email returns a rule that validates standard email formats.
func Email() ValidationRule[string] {
	return validator.Email()
}

// Phone returns a rule that validates phone number formats.
func Phone() ValidationRule[string] {
	return validator.Phone()
}

// UUID returns a rule that validates UUID v4 formats.
func UUID() ValidationRule[string] {
	return validator.UUID()
}

// Enum returns a rule that ensures the value is one of the provided allowed values.
// Note: For static typing, you may need to cast or use specific enum types.
func Enum[T comparable]() ValidationRule[T] {
	return validator.Enum[T]()
}

// NotEmpty returns a rule that ensures the string is not just whitespace.
func NotEmpty() ValidationRule[string] {
	return validator.NotEmpty()
}

// EqString returns a rule that ensures the string equals the expected value.
func EqString(expected string) ValidationRule[string] {
	return validator.Eq(expected)
}

// NeqString returns a rule that ensures the string does not equal the expected value.
func NeqString(expected string) ValidationRule[string] {
	return validator.Neq(expected)
}

// --- Numeric Rules ---

// RequiredNum returns a rule that ensures the numeric value is not zero.
func RequiredNum[T validator.Number]() ValidationRule[T] {
	return validator.Required[T]()
}

// Min returns a rule that ensures the number is greater than or equal to min.
func Min[T validator.Number](min T) ValidationRule[T] {
	return validator.Min(min)
}

// Max returns a rule that ensures the number is less than or equal to max.
func Max[T validator.Number](max T) ValidationRule[T] {
	return validator.Max(max)
}

// MinMax returns a rule that ensures the number is within the range [min, max].
func MinMax[T validator.Number](min, max T) ValidationRule[T] {
	return validator.MinMax(min, max)
}

// NotZero returns a rule that ensures the number is not zero.
func NotZero[T validator.Number]() ValidationRule[T] {
	return validator.NotZero[T]()
}

// Positive returns a rule that ensures the number is strictly positive (> 0).
func Positive[T validator.Number]() ValidationRule[T] {
	return validator.Positive[T]()
}

// Negative returns a rule that ensures the number is strictly negative (< 0).
func Negative[T validator.Number]() ValidationRule[T] {
	return validator.Negative[T]()
}

// Eq returns a rule that ensures the number equals the expected value.
func Eq[T validator.Number](expected T) ValidationRule[T] {
	return validator.Eq(expected)
}

// Neq returns a rule that ensures the number does not equal the expected value.
func Neq[T validator.Number](expected T) ValidationRule[T] {
	return validator.Neq(expected)
}

// Lt returns a rule that ensures the number is strictly less than the expected value.
func Lt[T validator.Number](expected T) ValidationRule[T] {
	return validator.Lt(expected)
}

// Lte returns a rule that ensures the number is less than or equal to the expected value.
func Lte[T validator.Number](expected T) ValidationRule[T] {
	return validator.Lte(expected)
}

// Gt returns a rule that ensures the number is strictly greater than the expected value.
func Gt[T validator.Number](expected T) ValidationRule[T] {
	return validator.Gt(expected)
}

// Gte returns a rule that ensures the number is greater than or equal to the expected value.
func Gte[T validator.Number](expected T) ValidationRule[T] {
	return validator.Gte(expected)
}

// --- Time Rules ---

// RequiredTime returns a rule that ensures the time.Time value is not zero.
func RequiredTime() ValidationRule[time.Time] {
	return validator.Required[time.Time]()
}

// TimeNotZero returns a rule that ensures the time.Time value is not the zero value.
func TimeNotZero() ValidationRule[time.Time] {
	return validator.TimeNotZero()
}

// After returns a rule that ensures the time is after the specified date.
// Layout follows time.Parse format (e.g., "2006-01-02").
func After(layout, date string) ValidationRule[time.Time] {
	return validator.After(layout, date)
}

// Before returns a rule that ensures the time is before the specified date.
// Layout follows time.Parse format (e.g., "2006-01-02").
func Before(layout, date string) ValidationRule[time.Time] {
	return validator.Before(layout, date)
}

// --- Duration Rules ---

// RequiredDuration returns a rule that ensures the duration is not zero.
func RequiredDuration() ValidationRule[time.Duration] {
	return validator.Required[time.Duration]()
}

// DurationNotZero returns a rule that ensures the duration is not zero.
func DurationNotZero() ValidationRule[time.Duration] {
	return validator.DurationNotZero()
}

// DurationMin returns a rule that ensures the duration is at least min.
// Min is parsed as a string (e.g., "1h", "30m").
func DurationMin(min string) ValidationRule[time.Duration] {
	return validator.DurationMin(min)
}

// DurationMax returns a rule that ensures the duration is at most max.
// Max is parsed as a string (e.g., "1h", "30m").
func DurationMax(max string) ValidationRule[time.Duration] {
	return validator.DurationMax(max)
}

// DurationRange returns a rule that ensures the duration is within [min, max].
func DurationRange(min, max time.Duration) ValidationRule[time.Duration] {
	return validator.DurationRange(min, max)
}

// --- Composers ---

// And combines multiple rules; all must pass.
func And[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return validator.And(rules...)
}

// Or combines multiple rules; at least one must pass.
func Or[T any](rules ...ValidationRule[T]) ValidationRule[T] {
	return validator.Or(rules...)
}

// IfThen applies a rule only if the condition function returns true.
func IfThen[T any](condition func(T) bool, rule ValidationRule[T]) ValidationRule[T] {
	return validator.IfThen(condition, rule)
}

// Custom creates a custom validation rule using a function.
func Custom[T any](fn func(T) *errs.ValidationError) ValidationRule[T] {
	return validator.Custom(fn)
}

// NewSliceValidator creates a validator for slices, allowing element-wise validation.
func NewSliceValidator[T any](fieldName string) *validator.SliceValidator[T] {
	return validator.NewSliceValidator[T](fieldName)
}
