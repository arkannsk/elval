package elval

import (
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
)

// Тестовый тип, реализующий OptionalExtractor
type testOptional[T any] struct {
	value T
	ok    bool
}

func (t testOptional[T]) Value() (T, bool) {
	return t.value, t.ok
}

// Тестовый тип, реализующий GetExtractor
type testGetOptional[T any] struct {
	value T
	ok    bool
}

func (t testGetOptional[T]) Get() (T, bool) {
	return t.value, t.ok
}

// Тестовый тип, реализующий OptionalExtractor с nil значением
type nilOptional[T any] struct{}

func (n nilOptional[T]) Value() (T, bool) {
	var zero T
	return zero, false
}

// Тестовый тип, который не реализует OptionalExtractor
type nonOptional struct{}

func TestGenericUnwrap(t *testing.T) {
	t.Run("Test with valid value", func(t *testing.T) {
		optional := testOptional[int]{value: 42, ok: true}
		result := Unwrap[int](optional)

		assert.True(t, result.IsPresent())
		value, ok := result.Value()
		assert.True(t, ok)
		assert.Equal(t, 42, value)
	})

	t.Run("Test with nil value", func(t *testing.T) {
		optional := testOptional[int]{value: 0, ok: false}
		result := Unwrap[int](optional)

		assert.False(t, result.IsPresent())
		value, ok := result.Value()
		assert.False(t, ok)
		assert.Equal(t, 0, value)
	})

	t.Run("Test with nil optional", func(t *testing.T) {
		var optional OptionalExtractor[int] = nil
		result := Unwrap[int](optional)

		assert.False(t, result.IsPresent())
		value, ok := result.Value()
		assert.False(t, ok)
		assert.Equal(t, 0, value)
	})

	t.Run("Test with nilOptional", func(t *testing.T) {
		optional := nilOptional[int]{}
		result := Unwrap[int](optional)

		assert.False(t, result.IsPresent())
		value, ok := result.Value()
		assert.False(t, ok)
		assert.Equal(t, 0, value)
	})

	t.Run("Test with string value", func(t *testing.T) {
		optional := testOptional[string]{value: "hello", ok: true}
		result := Unwrap[string](optional)

		assert.True(t, result.IsPresent())
		value, ok := result.Value()
		assert.True(t, ok)
		assert.Equal(t, "hello", value)
	})
}

func TestGenericUnwrapWithGet(t *testing.T) {
	t.Run("Test with valid value using Get", func(t *testing.T) {
		optional := testGetOptional[int]{value: 42, ok: true}
		result := UnwrapWithGet[int](optional)

		assert.True(t, result.IsPresent())
		value, ok := result.Value()
		assert.True(t, ok)
		assert.Equal(t, 42, value)
	})

	t.Run("Test with nil value using Get", func(t *testing.T) {
		optional := testGetOptional[int]{value: 0, ok: false}
		result := UnwrapWithGet[int](optional)

		assert.False(t, result.IsPresent())
		value, ok := result.Value()
		assert.False(t, ok)
		assert.Equal(t, 0, value)
	})
}

func TestGenericIsPresent(t *testing.T) {
	t.Run("Test IsPresent with present value", func(t *testing.T) {
		optional := testOptional[int]{value: 100, ok: true}
		result := Unwrap[int](optional)

		assert.True(t, result.IsPresent())
	})

	t.Run("Test IsPresent with absent value", func(t *testing.T) {
		optional := testOptional[int]{value: 0, ok: false}
		result := Unwrap[int](optional)

		assert.False(t, result.IsPresent())
	})
}

func TestGenericValue(t *testing.T) {
	t.Run("Test Value with present value", func(t *testing.T) {
		optional := testOptional[int]{value: 200, ok: true}
		result := Unwrap[int](optional)

		value, ok := result.Value()
		assert.True(t, ok)
		assert.Equal(t, 200, value)
	})

	t.Run("Test Value with absent value", func(t *testing.T) {
		optional := testOptional[int]{value: 0, ok: false}
		result := Unwrap[int](optional)

		value, ok := result.Value()
		assert.False(t, ok)
		assert.Equal(t, 0, value)
	})
}

