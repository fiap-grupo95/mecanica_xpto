package repository

// N:N relationship between Service and ServiceOrder
type ServiceServiceOrderModel struct {
	ServiceID      uint `gorm:"primaryKey"`
	ServiceOrderID uint `gorm:"primaryKey"`
}
