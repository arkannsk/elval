package sample

// User структура для тестирования простых полей
type User struct {
	// @evl:validate required
	// @evl:validate min-max:3,50
	Name string

	// @evl:validate required
	// @evl:validate pattern:email
	Email string
}

// Post структура для тестирования многострочных комментариев
type Post struct {
	/*
	   @evl:validate required
	   @evl:validate min-max:1,100
	*/
	Title string

	// @evl:validate optional
	// @evl:validate pattern:uuid
	AuthorID string
}

// Product структура для тестирования разных типов полей
type Product struct {
	// @evl:validate required
	// @evl:validate min-max:0,999.99
	Price float64

	// @evl:validate min-max:1,1000
	Quantity int

	// @evl:validate pattern:^[A-Z]{3}-\d+$  // Один слеш
	Code string
}
