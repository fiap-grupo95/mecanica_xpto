package repository

import (
	"mecanica_xpto/internal/domain/entities"
	"time"
)

type PaymentModel struct {
	ID             uint              `gorm:"primaryKey"`
	ServiceOrderID uint              `gorm:"unique;not null"`
	ServiceOrder   ServiceOrderModel `gorm:"foreignKey:ServiceOrderID;references:ID"`
	PaymentDate    time.Time         `gorm:"not null"`
}

func (pm *PaymentModel) ToDomain() entities.Payment {
	return entities.Payment{
		ID:           pm.ID,
		ServiceOrder: pm.ServiceOrder.ToDomain(),
		PaymentDate:  pm.PaymentDate,
	}
}
