package generics

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericStruct_OaSchema(t *testing.T) {
	s := (&GenericStruct{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Struct using generic wrappers", s.Description)
	assert.Len(t, s.Properties, 4)

	assertProps(t, s, "name", "string", "Optional name")
	assertProps(t, s, "age", "integer", "Optional age")
	assertProps(t, s, "data", "string", "Result of computation")

	itemsProp, ok := s.Properties["items"]
	require.True(t, ok)
	assert.Equal(t, "array", itemsProp.Type)
	assert.Contains(t, itemsProp.Items.Ref, "Item")
}

func TestItem_OaSchema(t *testing.T) {
	s := (&Item{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 2)

	assertProps(t, s, "id", "string")
	assertProps(t, s, "name", "string")
}

func TestWithCustomGeneric_OaSchema(t *testing.T) {
	s := (&WithCustomGeneric{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 1)

	valueProp, ok := s.Properties["value"]
	require.True(t, ok)
	assert.Equal(t, "string", valueProp.Type)
	assert.Equal(t, "Custom generic field (rewritten as string)", valueProp.Description)
}

func assertProps(t *testing.T, schema *oa.Schema, name, typ string, desc ...string) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, typ, prop.Type)
	if len(desc) > 0 {
		assert.Equal(t, desc[0], prop.Description)
	}
}
