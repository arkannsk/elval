package main

import (
	"fmt"
	"time"

	"github.com/arkannsk/elval"
)

// User структура пользователя
type User struct {
	Name      string
	Email     string
	Age       int
	Phone     string
	BirthDate time.Time
	Tags      []string
	Scores    []int
}

// Validate реализует валидацию структуры
func (u *User) Validate() error {
	var errs []error

	// Валидация Name
	if err := elval.ValidateField("Name", u.Name,
		elval.Required(),
		elval.MinLen(2),
		elval.MaxLen(50),
	); err != nil {
		errs = append(errs, err)
	}

	// Валидация Email
	if err := elval.ValidateField("Email", u.Email,
		elval.Required(),
		elval.Email(),
	); err != nil {
		errs = append(errs, err)
	}

	// Валидация Age
	if err := elval.ValidateField("Age", u.Age,
		elval.RequiredNum[int](),
		elval.Min(18),
		elval.Max(120),
	); err != nil {
		errs = append(errs, err)
	}

	// Валидация Phone (опциональный)
	if u.Phone != "" {
		if err := elval.ValidateField("Phone", u.Phone,
			elval.Phone(),
		); err != nil {
			errs = append(errs, err)
		}
	}

	// Валидация BirthDate
	if err := elval.ValidateField("BirthDate", u.BirthDate,
		elval.RequiredTime(),
		elval.After("2006-01-02", "1900-01-01"),
		elval.Before("2006-01-02", "2024-12-31"),
	); err != nil {
		errs = append(errs, err)
	}

	// Валидация Tags (слайс строк)
	tagsValidator := elval.NewSliceValidator[string]("Tags").
		Required().
		Min(1).
		Max(10)

	if err := tagsValidator.Validate(u.Tags); err != nil {
		errs = append(errs, err)
	}

	// Валидация Scores (слайс int)
	scoresValidator := elval.NewSliceValidator[int]("Scores").
		NotZero()

	if err := scoresValidator.Validate(u.Scores); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("ошибки валидации: %v", errs)
	}
	return nil
}

func main() {
	// Валидный пользователь
	user1 := User{
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       30,
		Phone:     "+1234567890",
		BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:      []string{"go", "developer"},
		Scores:    []int{100, 95, 90},
	}

	if err := user1.Validate(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Пользователь 1 валиден")
	}

	// Невалидный пользователь
	user2 := User{
		Name:      "J",
		Email:     "invalid",
		Age:       15,
		Phone:     "123",
		BirthDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:      []string{},
		Scores:    []int{},
	}

	if err := user2.Validate(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Пользователь 2 валиден")
	}
}
