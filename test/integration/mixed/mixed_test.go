package mixed

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneration(t *testing.T) {
	// Проверяем что user.gen.go существует
	userGen := filepath.Join("user.gen.go")
	_, err := os.Stat(userGen)
	assert.NoError(t, err, "user.gen.go должен существовать")

	// Проверяем что product.gen.go существует
	productGen := filepath.Join("product.gen.go")
	_, err = os.Stat(productGen)
	assert.NoError(t, err, "product.gen.go должен существовать")

	// Проверяем что common.gen.go НЕ существует
	commonGen := filepath.Join("common.gen.go")
	_, err = os.Stat(commonGen)
	assert.True(t, os.IsNotExist(err), "common.gen.go не должен существовать")
}

func TestUserValidation(t *testing.T) {
	// Проверяем что метод Validate существует у User
	user := User{
		Name:  "John",
		Email: "john@example.com",
		Age:   25,
	}

	err := user.Validate()
	assert.NoError(t, err)

	// Невалидный пользователь
	invalidUser := User{
		Name:  "J",
		Email: "invalid",
		Age:   15,
	}

	err = invalidUser.Validate()
	assert.Error(t, err)
}

func TestProductValidation(t *testing.T) {
	// Проверяем что метод Validate существует у Product
	product := Product{
		Status:   "active",
		Quantity: 10,
		Price:    99.99,
	}

	err := product.Validate()
	assert.NoError(t, err)

	// Невалидный продукт
	invalidProduct := Product{
		Status:   "inactive",
		Quantity: 0,
		Price:    0,
	}

	err = invalidProduct.Validate()
	assert.Error(t, err)
}

func TestNoAnnotationsStruct(t *testing.T) {
	// Проверяем что у NoAnnotations НЕТ метода Validate
	// Используем reflection для проверки
	noAnnot := NoAnnotations{ID: 1, Name: "test"}

	// Проверяем что метод Validate отсутствует
	_, ok := any(&noAnnot).(interface{ Validate() error })
	assert.False(t, ok, "NoAnnotations не должен иметь метод Validate")
}
