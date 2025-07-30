package entities

// N:N relationship between Service and ServiceOrder
type ServiceServiceOrder struct {
	ServiceID      uint `gorm:"primaryKey"`
	ServiceOrderID uint `gorm:"primaryKey"`
}