func TestPanicSafety(t *testing.T) {
	t.Run("Test panic safety with nil interface", func(t *testing.T) {
		// Тестирование, что не происходит паники при передаче nil
		var nilOptional OptionalExtractor[int] = nil
		result := Unwrap[int](nilOptional)

		assert.False(t, result.IsPresent())
		value, ok := result.Value()
		assert.False(t, ok)
		assert.Equal(t, 0, value)
	})
}

func TestGenericUnwrapWithMoOptional(t *testing.T) {
	t.Run("Test with mo.Some value", func(t *testing.T) {
		optional := mo.Some[int](42)
		// Для mo.Option мы используем UnwrapWithGet, так как она реализует GetExtractor
		result := UnwrapWithGet[int](optional)

		assert.True(t, result.IsPresent())
		value, ok := result.Value()
		assert.True(t, ok)
		assert.Equal(t, 42, value)
	})

	t.Run("Test with mo.None value", func(t *testing.T) {
		optional := mo.Option[int]{}
		result := UnwrapWithGet[int](optional)

		assert.False(t, result.IsPresent())
		value, ok := result.Value()
		assert.False(t, ok)
		assert.Equal(t, 0, value)
	})

	t.Run("Test with string value", func(t *testing.T) {
		optional := mo.Some[string]("hello")
		result := UnwrapWithGet[string](optional)

		assert.True(t, result.IsPresent())
		value, ok := result.Value()
		assert.True(t, ok)
		assert.Equal(t, "hello", value)
	})

	t.Run("Test with bool value", func(t *testing.T) {
		optional := mo.Some[bool](true)
		result := UnwrapWithGet[bool](optional)

		assert.True(t, result.IsPresent())
		value, ok := result.Value()
		assert.True(t, ok)
		assert.Equal(t, true, value)
	})

	t.Run("Test with struct value", func(t *testing.T) {
		type TestStruct struct {
			Field string
		}
		optional := mo.Some[TestStruct](TestStruct{Field: "test"})
		result := UnwrapWithGet[TestStruct](optional)

		assert.True(t, result.IsPresent())
		value, ok := result.Value()
		assert.True(t, ok)
		assert.Equal(t, "test", value.Field)
	})

	t.Run("Test with nil mo.Option", func(t *testing.T) {
		var optional mo.Option[int]
		result := UnwrapWithGet[int](optional)

		assert.False(t, result.IsPresent())
		value, ok := result.Value()
		assert.False(t, ok)
		assert.Equal(t, 0, value)
	})
}

func TestGenericIsPresentWithMoOptional(t *testing.T) {
	t.Run("Test IsPresent with mo.Some value", func(t *testing.T) {
		optional := mo.Some[int](100)
		result := UnwrapWithGet[int](optional)

		assert.True(t, result.IsPresent())
	})

	t.Run("Test IsPresent with mo.None value", func(t *testing.T) {
		optional := mo.Option[int]{}
		result := UnwrapWithGet[int](optional)

		assert.False(t, result.IsPresent())
	})
}

func TestGenericValueWithMoOptional(t *testing.T) {
	t.Run("Test Value with mo.Some value", func(t *testing.T) {
		optional := mo.Some[int](200)
		result := UnwrapWithGet[int](optional)

		value, ok := result.Value()
		assert.True(t, ok)
		assert.Equal(t, 200, value)
	})

	t.Run("Test Value with mo.None value", func(t *testing.T) {
		optional := mo.Option[int]{}
		result := UnwrapWithGet[int](optional)

		value, ok := result.Value()
		assert.False(t, ok)
		assert.Equal(t, 0, value)
	})
}
