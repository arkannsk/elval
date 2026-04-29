package multi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name      string
		user      User
		wantError bool
		errorMsg  string
	}{
		{
			name: "валидный пользователь",
			user: User{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
			wantError: false,
		},
		{
			name: "пустое имя",
			user: User{
				Name:  "",
				Email: "john@example.com",
				Age:   25,
			},
			wantError: true,
			errorMsg:  "Name",
		},
		{
			name: "невалидный email",
			user: User{
				Name:  "John Doe",
				Email: "invalid",
				Age:   25,
			},
			wantError: true,
			errorMsg:  "Email",
		},
		{
			name: "возраст меньше 18",
			user: User{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   16,
			},
			wantError: true,
			errorMsg:  "Age",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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
				Status:   "active",
				Quantity: 10,
				Price:    99.99,
			},
			wantError: false,
		},
		{
			name: "невалидный статус",
			product: Product{
				Status:   "inactive",
				Quantity: 10,
				Price:    99.99,
			},
			wantError: true,
			errorMsg:  "Status",
		},
		{
			name: "количество меньше 1",
			product: Product{
				Status:   "active",
				Quantity: 0,
				Price:    99.99,
			},
			wantError: true,
			errorMsg:  "Quantity",
		},
		{
			name: "цена меньше или равна 0",
			product: Product{
				Status:   "active",
				Quantity: 10,
				Price:    0,
			},
			wantError: true,
			errorMsg:  "Price",
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
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrderValidation(t *testing.T) {
	tests := []struct {
		name      string
		order     Order
		wantError bool
		errorMsg  string
	}{
		{
			name: "валидный заказ",
			order: Order{
				ID:     "ORD-001",
				UserID: 1,
				Total:  100.50,
				Items:  []string{"item1", "item2"},
			},
			wantError: false,
		},
		{
			name: "пустой ID",
			order: Order{
				ID:     "",
				UserID: 1,
				Total:  100.50,
				Items:  []string{"item1"},
			},
			wantError: true,
			errorMsg:  "ID",
		},
		{
			name: "UserID меньше 1",
			order: Order{
				ID:     "ORD-001",
				UserID: 0,
				Total:  100.50,
				Items:  []string{"item1"},
			},
			wantError: true,
			errorMsg:  "UserID",
		},
		{
			name: "Total меньше 0.01 (опциональный)",
			order: Order{
				ID:     "ORD-001",
				UserID: 1,
				Total:  0.005,
				Items:  []string{"item1"},
			},
			wantError: true,
			errorMsg:  "Total",
		},
		{
			name: "Total пустой (опциональный)",
			order: Order{
				ID:     "ORD-001",
				UserID: 1,
				Total:  0,
				Items:  []string{"item1"},
			},
			wantError: false, // optional поле, 0 допустим
		},
		{
			name: "Items больше 100",
			order: Order{
				ID:     "ORD-001",
				UserID: 1,
				Total:  100.50,
				Items:  make([]string, 101),
			},
			wantError: true,
			errorMsg:  "Items",
		},
		{
			name: "Items пустой (опциональный)",
			order: Order{
				ID:     "ORD-001",
				UserID: 1,
				Total:  100.50,
				Items:  []string{},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.Validate()
			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestNoGenerationForConfig(t *testing.T) {
	// Проверяем что у Config нет метода Validate
	// Используем отражение, чтобы проверить отсутствие метода
	config := Config{Host: "localhost", Port: 8080}

	// Проверяем что тип Config не имеет метода Validate
	_, ok := any(&config).(interface{ Validate() error })
	assert.False(t, ok, "Config не должен иметь метод Validate")
}
