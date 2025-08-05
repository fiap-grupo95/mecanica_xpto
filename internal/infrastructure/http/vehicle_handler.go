package http

type VehicleHandler struct {
	service VehicleServiceInterface
}

func NewVehicleHandler(service VehicleServiceInterface) *VehicleHandler {
	return &VehicleHandler{
		service: service,
	}
}

func (v VehicleHandler) GetVehicle() {
	// Handler logic to get a vehicle by ID
}
