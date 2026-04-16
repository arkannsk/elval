package openapi

//go:generate elval-gen gen -input . -openapi

import "time"

type User struct {
	// @evl:validate required
	// @evl:validate min:3
	// @evl:validate max:50
	// @oa:title "User Name"
	// @oa:description "Full name of the user"
	// @oa:example "John Doe"
	Name string

	// @evl:validate required
	// @evl:validate pattern:email
	// @oa:format email
	Email string

	// @evl:validate min:18
	// @evl:validate max:120
	Age int

	// @evl:validate optional
	// @evl:validate pattern:phone
	Phone string
}

type Product struct {
	// @evl:validate required
	// @evl:validate enum:active,inactive,archived
	// @oa:description "Product status"
	Status string

	// @evl:validate required
	// @evl:validate gt:0
	Price float64

	// @evl:validate optional
	// @evl:validate min:1
	// @evl:validate max:1000
	Quantity int

	// @evl:validate required
	// @evl:validate len:10
	Code string
}

type Order struct {
	// @evl:validate required
	ID string

	// @evl:validate required
	// @evl:validate before:2025-12-31
	// @evl:validate after:2020-01-01
	CreatedAt time.Time

	// @evl:validate required
	// @evl:validate min:1
	// @evl:validate max:100
	Items []string
}
