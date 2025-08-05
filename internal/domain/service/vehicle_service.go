package service

type VehicleServiceInterface interface {
	GetAllVehicles() ([]VehicleDTO, error)
	GetVehicleByID(id uint) (*VehicleDTO, error)
	GetVehicleByPlate(plate string) (*VehicleDTO, error)
	GetVehiclesByCustomerID(customerID uint) ([]VehicleDTO, error)
	CreateVehicle(vehicle VehicleDTO) (*VehicleDTO, error)
	UpdateVehicle(vehicle VehicleDTO) (*VehicleDTO, error)
	DeleteVehicle(id uint) error
}

type VehicleService struct {
	repo VehicleRepositoryInterface
}

func NewVehicleService(repo VehicleRepositoryInterface) VehicleServiceInterface {
	return &VehicleService{repo: repo}
}

func (s *VehicleService) GetAllVehicles() ([]VehicleDTO, error) {
	return s.repo.FindAll()
}
func (s *VehicleService) GetVehicleByID(id uint) (*VehicleDTO, error) {
	return s.repo.FindByID(id)
}
func (s *VehicleService) GetVehicleByPlate(plate string) (*VehicleDTO, error) {
	voPlate, err := valueobject.NewPlate(plate)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByPlate(*voPlate)
}
func (s *VehicleService) GetVehiclesByCustomerID(customerID uint) ([]VehicleDTO, error) {
	return s.repo.FindByCustomerID(customerID)
}
func (s *VehicleService) CreateVehicle(vehicle VehicleDTO) (*VehicleDTO, error) {
	return s.repo.Create(vehicle)
}

func (s *VehicleService) UpdateVehicle(vehicle VehicleDTO) (*VehicleDTO, error) {
	return s.repo.Update(vehicle)
}
func (s *VehicleService) DeleteVehicle(id uint) error {
	return s.repo.Delete(id)
}
