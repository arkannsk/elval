package annotations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDecoratorsFromTexts(t *testing.T) {
	tests := []struct {
		name     string
		texts    []string
		expected []Decorator
	}{
		{
			name:  "Single decorator",
			texts: []string{"@evl:decor ctx-get UserID"},
			expected: []Decorator{
				{Type: "ctx-get", Params: []string{"UserID"}, Raw: "@evl:decor ctx-get UserID"},
			},
		},
		{
			name:  "Decorator with multiple params",
			texts: []string{"@evl:decor custom a,b,c"},
			expected: []Decorator{
				{Type: "custom", Params: []string{"a", "b", "c"}, Raw: "@evl:decor custom a,b,c"},
			},
		},
		{
			name:  "Multiple decorators",
			texts: []string{"@evl:decor trim", "@evl:decor upper"},
			expected: []Decorator{
				{Type: "trim", Params: nil, Raw: "@evl:decor trim"},
				{Type: "upper", Params: nil, Raw: "@evl:decor upper"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseDecoratorsFromTexts(tt.texts)
			assert.Equal(t, tt.expected, result)
		})
	}
}
