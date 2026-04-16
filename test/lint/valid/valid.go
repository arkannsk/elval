// test/lint/valid/valid.go
package valid

type ValidUser struct {
	// @evl:validate required
	// @evl:validate min:3
	// @evl:validate max:50
	Name string

	// @evl:validate required
	// @evl:validate pattern:email
	Email string

	// @evl:validate min:18
	// @evl:validate max:120
	Age int

	// @evl:validate optional
	// @evl:validate enum:admin,user,guest
	Role string

	// @evl:validate not-zero
	Tags []string

	// @evl:validate min:0.01
	// @evl:validate max:999.99
	Price float64
}
