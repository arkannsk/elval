package validator

import (
	"testing"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequired(t *testing.T) {
	rule := Required[string]()

	t.Run("непустая строка", func(t *testing.T) {
		err := rule("hello")
		require.Nil(t, err)
	})

	t.Run("пустая строка", func(t *testing.T) {
		err := rule("")
		require.NotNil(t, err)
	})

	t.Run("число не ноль", func(t *testing.T) {
		ruleInt := Required[int]()
		err := ruleInt(42)
		require.Nil(t, err)
	})

	t.Run("число ноль", func(t *testing.T) {
		ruleInt := Required[int]()
		err := ruleInt(0)
		require.NotNil(t, err)
	})
}

func TestOptional(t *testing.T) {
	rule := Optional[string]()

	err := rule("any value")
	require.Nil(t, err)

	err = rule("")
	require.Nil(t, err)
}

func TestCustom(t *testing.T) {
	customRule := Custom(func(s string) *errs.ValidationError {
		if s != "secret" {
			return errs.NewValidationError("", "", "")
		}
		return nil
	})

	err := customRule("secret")
	require.Nil(t, err)

	err = customRule("wrong")
	assert.Error(t, err)
}
