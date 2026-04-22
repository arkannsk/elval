package parser

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExprToString(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "Ident simple",
			expr:     &ast.Ident{Name: "string"},
			expected: "string",
		},
		{
			name:     "Ident struct",
			expr:     &ast.Ident{Name: "User"},
			expected: "User",
		},

		{
			name: "SelectorExpr time.Time",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "time"},
				Sel: &ast.Ident{Name: "Time"},
			},
			expected: "time.Time",
		},
		{
			name: "SelectorExpr geojson.Point",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "geojson"},
				Sel: &ast.Ident{Name: "Point"},
			},
			expected: "geojson.Point",
		},

		{
			name: "Pointer to Ident",
			expr: &ast.StarExpr{
				X: &ast.Ident{Name: "User"},
			},
			expected: "*User",
		},
		{
			name: "Pointer to SelectorExpr",
			expr: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "time"},
					Sel: &ast.Ident{Name: "Time"},
				},
			},
			expected: "*time.Time",
		},
		{
			name: "Slice of Ident",
			expr: &ast.ArrayType{
				Elt: &ast.Ident{Name: "string"},
			},
			expected: "[]string",
		},
		{
			name: "Slice of SelectorExpr",
			expr: &ast.ArrayType{
				Elt: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "geojson"},
					Sel: &ast.Ident{Name: "Feature"},
				},
			},
			expected: "[]geojson.Feature",
		},
		{
			name: "Slice of Pointer",
			expr: &ast.ArrayType{
				Elt: &ast.StarExpr{
					X: &ast.Ident{Name: "User"},
				},
			},
			expected: "[]*User",
		},

		{
			name: "Array fixed length int",
			expr: &ast.ArrayType{
				Len: &ast.BasicLit{Kind: token.INT, Value: "2"},
				Elt: &ast.Ident{Name: "float64"},
			},
			expected: "[2]float64",
		},
		{
			name: "Array fixed length selector",
			expr: &ast.ArrayType{
				Len: &ast.BasicLit{Kind: token.INT, Value: "3"},
				Elt: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "geojson"},
					Sel: &ast.Ident{Name: "Coordinate"},
				},
			},
			expected: "[3]geojson.Coordinate",
		},

		{
			name: "Slice of Arrays [][]float64",
			expr: &ast.ArrayType{
				Elt: &ast.ArrayType{
					Len: nil, // Slice
					Elt: &ast.Ident{Name: "float64"},
				},
			},
			expected: "[][]float64",
		},
		{
			name: "Array of Slices [2][]int",
			expr: &ast.ArrayType{
				Len: &ast.BasicLit{Kind: token.INT, Value: "2"},
				Elt: &ast.ArrayType{
					Len: nil,
					Elt: &ast.Ident{Name: "int"},
				},
			},
			expected: "[2][]int",
		},
		{
			name: "Generic single arg Option[string]",
			expr: &ast.IndexExpr{
				X:     &ast.Ident{Name: "Option"},
				Index: &ast.Ident{Name: "string"},
			},
			expected: "Option[string]",
		},
		{
			name: "Generic multi args Map[K,V]",
			expr: &ast.IndexListExpr{
				X: &ast.Ident{Name: "Map"},
				Indices: []ast.Expr{
					&ast.Ident{Name: "string"},
					&ast.Ident{Name: "int"},
				},
			},
			expected: "Map[string, int]",
		},
		{
			name: "Generic with complex type Option[geojson.Point]",
			expr: &ast.IndexExpr{
				X: &ast.Ident{Name: "Option"},
				Index: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "geojson"},
					Sel: &ast.Ident{Name: "Point"},
				},
			},
			expected: "Option[geojson.Point]",
		},
		{
			name: "Simple map",
			expr: &ast.MapType{
				Key:   &ast.Ident{Name: "string"},
				Value: &ast.Ident{Name: "any"},
			},
			expected: "map[string]any",
		},
		{
			name: "Complex map value",
			expr: &ast.MapType{
				Key: &ast.Ident{Name: "string"},
				Value: &ast.ArrayType{
					Elt: &ast.Ident{Name: "int"},
				},
			},
			expected: "map[string][]int",
		},
		{
			name: "Parenthesized type",
			expr: &ast.ParenExpr{
				X: &ast.Ident{Name: "string"},
			},
			expected: "(string)",
		},

		{
			name:     "Unknown expr type (nil check simulation via interface cast if needed, but here we test default)",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := exprToString(tt.expr)
			assert.Equal(t, tt.expected, result, "exprToString mismatch")
		})
	}
}
