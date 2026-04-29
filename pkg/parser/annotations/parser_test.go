package annotations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOaAnnotationsFromTexts(t *testing.T) {
	tests := []struct {
		name     string
		texts    []string
		expected []OaAnnotation
	}{
		{
			name:     "Empty texts",
			texts:    []string{},
			expected: nil,
		},
		{
			name: "Single description",
			texts: []string{
				"@oa:description User entity",
			},
			expected: []OaAnnotation{
				{Type: "description", Value: "User entity"},
			},
		},
		{
			name: "Description with colon separator",
			texts: []string{
				"@oa:description: User entity",
			},
			expected: []OaAnnotation{
				{Type: "description", Value: "User entity"},
			},
		},
		{
			name: "Multiple descriptions joined by newline",
			texts: []string{
				"@oa:description First line",
				"@oa:description Second line",
			},
			expected: []OaAnnotation{
				{Type: "description", Value: "First line\nSecond line"},
			},
		},
		{
			name: "Description and other annotation",
			texts: []string{
				"@oa:description User entity",
				"@oa:discriminator type",
			},
			expected: []OaAnnotation{
				{Type: "description", Value: "User entity"},
				{Type: "discriminator", Value: "type"},
			},
		},
		{
			name: "Value with quotes",
			texts: []string{
				`@oa:description "User entity"`,
			},
			expected: []OaAnnotation{
				{Type: "description", Value: "User entity"},
			},
		},
		{
			name: "Mixed separators",
			texts: []string{
				"@oa:description First part",
				"@oa:discriminator:type",
			},
			expected: []OaAnnotation{
				{Type: "description", Value: "First part"},
				{Type: "discriminator", Value: "type"},
			},
		},
		{
			name: "Flag without value",
			texts: []string{
				"@oa:deprecated",
			},
			expected: []OaAnnotation{
				{Type: "deprecated", Value: ""},
			},
		},
		{
			name: "Invalid format ignored",
			texts: []string{
				"@oa:invalid_format_without_space_or_colon",
			},
			expected: []OaAnnotation{
				{Type: "invalid_format_without_space_or_colon", Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseOaAnnotationsFromTexts(tt.texts)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseDirectivesFromTexts(t *testing.T) {
	tests := []struct {
		name     string
		texts    []string
		expected []Directive
	}{
		{
			name:     "Empty texts",
			texts:    []string{},
			expected: nil,
		},
		{
			name: "Simple required",
			texts: []string{
				"@evl:validate required",
			},
			expected: []Directive{
				{Type: "required", Raw: "@evl:validate required", Params: nil},
			},
		},
		{
			name: "Min with param",
			texts: []string{
				"@evl:validate min:5",
			},
			expected: []Directive{
				{Type: "min", Raw: "@evl:validate min:5", Params: []string{"5"}},
			},
		},
		{
			name: "Pattern directive",
			texts: []string{
				"@evl:validate pattern:^\\d+$",
			},
			expected: []Directive{
				{Type: "pattern", Raw: "@evl:validate pattern:^\\d+$", Params: []string{"^\\d+$"}},
			},
		},
		{
			name: "Multiple params",
			texts: []string{
				"@evl:validate oneof:a,b,c",
			},
			expected: []Directive{
				{Type: "oneof", Raw: "@evl:validate oneof:a,b,c", Params: []string{"a", "b", "c"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseDirectivesFromTexts(tt.texts)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCleanComment(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"// Hello", "Hello"},
		{"/* Hello */", "Hello"},
		{"// @oa:description Test", "@oa:description Test"},
		{"  //   Trimmed  ", "Trimmed"},
		{"No comment marker", "No comment marker"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, CleanComment(tt.input))
		})
	}
}

func TestTrimQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"value"`, "value"},
		{"'value'", "value"},
		{"no quotes", "no quotes"},
		{`" mixed spaces "`, "mixed spaces"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, trimQuotes(tt.input))
		})
	}
}
