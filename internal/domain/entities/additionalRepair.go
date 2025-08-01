package entities

type AdditionalRepairStatus struct {
	ID                uint               `json:"id"`
	Description       string             `json:"description"`
	AdditionalRepairs []AdditionalRepair `json:"additional_repairs,omitempty"`
}

type AdditionalRepair struct {
	ID             uint                   `json:"id"`
	ServiceOrderID uint                   `json:"service_order_id"`
	ServiceOrder   ServiceOrder           `json:"service_order,omitempty"`
	ServiceID      uint                   `json:"service_id"`
	Service        Service                `json:"service,omitempty"`
	PartsSupplyID  uint                   `json:"parts_supply_id"`
	PartsSupply    PartsSupply            `json:"parts_supply,omitempty"`
	ARStatusID     uint                   `json:"ar_status_id"`
	ARStatus       AdditionalRepairStatus `json:"ar_status,omitempty"`
}
