//go:build integration

package local_generic

import (
	"testing"

	"github.com/arkannsk/elval/test/integration/local_generic/model"
	"github.com/stretchr/testify/require"
)

func TestProduct_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		product Product
		wantErr bool
	}{
		{
			name: "valid product with all fields",
			product: Product{
				Name:    "Widget",
				Price:   19.99,
				InStock: true,
				Tags:    []model.Option[string]{model.Some("new"), model.Some("sale")},
				Rating:  model.Some(4.5),
			},
			wantErr: false,
		},
		{
			name: "missing required name",
			product: Product{
				Name:    "", // пусто
				Price:   10.0,
				InStock: true,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			product: Product{
				Name:    "Freebie",
				Price:   -1.0, // < 0
				InStock: true,
			},
			wantErr: true,
		},
		{
			name: "optional rating absent — should pass",
			product: Product{
				Name:    "Basic",
				Price:   5.0,
				InStock: true,
				// Rating не задан
			},
			wantErr: false,
		},
		{
			name: "slice with invalid option value",
			product: Product{
				Name:    "Tagged",
				Price:   15.0,
				InStock: true,
				Tags: []model.Option[string]{
					model.Some("valid"),
					model.Some(""), // пустая строка может нарушать min_len, если есть правило
				},
			},
			// Настройте ожидание в зависимости от ваших правил для []Option[string]
			wantErr: false, // или true, если есть валидация на элементы
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.product.Validate()

			if tt.wantErr {
				require.Error(t, err)
				// При необходимости проверьте errs.ValidationError
			} else {
				require.NoError(t, err)
			}
		})
	}
}
