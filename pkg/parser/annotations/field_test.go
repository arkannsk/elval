package annotations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessFieldAnnotations(t *testing.T) {
	tests := []struct {
		name        string
		annotations []OaAnnotation
		expected    FieldAnnotationResult
	}{
		{
			name:        "Empty annotations",
			annotations: []OaAnnotation{},
			expected: FieldAnnotationResult{
				Remaining: []OaAnnotation{},
			},
		},
		{
			name: "Rewrite ref",
			annotations: []OaAnnotation{
				{Type: "rewrite.ref", Value: "#/components/schemas/User"},
			},
			expected: FieldAnnotationResult{
				RewriteRef: "#/components/schemas/User",
				Remaining:  []OaAnnotation{},
			},
		},
		{
			name: "Rewrite type",
			annotations: []OaAnnotation{
				{Type: "rewrite.type", Value: "integer"},
			},
			expected: FieldAnnotationResult{
				RewriteType: "integer",
				Remaining:   []OaAnnotation{},
			},
		},
		{
			name: "Ignore field",
			annotations: []OaAnnotation{
				{Type: "ignore", Value: ""},
			},
			expected: FieldAnnotationResult{
				IsIgnored: true,
				Remaining: []OaAnnotation{},
			},
		},
		{
			name: "OA In path",
			annotations: []OaAnnotation{
				{Type: "in", Value: "path id"},
			},
			expected: FieldAnnotationResult{
				OaIn:        "path",
				OaParamName: "id",
				Remaining:   []OaAnnotation{},
			},
		},
		{
			name: "OA In query only",
			annotations: []OaAnnotation{
				{Type: "in", Value: "query"},
			},
			expected: FieldAnnotationResult{
				OaIn:        "query",
				OaParamName: "",
				Remaining:   []OaAnnotation{},
			},
		},
		{
			name: "Remaining annotations",
			annotations: []OaAnnotation{
				{Type: "description", Value: "User ID"},
				{Type: "format", Value: "uuid"},
			},
			expected: FieldAnnotationResult{
				Remaining: []OaAnnotation{
					{Type: "description", Value: "User ID"},
					{Type: "format", Value: "uuid"},
				},
			},
		},
		{
			name: "Complex mix",
			annotations: []OaAnnotation{
				{Type: "rewrite.type", Value: "string"},
				{Type: "in", Value: "header X-Request-ID"},
				{Type: "description", Value: "Request ID"},
			},
			expected: FieldAnnotationResult{
				RewriteType: "string",
				OaIn:        "header",
				OaParamName: "X-Request-ID",
				Remaining: []OaAnnotation{
					{Type: "description", Value: "Request ID"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessFieldAnnotations(tt.annotations)
			assert.Equal(t, tt.expected.RewriteRef, result.RewriteRef)
			assert.Equal(t, tt.expected.RewriteType, result.RewriteType)
			assert.Equal(t, tt.expected.IsIgnored, result.IsIgnored)
			assert.Equal(t, tt.expected.OaIn, result.OaIn)
			assert.Equal(t, tt.expected.OaParamName, result.OaParamName)
			assert.Equal(t, tt.expected.Remaining, result.Remaining)
		})
	}
}

func TestRewriteTypeToOa(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"string", "string"},
		{"text", "string"},
		{"bool", "boolean"},
		{"int", "integer"},
		{"float", "number"},
		{"[]", "array"},
		{"unknown", "object"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, RewriteTypeToOa(tt.input))
		})
	}
}
