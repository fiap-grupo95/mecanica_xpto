package repository

import (
	"mecanica_xpto/internal/domain/model/entities"
	"time"
)

type PaymentDTO struct {
	ID             uint            `gorm:"primaryKey"`
	ServiceOrderID uint            `gorm:"unique;not null"`
	ServiceOrder   ServiceOrderDTO `gorm:"foreignKey:ServiceOrderID;references:ID"`
	PaymentDate    time.Time       `gorm:"not null"`
}

func (pm *PaymentDTO) ToDomain() entities.Payment {
	return entities.Payment{
		ID:           pm.ID,
		ServiceOrder: pm.ServiceOrder.ToDomain(),
		PaymentDate:  pm.PaymentDate,
	}
}
