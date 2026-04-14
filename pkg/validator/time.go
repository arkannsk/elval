package validator

import (
	"fmt"
	"time"
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

	return func(value time.Time) error {
		if value.IsZero() {
			return nil // zero time пропускаем, используйте Required если нужно
		}
		if value.Before(target) || value.Equal(target) {
			return NewValidationError("after", "дата должна быть после %s", target.Format("2006-01-02"))
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

	return func(value time.Time) error {
		if value.IsZero() {
			return nil
		}
		if value.After(target) || value.Equal(target) {
			return NewValidationError("before", "дата должна быть до %s", target.Format("2006-01-02"))
		}
		return nil
	}
}

// TimeNotZero проверяет что время не нулевое (не IsZero())
func TimeNotZero() ValidationRule[time.Time] {
	return func(value time.Time) error {
		if value.IsZero() {
			return NewValidationError("not_zero", "дата не может быть нулевой")
		}
		return nil
	}
}

// DurationMin проверяет минимальную длительность
func DurationMin(min time.Duration) ValidationRule[time.Duration] {
	return func(value time.Duration) error {
		if value < min {
			return NewValidationError("min", "длительность должна быть не менее %v", min)
		}
		return nil
	}
}

// DurationMax проверяет максимальную длительность
func DurationMax(max time.Duration) ValidationRule[time.Duration] {
	return func(value time.Duration) error {
		if value > max {
			return NewValidationError("max", "длительность должна быть не более %v", max)
		}
		return nil
	}
}

// DurationRange проверяет диапазон длительности
func DurationRange(min, max time.Duration) ValidationRule[time.Duration] {
	return func(value time.Duration) error {
		if value < min {
			return NewValidationError("min", "длительность должна быть не менее %v", min)
		}
		if value > max {
			return NewValidationError("max", "длительность должна быть не более %v", max)
		}
		return nil
	}
}

// DurationNotZero проверяет что длительность не нулевая
func DurationNotZero() ValidationRule[time.Duration] {
	return func(value time.Duration) error {
		if value == 0 {
			return NewValidationError("not_zero", "длительность не может быть нулевой")
		}
		return nil
	}
}
