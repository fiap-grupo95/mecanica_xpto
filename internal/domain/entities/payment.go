package entities

import "time"

type Payment struct {
	ID             uint         `gorm:"primaryKey"`
	ServiceOrderID uint         `gorm:"unique;not null"`
	ServiceOrder   ServiceOrder `gorm:"foreignKey:ServiceOrderID;references:ID"`
	PaymentDate    time.Time    `gorm:"not null"`
}
