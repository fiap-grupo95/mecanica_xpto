package entities

import "mecanica_xpto/internal/domain/entities/valueobject"

type Customer struct {
	ID            uint                 `json:"id"`
	UserID        uint                 `json:"user_id"`
	User          *User                `json:"user,omitempty"`
	CPF_CNPJ      valueobject.CPF_CNPJ `json:"CPF_CNPJ"`
	PhoneNumber   string               `json:"phone_number"`
	FullName      string               `json:"full_name"`
	Vehicles      []Vehicle            `json:"vehicles,omitempty"`
	ServiceOrders []ServiceOrder       `json:"service_orders,omitempty"`
}
