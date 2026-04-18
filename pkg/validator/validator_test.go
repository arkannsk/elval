package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidator(t *testing.T) {
	v := New[string]("test")
	assert.NotNil(t, v)
	assert.Equal(t, "test", v.fieldName)
	assert.Empty(t, v.rules)
}

func TestAddRule(t *testing.T) {
	v := New[string]("test")
	v.AddRule(Required[string]())

	assert.Len(t, v.rules, 1)
}

func TestValidate(t *testing.T) {
	t.Run("успешная валидация", func(t *testing.T) {
		v := New[string]("name").
			AddRule(Required[string]()).
			AddRule(LenRange(3, 10))

		err := v.Validate("John")
		require.Nil(t, err)
	})

	t.Run("ошибка валидации", func(t *testing.T) {
		v := New[string]("name").
			AddRule(Required[string]()).
			AddRule(LenRange(3, 10))

		err := v.Validate("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field: name")
	})
}

func TestValidateFunc(t *testing.T) {
	err := ValidateFunc("email", "test@example.com", Required[string](), Email())
	require.Nil(t, err)

	err = ValidateFunc("email", "invalid", Required[string](), Email())
	assert.Error(t, err)
}
