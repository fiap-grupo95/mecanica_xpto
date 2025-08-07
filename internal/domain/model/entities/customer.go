package entities

import (
	"mecanica_xpto/internal/domain/model/valueobject"
)

type Customer struct {
	ID          uint                `json:"id"`
	User        *User               `json:"user,omitempty"`
	Document    valueobject.CpfCnpj `json:"document"`
	PhoneNumber string              `json:"phone_number"`
	FullName    string              `json:"full_name"`
	//Vehicles      []Vehicle            `json:"vehicles,omitempty"`
	//ServiceOrders []ServiceOrder       `json:"service_orders,omitempty"`
}
