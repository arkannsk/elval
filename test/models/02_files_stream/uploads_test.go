package files_streams

import (
	"testing"

	oa "github.com/arkannsk/elval/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandardFileTypes_OaSchema(t *testing.T) {
	s := (&StandardFileTypes{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Upload request with standard file/stream types", s.Description)
	assert.Len(t, s.Properties, 5)

	assertFileProp(t, s, "avatar", "binary", "Avatar as *os.File")
	assertFileProp(t, s, "payload", "binary", "Payload as io.Reader")
	assertFileProp(t, s, "thumbnail", "binary", "Thumbnail as io.ReadCloser")
	assertFileProp(t, s, "data", "byte", "Data as base64 []byte")
	assertFileProp(t, s, "attachment", "binary", "Multipart file")
}

func TestCustomWithAnnotations_OaSchema(t *testing.T) {
	s := (&CustomWithAnnotations{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 3)

	assertFileProp(t, s, "avatar", "binary", "Custom file reader")
	assertFileProp(t, s, "payload", "binary", "Raw stream from custom type")

	dataProp, ok := s.Properties["data"]
	require.True(t, ok)
	assert.Equal(t, "object", dataProp.Type)
	assert.Equal(t, "byte", dataProp.Format)
	assert.Contains(t, dataProp.Ref, "CustomBuffer")
}

func TestMixedRequest_OaSchema(t *testing.T) {
	s := (&MixedRequest{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 5)

	assertFileProp(t, s, "file", "binary", "Main file")
	assertProps(t, s, "userid", "string", "User ID from form")
	assertProps(t, s, "description", "string", "Optional description")

	tagsProp, ok := s.Properties["tags"]
	require.True(t, ok)
	assert.Equal(t, "array", tagsProp.Type)
	assert.Equal(t, "string", tagsProp.Items.Type)

	catsProp, ok := s.Properties["categories"]
	require.True(t, ok)
	assert.Equal(t, "array", catsProp.Type)
	assert.Equal(t, float64(1), *catsProp.Minimum)
	assert.Contains(t, s.Required, "categories")
}

func TestCustomStream_OaSchema(t *testing.T) {
	s := (&CustomStream{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 0)
}

func TestCustomBuffer_OaSchema(t *testing.T) {
	s := (&CustomBuffer{}).OaSchema()
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
	assert.Len(t, s.Properties, 0)
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

func assertFileProp(t *testing.T, schema *oa.Schema, name, format, desc string) {
	prop, ok := schema.Properties[name]
	require.True(t, ok, "property %q should exist", name)
	assert.Equal(t, "string", prop.Type)
	assert.Equal(t, format, prop.Format)
	assert.Equal(t, desc, prop.Description)
}
