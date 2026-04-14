package multi

//go:generate elval-gen -input .

// User структура пользователя
type User struct {
	// @evl:validate required
	// @evl:validate min:2
	// @evl:validate max:50
	Name string

	// @evl:validate required
	// @evl:validate pattern:email
	Email string

	// @evl:validate min:18
	// @evl:validate max:120
	Age int
}

// Product структура продукта
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

// Order структура заказа
type Order struct {
	// @evl:validate required
	ID string

	// @evl:validate required
	// @evl:validate min:1
	UserID int

	// @evl:validate optional
	// @evl:validate min:0.01
	Total float64

	// @evl:validate optional
	// @evl:validate min:1
	// @evl:validate max:100
	Items []string
}

// Config структура конфигурации (без аннотаций - не должна генерироваться)
type Config struct {
	Host string
	Port int
}
