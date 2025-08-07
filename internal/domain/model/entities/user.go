package entities

import (
	"mecanica_xpto/internal/domain/model/valueobject"
	"time"
)

type User struct {
	ID        uint                 `json:"id"`
	Email     string               `json:"email"`
	Password  string               `json:"password"`
	UserType  valueobject.UserType `json:"user_type,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	DeletedAt *time.Time           `json:"deleted_at,omitempty"`
	//Customer  *Customer            `json:"customer,omitempty"`
}
