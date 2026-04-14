package product

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductValidation(t *testing.T) {
	tests := []struct {
		name      string
		product   Product
		wantError bool
		errorMsg  string
	}{
		{
			name: "валидный продукт",
			product: Product{
				Status:   "active", // добавили
				Quantity: 10,
				Price:    99.99,
				Age:      25,
				Role:     "admin",
				Score:    100,
				Discount: 50,
				Tax:      20,
			},
			wantError: false,
		},
		{
			name: "невалидный статус (eq:active)",
			product: Product{
				Status:   "inactive",
				Quantity: 10,
				Price:    99.99,
				Age:      25,
				Score:    100,
				Discount: 50,
				Tax:      20,
			},
			wantError: true,
			errorMsg:  "Status",
		},
		{
			name: "цена меньше или равна 0 (gt:0)",
			product: Product{
				Status:   "active",
				Quantity: 10,
				Price:    0,
				Age:      25,
				Score:    100,
				Discount: 50,
				Tax:      20,
			},
			wantError: true,
			errorMsg:  "Price",
		},
		{
			name: "количество меньше 1 (min:1)",
			product: Product{
				Status:   "active",
				Quantity: 0,
				Price:    99.99,
				Age:      25,
				Score:    100,
				Discount: 50,
				Tax:      20,
			},
			wantError: true,
			errorMsg:  "Quantity",
		},
		{
			name: "количество больше 100 (max:100)",
			product: Product{
				Status:   "active",
				Quantity: 101,
				Price:    99.99,
				Age:      25,
				Score:    100,
				Discount: 50,
				Tax:      20,
			},
			wantError: true,
			errorMsg:  "Quantity",
		},
		{
			name: "возраст меньше 18 (gte:18)",
			product: Product{
				Status:   "active",
				Quantity: 10,
				Price:    99.99,
				Age:      16,
				Score:    100,
				Discount: 50,
				Tax:      20,
			},
			wantError: true,
			errorMsg:  "Age",
		},
		{
			name: "score равен 0 (neq:0)",
			product: Product{
				Status:   "active",
				Quantity: 10,
				Price:    99.99,
				Age:      25,
				Score:    0,
				Discount: 50,
				Tax:      20,
			},
			wantError: true,
			errorMsg:  "Score",
		},
		{
			name: "discount не меньше 100 (lt:100)",
			product: Product{
				Status:   "active",
				Quantity: 10,
				Price:    99.99,
				Age:      25,
				Score:    100,
				Discount: 100,
				Tax:      20,
			},
			wantError: true,
			errorMsg:  "Discount",
		},
		{
			name: "tax больше 50 (lte:50)",
			product: Product{
				Status:   "active",
				Quantity: 10,
				Price:    99.99,
				Age:      25,
				Score:    100,
				Discount: 50,
				Tax:      51,
			},
			wantError: true,
			errorMsg:  "Tax",
		},
		{
			name: "role опциональный - пустой",
			product: Product{
				Status:   "active",
				Quantity: 10,
				Price:    99.99,
				Age:      25,
				Role:     "",
				Score:    100,
				Discount: 50,
				Tax:      20,
			},
			wantError: false,
		},
		{
			name: "role опциональный - невалидное значение",
			product: Product{
				Status:   "active",
				Quantity: 10,
				Price:    99.99,
				Age:      25,
				Role:     "user",
				Score:    100,
				Discount: 50,
				Tax:      20,
			},
			wantError: true,
			errorMsg:  "Role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				t.Logf("Ошибка: %v", err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductValidationBoundary(t *testing.T) {
	t.Run("граничные значения", func(t *testing.T) {
		product := Product{
			Status:   "active",
			Quantity: 1,
			Price:    0.01,
			Age:      18,
			Score:    1,
			Discount: 99,
			Tax:      0,
		}
		err := product.Validate()
		assert.NoError(t, err)
	})

	t.Run("максимальные значения", func(t *testing.T) {
		product := Product{
			Status:   "active",
			Quantity: 100,
			Price:    999999.99,
			Age:      120,
			Score:    999999,
			Discount: 99,
			Tax:      50,
		}
		err := product.Validate()
		assert.NoError(t, err)
	})
}
