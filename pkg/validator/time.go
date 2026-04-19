package validator

import (
	"fmt"
	"time"

	"github.com/arkannsk/elval/pkg/errs"
)

// After проверяет что время после указанной даты
func After(layout, date string) ValidationRule[time.Time] {
	// Парсим целевую дату
	target, err := time.Parse(layout, date)
	if err != nil {
		// Если layout не указан, пробуем стандартные форматы
		formats := []string{
			time.RFC3339,
			"2006-01-02",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05+07:00",
		}
		for _, f := range formats {
			if target, err = time.Parse(f, date); err == nil {
				break
			}
		}
		if err != nil {
			panic(fmt.Sprintf("invalid date format: %s", date))
		}
	}

	return func(value time.Time) *errs.ValidationError {
		if value.IsZero() {
			return nil // zero time пропускаем, используйте Required если нужно
		}
		if value.Before(target) || value.Equal(target) {
			return errs.NewValidationError("after", "date must be after %s", target.Format("2006-01-02"))
		}
		return nil
	}
}

// Before проверяет что время до указанной даты
func Before(layout, date string) ValidationRule[time.Time] {
	target, err := time.Parse(layout, date)
	if err != nil {
		formats := []string{
			time.RFC3339,
			"2006-01-02",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05+07:00",
		}
		for _, f := range formats {
			if target, err = time.Parse(f, date); err == nil {
				break
			}
		}
		if err != nil {
			panic(fmt.Sprintf("invalid date format: %s", date))
		}
	}

	return func(value time.Time) *errs.ValidationError {
		if value.IsZero() {
			return nil
		}
		if value.After(target) || value.Equal(target) {
			return errs.NewValidationError("before", "date must be before %s", target.Format("2006-01-02"))
		}
		return nil
	}
}

// TimeNotZero проверяет что время не нулевое (не IsZero())
func TimeNotZero() ValidationRule[time.Time] {
	return func(value time.Time) *errs.ValidationError {
		if value.IsZero() {
			return errs.NewValidationError("not_zero", "date must be not zero")
		}
		return nil
	}
}

// DurationMin проверяет минимальную длительность
func DurationMin(min time.Duration) ValidationRule[time.Duration] {
	return func(value time.Duration) *errs.ValidationError {
		if value < min {
			return errs.NewValidationError("min", "duration min: %v", min)
		}
		return nil
	}
}

// DurationMax проверяет максимальную длительность
func DurationMax(max time.Duration) ValidationRule[time.Duration] {
	return func(value time.Duration) *errs.ValidationError {
		if value > max {
			return errs.NewValidationError("max", "max date: %v", max)
		}
		return nil
	}
}

// DurationRange проверяет диапазон длительности
func DurationRange(min, max time.Duration) ValidationRule[time.Duration] {
	return func(value time.Duration) *errs.ValidationError {
		if value < min {
			return errs.NewValidationError("min", "duration must be lt %v", min)
		}
		if value > max {
			return errs.NewValidationError("max", "duration must be gt %v", max)
		}
		return nil
	}
}

// DurationNotZero проверяет что длительность не нулевая
func DurationNotZero() ValidationRule[time.Duration] {
	return func(value time.Duration) *errs.ValidationError {
		if value == 0 {
			return errs.NewValidationError("not_zero", "date must be not nil")
		}
		return nil
	}
}
