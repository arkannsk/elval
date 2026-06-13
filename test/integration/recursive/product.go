// product.go
package integration

type Node struct {
	// @evl:validate required
	// @evl:validate max:50
	Value string

	// @evl:validate optional
	Next *Node
}
