package file_upload

import (
	"io"
	"os"
)

// @oa:ignore
type CustomReader struct{}

func (r CustomReader) Read(p []byte) (n int, err error) { return 0, nil }

type UploadRequest struct {
	// @oa:description User avatar
	Avatar *os.File

	// @oa:description Raw data stream
	Payload io.Reader

	// @oa:description Base64 encoded
	Data []byte

	// @oa:file
	// @oa:description "Custom avatar reader"
	Custom *CustomReader

	// @oa:stream
	// @oa:description "Custom avatar reader"
	Stream *CustomReader
}
