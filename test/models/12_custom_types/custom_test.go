package customtypes

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithAliases_OaSchema(t *testing.T) {
	s := (&WithAliases{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Struct with type aliases", s.Description)
	assert.Len(t, s.Properties, 2)

	nameProp, ok := s.Properties["name"]
	require.True(t, ok)
	assert.Equal(t, "Aliased string", nameProp.Description)

	emailProp, ok := s.Properties["email"]
	require.True(t, ok)
	assert.Equal(t, float64(3), *emailProp.Minimum)
	assert.Regexp(t, `@`, emailProp.Pattern)
}

func TestCustomReader_OaSchema(t *testing.T) {
	s := (&CustomReader{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 1)

	dataProp, ok := s.Properties["data"]
	require.True(t, ok)
	assert.Equal(t, "string", dataProp.Type)
	assert.Equal(t, "byte", dataProp.Format)
}

func TestWithCustomReader_OaSchema(t *testing.T) {
	s := (&WithCustomReader{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 1)

	fileProp, ok := s.Properties["file"]
	require.True(t, ok)
	assert.Equal(t, "string", fileProp.Type)
	assert.Equal(t, "binary", fileProp.Format)
}

func TestEmbedStruct_OaSchema(t *testing.T) {
	s := (&EmbedStruct{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 2)

	assertProps(t, s, "id", "string", "Embedded ID")
	assertProps(t, s, "name", "string", "Embedded name")
}

func TestWithEmbed_OaSchema(t *testing.T) {
	s := (&WithEmbed{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 2)

	embedProp, ok := s.Properties["embedstruct"]
	require.True(t, ok)
	assert.Equal(t, "object", embedProp.Type)
	assert.Contains(t, embedProp.Ref, "EmbedStruct")

	assertProps(t, s, "extra", "string", "Additional field")
}

func assertProps(t *testing.T, schema *oa.Schema, name, typ, desc string) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, typ, prop.Type)
	assert.Equal(t, desc, prop.Description)
}
