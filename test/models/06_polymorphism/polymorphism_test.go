package polymorphism

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShape_OaSchema(t *testing.T) {
	s := (&Shape{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "Geometric shape (polymorphic)", s.Description)
	assert.Len(t, s.OneOf, 2)
	assert.Contains(t, s.OneOf[0].Ref, "CircleShape")
	assert.Contains(t, s.OneOf[1].Ref, "RectangleShape")

	require.NotNil(t, s.Discriminator)
	assert.Equal(t, "type", s.Discriminator.PropertyName)
	assert.Equal(t, map[string]string{
		"circle":    "CircleShape",
		"rectangle": "RectangleShape",
	}, s.Discriminator.Mapping)
}

func TestCircleShape_OaSchema(t *testing.T) {
	s := (&CircleShape{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Circle geometry", s.Description)
	assert.Len(t, s.Properties, 2)

	typeProp, ok := s.Properties["type"]
	require.True(t, ok)
	assert.Equal(t, "string", typeProp.Type)
	assert.ElementsMatch(t, []any{"circle"}, typeProp.Enum)
	assert.Contains(t, s.Required, "type")

	assertProps(t, s, "radius", "number", "Radius in meters")
}

func TestRectangleShape_OaSchema(t *testing.T) {
	s := (&RectangleShape{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 3)

	typeProp, ok := s.Properties["type"]
	require.True(t, ok)
	assert.Equal(t, "string", typeProp.Type)
	assert.ElementsMatch(t, []any{"rectangle"}, typeProp.Enum)
	assert.Contains(t, s.Required, "type")

	assertProps(t, s, "width", "number", "Width in meters")
	assertProps(t, s, "height", "number", "Height in meters")
}

func TestContainer_OaSchema(t *testing.T) {
	s := (&Container{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 1)

	shapeProp, ok := s.Properties["shape"]
	require.True(t, ok)
	assert.Equal(t, "object", shapeProp.Type)
	assert.Contains(t, shapeProp.Ref, "Shape")
}

func TestOneOfExample_OaSchema(t *testing.T) {
	s := (&OneOfExample{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 1)

	valueProp, ok := s.Properties["value"]
	require.True(t, ok)
	assert.Len(t, valueProp.OneOf, 2)
	assert.Contains(t, valueProp.OneOf[0].Ref, "StringValue")
	assert.Contains(t, valueProp.OneOf[1].Ref, "NumberValue")
}

func TestStringValue_OaSchema(t *testing.T) {
	s := (&StringValue{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 1)

	assertProps(t, s, "value", "string")
}

func TestNumberValue_OaSchema(t *testing.T) {
	s := (&NumberValue{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 1)

	assertProps(t, s, "value", "number")
}

func assertProps(t *testing.T, schema *oa.Schema, name, typ string, desc ...string) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, typ, prop.Type)
	if len(desc) > 0 {
		assert.Equal(t, desc[0], prop.Description)
	}
}
