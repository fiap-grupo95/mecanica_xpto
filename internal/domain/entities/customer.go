package entities

// 1:1 relationship between Customer and User
// 1:N relationship between Customer and Vehicle
// 1:N relationship between Customer and ServiceOrder
type Customer struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        uint   `gorm:"unique;not null"`
	User          *User  `gorm:"foreignKey:UserID;references:ID"`
	Document      string `gorm:"size:20;not null"`
	PhoneNumber   string `gorm:"size:20;not null"`
	FullName      string `gorm:"size:100;not null"`
	Vehicles      []Vehicle
	ServiceOrders []ServiceOrder
}
