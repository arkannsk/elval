package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldValidator_String(t *testing.T) {
	fv := NewFieldValidator[string]("username")
	fv.AddRequired()
	fv.AddMinMax([]string{"3", "20"})

	assert.NoError(t, fv.Validate("john"))
	assert.Error(t, fv.Validate(""))
	assert.Error(t, fv.Validate("jo"))
}

func TestFieldValidator_Int(t *testing.T) {
	fv := NewFieldValidator[int]("age")
	fv.AddMinMax([]string{"18", "99"})

	assert.NoError(t, fv.Validate(25))
	assert.Error(t, fv.Validate(16))
	assert.Error(t, fv.Validate(100))
}

func TestFieldValidator_Float64(t *testing.T) {
	fv := NewFieldValidator[float64]("price")
	fv.AddMinMax([]string{"0.01", "999.99"})

	assert.NoError(t, fv.Validate(50.5))
	assert.Error(t, fv.Validate(0.0))
	assert.Error(t, fv.Validate(1000.0))
}

func TestFieldValidator_RequiredWithTypes(t *testing.T) {
	t.Run("required для string", func(t *testing.T) {
		fv := NewFieldValidator[string]("name")
		fv.AddRequired()

		assert.NoError(t, fv.Validate("John"))
		assert.Error(t, fv.Validate(""))
	})

	t.Run("required для int", func(t *testing.T) {
		fv := NewFieldValidator[int]("age")
		fv.AddRequired()

		assert.NoError(t, fv.Validate(25))
		assert.Error(t, fv.Validate(0))
	})

	t.Run("required для bool", func(t *testing.T) {
		fv := NewFieldValidator[bool]("active")
		fv.AddRequired()

		assert.NoError(t, fv.Validate(true))
		assert.Error(t, fv.Validate(false))
	})
}

func TestFieldValidator_Pattern(t *testing.T) {
	t.Run("email pattern", func(t *testing.T) {
		fv := NewFieldValidator[string]("email")
		fv.AddPattern([]string{"email"})

		assert.NoError(t, fv.Validate("test@example.com"))
		assert.Error(t, fv.Validate("invalid"))
	})

	t.Run("custom regexp", func(t *testing.T) {
		fv := NewFieldValidator[string]("code")
		fv.AddPattern([]string{`^[A-Z]{3}-\d+$`})

		assert.NoError(t, fv.Validate("ABC-123"))
		assert.Error(t, fv.Validate("abc-123"))
	})
}

func TestFieldValidator_Combined(t *testing.T) {
	fv := NewFieldValidator[string]("username")
	fv.AddRequired()
	fv.AddMinMax([]string{"3", "20"})
	fv.AddPattern([]string{`^[a-z]+$`})

	// Все правила соблюдены
	assert.NoError(t, fv.Validate("john"))

	// Пустое значение
	assert.Error(t, fv.Validate(""))

	// Слишком короткое
	assert.Error(t, fv.Validate("jo"))

	// Содержит заглавные буквы
	assert.Error(t, fv.Validate("John"))
}

func TestFieldValidator_Errors(t *testing.T) {
	t.Run("невалидные параметры min-max", func(t *testing.T) {
		fv := NewFieldValidator[int]("age")
		fv.AddMinMax([]string{"invalid"})

		err := fv.Validate(25)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибки конфигурации")
	})

	t.Run("pattern для int - ошибка", func(t *testing.T) {
		fv := NewFieldValidator[int]("age")
		fv.AddPattern([]string{"email"})

		err := fv.Validate(25)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "поддерживается только для string")
	})
}
