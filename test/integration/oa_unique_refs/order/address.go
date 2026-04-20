package order

// Address представляет адрес доставки заказа
// ⚠️ То же имя, что и user.Address, но другая семантика
type Address struct {
	// @oa:example "Warehouse A"
	Warehouse string
	// @oa:example "Zone B-12"
	Zone string
	// @oa:example "12:00-18:00"
	DeliveryWindow string
}
