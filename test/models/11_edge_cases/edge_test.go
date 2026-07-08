package edgecases

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyStruct_OaSchema(t *testing.T) {
	s := (&EmptyStruct{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Empty struct (no fields)", s.Description)
	assert.Len(t, s.Properties, 0)
}

func TestPointerChain_OaSchema(t *testing.T) {
	s := (&PointerChain{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 3)

	assertProps(t, s, "level1", "string", "String pointer")
	_, ok := s.Properties["level2"]
	require.True(t, ok)
	level3Prop, ok := s.Properties["level3"]
	require.True(t, ok)
	assert.Equal(t, "object", level3Prop.Type)
	assert.Contains(t, level3Prop.Ref, "Nested")
}

func TestCircularRefA_OaSchema(t *testing.T) {
	s := (&CircularRefA{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 1)

	bProp, ok := s.Properties["b"]
	require.True(t, ok)
	assert.Contains(t, bProp.Ref, "CircularRefB")
}

func TestCircularRefB_OaSchema(t *testing.T) {
	s := (&CircularRefB{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 1)

	aProp, ok := s.Properties["a"]
	require.True(t, ok)
	assert.Contains(t, aProp.Ref, "CircularRefA")
}

func TestWithInterface_OaSchema(t *testing.T) {
	s := (&WithInterface{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 1)

	dataProp, ok := s.Properties["data"]
	require.True(t, ok)
	assert.Len(t, dataProp.OneOf, 2)
	assert.Contains(t, dataProp.OneOf[0].Ref, "StringEdgeValue")
	assert.Contains(t, dataProp.OneOf[1].Ref, "NumberEdgeValue")
}

func TestWithUntypedNil_OaSchema(t *testing.T) {
	s := (&WithUntypedNil{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 3)

	assertProps(t, s, "optional", "string", "Nullable string")

	childProp, ok := s.Properties["child"]
	require.True(t, ok)
	assert.Contains(t, childProp.Ref, "Nested")
}

func assertProps(t *testing.T, schema *oa.Schema, name, typ, desc string) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, typ, prop.Type)
	assert.Equal(t, desc, prop.Description)
}
