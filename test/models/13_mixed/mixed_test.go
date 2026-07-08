package mixed

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMegaStruct_OaSchema(t *testing.T) {
	s := (&MegaStruct{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Comprehensive example with all features", s.Description)
	assert.Len(t, s.Properties, 10)

	// name: required, min/max length
	nameProp, ok := s.Properties["name"]
	require.True(t, ok)
	assert.Equal(t, "string", nameProp.Type)
	assert.Equal(t, int64(2), *nameProp.MinLength)
	assert.Equal(t, int64(50), *nameProp.MaxLength)
	assert.Contains(t, s.Required, "name")

	// avatar: file
	avatarProp, ok := s.Properties["avatar"]
	require.True(t, ok)
	assert.Equal(t, "string", avatarProp.Type)
	assert.Equal(t, "binary", avatarProp.Format)
	assert.Contains(t, s.Required, "avatar")

	// thumbnail: byte
	thumbProp, ok := s.Properties["thumbnail"]
	require.True(t, ok)
	assert.Equal(t, "byte", thumbProp.Format)

	// status: enum
	statusProp, ok := s.Properties["status"]
	require.True(t, ok)
	assert.ElementsMatch(t, []any{"active", "inactive"}, statusProp.Enum)
	assert.Contains(t, s.Required, "status")

	// variant: oneOf
	variantProp, ok := s.Properties["variant"]
	require.True(t, ok)
	assert.Len(t, variantProp.OneOf, 2)
	assert.Contains(t, variantProp.OneOf[0].Ref, "UserVariant")
	assert.Contains(t, variantProp.OneOf[1].Ref, "AdminVariant")

	// discriminator
	require.NotNil(t, s.Discriminator)
	assert.Equal(t, "kind", s.Discriminator.PropertyName)
	assert.Equal(t, map[string]string{
		"admin": "AdminVariant",
		"user":  "UserVariant",
	}, s.Discriminator.Mapping)
}

func TestMegaStruct_OaParams(t *testing.T) {
	params := (&MegaStruct{}).OaParams()
	require.Len(t, params, 3)

	paramMap := make(map[string]*oa.Parameter)
	for _, p := range params {
		paramMap[p.Name] = p
	}

	id, ok := paramMap["ID"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamPath, id.In)
	assert.True(t, id.Required)

	incDel, ok := paramMap["IncludeDeleted"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamQuery, incDel.In)
	assert.False(t, incDel.Required)

	apiVer, ok := paramMap["APIVersion"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamHeader, apiVer.In)
}

func TestUserVariant_OaSchema(t *testing.T) {
	s := (&UserVariant{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 2)
	assert.Contains(t, s.Required, "kind")
}

func TestAdminVariant_OaSchema(t *testing.T) {
	s := (&AdminVariant{}).OaSchema()
	require.NotNil(t, s)
	assert.Len(t, s.Properties, 2)
	assert.Contains(t, s.Required, "kind")

	levelProp, ok := s.Properties["admin_level"]
	require.True(t, ok)
	assert.Equal(t, "integer", levelProp.Type)
}
