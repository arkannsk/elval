package validators

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllStringValidators_OaSchema(t *testing.T) {
	s := (&AllStringValidators{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 8)

	nameProp, ok := s.Properties["name"]
	require.True(t, ok)
	assert.Equal(t, "string", nameProp.Type)
	assert.Equal(t, int64(3), *nameProp.MinLength)
	assert.Equal(t, int64(50), *nameProp.MaxLength)
	assert.Contains(t, s.Required, "name")

	emailProp, ok := s.Properties["email"]
	require.True(t, ok)
	assert.Equal(t, "string", emailProp.Type)
	assert.Regexp(t, `@`, emailProp.Pattern)

	phoneProp, ok := s.Properties["phone"]
	require.True(t, ok)
	assert.Equal(t, "phone", phoneProp.Format)

	idProp, ok := s.Properties["id"]
	require.True(t, ok)
	assert.Equal(t, "uuid", idProp.Format)

	customProp, ok := s.Properties["custompattern"]
	require.True(t, ok)
	assert.Equal(t, "^https:/", customProp.Pattern)
}

func TestAllNumericValidators_OaSchema(t *testing.T) {
	s := (&AllNumericValidators{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 8)

	percentProp, ok := s.Properties["percent"]
	require.True(t, ok)
	assert.Equal(t, "integer", percentProp.Type)
	assert.Equal(t, float64(0), *percentProp.Minimum)
	assert.Equal(t, float64(100), *percentProp.Maximum)

	posProp, ok := s.Properties["positive"]
	require.True(t, ok)
	assert.True(t, posProp.ExclusiveMinimum)

	countProp, ok := s.Properties["count"]
	require.True(t, ok)
	assert.True(t, countProp.ExclusiveMaximum)

	amountProp, ok := s.Properties["amount"]
	require.True(t, ok)
	assert.Equal(t, "number", amountProp.Type)
}

func TestAllEnumAndSliceValidators_OaSchema(t *testing.T) {
	s := (&AllEnumAndSliceValidators{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 3)

	statusProp, ok := s.Properties["status"]
	require.True(t, ok)
	assert.ElementsMatch(t, []any{"active", "inactive", "pending"}, statusProp.Enum)

	tagsProp, ok := s.Properties["tags"]
	require.True(t, ok)
	assert.Equal(t, "array", tagsProp.Type)
	assert.Equal(t, float64(1), *tagsProp.Minimum)
	assert.Equal(t, float64(10), *tagsProp.Maximum)
	assert.Contains(t, s.Required, "tags")

	fixedProp, ok := s.Properties["fixedlist"]
	require.True(t, ok)
	assert.Equal(t, int64(3), *fixedProp.MinLength)
	assert.Equal(t, int64(3), *fixedProp.MaxLength)
}

func TestDateAndDurationValidators_OaSchema(t *testing.T) {
	s := (&DateAndDurationValidators{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 3)

	assertProps(t, s, "startdate", "string", "Date after threshold")
	assertProps(t, s, "enddate", "string", "Date before deadline")
	assertProps(t, s, "timestamp", "string", "Non-zero timestamp")
}

func assertProps(t *testing.T, schema *oa.Schema, name, typ, desc string) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, typ, prop.Type)
	assert.Equal(t, desc, prop.Description)
}
