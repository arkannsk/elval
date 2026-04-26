package main

import (
	"fmt"
	"time"

	"github.com/arkannsk/elval"
)

func main() {
	// Строки
	err := elval.Validate("john",
		elval.Required(),
		elval.MinLen(3),
		elval.MaxLen(20),
	)
	fmt.Println("String:", err)

	// Числа
	err = elval.Validate(25,
		elval.RequiredNum[int](),
		elval.Min(18),
		elval.Max(99),
	)
	fmt.Println("Int:", err)

	// Сравнения
	err = elval.Validate(75,
		elval.Gt(0),
		elval.Lte(100),
	)
	fmt.Println("Comparison:", err)

	// Email
	err = elval.Validate("test@example.com",
		elval.Required(),
		elval.Email(),
	)
	fmt.Println("Email:", err)

	// Time
	err = elval.Validate(time.Now(),
		elval.RequiredTime(),
		elval.After("2006-01-02", "2020-01-01"),
	)
	fmt.Println("Time:", err)

	// Duration
	err = elval.Validate(5*time.Second,
		elval.RequiredDuration(),
		elval.DurationMin("1s"),
		elval.DurationMax("10s"),
	)
	fmt.Println("Duration:", err)

}
