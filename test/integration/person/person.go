package person

import "time"

type Person struct {
	// @evl:validate required
	// @evl:validate min:2
	// @evl:validate max:100
	Name string

	// @evl:validate required
	// @evl:validate pattern:email
	Email string

	// @evl:validate required
	// @evl:validate min:18
	// @evl:validate max:120
	Age int

	// @evl:validate optional
	// @evl:validate pattern:phone
	Phone string

	// @evl:validate required
	// @evl:validate after:1900-01-01
	// @evl:validate before:2024-12-31
	BirthDate time.Time

	// @evl:validate required
	// @evl:validate min:1
	// @evl:validate max:10
	Tags []string

	// @evl:validate optional
	// @evl:validate not-zero
	Scores []int
}
