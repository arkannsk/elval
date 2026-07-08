package parser

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helperField создаёт ast.Field с заданным тегом для тестов
func helperField(tag string) *ast.Field {
	if tag == "" {
		return &ast.Field{}
	}
	return &ast.Field{Tag: &ast.BasicLit{Value: "`" + tag + "`"}}
}

func TestExtractSerializedName(t *testing.T) {
	tests := []struct {
		name      string
		tag       string
		goField   string
		expected  string
	}{
		// json tag
		{"json snake_case", `json:"user_id"`, "UserID", "user_id"},
		{"json camelCase", `json:"firstName"`, "FirstName", "firstName"},
		{"json with omitempty", `json:"email,omitempty"`, "Email", "email"},
		{"json dash (ignore)", `json:"-"`, "Ignored", ""},
		// yaml tag (no json)
		{"yaml only", `yaml:"user_name"`, "UserName", "user_name"},
		{"yaml with omitempty", `yaml:"name,omitempty"`, "Name", "name"},
		// xml tag (no json, no yaml)
		{"xml only", `xml:"record_id"`, "RecordID", "record_id"},
		// priority: json > yaml > xml
		{"json wins over yaml", `json:"full_name" yaml:"full-name"`, "FullName", "full_name"},
		{"json wins over all", `json:"title" yaml:"record_title" xml:"RecordTitle"`, "RecordTitle", "title"},
		{"yaml wins over xml", `yaml:"name" xml:"Name"`, "Name", "name"},
		// no tag — fallback
		{"no tag", "", "LastName", "lastname"},
		{"no tag PascalCase", "", "FirstName", "firstname"},
		// empty json value
		{"json empty value (falls back to Go name)", `json:""`, "Empty", "empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := helperField(tt.tag)
			result := extractSerializedName(field, tt.goField)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseStructTags(t *testing.T) {
	tests := []struct {
		name     string
		tagStr   string
		expected map[string]string
	}{
		{
			name:   "single tag",
			tagStr: `json:"user_id"`,
			expected: map[string]string{
				"json": "user_id",
			},
		},
		{
			name:   "multiple tags",
			tagStr: `json:"user_id" yaml:"user_name" xml:"UserID"`,
			expected: map[string]string{
				"json": "user_id",
				"yaml": "user_name",
				"xml":  "UserID",
			},
		},
		{
			name:   "tag with omitempty",
			tagStr: `json:"email,omitempty"`,
			expected: map[string]string{
				"json": "email,omitempty",
			},
		},
		{
			name:   "dash tag",
			tagStr: `json:"-"`,
			expected: map[string]string{
				"json": "-",
			},
		},
		{
			name:   "empty string",
			tagStr: ``,
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStructTags(tt.tagStr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
