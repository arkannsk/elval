package nested_structs

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddress_OaSchema(t *testing.T) {
	s := (&Address{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Postal address", s.Description)
	assert.Len(t, s.Properties, 4)

	assertProps(t, s, "street", "string", "Street address", "123 Main St")
	assertProps(t, s, "city", "string", "City name")
	assertProps(t, s, "zipcode", "string", "Postal code", "12345")
	assertProps(t, s, "country", "string", "Country code (ISO 3166-1 alpha-2)", "US")
}

func TestUserWithAddress_OaSchema(t *testing.T) {
	s := (&UserWithAddress{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "User with nested address", s.Description)
	assert.Len(t, s.Properties, 4)

	assertProps(t, s, "id", "string", "User ID")
	assertProps(t, s, "name", "string", "User name")

	billingProp, ok := s.Properties["billing"]
	require.True(t, ok)
	assert.Equal(t, "object", billingProp.Type)
	assert.Contains(t, billingProp.Ref, "Address")

	shippingProp, ok := s.Properties["shipping"]
	require.True(t, ok)
	assert.Equal(t, "object", shippingProp.Type)
	assert.Contains(t, shippingProp.Ref, "Address")
}

func TestRecursiveNode_OaSchema(t *testing.T) {
	s := (&RecursiveNode{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 3)

	assertProps(t, s, "value", "string", "Node value")

	childrenProp, ok := s.Properties["children"]
	require.True(t, ok)
	assert.Equal(t, "array", childrenProp.Type)
	assert.Contains(t, childrenProp.Items.Ref, "RecursiveNode")

	parentProp, ok := s.Properties["parent"]
	require.True(t, ok)
	assert.Equal(t, "object", parentProp.Type)
	assert.Contains(t, parentProp.Ref, "RecursiveNode")
}

func TestDeepNesting_OaSchema(t *testing.T) {
	s := (&DeepNesting{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 1)
	_, ok := s.Properties["level1"]
	require.True(t, ok)
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
