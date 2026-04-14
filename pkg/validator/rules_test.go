package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	rule := Required[string]()

	t.Run("непустая строка", func(t *testing.T) {
		err := rule("hello")
		assert.NoError(t, err)
	})

	t.Run("пустая строка", func(t *testing.T) {
		err := rule("")
		assert.ErrorIs(t, err, ErrRequired)
	})

	t.Run("число не ноль", func(t *testing.T) {
		ruleInt := Required[int]()
		err := ruleInt(42)
		assert.NoError(t, err)
	})

	t.Run("число ноль", func(t *testing.T) {
		ruleInt := Required[int]()
		err := ruleInt(0)
		assert.ErrorIs(t, err, ErrRequired)
	})
}

func TestOptional(t *testing.T) {
	rule := Optional[string]()

	err := rule("any value")
	assert.NoError(t, err)

	err = rule("")
	assert.NoError(t, err)
}

func TestCustom(t *testing.T) {
	customRule := Custom(func(s string) error {
		if s != "secret" {
			return assert.AnError
		}
		return nil
	})

	err := customRule("secret")
	assert.NoError(t, err)

	err = customRule("wrong")
	assert.Error(t, err)
}
