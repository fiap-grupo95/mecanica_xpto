package entities

import (
	"mecanica_xpto/internal/domain/entities/valueobject"
	"time"
)

type ServiceOrderStatus struct {
	ID            uint               `json:"id"`
	Description   valueobject.Status `json:"description"`
	ServiceOrders []ServiceOrder     `json:"service_orders,omitempty"`
}

type ServiceOrder struct {
	ID                   uint               `json:"id"`
	Customer             Customer           `json:"customer,omitempty"`
	Vehicle              Vehicle            `json:"vehicle,omitempty"`
	ServiceOrderStatus   valueobject.Status `json:"service_order_status"`
	Estimate             float64            `json:"estimate"`
	StartedExecutionDate time.Time          `json:"started_execution_date"`
	FinalExecutionDate   time.Time          `json:"final_execution_date"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
	AdditionalRepairs    []AdditionalRepair `json:"additional_repairs,omitempty"`
	Payment              *Payment           `json:"payment,omitempty"`
	PartsSupplies        []PartsSupply      `json:"parts_supplies,omitempty"`
	Services             []Service          `json:"services,omitempty"`
}
