package invalid_directives

import "time"

//go:generate go run ../../../cmd/elval-gen generate -input .

// MixedValidationErrors demonstrates various directive validation scenarios.
type MixedValidationErrors struct {
	// Valid: no diagnostics
	// @evl:validate required
	Name string `json:"name"`

	// WARNING: unknown directive
	// @evl:validate phone
	Phone string `json:"phone"`

	// WARNING: type mismatch (pattern on int)
	// @evl:validate pattern:email
	Score int `json:"score"`

	// WARNING: missing parameter for min
	// @evl:validate min
	Age int `json:"age"`

	// WARNING: invalid duration format (needs unit like ms, s, h)
	// @evl:validate min:500
	Timeout time.Duration `json:"timeout"`

	// WARNING: deprecated directive
	// @evl:validate min-max:1,100
	Range int `json:"range"`

	// WARNING: enum requires at least one value
	// @evl:validate
	Status string `json:"status"`

	// Valid: pointer + optional + min
	// @evl:validate optional
	// @evl:validate min:3
	Bio *string `json:"bio,omitempty"`
}
