package invalid

// testdata/invalid/invalid.go

type InvalidStruct struct {
	// @evl:validate required
	// @evl:validate min:10
	// @evl:validate max:5
	// Ошибка: min (10) > max (5)
	Name string

	// @evl:validate required
	// @evl:validate min-max:10,5  // <- эта ошибка не ловится
	Age int

	// @evl:validate pattern:email
	Count int

	// @evl:validate unknown
	Field string

	// @evl:validate min:invalid
	InvalidNumber int

	// @evl:validate required
	// @evl:validate min:1
	// @evl:validate max:10
	// @evl:validate pattern:email
	Tags []string
}
