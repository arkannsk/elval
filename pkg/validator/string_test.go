package validator

import (
	"testing"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinLen(t *testing.T) {
	rule := MinLen(3)

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"достаточно длинная", "abc", false},
		{"ровно минимум", "abc", false},
		{"слишком короткая", "ab", true},
		{"пустая строка", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestMaxLen(t *testing.T) {
	rule := MaxLen(5)

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"в пределах лимита", "abc", false},
		{"ровно лимит", "abcde", false},
		{"превышает лимит", "abcdef", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestLenRange(t *testing.T) {
	rule := LenRange(3, 10)

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"в диапазоне", "hello", false},
		{"минимальная длина", "abc", false},
		{"максимальная длина", "abcdefghij", false},
		{"меньше минимума", "ab", true},
		{"больше максимума", "abcdefghijk", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestMatchRegexp(t *testing.T) {
	rule := MatchRegexp(`^[A-Z][a-z]+$`)

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"подходит", "John", false},
		{"начинается с заглавной", "John", false},
		{"не подходит", "john", true},
		{"все заглавные", "JOHN", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	rule := Email()

	tests := []struct {
		name      string
		email     string
		wantError bool
	}{
		{"валидный email", "user@example.com", false},
		{"с точкой в домене", "user@mail.co.uk", false},
		{"с плюсом", "user+tag@example.com", false},
		{"без @", "userexample.com", true},
		{"без домена", "user@", true},
		{"без имени", "@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.email)
			if tt.wantError {
				require.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrInvalidEmail)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestPhone(t *testing.T) {
	rule := Phone()

	tests := []struct {
		name      string
		phone     string
		wantError bool
	}{
		{"с плюсом", "+1234567890", false},
		{"без плюса", "1234567890", false},
		{"минимальная длина 8 цифр", "12345678", false},
		{"максимальная длина 15 цифр", "123456789012345", false},
		{"слишком короткий (7 цифр)", "1234567", true},
		{"слишком длинный (16 цифр)", "1234567890123456", true},
		{"содержит буквы", "123abc7890", true},
		{"с пробелами", "123 456 7890", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.phone)
			if tt.wantError {
				require.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrInvalidPhone)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestUUID(t *testing.T) {
	rule := UUID()

	tests := []struct {
		name      string
		uuid      string
		wantError bool
	}{
		{"валидный UUID", "123e4567-e89b-12d3-a456-426614174000", false},
		{"невалидный", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.uuid)
			if tt.wantError {
				require.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrInvalidUUID)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
