package generics

// ValidatedOption — дженерик с валидацией на внутреннем типе
// @oa:description "Option with validation on inner value"
type ValidatedOption[T any] struct {
	// @evl:validate optional
	value    T
	hasValue bool
}

// UserRequest — структура с валидированными дженерик-полями
// @oa:description "Request with validated generic fields"
type UserRequest struct {
	// @evl:validate required
	// @oa:description "User ID (required string in Option)"
	ID ValidatedOption[string]

	// @evl:validate optional
	// @evl:validate min:0
	// @evl:validate max:120
	// @oa:description "Age (optional int in Option)"
	Age ValidatedOption[int]

	// @oa:description "Tags (slice of strings)"
	Tags []string
}
