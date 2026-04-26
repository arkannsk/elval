package user_http_request

//go:generate elval-gen gen -input . -openapi

// GetUserRequest запрос для получения пользователя
type GetUserRequest struct {
	// @oa:in path id
	// @oa:format uuid
	// @evl:validate required
	// @evl:validate pattern:uuid
	ID string

	// @oa:in query fields
	// @oa:description Comma-separated list of fields to include
	Fields []string

	// @oa:in header X-Request-ID
	// @oa:description Unique request identifier for tracing
	RequestID string

	// Обычное поле (для body или просто часть структуры)
	// @oa:description Internal user role filter
	Role string
}
