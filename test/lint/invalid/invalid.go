package invalid

import "time"

type InvalidUser struct {
	// @evl:validate required
	// @evl:validate unknown
	Name string

	// @evl:validate min
	Age int

	// @evl:validate max:abc
	Score int

	// @evl:validate pattern:invalid[regexp
	Email string

	// @evl:validate enum:
	Role string

	// @evl:validate after:2024-13-01
	CreatedAt time.Time

	// @evl:validate min:3
	// @evl:validate max:2
	// @evl:validate len:abc
	Tags []string

	// @evl:validate min:3
	// @evl:validate max:2
	Items []string

	// @evl:validate min:1.5
	// @evl:validate max:10.5
	Price float64
}
