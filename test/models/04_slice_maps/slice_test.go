package collections

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSliceVariations_OaSchema(t *testing.T) {
	s := (&SliceVariations{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Struct with various slice types", s.Description)
	assert.Len(t, s.Properties, 5)

	tagsProp, ok := s.Properties["tags"]
	require.True(t, ok)
	assert.Equal(t, "array", tagsProp.Type)
	assert.Equal(t, "string", tagsProp.Items.Type)

	idsProp, ok := s.Properties["ids"]
	require.True(t, ok)
	assert.Equal(t, "array", idsProp.Type)
	assert.Equal(t, "integer", idsProp.Items.Type)

	itemsProp, ok := s.Properties["items"]
	require.True(t, ok)
	assert.Equal(t, "array", itemsProp.Type)
	assert.Contains(t, itemsProp.Items.Ref, "Item")

	ptrItemsProp, ok := s.Properties["ptritems"]
	require.True(t, ok)
	assert.Equal(t, "array", ptrItemsProp.Type)
	assert.Contains(t, ptrItemsProp.Items.Ref, "Item")
}

func TestItem_OaSchema(t *testing.T) {
	s := (&Item{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 2)

	assertProps(t, s, "id", "string", "Item ID")
	assertProps(t, s, "name", "string", "Item name")
}

func TestMapVariations_OaSchema(t *testing.T) {
	s := (&MapVariations{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 4)

	metadataProp, ok := s.Properties["metadata"]
	require.True(t, ok)
	assert.Equal(t, "object", metadataProp.Type)
	assert.Equal(t, "String-to-string map", metadataProp.Description)
	assert.Equal(t, "{\"key\": \"value\"}", metadataProp.Example)
	assertProps(t, s, "counts", "object", "String-to-int map")
	assertProps(t, s, "items", "object", "String-to-struct map")
	assertProps(t, s, "complex", "object", "Nested map")
}

func TestArrayFixed_OaSchema(t *testing.T) {
	s := (&ArrayFixed{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 2)

	coordsProp, ok := s.Properties["coords"]
	require.True(t, ok)
	assert.Equal(t, "array", coordsProp.Type)
	assert.Equal(t, "number", coordsProp.Items.Type)

	tagsProp, ok := s.Properties["tags"]
	require.True(t, ok)
	assert.Equal(t, "array", tagsProp.Type)
	assert.Equal(t, "string", tagsProp.Items.Type)
}

func assertProps(t *testing.T, schema *oa.Schema, name, typ, desc string, example ...any) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, typ, prop.Type)
	assert.Equal(t, desc, prop.Description)
	if len(example) > 0 {
		assert.Equal(t, example[0], prop.Example)
	}
}
