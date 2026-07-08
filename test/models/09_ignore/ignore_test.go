package ignore

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithIgnoredField_OaSchema(t *testing.T) {
	s := (&WithIgnoredField{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Struct with ignored field", s.Description)
	assert.Len(t, s.Properties, 1)

	assertProps(t, s, "public", "string", "Public field")

	// Ignored field should not appear
	_, ok := s.Properties["secret"]
	assert.False(t, ok, "ignored field should not be in schema")
}

func TestWithOverride_OaSchema(t *testing.T) {
	s := (&WithOverride{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 1)

	prop, ok := s.Properties["file"]
	require.True(t, ok)
	assert.Equal(t, "string", prop.Type)
	assert.Equal(t, "binary", prop.Format)
}

func TestOnlyIgnoredFields_OaSchema(t *testing.T) {
	s := (&OnlyIgnoredFields{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Struct with all fields ignored", s.Description)
	assert.Len(t, s.Properties, 0)
}

func assertProps(t *testing.T, schema *oa.Schema, name, typ, desc string) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, typ, prop.Type)
	assert.Equal(t, desc, prop.Description)
}
