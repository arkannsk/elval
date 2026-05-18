package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderStatusValidation(t *testing.T) {
	tests := []struct {
		name      string
		status    string
		wantError bool
	}{
		{"валидный статус: pending", "pending", false},
		{"валидный статус: processing", "processing", false},
		{"валидный статус: shipped", "shipped", false},
		{"валидный статус: delivered", "delivered", false},
		{"валидный статус: cancelled", "cancelled", false},
		{"невалидный статус: unknown", "unknown", true},
		{"невалидный статус: active", "active", true},
		{"пустой статус", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := Order{
				Status:   tt.status,
				Priority: 1,
				Size:     "small",
			}
			err := order.Validate()
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestOrderPriorityValidation(t *testing.T) {
	tests := []struct {
		name      string
		priority  int
		wantError bool
	}{
		{"валидный приоритет: 1", 1, false},
		{"валидный приоритет: 2", 2, false},
		{"валидный приоритет: 3", 3, false},
		{"валидный приоритет: 4", 4, false},
		{"валидный приоритет: 5", 5, false},
		{"невалидный приоритет: 0", 0, true},
		{"невалидный приоритет: 6", 6, true},
		{"невалидный приоритет: -1", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := Order{
				Status:   "pending",
				Priority: tt.priority,
				Size:     "small",
			}
			err := order.Validate()
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestOrderSizeOptional(t *testing.T) {
	tests := []struct {
		name      string
		size      string
		wantError bool
	}{
		{"опциональное поле - пустое", "", false},
		{"валидный размер: small", "small", false},
		{"валидный размер: medium", "medium", false},
		{"валидный размер: large", "large", false},
		{"невалидный размер: xl", "xl", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := Order{
				Status:   "pending",
				Priority: 1,
				Size:     tt.size,
			}
			err := order.Validate()
			if tt.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "Size")
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestUserRoleValidation(t *testing.T) {
	user := User{
		Role:  "admin",
		Level: 1,
	}
	require.Nil(t, user.Validate())

	user.Role = "superuser"
	require.Error(t, user.Validate())
	assert.Contains(t, user.Validate().Error(), "Role")
}

func TestUserLevelValidation(t *testing.T) {
	tests := []struct {
		level     int
		wantError bool
	}{
		{1, false},
		{2, false},
		{3, false},
		{0, true},
		{4, true},
		{-1, true},
	}

	for _, tt := range tests {
		user := User{
			Role:  "admin",
			Level: tt.level,
		}
		err := user.Validate()
		if tt.wantError {
			require.Error(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}
