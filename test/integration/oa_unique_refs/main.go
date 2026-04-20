package main

import (
	"github.com/arkannsk/elval/test/integration/oa_unique_refs/order"
	"github.com/arkannsk/elval/test/integration/oa_unique_refs/user"
)

//go:generate elval-gen gen -input ./user -openapi
//go:generate elval-gen gen -input ./order -openapi

// APIResponse структура, которая использует оба типа с конфликтом имён
// @oa:oneOf "UserResponse,OrderResponse"
type APIResponse struct {
	// Union-тип: может быть ответом пользователя или заказа
	// @oa:oneOf "user.User,order.Order"
	Data any
}

// UserResponse — обёртка для ответа с пользователем
type UserResponse struct {
	user.User
	// @oa:example "success"
	Status string
}

// OrderResponse — обёртка для ответа с заказом
type OrderResponse struct {
	order.Order
	// @oa:example "shipped"
	Status string
}
