package entities

type Customer struct {
	ID            uint           `json:"id"`
	UserID        uint           `json:"user_id"`
	User          *User          `json:"user,omitempty"`
	Document      string         `json:"document"`
	PhoneNumber   string         `json:"phone_number"`
	FullName      string         `json:"full_name"`
	Vehicles      []Vehicle      `json:"vehicles,omitempty"`
	ServiceOrders []ServiceOrder `json:"service_orders,omitempty"`
}
