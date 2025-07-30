package entities

// 1:N relationship between AdditionalRepair and AdditionalRepairStatus
type AdditionalRepairStatus struct {
	ID                uint   `gorm:"primaryKey"`
	Description       string `gorm:"size:50;not null"`
	AdditionalRepairs []AdditionalRepair
}

type AdditionalRepair struct {
	ID             uint                   `gorm:"primaryKey"`
	ServiceOrderID uint                   `gorm:"not null"`
	ServiceOrder   ServiceOrder           `gorm:"foreignKey:ServiceOrderID"`
	ServiceID      uint                   `gorm:"not null"`
	Service        Service                `gorm:"foreignKey:ServiceID"`
	PartsSupplyID  uint                   `gorm:"not null"`
	PartsSupply    PartsSupply            `gorm:"foreignKey:PartsSupplyID"`
	ARStatusID     uint                   `gorm:"not null"`
	ARStatus       AdditionalRepairStatus `gorm:"foreignKey:ARStatusID"`
}
