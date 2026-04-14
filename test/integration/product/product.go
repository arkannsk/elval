package product

type Product struct {
	// @evl:validate required
	// @evl:validate eq:"active"
	Status string

	// @evl:validate required
	// @evl:validate min:1
	// @evl:validate max:100
	Quantity int

	// @evl:validate required
	// @evl:validate gt:0
	Price float64

	// @evl:validate required
	// @evl:validate gte:18
	Age int

	// @evl:validate optional
	// @evl:validate eq:"admin"
	Role string

	// @evl:validate neq:0
	Score int

	// @evl:validate lt:100
	Discount int

	// @evl:validate lte:50
	Tax int
}
