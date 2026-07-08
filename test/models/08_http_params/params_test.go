package httpparams

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryParams_OaSchema(t *testing.T) {
	s := (&QueryParams{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Request with query parameters", s.Description)
	assert.Len(t, s.Properties, 1)

	assertProps(t, s, "body_field", "string", "Request body field")
}

func TestQueryParams_OaParams(t *testing.T) {
	params := (&QueryParams{}).OaParams()
	require.Len(t, params, 4)

	paramMap := make(map[string]*oa.Parameter)
	for _, p := range params {
		paramMap[p.Name] = p
	}

	q, ok := paramMap["Query"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamQuery, q.In)
	assert.False(t, q.Required)
	assert.Equal(t, "string", q.Schema.Type)

	page, ok := paramMap["Page"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamQuery, page.In)
	assert.Equal(t, "integer", page.Schema.Type)

	limit, ok := paramMap["Limit"]
	require.True(t, ok)
	assert.Equal(t, "integer", limit.Schema.Type)

	status, ok := paramMap["Status"]
	require.True(t, ok)
	assert.Equal(t, "string", status.Schema.Type)
	assert.ElementsMatch(t, []any{"active", "inactive", "pending"}, status.Schema.Enum)
}

func TestPathParams_OaParams(t *testing.T) {
	params := (&PathParams{}).OaParams()
	require.Len(t, params, 2)

	paramMap := make(map[string]*oa.Parameter)
	for _, p := range params {
		paramMap[p.Name] = p
	}

	uid, ok := paramMap["userId"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamPath, uid.In)
	assert.True(t, uid.Required)
	assert.Equal(t, "string", uid.Schema.Type)

	rid, ok := paramMap["resource_id"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamPath, rid.In)
	assert.True(t, rid.Required)
	assert.Equal(t, "integer", rid.Schema.Type)
}

func TestHeaderParams_OaParams(t *testing.T) {
	params := (&HeaderParams{}).OaParams()
	require.Len(t, params, 2)

	paramMap := make(map[string]*oa.Parameter)
	for _, p := range params {
		paramMap[p.Name] = p
	}

	apiKey, ok := paramMap["X-API-Key"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamHeader, apiKey.In)
	assert.False(t, apiKey.Required)

	reqID, ok := paramMap["request-id"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamHeader, reqID.In)
}

func TestMixedParams_OaParams(t *testing.T) {
	params := (&MixedParams{}).OaParams()
	require.Len(t, params, 4)

	paramMap := make(map[string]*oa.Parameter)
	for _, p := range params {
		paramMap[p.Name] = p
	}

	id, ok := paramMap["id"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamPath, id.In)
	assert.True(t, id.Required)

	filter, ok := paramMap["filter"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamQuery, filter.In)

	auth, ok := paramMap["X-Auth-Token"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamHeader, auth.In)

	session, ok := paramMap["session_id"]
	require.True(t, ok)
	assert.Equal(t, oa.ParamCookie, session.In)
}

func assertProps(t *testing.T, schema *oa.Schema, name, typ, desc string) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, typ, prop.Type)
	assert.Equal(t, desc, prop.Description)
}
