package structtags

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonTagModel_OaSchema(t *testing.T) {
	s := (&JsonTagModel{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "Struct with various json tags", s.Description)
	assert.Len(t, s.Properties, 4)

	assertProps(t, s, "user_id", "string", "Field with snake_case json tag")
	assertProps(t, s, "firstName", "string", "Field with camelCase json tag")
	assertProps(t, s, "email", "string", "Field with json omitempty")
	assertProps(t, s, "notag", "string", "Field without json tag (fallback to lowercase)")
}

func TestYamlTagModel_OaSchema(t *testing.T) {
	s := (&YamlTagModel{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 2)

	// yaml-only tag should be used when no json tag exists
	assertProps(t, s, "user_name", "string", "Only yaml tag (no json)")
	// json takes priority over yaml
	assertProps(t, s, "full_name", "string", "Both json and yaml — json wins")
}

func TestXmlTagModel_OaSchema(t *testing.T) {
	s := (&XmlTagModel{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 2)

	// xml-only tag used when no json/yaml
	assertProps(t, s, "record_id", "string", "Only xml tag (no json, no yaml)")
	// json wins over xml
	assertProps(t, s, "title", "string", "All three tags — json wins")
}

func TestMixedTagsModel_OaSchema(t *testing.T) {
	s := (&MixedTagsModel{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 4)

	assertProps(t, s, "first_name", "string", "json tag")
	assertProps(t, s, "lastname", "string", "no tag — fallback")
	assertProps(t, s, "age", "integer", "yaml only")
	assertProps(t, s, "address", "string", "json with omitempty")
}

func assertProps(t *testing.T, schema *oa.Schema, name, typ, desc string) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, typ, prop.Type)
	assert.Equal(t, desc, prop.Description)
}
