package user

//go:generate elval-gen gen -input . -openapi

// User представляет пользователя системы
// @oa:description "Полная информация о пользователе"
type User struct {
	// @oa:example "user-123"
	ID string

	// @oa:example "john@example.com"
	Email string

	// Вложенная структура из того же пакета → ref без префикса
	// @oa:example {"street":"Lenina 10","zipCode":"101000"}
	Address Address
}
