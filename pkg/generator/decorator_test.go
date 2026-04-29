package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDecoratorCode(t *testing.T) {
	tests := []struct {
		name      string
		decType   string
		paramName string
		fieldName string
		expected  string
	}{
		{
			name:      "ctx-get",
			decType:   "ctx-get",
			paramName: "UserID",
			fieldName: "ID",
			expected: `if val := ctx.Value("UserID"); val != nil {
	if str, ok := val.(string); ok {
		v.ID = str
	}
}`,
		},
		{
			name:      "env-get",
			decType:   "env-get",
			paramName: "APP_PORT",
			fieldName: "Port",
			expected:  `v.Port = os.Getenv("APP_PORT")`,
		},
		{
			name:      "trim",
			decType:   "trim",
			paramName: "",
			fieldName: "Name",
			expected:  `v.Name = strings.TrimSpace(v.Name)`,
		},
		{
			name:      "unknown type",
			decType:   "unknown",
			paramName: "",
			fieldName: "X",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateDecoratorCode(tt.decType, tt.paramName, tt.fieldName)
			assert.Equal(t, tt.expected, result)
		})
	}
}
