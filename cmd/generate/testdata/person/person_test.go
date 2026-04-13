package person

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersonValidation(t *testing.T) {
	tests := []struct {
		name      string
		person    Person
		wantError bool
		errorMsg  string
	}{
		{
			name: "валидный человек",
			person: Person{
				Name:      "John Doe",
				Email:     "john@example.com",
				Age:       30,
				Phone:     "+1234567890",
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:      []string{"go", "developer"},
				Scores:    []int{100, 95},
			},
			wantError: false,
		},
		{
			name: "пустое имя",
			person: Person{
				Name:      "",
				Email:     "john@example.com",
				Age:       30,
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:      []string{"go"},
				Scores:    []int{100},
			},
			wantError: true,
			errorMsg:  "Name",
		},
		{
			name: "слишком короткое имя",
			person: Person{
				Name:      "J",
				Email:     "john@example.com",
				Age:       30,
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:      []string{"go"},
				Scores:    []int{100},
			},
			wantError: true,
			errorMsg:  "Name",
		},
		{
			name: "невалидный email",
			person: Person{
				Name:      "John Doe",
				Email:     "invalid",
				Age:       30,
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:      []string{"go"},
				Scores:    []int{100},
			},
			wantError: true,
			errorMsg:  "Email",
		},
		{
			name: "возраст меньше 18",
			person: Person{
				Name:      "John Doe",
				Email:     "john@example.com",
				Age:       16,
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:      []string{"go"},
				Scores:    []int{100},
			},
			wantError: true,
			errorMsg:  "Age",
		},
		{
			name: "пустые теги (required)",
			person: Person{
				Name:      "John Doe",
				Email:     "john@example.com",
				Age:       30,
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:      []string{},
				Scores:    []int{100},
			},
			wantError: true,
			errorMsg:  "Tags",
		},
		{
			name: "nil теги (required)",
			person: Person{
				Name:      "John Doe",
				Email:     "john@example.com",
				Age:       30,
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:      nil,
				Scores:    []int{100},
			},
			wantError: true,
			errorMsg:  "Tags",
		},
		{
			name: "пустые Scores (not-zero)",
			person: Person{
				Name:      "John Doe",
				Email:     "john@example.com",
				Age:       30,
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:      []string{"go"},
				Scores:    []int{},
			},
			wantError: true,
			errorMsg:  "Scores",
		},
		{
			name: "дата рождения в будущем (before)",
			person: Person{
				Name:      "John Doe",
				Email:     "john@example.com",
				Age:       30,
				BirthDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:      []string{"go"},
				Scores:    []int{100},
			},
			wantError: true,
			errorMsg:  "BirthDate",
		},
		{
			name: "невалидный телефон",
			person: Person{
				Name:      "John Doe",
				Email:     "john@example.com",
				Age:       30,
				Phone:     "invalid",
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:      []string{"go"},
				Scores:    []int{100},
			},
			wantError: true,
			errorMsg:  "Phone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.person.Validate()

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				t.Logf("Ошибка: %v", err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPersonValidationEdgeCases(t *testing.T) {
	t.Run("максимальное количество тегов", func(t *testing.T) {
		tags := make([]string, 10)
		for i := range tags {
			tags[i] = "tag"
		}

		person := Person{
			Name:      "John Doe",
			Email:     "john@example.com",
			Age:       30,
			BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Tags:      tags,
			Scores:    []int{100},
		}

		err := person.Validate()
		assert.NoError(t, err)
	})

	t.Run("слишком много тегов", func(t *testing.T) {
		tags := make([]string, 11)
		for i := range tags {
			tags[i] = "tag"
		}

		person := Person{
			Name:      "John Doe",
			Email:     "john@example.com",
			Age:       30,
			BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Tags:      tags,
			Scores:    []int{100},
		}

		err := person.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Tags")
	})
}
