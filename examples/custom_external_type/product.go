package main

import "github.com/samber/mo"

//go:generate elval-gen generate -input .

type User struct {
	// @evl:validate required
	// @evl:validate x-option-present
	Name mo.Option[string]

	// @evl:validate required
	// @evl:validate x-option-value-min:3
	// @evl:validate x-option-value-max:50
	Email mo.Option[string]

	// @evl:validate x-option-absent
	Phone mo.Option[string]
}
