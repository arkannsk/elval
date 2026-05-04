package file_upload

import (
	"io"
	"os"
)

type UploadRequest struct {
	// @oa:description User avatar
	Avatar *os.File

	// @oa:description Raw data stream
	Payload io.Reader

	// @oa:description Base64 encoded
	Data []byte
}
