package sample

// User представляет пользователя системы
type User struct {
	// @evl:validate required
	// @evl:validate min:3
	// @evl:validate max:50
	Name string

	// @evl:validate required
	// @evl:validate pattern:email
	Email string

	// @evl:validate min:18
	// @evl:validate max:99
	Age int

	// @evl:validate optional
	// @evl:validate pattern:phone
	Phone string
}

// Post представляет пост в блоге
type Post struct {
	/*
		@evl:validate required
		@evl:validate min:1
		@evl:validate max:200
	*/
	Title string

	// @evl:validate optional
	// @evl:validate pattern:uuid
	AuthorID string

	// @evl:validate required
	// @evl:validate min:1
	// @evl:validate max:10
	Tags []string
}

// Product представляет товар
type Product struct {
	// @evl:validate required
	// @evl:validate min:0.01
	// @evl:validate max:9999.99
	Price float64

	// @evl:validate required
	// @evl:validate min:1
	// @evl:validate max:10000
	Quantity int

	// @evl:validate required
	// @evl:validate pattern:^[A-Z]{3}-\d+$
	Code string

	// @evl:validate optional
	// @evl:validate len:3
	Coordinates []float64
}

// Config представляет конфигурацию
type Config struct {
	// @evl:validate required
	// @evl:validate min:1024
	// @evl:validate max:65535
	Port int

	// @evl:validate optional
	// @evl:validate pattern:^[a-z]+$
	Environment string

	// @evl:validate required
	// @evl:validate not-zero
	Features []string
}
