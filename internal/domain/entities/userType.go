package entities

import "mecanica_xpto/internal/domain/entities/valueobject"

type UserType struct {
	ID    uint                 `json:"id"`
	Type  valueobject.UserType `json:"type"`
	Users []User               `json:"users,omitempty"`
}
