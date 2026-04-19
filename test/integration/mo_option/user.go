package main

import (
	"github.com/samber/mo"
)

//go:generate elval-gen generate -input .

// User демонстрирует валидацию внешних дженериков (mo.Option)
type User struct {
	// @evl:validate required
	// @evl:validate min:3
	// @evl:validate max:50
	Name mo.Option[string]

	// @evl:validate required
	// @evl:validate pattern:email
	Email mo.Option[string]

	// @evl:validate optional
	// @evl:validate min:18
	// @evl:validate max:120
	Age mo.Option[int]

	// @evl:validate optional
	// @evl:validate x-strong-password
	Password mo.Option[string]

	// @evl:validate required
	// @evl:validate min:1
	// @evl:validate max:10
	Tag mo.Option[string]
}
