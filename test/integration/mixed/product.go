package mixed

type Product struct {
	// @evl:validate required
	// @evl:validate eq:"active"
	Status string

	// @evl:validate required
	// @evl:validate min:1
	// @evl:validate max:100
	Quantity int

	// @evl:validate gt:0
	Price float64
}
