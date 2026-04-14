package main

import (
	"context"
	"net/http"
	"testing"
)

func TestUserWithoutContext(t *testing.T) {
	user := &User{}

	// Пустой контекст - декораторы не сработают
	if err := user.Decorate(context.Background()); err != nil {
		t.Errorf("Decorate error: %v", err)
	}

	// ID и Role останутся пустыми - валидация должна упасть
	if err := user.Validate(); err == nil {
		t.Error("Expected validation error, got nil")
	} else {
		t.Logf("Expected error: %v", err)
	}
}

func TestUserWithInvalidRole(t *testing.T) {
	user := &User{}

	ctx := context.Background()
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-User-Role", "superadmin") // невалидная роль
	ctx = context.WithValue(ctx, "http.request", req)
	ctx = context.WithValue(ctx, "user_id", "12345")

	user.Decorate(ctx)

	if err := user.Validate(); err == nil {
		t.Error("Expected validation error for invalid role, got nil")
	}
}
