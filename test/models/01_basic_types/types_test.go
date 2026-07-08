package basic_types

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimplePrimitives_OaSchema(t *testing.T) {
	s := (&SimplePrimitives{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Struct with all primitive Go types", s.Description)
	assert.Len(t, s.Properties, 6)

	assertProps(t, s, "name", "string", "String field", "hello")
	assertProps(t, s, "count", "integer", "Integer field", 42)
	assertProps(t, s, "bignumber", "integer", "64-bit integer", 9223372036854775807)
	assertProps(t, s, "rate", "number", "Float value", 3.14)
	assertProps(t, s, "active", "boolean", "Boolean flag", true)
	assertProps(t, s, "createdat", "string", "Timestamp", "2024-01-15T10:30:00Z")
}

func TestWithPointers_OaSchema(t *testing.T) {
	s := (&WithPointers{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Struct with pointer fields (nullable in OpenAPI)", s.Description)
	assert.Len(t, s.Properties, 3)

	assertProps(t, s, "name", "string", "Optional name")
	assertProps(t, s, "count", "integer", "Optional count")
	assertProps(t, s, "active", "boolean", "Optional flag")
}

func TestWithDefaults_OaSchema(t *testing.T) {
	s := (&WithDefaults{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Struct with default values", s.Description)
	assert.Len(t, s.Properties, 3)

	assertProps(t, s, "status", "string", "Status with default")
	assertProps(t, s, "limit", "integer", "Limit with default")
	assertProps(t, s, "enabled", "boolean", "Enabled with default")
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
