// examples/custom/main.go
package main

import (
	"fmt"
)

func main() {
	// Валидный продукт
	product1 := Product{
		Color: "red",
		Count: 42,
	}

	if err := product1.Validate(); err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
	} else {
		fmt.Println("✅ Продукт 1 валиден!")
	}

	// Невалидный продукт (неправильный цвет)
	product2 := Product{
		Color: "purple",
		Count: 42,
	}

	if err := product2.Validate(); err != nil {
		fmt.Printf("❌ Продукт 2: %v\n", err)
	}

	// Невалидный продукт (нечетное число)
	product3 := Product{
		Color: "red",
		Count: 43,
	}

	if err := product3.Validate(); err != nil {
		fmt.Printf("❌ Продукт 3: %v\n", err)
	}

	// Невалидный продукт (оба поля невалидны)
	product4 := Product{
		Color: "purple",
		Count: 43,
	}

	if err := product4.Validate(); err != nil {
		fmt.Printf("❌ Продукт 4: %v\n", err)
	}
}
