package service

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/domain/repository"
)

type VehicleServiceInterface interface {
	GetAllVehicles() ([]entities.Vehicle, error)
	GetVehicleByID(id uint) (*entities.Vehicle, error)
	GetVehicleByPlate(plate string) (*entities.Vehicle, error)
	GetVehiclesByCustomerID(customerID uint) ([]entities.Vehicle, error)
	CreateVehicle(vehicle entities.Vehicle) (*entities.Vehicle, error)
	UpdateVehicle(vehicle entities.Vehicle) (*entities.Vehicle, error)
	DeleteVehicle(id uint) error
}

type VehicleService struct {
	repo repository.VehicleRepositoryInterface
}

func NewVehicleService(repo repository.VehicleRepositoryInterface) VehicleServiceInterface {
	return &VehicleService{repo: repo}
}

func (s *VehicleService) GetAllVehicles() ([]entities.Vehicle, error) {
	var vehiclesList []entities.Vehicle

	vehicles, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, v := range vehicles {
		vehiclesList = append(vehiclesList, *v.ToDomain())
	}
	return vehiclesList, nil
}
func (s *VehicleService) GetVehicleByID(id uint) (*entities.Vehicle, error) {
	vehicle, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return vehicle.ToDomain(), nil
}
func (s *VehicleService) GetVehicleByPlate(plate string) (*entities.Vehicle, error) {
	voPlate := valueobject.ParsePlate(plate)
	vehicle, err := s.repo.FindByPlate(voPlate)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, nil
	}
	return vehicle.ToDomain(), nil
}
func (s *VehicleService) GetVehiclesByCustomerID(customerID uint) ([]entities.Vehicle, error) {
	vehicles, err := s.repo.FindByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	var vehiclesList []entities.Vehicle
	for _, v := range vehicles {
		vehiclesList = append(vehiclesList, *v.ToDomain())
	}
	if len(vehiclesList) > 0 {
		return vehiclesList, nil
	}
	// If no vehicles found, return an empty slice instead of nil
	if vehiclesList == nil {
		vehiclesList = []entities.Vehicle{}
	}
	return vehiclesList, nil
}
func (s *VehicleService) CreateVehicle(vehicle entities.Vehicle) (*entities.Vehicle, error) {
	result, err := s.repo.Create(vehicle)
	if err != nil {
		return nil, err
	}
	return result.ToDomain(), nil
}

func (s *VehicleService) UpdateVehicle(vehicle entities.Vehicle) (*entities.Vehicle, error) {
	result, err := s.repo.Update(vehicle)
	if err != nil {
		return nil, err
	}
	return result.ToDomain(), nil
}

func (s *VehicleService) DeleteVehicle(id uint) error {
	return s.repo.Delete(id)
}
