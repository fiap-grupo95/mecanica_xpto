package entities

import (
	"mecanica_xpto/internal/domain/model/valueobject"
	"time"
)

type ServiceOrder struct {
	ID                   uint                           `json:"id"`
	CustomerID           uint                           `json:"customer_id"`
	Customer             Customer                       `json:"customer,omitempty"`
	VehicleID            uint                           `json:"vehicle_id"`
	Vehicle              Vehicle                        `json:"vehicle,omitempty"`
	ServiceOrderStatus   valueobject.ServiceOrderStatus `json:"service_order_status"`
	Estimate             float64                        `json:"estimate"`
	StartedExecutionDate *time.Time                     `json:"started_execution_date"`
	FinalExecutionDate   *time.Time                     `json:"final_execution_date"`
	CreatedAt            *time.Time                     `json:"created_at"`
	UpdatedAt            *time.Time                     `json:"updated_at"`
	AdditionalRepairs    []AdditionalRepair             `json:"additional_repairs,omitempty"`
	Payment              *Payment                       `json:"payment,omitempty"`
	PartsSupplies        []PartsSupply                  `json:"parts_supplies,omitempty"`
	Services             []Service                      `json:"services,omitempty"`
}
