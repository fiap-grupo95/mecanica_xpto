package entities

import (
	"mecanica_xpto/internal/domain/model/valueobject"
)

type AdditionalRepairStatus struct {
	ID                uint               `json:"id" validate:"required"`
	Description       string             `json:"description"`
	AdditionalRepairs []AdditionalRepair `json:"additional_repairs,omitempty"`
}

type AdditionalRepair struct {
	ID           uint               `json:"id"`
	ServiceOrder ServiceOrder       `json:"service_order,omitempty"`
	Service      Service            `json:"service,omitempty"`
	PartsSupply  PartsSupply        `json:"parts_supply,omitempty"`
	ARStatus     valueobject.Status `json:"ar_status,omitempty"`
}
