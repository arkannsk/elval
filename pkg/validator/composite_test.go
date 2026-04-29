package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnd(t *testing.T) {
	rule := And(
		Required[string](),
		LenRange(3, 10),
		MatchRegexp(`^[a-z]+$`),
	)

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"все правила соблюдены", "john", false},
		{"пустая строка", "", true},
		{"слишком короткая", "jo", true},
		{"с заглавными", "John", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestOr(t *testing.T) {
	rule := Or(
		Email(),
		Phone(),
	)

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"email", "test@example.com", false},
		{"phone", "+1234567890", false},
		{"ни то ни другое", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestIfThen(t *testing.T) {
	// Если значение не пустое, то проверяем email
	rule := IfThen(
		func(s string) bool { return s != "" },
		Email(),
	)

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"пустое - пропускаем", "", false},
		{"не пустое и валидное", "test@example.com", false},
		{"не пустое и невалидное", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule(tt.value)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
