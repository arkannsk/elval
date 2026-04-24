package rewrite

import "encoding/json"

//go:generate elval-gen gen -input . -openapi

type UploadRequest struct {
	// Технически это json.RawMessage или []byte,
	// но в документации мы хотим видеть просто string
	// @oa:rewrite.type string
	// @oa:description "JSON payload as a base64 encoded string"
	Payload json.RawMessage `json:"payload"`

	// Обычное поле
	UserID string `json:"user_id"`
}