package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name      string
		user      User
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid user",
			user: User{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   new(25),
				Role:  "admin",
			},
			wantError: false,
		},
		{
			name: "invalid ID - not uuid",
			user: User{
				ID:    "123",
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  "admin",
			},
			wantError: true,
			errorMsg:  "ID",
		},
		{
			name: "invalid ID - empty",
			user: User{
				ID:    "",
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  "admin",
			},
			wantError: true,
			errorMsg:  "ID",
		},
		{
			name: "invalid name - too short",
			user: User{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Name:  "Jo",
				Email: "john@example.com",
				Role:  "admin",
			},
			wantError: true,
			errorMsg:  "Name",
		},
		{
			name: "invalid email",
			user: User{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Name:  "John Doe",
				Email: "invalid",
				Role:  "admin",
			},
			wantError: true,
			errorMsg:  "Email",
		},
		{
			name: "invalid role",
			user: User{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  "superuser",
			},
			wantError: true,
			errorMsg:  "Role",
		},
		{
			name: "age optional - nil",
			user: User{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   nil,
				Role:  "admin",
			},
			wantError: false,
		},
		{
			name: "age too low",
			user: User{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   ptr(16),
				Role:  "admin",
			},
			wantError: true,
			errorMsg:  "Age",
		},
		{
			name: "age too high",
			user: User{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   ptr(150),
				Role:  "admin",
			},
			wantError: true,
			errorMsg:  "Age",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				t.Logf("Error: %v", err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func ptr(i int) *int {
	return &i
}
