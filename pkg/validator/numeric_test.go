package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMin(t *testing.T) {
	rule := Min(18)

	tests := []struct {
		name      string
		value     int
		wantError bool
	}{
		{"больше минимума", 25, false},
		{"равно минимуму", 18, false},
		{"меньше минимума", 16, true},
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

func TestMax(t *testing.T) {
	rule := Max(100)

	tests := []struct {
		name      string
		value     int
		wantError bool
	}{
		{"меньше максимума", 50, false},
		{"равно максимуму", 100, false},
		{"больше максимума", 150, true},
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

func TestMinMax(t *testing.T) {
	rule := MinMax(18, 99)

	tests := []struct {
		name      string
		value     int
		wantError bool
	}{
		{"в диапазоне", 25, false},
		{"нижняя граница", 18, false},
		{"верхняя граница", 99, false},
		{"ниже диапазона", 16, true},
		{"выше диапазона", 100, true},
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

func TestPositive(t *testing.T) {
	rule := Positive[int]()

	tests := []struct {
		name      string
		value     int
		wantError bool
	}{
		{"положительное", 42, false},
		{"ноль", 0, true},
		{"отрицательное", -5, true},
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

func TestNegative(t *testing.T) {
	rule := Negative[int]()

	tests := []struct {
		name      string
		value     int
		wantError bool
	}{
		{"отрицательное", -42, false},
		{"ноль", 0, true},
		{"положительное", 5, true},
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

func TestNotZero(t *testing.T) {
	rule := NotZero[int]()

	tests := []struct {
		name      string
		value     int
		wantError bool
	}{
		{"не ноль", 42, false},
		{"ноль", 0, true},
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
