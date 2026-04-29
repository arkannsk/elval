package elval

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidate_StringRules(t *testing.T) {
	t.Run("Required", func(t *testing.T) {
		require.Error(t, Validate("", Required()))
		require.Nil(t, Validate("not empty", Required()))
	})

	t.Run("MinLen", func(t *testing.T) {
		require.Error(t, Validate("ab", MinLen(3)))
		require.Nil(t, Validate("abc", MinLen(3)))
	})

	t.Run("MaxLen", func(t *testing.T) {
		require.Error(t, Validate("abcd", MaxLen(3)))
		require.Nil(t, Validate("abc", MaxLen(3)))
	})

	t.Run("Email", func(t *testing.T) {
		require.Error(t, Validate("invalid-email", Email()))
		require.Nil(t, Validate("test@example.com", Email()))
	})

	t.Run("UUID", func(t *testing.T) {
		require.Error(t, Validate("not-a-uuid", UUID()))
		require.Nil(t, Validate("550e8400-e29b-41d4-a716-446655440000", UUID()))
	})
}

func TestValidate_NumericRules(t *testing.T) {
	t.Run("Min", func(t *testing.T) {
		require.Error(t, Validate[int](5, Min[int](10)))
		require.Nil(t, Validate[int](10, Min[int](10)))
		require.Nil(t, Validate[int](15, Min[int](10)))
	})

	t.Run("Max", func(t *testing.T) {
		require.Error(t, Validate[int](15, Max[int](10)))
		require.Nil(t, Validate[int](10, Max[int](10)))
	})

	t.Run("Positive", func(t *testing.T) {
		require.Error(t, Validate[int](-1, Positive[int]()))
		require.Error(t, Validate[int](0, Positive[int]()))
		require.Nil(t, Validate[int](1, Positive[int]()))
	})

	t.Run("MinMax", func(t *testing.T) {
		require.Error(t, Validate[float64](5.0, MinMax[float64](10.0, 20.0)))
		require.Nil(t, Validate[float64](15.0, MinMax[float64](10.0, 20.0)))
	})
}

func TestValidate_TimeRules(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	t.Run("After", func(t *testing.T) {
		layout := "2006-01-02"
		todayStr := now.Format(layout)

		require.Error(t, Validate(past, After(layout, todayStr)))
		require.Nil(t, Validate(future, After(layout, todayStr)))
	})
}

func TestValidate_Composers(t *testing.T) {
	t.Run("And - Basic", func(t *testing.T) {
		// Both must pass
		require.Nil(t, Validate("valid", And(Required(), MinLen(3))))
		// One fails (Required fails first)
		err := Validate("", And(Required(), MinLen(3)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "required")
	})

	t.Run("And - Second rule fails", func(t *testing.T) {
		// First passes, second fails
		err := Validate("ab", And(Required(), MinLen(3)))
		require.Error(t, err)
	})

	t.Run("Or - At least one passes", func(t *testing.T) {
		// EqString("") passes for empty string
		require.Nil(t, Validate("", Or(Required(), EqString(""))))

		// Required passes for non-empty
		require.Nil(t, Validate("x", Or(Required(), EqString("y"))))
	})

	t.Run("Or - All fail", func(t *testing.T) {
		// Both fail
		err := Validate("x", Or(EqString("a"), EqString("b")))
		require.Error(t, err)
		// ErrorIs может быть сложным с Or, проверим сообщение или факт ошибки
		require.Contains(t, err.Error(), "eq")
	})

	t.Run("IfThen - Condition True", func(t *testing.T) {
		// Condition is true, rule applies and fails
		err := Validate(5, IfThen(func(i int) bool { return i > 0 }, Min[int](10)))
		require.Error(t, err)
	})

	t.Run("IfThen - Condition False", func(t *testing.T) {
		// Condition is false, rule skipped, no error
		require.Nil(t, Validate(5, IfThen(func(i int) bool { return i < 0 }, Min[int](10))))
	})

	t.Run("Nested Composers", func(t *testing.T) {
		// Complex chain: And(Or(...), IfThen(...))
		// Value "hello":
		// 1. Or(Required(), EqString("")) -> Passes (Required is true)
		// 2. IfThen(len > 3, MinLen(5)) -> Condition true, Rule fails (len 5 == 5? No, MinLen(5) means >=5. "hello" is 5. So it passes.)
		// Let's change to MinLen(6) to make it fail.

		rule := And(
			Or(Required(), EqString("")),
			IfThen(func(s string) bool { return len(s) > 3 }, MinLen(6)),
		)

		// "hello" has length 5. Condition len>3 is true. MinLen(6) fails.
		err := Validate("hello", rule)
		require.Error(t, err)

		// "hi" has length 2. Condition len>3 is false. Rule skipped. Passes.
		require.Nil(t, Validate("hi", rule))
	})
}

func TestValidate_Custom_Chains(t *testing.T) {
	customErr := &ValidationError{
		Field:   "custom",
		Rule:    "custom-rule",
		Message: "custom error message",
	}

	rule := Custom[string](func(s string) *ValidationError {
		if s == "bad" {
			return customErr
		}
		return nil
	})

	t.Run("Custom in And", func(t *testing.T) {
		// Custom fails
		err := Validate("bad", And(Required(), rule))
		require.Error(t, err)
		require.Contains(t, err.Error(), "custom error message")

		// Custom passes, but Required fails
		err = Validate("", And(Required(), rule))
		require.Error(t, err)
		require.Contains(t, err.Error(), "required")
	})

	t.Run("Custom in Or", func(t *testing.T) {
		// Custom fails, but Required passes (if we used Required as other part, but here both are custom-like logic or mixed)
		// Let's use Or(Custom, EqString("ok"))
		orRule := Or(rule, EqString("ok"))

		// "bad" -> Custom fails, EqString("ok") fails. Error expected.
		err := Validate("bad", orRule)
		require.Error(t, err)

		// "ok" -> Custom fails, EqString("ok") passes. No error.
		require.Nil(t, Validate("ok", orRule))
	})
}

func TestValidate_Custom_Single(t *testing.T) {
	customErr := &ValidationError{
		Field:   "custom",
		Rule:    "custom-rule",
		Message: "custom error message",
	}

	rule := Custom[string](func(s string) *ValidationError {
		if s == "bad" {
			return customErr
		}
		return nil
	})

	err := Validate("bad", rule)
	require.Error(t, err)
	require.Contains(t, err.Error(), "custom error message")

	require.Nil(t, Validate("good", rule))
}
