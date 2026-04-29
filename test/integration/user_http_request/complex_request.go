package user_http_request

import "time"

// ComplexRequest демонстрирует все возможности парсинга HTTP параметров
type ComplexRequest struct {
	// --- Path Parameters ---

	// @oa:in path id
	// @oa:format uuid
	// @evl:validate required pattern:uuid
	UserID string

	// @oa:in path version
	// @evl:validate required min:1 max:3
	Version int

	// --- Query Parameters ---

	// @oa:in query page
	// @evl:validate required min:1
	Page int

	// @oa:in query limit
	// @evl:validate min:1 max:100
	Limit int

	// @oa:in query ids
	// @oa:description List of user IDs to fetch
	// @evl:validate required
	IDs []int

	// @oa:in query tags
	// @oa:description Filter by tags
	Tags []string

	// @oa:in query active
	// @oa:description Include only active users
	Active bool

	// @oa:in query created_after
	// @oa:format date-time
	// @oa:description Filter by creation time
	CreatedAfter time.Time

	// @oa:in query score
	// @oa:description Filter by minimum score
	Score float64

	// --- Header Parameters ---

	// @oa:in header X-Request-ID
	// @oa:description Unique request ID for tracing
	RequestID string

	// @oa:in header X-Tenant-ID
	// @oa:description Tenant identifier
	TenantID int64

	// @oa:in header X-Rate-Limit
	// @oa:description Custom rate limit override
	RateLimit uint32

	// --- Body / Other Fields (Ignored by ParseRequest if no @oa:in) ---

	// @oa:description Internal role filter (not from HTTP params)
	Role string

	// @oa:description Metadata map
	Metadata map[string]string
}
