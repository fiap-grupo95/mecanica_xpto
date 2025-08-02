package entities

import (
	"time"
)

type Payment struct {
	ID           uint         `json:"id"`
	ServiceOrder ServiceOrder `json:"service_order,omitempty"`
	PaymentDate  time.Time    `json:"payment_date"`
}
