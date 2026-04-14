package main

import (
	"fmt"
	"strings"

	"github.com/arkannsk/elval/pkg/validator"
	"github.com/samber/mo"
)

func main() {
	// Валидный пользователь
	user1 := User{
		Name:  mo.Some("John Doe"),
		Email: mo.Some("john@example.com"),
		Phone: mo.None[string](),
	}

	if err := user1.Validate(); err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
	} else {
		fmt.Println("✅ Пользователь 1 валиден!")
	}

	// Невалидный пользователь - Name пустой
	user2 := User{
		Name:  mo.None[string](),
		Email: mo.Some("john@example.com"),
		Phone: mo.None[string](),
	}

	if err := user2.Validate(); err != nil {
		fmt.Printf("❌ Пользователь 2: %v\n", err)
	}

	// Невалидный пользователь - Email слишком короткий
	user3 := User{
		Name:  mo.Some("John Doe"),
		Email: mo.Some("jo"),
		Phone: mo.None[string](),
	}

	if err := user3.Validate(); err != nil {
		fmt.Printf("❌ Пользователь 3: %v\n", err)
	}

	// Невалидный пользователь - Email слишком длинный
	user4 := User{
		Name:  mo.Some("John Doe"),
		Email: mo.Some("verylongemail@example.com" + strings.Repeat("a", 100)),
		Phone: mo.None[string](),
	}

	if err := user4.Validate(); err != nil {
		fmt.Printf("❌ Пользователь 4: %v\n", err)
	}

	// Невалидный пользователь - Phone должен быть пустым
	user5 := User{
		Name:  mo.Some("John Doe"),
		Email: mo.Some("john@example.com"),
		Phone: mo.Some("1234567890"),
	}

	if err := user5.Validate(); err != nil {
		fmt.Printf("❌ Пользователь 5: %v\n", err)
	}

	type Config struct {
		// @evl:validate x-option-value-eq:production
		Environment mo.Option[string]
	}

	config := Config{
		Environment: mo.Some("production"),
	}

	if err := validator.ValidateCustom("x-option-value-eq", config.Environment, "production"); err != nil {
		fmt.Printf("❌ Config: %v\n", err)
	} else {
		fmt.Println("✅ Config валиден!")
	}
}
