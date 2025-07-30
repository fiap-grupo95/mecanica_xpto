package entities

// N:N relationship between PartsSupply and ServiceOrder
type PartsSupplyServiceOrder struct {
	PartsSupplyID  uint `gorm:"primaryKey"`
	ServiceOrderID uint `gorm:"primaryKey"`
	Quantity       int  `gorm:"not null;default:1"`
}
