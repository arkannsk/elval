# Utils Package

Package utils provides helper functions for ElVal code generation.

## Functions

## Usage

#### func  NewSliceValidator

```go
func NewSliceValidator[T any](fieldName string) *validator.SliceValidator[T]
```
NewSliceValidator creates a validator for slices, allowing element-wise
validation.

#### func  Validate

```go
func Validate[T any](value T, rules ...ValidationRule[T]) error
```
Validate validates a single value against a list of rules. It returns an error
if any rule fails.

Example:

    err := elval.Validate("user@example.com", elval.Email())

#### func  ValidateField

```go
func ValidateField[T any](fieldName string, value T, rules ...ValidationRule[T]) error
```
ValidateField validates a field with a specific name. The fieldName is used in
error messages to identify which field failed.

Example:

    err := elval.ValidateField("Email", user.Email, elval.Email())

#### type Generic

```go
type Generic[T any] struct {
}
```

Generic — универсальная обёртка для извлечения значения из любого Optional-типа.

#### func  Unwrap

```go
func Unwrap[T any](src any) Generic[T]
```
Unwrap извлекает значение из произвольной обёртки. Поддерживает методы: Value(),
Get(), Unwrap(), Ok(), GetOrZero()

#### func (Generic[T]) IsPresent

```go
func (g Generic[T]) IsPresent() bool
```

#### func (Generic[T]) Value

```go
func (g Generic[T]) Value() (T, bool)
```

#### type ValidationError

```go
type ValidationError = errs.ValidationError
```

ValidationError represents a validation failure.

#### type ValidationRule

```go
type ValidationRule[T any] = validator.ValidationRule[T]
```

ValidationRule defines a generic rule for validating a value of type T.

#### func  After

```go
func After(layout, date string) ValidationRule[time.Time]
```
After returns a rule that ensures the time is after the specified date. Layout
follows time.Parse format (e.g., "2006-01-02").

#### func  And

```go
func And[T any](rules ...ValidationRule[T]) ValidationRule[T]
```
And combines multiple rules; all must pass.

#### func  Before

```go
func Before(layout, date string) ValidationRule[time.Time]
```
Before returns a rule that ensures the time is before the specified date. Layout
follows time.Parse format (e.g., "2006-01-02").

#### func  Custom

```go
func Custom[T any](fn func(T) *errs.ValidationError) ValidationRule[T]
```
Custom creates a custom validation rule using a function.

#### func  DurationMax

```go
func DurationMax(max string) ValidationRule[time.Duration]
```
DurationMax returns a rule that ensures the duration is at most max. Max is
parsed as a string (e.g., "1h", "30m").

#### func  DurationMin

```go
func DurationMin(min string) ValidationRule[time.Duration]
```
DurationMin returns a rule that ensures the duration is at least min. Min is
parsed as a string (e.g., "1h", "30m").

#### func  DurationNotZero

```go
func DurationNotZero() ValidationRule[time.Duration]
```
DurationNotZero returns a rule that ensures the duration is not zero.

#### func  DurationRange

```go
func DurationRange(min, max time.Duration) ValidationRule[time.Duration]
```
DurationRange returns a rule that ensures the duration is within [min, max].

#### func  Email

```go
func Email() ValidationRule[string]
```
Email returns a rule that validates standard email formats.

#### func  Enum

```go
func Enum[T comparable]() ValidationRule[T]
```
Enum returns a rule that ensures the value is one of the provided allowed
values. Note: For static typing, you may need to cast or use specific enum
types.

#### func  Eq

```go
func Eq[T validator.Number](expected T) ValidationRule[T]
```
Eq returns a rule that ensures the number equals the expected value.

#### func  EqString

```go
func EqString(expected string) ValidationRule[string]
```
EqString returns a rule that ensures the string equals the expected value.

#### func  Gt

```go
func Gt[T validator.Number](expected T) ValidationRule[T]
```
Gt returns a rule that ensures the number is strictly greater than the expected
value.

