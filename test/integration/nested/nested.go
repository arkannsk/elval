package nested

//go:generate elval-gen generate -input .

type Address struct {
	// @evl:validate required
	// @evl:validate min:2
	City string

	// @evl:validate required
	// @evl:validate min:5
	Street string

	// @evl:validate optional
	// @evl:validate pattern:^[0-9]{5,6}$
	ZipCode string
}

type User struct {
	// @evl:validate required
	// @evl:validate min:2
	Name string

	// @evl:validate required
	// @evl:validate pattern:email
	Email string

	Address Address

	// @evl:validate optional
	BillingAddress *Address
}

type Company struct {
	// @evl:validate required
	Name string

	// @evl:validate required
	// @evl:validate min:1
	Addresses []Address

	// @evl:validate optional
	Users []User
}
