package nested

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNestedStructValidation(t *testing.T) {
	t.Run("валидные вложенные структуры", func(t *testing.T) {
		user := User{
			Name:  "John Doe",
			Email: "john@example.com",
			Address: Address{
				City:    "New York",
				Street:  "5th Avenue",
				ZipCode: "10001",
			},
		}
		err := user.Validate()
		assert.NoError(t, err)
	})

	t.Run("невалидный City в Address", func(t *testing.T) {
		user := User{
			Name:  "John Doe",
			Email: "john@example.com",
			Address: Address{
				City:   "N", // слишком короткий
				Street: "5th Avenue",
			},
		}
		err := user.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Address")
		assert.Contains(t, err.Error(), "City")
	})

	t.Run("nil указатель на структуру (optional)", func(t *testing.T) {
		user := User{
			Name:           "John Doe",
			Email:          "john@example.com",
			Address:        Address{City: "New York", Street: "5th Avenue"},
			BillingAddress: nil, // optional - допустимо
		}
		err := user.Validate()
		assert.NoError(t, err)
	})
}

func TestSliceOfStructsValidation(t *testing.T) {
	t.Run("валидный слайс структур", func(t *testing.T) {
		company := Company{
			Name: "Tech Corp",
			Addresses: []Address{
				{City: "New York", Street: "5th Avenue"},
				{City: "Boston", Street: "Main Street"},
			},
		}
		err := company.Validate()
		assert.NoError(t, err)
	})

	t.Run("пустой слайс (required)", func(t *testing.T) {
		company := Company{
			Name:      "Tech Corp",
			Addresses: []Address{},
		}
		err := company.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Addresses")
	})

	t.Run("невалидный элемент в слайсе", func(t *testing.T) {
		company := Company{
			Name: "Tech Corp",
			Addresses: []Address{
				{City: "New York", Street: "5th Avenue"},
				{City: "B", Street: "Main Street"}, // слишком короткий City
			},
		}
		err := company.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Addresses")
	})
}
