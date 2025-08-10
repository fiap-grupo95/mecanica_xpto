package entities

import (
	"mecanica_xpto/internal/domain/model/valueobject"
	"time"
)

type AdditionalRepair struct {
	ID            uint                               `json:"id"`
	ARStatus      valueobject.AdditionalRepairStatus `json:"ar_status,omitempty"`
	Estimate      float64                            `json:"estimate"`
	CreatedAt     time.Time                          `json:"created_at"`
	UpdatedAt     time.Time                          `json:"updated_at"`
	PartsSupplies []PartsSupply                      `json:"parts_supplies,omitempty"`
	Services      []Service                          `json:"services,omitempty"`
}
