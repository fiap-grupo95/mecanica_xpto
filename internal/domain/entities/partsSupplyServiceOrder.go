package entities

type PartsSupplyServiceOrder struct {
	PartsSupplyID  uint `json:"parts_supply_id"`
	ServiceOrderID uint `json:"service_order_id"`
	Quantity       int  `json:"quantity"`
}
