package validator

import (
	"fmt"
	"time"

	"github.com/arkannsk/elval/pkg/errs"
)

// parseTime пытается распарсить дату по нескольким форматам
func parseTime(date string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05+07:00",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, date); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s (supported: RFC3339, 2006-01-02)", date)
}

// After проверяет, что время строго после указанной даты
func After(layout, date string) ValidationRule[time.Time] {
	target, err := parseTime(date)
	if err != nil {
		// Возвращаем правило, которое всегда будет ошибкой
		return func(value time.Time) *errs.ValidationError {
			return errs.NewValidationError(
				"",
				"after_config_error",
				"validation rule 'After' has invalid target date configuration: %v", err,
			)
		}
	}

	return func(value time.Time) *errs.ValidationError {
		if value.IsZero() {
			return nil // Zero time пропускаем, используйте Required если нужно
		}
		if !value.After(target) {
			return errs.NewValidationError("", "after", "date must be after %s", target.Format(time.RFC3339))
		}
		return nil
	}
}

// Before проверяет, что время строго до указанной даты
func Before(layout, date string) ValidationRule[time.Time] {
	target, err := parseTime(date)
	if err != nil {
		// Возвращаем правило, которое всегда будет ошибкой
		return func(value time.Time) *errs.ValidationError {
			return errs.NewValidationError(
				"",
				"before_config_error",
				"validation rule 'Before' has invalid target date configuration: %v", err,
			)
		}
	}

	return func(value time.Time) *errs.ValidationError {
		if value.IsZero() {
			return nil
		}
		if !value.Before(target) {
			return errs.NewValidationError("", "before", "date must be before %s", target.Format(time.RFC3339))
		}
		return nil
	}
}

// TimeNotZero проверяет, что время не нулевое (не IsZero())
func TimeNotZero() ValidationRule[time.Time] {
	return func(value time.Time) *errs.ValidationError {
		if value.IsZero() {
			return errs.NewValidationError("", "not_zero", "date must not be zero")
		}
		return nil
	}
}

// DurationMin проверяет минимальную длительность
func DurationMin(minStr string) ValidationRule[time.Duration] {
	min, err := time.ParseDuration(minStr)
	if err != nil {
		// Возвращаем правило, которое всегда будет ошибкой
		return func(value time.Duration) *errs.ValidationError {
			return errs.NewValidationError(
				"",
				"min_duration_config_error",
				"validation rule 'DurationMin' has invalid duration configuration: %v", err,
			)
		}
	}

	return func(value time.Duration) *errs.ValidationError {
		if value < min {
			return errs.NewValidationError("", "min_duration", "duration must be at least %s", min)
		}
		return nil
	}
}

// DurationMax проверяет максимальную длительность
func DurationMax(maxStr string) ValidationRule[time.Duration] {
	max, err := time.ParseDuration(maxStr)
	if err != nil {
		// Возвращаем правило, которое всегда будет ошибкой
		return func(value time.Duration) *errs.ValidationError {
			return errs.NewValidationError(
				"",
				"max_duration_config_error",
				"validation rule 'DurationMax' has invalid duration configuration: %v", err,
			)
		}
	}

	return func(value time.Duration) *errs.ValidationError {
		if value > max {
			return errs.NewValidationError("", "max_duration", "duration must be at most %s", max)
		}
		return nil
	}
}

// DurationRange проверяет диапазон длительности [min, max]
func DurationRange(min, max time.Duration) ValidationRule[time.Duration] {
	return func(value time.Duration) *errs.ValidationError {
		if value < min {
			return errs.NewValidationError("", "min_duration", "duration must be at least %s", min)
		}
		if value > max {
			return errs.NewValidationError("", "max_duration", "duration must be at most %s", max)
		}
		return nil
	}
}

// DurationNotZero проверяет, что длительность не нулевая
func DurationNotZero() ValidationRule[time.Duration] {
	return func(value time.Duration) *errs.ValidationError {
		if value == 0 {
			return errs.NewValidationError("", "not_zero", "duration must not be zero")
		}
		return nil
	}
}
