package order

//go:generate elval-gen gen -input . -openapi

// Order представляет заказ в системе
// @oa:description "Информация о заказе"
type Order struct {
	// @oa:example "order-456"
	ID string

	// @oa:example 2999.99
	Total float64

	// Вложенная структура из ДРУГОГО пакета → ref с префиксом пакета
	// @oa:example {"warehouse":"Main","zone":"A-1","deliveryWindow":"09:00-12:00"}
	ShippingAddress Address
}
