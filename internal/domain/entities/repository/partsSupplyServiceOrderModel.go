package repository

// N:N relationship between PartsSupply and ServiceOrder
type PartsSupplyServiceOrderModel struct {
	PartsSupplyID  uint `gorm:"primaryKey"`
	ServiceOrderID uint `gorm:"primaryKey"`
	Quantity       int  `gorm:"not null;default:1"`
}
