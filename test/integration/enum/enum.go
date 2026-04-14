// test/integration/enum/enum.go
package enum

//go:generate elval-gen -input .

type Order struct {
	// @evl:validate required
	// @evl:validate enum:pending,processing,shipped,delivered,cancelled
	Status string

	// @evl:validate required
	// @evl:validate enum:1,2,3,4,5
	Priority int

	// @evl:validate optional
	// @evl:validate enum:small,medium,large
	Size string
}

type User struct {
	// @evl:validate required
	// @evl:validate enum:admin,moderator,user
	Role string

	// @evl:validate required
	// @evl:validate enum:1,2,3
	Level int
}
