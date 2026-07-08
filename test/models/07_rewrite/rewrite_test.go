package rewrite

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithRewriteType_OaSchema(t *testing.T) {
	s := (&WithRewriteType{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Struct with type rewriting", s.Description)
	assert.Len(t, s.Properties, 3)

	assertProps(t, s, "rawjson", "string", "JSON payload documented as string")
	assertProps(t, s, "customid", "integer", "Custom ID type documented as integer")

	customListProp, ok := s.Properties["customlist"]
	require.True(t, ok)
	assert.Equal(t, "array", customListProp.Type)
}

func TestWithRewriteRef_OaSchema(t *testing.T) {
	s := (&WithRewriteRef{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 2)

	extProp, ok := s.Properties["external"]
	require.True(t, ok)
	assert.Equal(t, "object", extProp.Type)
	assert.Contains(t, extProp.Ref, "ExternalSchema")

	metaProp, ok := s.Properties["meta"]
	require.True(t, ok)
	assert.Contains(t, metaProp.Ref, "CommonMetadata")
}

func TestCommonMetadata_OaSchema(t *testing.T) {
	s := (&CommonMetadata{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Reusable metadata schema", s.Description)
	assert.Len(t, s.Properties, 2)

	assertProps(t, s, "created_at", "string", "Creation timestamp")
	assertProps(t, s, "updated_at", "string", "Last update timestamp")
}

func assertProps(t *testing.T, schema *oa.Schema, name, typ, desc string) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, typ, prop.Type)
	assert.Equal(t, desc, prop.Description)
}
