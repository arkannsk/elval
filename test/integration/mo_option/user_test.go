package main

import (
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
)

func validUser() User {
	return User{
		Name:     mo.Some("John Doe"),
		Email:    mo.Some("john@example.com"),
		Age:      mo.Some(25),
		Password: mo.Some("SecureP@ss1"),
		Tag:      mo.Some("admin"),
	}
}

func TestUser_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{"valid user", validUser(), false},
		{"optional fields absent", User{Name: mo.Some("Test"), Email: mo.Some("a@b.com"), Tag: mo.Some("user")}, false},

		{"missing required name", User{Email: mo.Some("a@b.com"), Tag: mo.Some("user")}, true},
		{"missing required email", User{Name: mo.Some("Test"), Tag: mo.Some("user")}, true},
		{"name too short", User{Name: mo.Some("Jo"), Email: mo.Some("a@b.com"), Tag: mo.Some("user")}, true},
		{"name too long", User{Name: mo.Some("x"), Email: mo.Some("a@b.com"), Tag: mo.Some("user")}, true}, // min:3, max:50 (тест на min)
		{"invalid email", User{Name: mo.Some("Test"), Email: mo.Some("not-email"), Tag: mo.Some("user")}, true},
		{"age too low", User{Name: mo.Some("Test"), Email: mo.Some("a@b.com"), Age: mo.Some(17), Tag: mo.Some("user")}, true},
		{"age too high", User{Name: mo.Some("Test"), Email: mo.Some("a@b.com"), Age: mo.Some(121), Tag: mo.Some("user")}, true},
		{"weak password", User{Name: mo.Some("Test"), Email: mo.Some("a@b.com"), Password: mo.Some("123"), Tag: mo.Some("user")}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.user.Validate()

			if tt.wantErr {
				require.Error(t, err, "ожидалась ошибка валидации")
			} else {
				require.NoError(t, err, "ошибки быть не должно")
			}
		})
	}
}

// Тест на работу elval.Unwrap с mo.Option (критично для универсальности)
func TestMoOption_Integration(t *testing.T) {
	t.Parallel()

	t.Run("None skips validation", func(t *testing.T) {
		u := User{
			Name:     mo.Some("Test"),
			Email:    mo.Some("a@b.com"),
			Age:      mo.Option[int]{}, // None → optional, должно пройти
			Password: mo.Option[string]{},
			Tag:      mo.Some("user"),
		}
		require.NoError(t, u.Validate(), "опциональное поле без значения не должно вызывать ошибку")
	})

	t.Run("Some triggers validation", func(t *testing.T) {
		u := User{
			Name:     mo.Some("Test"),
			Email:    mo.Some("a@b.com"),
			Age:      mo.Some(10), // < min:18 → ошибка
			Password: mo.Option[string]{},
			Tag:      mo.Some("user"),
		}
		require.Error(t, u.Validate(), "валидация должна сработать для Some с неверным значением")
	})
}