#### func  Gte

```go
func Gte[T validator.Number](expected T) ValidationRule[T]
```
Gte returns a rule that ensures the number is greater than or equal to the
expected value.

#### func  IfThen

```go
func IfThen[T any](condition func(T) bool, rule ValidationRule[T]) ValidationRule[T]
```
IfThen applies a rule only if the condition function returns true.

#### func  LenRange

```go
func LenRange(min, max int) ValidationRule[string]
```
LenRange returns a rule that ensures the string length is between min and max.

#### func  Lt

```go
func Lt[T validator.Number](expected T) ValidationRule[T]
```
Lt returns a rule that ensures the number is strictly less than the expected
value.

#### func  Lte

```go
func Lte[T validator.Number](expected T) ValidationRule[T]
```
Lte returns a rule that ensures the number is less than or equal to the expected
value.

#### func  MatchRegexp

```go
func MatchRegexp(pattern string) ValidationRule[string]
```
MatchRegexp returns a rule that ensures the string matches the given regex
pattern.

#### func  Max

```go
func Max[T validator.Number](max T) ValidationRule[T]
```
Max returns a rule that ensures the number is less than or equal to max.

#### func  MaxLen

```go
func MaxLen(max int) ValidationRule[string]
```
MaxLen returns a rule that ensures the string length is at most max.

#### func  Min

```go
func Min[T validator.Number](min T) ValidationRule[T]
```
Min returns a rule that ensures the number is greater than or equal to min.

#### func  MinLen

```go
func MinLen(min int) ValidationRule[string]
```
MinLen returns a rule that ensures the string length is at least min.

#### func  MinMax

```go
func MinMax[T validator.Number](min, max T) ValidationRule[T]
```
MinMax returns a rule that ensures the number is within the range [min, max].

#### func  Negative

```go
func Negative[T validator.Number]() ValidationRule[T]
```
Negative returns a rule that ensures the number is strictly negative (< 0).

#### func  Neq

```go
func Neq[T validator.Number](expected T) ValidationRule[T]
```
Neq returns a rule that ensures the number does not equal the expected value.

#### func  NeqString

```go
func NeqString(expected string) ValidationRule[string]
```
NeqString returns a rule that ensures the string does not equal the expected
value.

#### func  NotEmpty

```go
func NotEmpty() ValidationRule[string]
```
NotEmpty returns a rule that ensures the string is not just whitespace.

#### func  NotZero

```go
func NotZero[T validator.Number]() ValidationRule[T]
```
NotZero returns a rule that ensures the number is not zero.

#### func  Or

```go
func Or[T any](rules ...ValidationRule[T]) ValidationRule[T]
```
Or combines multiple rules; at least one must pass.

#### func  Phone

```go
func Phone() ValidationRule[string]
```
Phone returns a rule that validates phone number formats.

#### func  Positive

```go
func Positive[T validator.Number]() ValidationRule[T]
```
Positive returns a rule that ensures the number is strictly positive (> 0).

#### func  Required

```go
func Required() ValidationRule[string]
```
Required returns a rule that ensures the string is not empty.

#### func  RequiredDuration

```go
func RequiredDuration() ValidationRule[time.Duration]
```
RequiredDuration returns a rule that ensures the duration is not zero.

#### func  RequiredNum

```go
func RequiredNum[T validator.Number]() ValidationRule[T]
```
RequiredNum returns a rule that ensures the numeric value is not zero.

#### func  RequiredTime

```go
func RequiredTime() ValidationRule[time.Time]
```
RequiredTime returns a rule that ensures the time.Time value is not zero.

#### func  TimeNotZero

```go
func TimeNotZero() ValidationRule[time.Time]
```
TimeNotZero returns a rule that ensures the time.Time value is not the zero
value.

#### func  UUID

```go
func UUID() ValidationRule[string]
```
UUID returns a rule that validates UUID v4 formats.
