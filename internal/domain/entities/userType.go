package entities

type UserType struct {
	ID    uint   `gorm:"primaryKey"`
	Type  string `gorm:"size:50;not null"`
	Users []User
}
