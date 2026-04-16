package main

//go:generate elval-gen generate -input .

type Product struct {
	// @evl:validate required
	// @evl:validate x-color
	Color string

	// @evl:validate required
	// @evl:validate x-even
	Count int

	// @evl:validate x-between:10,90
	Score int

	// @evl:validate x-contains:important
	Description string
}
