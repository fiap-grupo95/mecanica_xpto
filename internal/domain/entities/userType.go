package entities

type UserType struct {
	ID    uint   `json:"id"`
	Type  string `json:"type"`
	Users []User `json:"users,omitempty"`
}
