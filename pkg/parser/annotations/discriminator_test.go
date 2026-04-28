package annotations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockTarget реализует DiscriminatorTarget для тестов
type mockTarget struct {
	discriminator *OaDiscriminator
	oneOf         []string
	oneOfRefs     []string
	anyOf         []string
	anyOfRefs     []string
}

func (m *mockTarget) GetDiscriminator() *OaDiscriminator  { return m.discriminator }
func (m *mockTarget) SetDiscriminator(d *OaDiscriminator) { m.discriminator = d }
func (m *mockTarget) GetOaOneOf() []string                { return m.oneOf }
func (m *mockTarget) SetOaOneOf(v []string)               { m.oneOf = v }
func (m *mockTarget) GetOaOneOfRefs() []string            { return m.oneOfRefs }
func (m *mockTarget) SetOaOneOfRefs(v []string)           { m.oneOfRefs = v }
func (m *mockTarget) GetOaAnyOf() []string                { return m.anyOf }
func (m *mockTarget) SetOaAnyOf(v []string)               { m.anyOf = v }
func (m *mockTarget) GetOaAnyOfRefs() []string            { return m.anyOfRefs }
func (m *mockTarget) SetOaAnyOfRefs(v []string)           { m.anyOfRefs = v }

func TestExtractDiscriminatorData(t *testing.T) {
	tests := []struct {
		name          string
		annotations   []OaAnnotation
		expectedProp  string
		expectedMap   map[string]string
		expectedOneOf []string
		expectedAnyOf []string
	}{
		{
			name:          "No discriminator",
			annotations:   []OaAnnotation{{Type: "description", Value: "Test"}},
			expectedProp:  "",
			expectedMap:   nil,
			expectedOneOf: nil,
			expectedAnyOf: nil,
		},
		{
			name: "Simple discriminator",
			annotations: []OaAnnotation{
				{Type: "discriminator.propertyName", Value: "type"},
			},
			expectedProp:  "type",
			expectedMap:   map[string]string{},
			expectedOneOf: nil,
			expectedAnyOf: nil,
		},
		{
			name: "Discriminator with mapping",
			annotations: []OaAnnotation{
				{Type: "discriminator.propertyName", Value: "kind"},
				{Type: "discriminator.mapping", Value: "dog:#/components/schemas/Dog"},
				{Type: "discriminator.mapping", Value: "cat:#/components/schemas/Cat"},
			},
			expectedProp: "kind",
			expectedMap: map[string]string{
				"dog": "#/components/schemas/Dog",
				"cat": "#/components/schemas/Cat",
			},
			expectedOneOf: nil,
			expectedAnyOf: nil,
		},
		{
			name: "OneOf and AnyOf",
			annotations: []OaAnnotation{
				{Type: "oneOf", Value: "Dog,Cat"},
				{Type: "anyOf-ref", Value: "Bird,Fish"},
			},
			expectedProp:  "",
			expectedMap:   nil,
			expectedOneOf: []string{"Dog", "Cat"},
			expectedAnyOf: nil,
		},
		{
			name: "Quotes in values",
			annotations: []OaAnnotation{
				{Type: "discriminator.propertyName", Value: `"type"`},
				{Type: "discriminator.mapping", Value: `"dog":"#/components/schemas/Dog"`},
			},
			expectedProp: "type",
			expectedMap: map[string]string{
				"dog": "#/components/schemas/Dog",
			},
			expectedOneOf: nil,
			expectedAnyOf: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := &mockTarget{}

			ExtractDiscriminatorData(target, tt.annotations)

			if tt.expectedProp != "" {
				assert.NotNil(t, target.discriminator)
				assert.Equal(t, tt.expectedProp, target.discriminator.PropertyName)
				assert.Equal(t, tt.expectedMap, target.discriminator.Mapping)
			} else {
				assert.Nil(t, target.discriminator)
			}

			assert.Equal(t, tt.expectedOneOf, target.oneOf)
			assert.Equal(t, tt.expectedAnyOf, target.anyOf)
		})
	}
}

func TestParseList(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"a,b,c", []string{"a", "b", "c"}},
		{" a , b ", []string{"a", "b"}},
		{"single", []string{"single"}},
	}

	for _, tt := range tests {
		result := parseList(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}
