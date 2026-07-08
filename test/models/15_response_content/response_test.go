package response_content

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserResponse_OaResponses(t *testing.T) {
	resp := (&UserResponse{}).OaResponses()
	require.NotNil(t, resp)
	assert.Len(t, resp, 1)

	r200, ok := resp[200]
	require.True(t, ok)
	assert.Equal(t, "github.com/arkannsk/elval/test/models/15_response_content.UserResponse", r200.Description)
	require.NotNil(t, r200.Content)
	assert.Len(t, r200.Content, 2)
	assert.Contains(t, r200.Content, "application/json")
	assert.Contains(t, r200.Content, "application/xml")
}

func TestErrorResponse_OaResponses(t *testing.T) {
	resp := (&ErrorResponse{}).OaResponses()
	require.NotNil(t, resp)
	assert.Len(t, resp, 1)

	r400, ok := resp[400]
	require.True(t, ok)
	require.NotNil(t, r400.Content)
	assert.Len(t, r400.Content, 1)
	assert.Contains(t, r400.Content, "application/json")
}

func TestCreateUserResponse_OaResponses(t *testing.T) {
	resp := (&CreateUserResponse{}).OaResponses()
	require.NotNil(t, resp)
	assert.Len(t, resp, 2)

	r200, ok := resp[200]
	require.True(t, ok)
	require.NotNil(t, r200.Content)
	assert.Len(t, r200.Content, 2)
	assert.Contains(t, r200.Content, "application/json")
	assert.Contains(t, r200.Content, "application/xml")

	r201, ok := resp[201]
	require.True(t, ok)
	require.NotNil(t, r201.Content)
	assert.Len(t, r201.Content, 1)
	assert.Contains(t, r201.Content, "application/json")
}

func TestNoMediaTypes_OaResponses(t *testing.T) {
	resp := (&NoMediaTypes{}).OaResponses()
	require.NotNil(t, resp)
	assert.Len(t, resp, 1)

	r204, ok := resp[204]
	require.True(t, ok)
	// No media types — empty content map
	assert.Len(t, r204.Content, 0)
}

func TestMultipleResponses_OaResponses(t *testing.T) {
	resp := (&MultipleResponses{}).OaResponses()
	require.NotNil(t, resp)
	assert.Len(t, resp, 3)

	r200, ok := resp[200]
	require.True(t, ok)
	require.NotNil(t, r200.Content)
	assert.Contains(t, r200.Content, "application/json")

	r401, ok := resp[401]
	require.True(t, ok)
	assert.Contains(t, r401.Content, "application/json")

	r500, ok := resp[500]
	require.True(t, ok)
	assert.Contains(t, r500.Content, "application/json")
}

func TestNoContentType_NoOaResponses(t *testing.T) {
	// NoContentType has no @oa:response annotation, so no OaResponses method is generated.
	// This test only verifies the schema is correct.
	s := (&NoContentType{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "Struct without response annotation", s.Description)
	assert.Len(t, s.Properties, 1)
}
