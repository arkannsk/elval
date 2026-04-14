package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceValidator_Required(t *testing.T) {
	t.Run("не nil слайс проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.Required()

		err := sv.Validate([]string{"a", "b"})
		assert.NoError(t, err)
	})

	t.Run("nil слайс не проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.Required()

		var tags []string
		err := sv.Validate(tags)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "слайс не может быть nil")
	})

	t.Run("пустой слайс проходит (не required)", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")

		err := sv.Validate([]string{})
		assert.NoError(t, err)
	})
}

func TestSliceValidator_NotZero(t *testing.T) {
	t.Run("не пустой слайс проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.NotZero()

		err := sv.Validate([]string{"a"})
		assert.NoError(t, err)
	})

	t.Run("пустой слайс не проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.NotZero()

		err := sv.Validate([]string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не может быть пустым")
	})

	t.Run("nil слайс не проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.NotZero()

		var tags []string
		err := sv.Validate(tags)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не может быть пустым")
	})
}

func TestSliceValidator_Min(t *testing.T) {
	t.Run("размер >= min проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.Min(2)

		assert.NoError(t, sv.Validate([]string{"a", "b"}))
		assert.NoError(t, sv.Validate([]string{"a", "b", "c"}))
	})

	t.Run("размер < min не проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.Min(3)

		err := sv.Validate([]string{"a", "b"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "минимальный размер 3")
	})

	t.Run("nil слайс не проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.Min(1)

		var tags []string
		err := sv.Validate(tags)
		assert.Error(t, err)
	})
}

func TestSliceValidator_Max(t *testing.T) {
	t.Run("размер <= max проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.Max(3)

		assert.NoError(t, sv.Validate([]string{"a"}))
		assert.NoError(t, sv.Validate([]string{"a", "b"}))
		assert.NoError(t, sv.Validate([]string{"a", "b", "c"}))
	})

	t.Run("размер > max не проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.Max(2)

		err := sv.Validate([]string{"a", "b", "c"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "максимальный размер 2")
	})
}

func TestSliceValidator_Len(t *testing.T) {
	t.Run("точный размер проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.Len(3)

		err := sv.Validate([]string{"a", "b", "c"})
		assert.NoError(t, err)
	})

	t.Run("неправильный размер не проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")
		sv.Len(3)

		err := sv.Validate([]string{"a", "b"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ожидался размер 3, получен 2")
	})
}

func TestSliceValidator_Each(t *testing.T) {
	t.Run("валидация каждого элемента", func(t *testing.T) {
		// Создаем валидатор для строк
		elemValidator := New[string]("tag").
			AddRule(Required[string]()).
			AddRule(MinLen(2))

		sv := NewSliceValidator[string]("tags").
			Each(elemValidator)

		tags := []string{"go", "rust", "python"}
		err := sv.Validate(tags)
		assert.NoError(t, err)
	})

	t.Run("ошибка валидации элемента", func(t *testing.T) {
		elemValidator := New[string]("tag").
			AddRule(Required[string]()).
			AddRule(MinLen(3))

		sv := NewSliceValidator[string]("tags").
			Each(elemValidator)

		tags := []string{"go", "rust"} // "go" слишком короткий
		err := sv.Validate(tags)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tags[0]")
		assert.Contains(t, err.Error(), "минимальная длина 3")
	})

	t.Run("пустой слайс с Each", func(t *testing.T) {
		elemValidator := New[string]("tag").
			AddRule(Required[string]()).
			AddRule(MinLen(2))

		sv := NewSliceValidator[string]("tags").
			Each(elemValidator)

		// Пустой слайс не вызывает валидацию элементов
		err := sv.Validate([]string{})
		assert.NoError(t, err)
	})
}

func TestSliceValidator_Combined(t *testing.T) {
	t.Run("комбинация всех правил", func(t *testing.T) {
		elemValidator := New[string]("tag").
			AddRule(Required[string]()).
			AddRule(MinLen(2))

		sv := NewSliceValidator[string]("tags").
			Required().
			NotZero().
			Min(2).
			Max(5).
			Each(elemValidator)

		tags := []string{"go", "rust", "python"}
		err := sv.Validate(tags)
		assert.NoError(t, err)
	})

	t.Run("ошибка размера", func(t *testing.T) {
		elemValidator := New[string]("tag").
			AddRule(MinLen(2))

		sv := NewSliceValidator[string]("tags").
			Min(3).
			Max(5).
			Each(elemValidator)

		tags := []string{"go", "rust"}
		err := sv.Validate(tags)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "минимальный размер 3")
	})
}

func TestSliceValidator_WithDifferentTypes(t *testing.T) {
	t.Run("слайс int", func(t *testing.T) {
		elemValidator := New[int]("number").
			AddRule(Min(1)).
			AddRule(Max(100))

		sv := NewSliceValidator[int]("numbers").
			Min(1).
			Each(elemValidator)

		numbers := []int{10, 20, 30}
		err := sv.Validate(numbers)
		assert.NoError(t, err)

		invalidNumbers := []int{0, 150}
		err = sv.Validate(invalidNumbers)
		assert.Error(t, err)
	})

	t.Run("слайс float64", func(t *testing.T) {
		elemValidator := New[float64]("price").
			AddRule(Min(0.01)).
			AddRule(Max(99.99))

		sv := NewSliceValidator[float64]("prices").
			Each(elemValidator)

		prices := []float64{10.5, 20.0, 30.75}
		err := sv.Validate(prices)
		assert.NoError(t, err)
	})

	t.Run("слайс структур", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		// Создаем валидатор для структуры
		personValidator := New[Person]("person").
			AddRule(Custom(func(p Person) error {
				if p.Name == "" {
					return assert.AnError
				}
				if p.Age < 18 {
					return assert.AnError
				}
				return nil
			}))

		sv := NewSliceValidator[Person]("people").
			Each(personValidator)

		people := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}
		err := sv.Validate(people)
		assert.NoError(t, err)
	})
}

func TestSliceValidator_Chaining(t *testing.T) {
	t.Run("проверка цепочки вызовов", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags").
			Required().
			Min(1).
			Max(10).
			NotZero()

		assert.NotNil(t, sv)

		tags := []string{"a", "b", "c"}
		err := sv.Validate(tags)
		assert.NoError(t, err)
	})
}

func TestSliceValidator_EmptyRules(t *testing.T) {
	t.Run("без правил всегда проходит", func(t *testing.T) {
		sv := NewSliceValidator[string]("tags")

		err := sv.Validate(nil)
		assert.NoError(t, err)

		err = sv.Validate([]string{})
		assert.NoError(t, err)

		err = sv.Validate([]string{"a", "b"})
		assert.NoError(t, err)
	})
}
