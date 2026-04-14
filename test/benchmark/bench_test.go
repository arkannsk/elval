package benchmark

import (
	"testing"
	"time"

	evl "github.com/arkannsk/elval/pkg/validator"
	"github.com/go-playground/validator/v10"
)

type User struct {
	Name      string
	Email     string
	Age       int
	Phone     string
	BirthDate time.Time
}

// Ручная валидация
func (u *User) ValidateManual() error {
	var errs []error

	if err := evl.ValidateFunc("Name", u.Name,
		evl.Required[string](),
		evl.MinLen(2),
		evl.MaxLen(50),
	); err != nil {
		errs = append(errs, err)
	}

	if err := evl.ValidateFunc("Email", u.Email,
		evl.Required[string](),
		evl.Email(),
	); err != nil {
		errs = append(errs, err)
	}

	if err := evl.ValidateFunc("Age", u.Age,
		evl.Required[int](),
		evl.Min(18),
		evl.Max(120),
	); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// go-playground/validator теги
type UserPlayground struct {
	Name      string    `validate:"required,min=2,max=50"`
	Email     string    `validate:"required,email"`
	Age       int       `validate:"required,min=18,max=120"`
	BirthDate time.Time `validate:"required"`
}

func BenchmarkElvalManual(b *testing.B) {
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.ValidateManual()
	}
}

func BenchmarkElvalGenerated(b *testing.B) {
	user := &UserGen{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.Validate()
	}
}

func BenchmarkPlayground(b *testing.B) {
	validate := validator.New()
	user := &UserPlayground{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validate.Struct(user)
	}
}

func BenchmarkPlaygroundWithCache(b *testing.B) {
	validate := validator.New()
	user := &UserPlayground{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	validate.Struct(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validate.Struct(user)
	}
}
