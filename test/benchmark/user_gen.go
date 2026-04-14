package benchmark

//go:generate elval-gen -input .

type UserGen struct {
	// @evl:validate required
	// @evl:validate min:2
	// @evl:validate max:50
	Name string

	// @evl:validate required
	// @evl:validate pattern:email
	Email string

	// @evl:validate required
	// @evl:validate min:18
	// @evl:validate max:120
	Age int
}
