package user

// Address представляет почтовый адрес пользователя
type Address struct {
	// @oa:example "Tverskaya 1"
	Street string
	// @oa:example "123456"
	ZipCode string
}
