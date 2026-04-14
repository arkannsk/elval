package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Устанавливаем переменную окружения
	os.Setenv("APP_ENV", "production")

	user := &User{}

	// Создаём контекст с данными
	ctx := context.WithValue(context.Background(), "user_id", "12345")

	// Добавляем HTTP запрос в контекст
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-User-Role", "admin")
	ctx = context.WithValue(ctx, "http.request", req)

	// Применяем декораторы
	if err := user.Decorate(ctx); err != nil {
		fmt.Printf("❌ Decorator error: %v\n", err)
		return
	}

	fmt.Println("📝 После декораторов:")
	fmt.Printf("  ID: %s\n", user.ID)
	fmt.Printf("  Role: %s\n", user.Role)
	fmt.Printf("  Environment: %s\n", user.Environment)
	fmt.Printf("  CreatedAt: %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  RequestID: %s\n", user.RequestID)

	// Валидируем
	if err := user.Validate(); err != nil {
		fmt.Printf("❌ Validation error: %v\n", err)
		return
	}

	fmt.Println("\n✅ User валиден!")
}
