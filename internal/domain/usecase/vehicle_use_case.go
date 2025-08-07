package usecase

import (
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/domain/repository/vehicles"
)

type VehicleServiceInterface interface {
	GetAllVehicles() ([]entities.Vehicle, error)
	GetVehicleByID(id uint) (*entities.Vehicle, error)
	GetVehicleByPlate(plate string) (*entities.Vehicle, error)
	GetVehiclesByCustomerID(customerID uint) ([]entities.Vehicle, error)
	CreateVehicle(vehicle entities.Vehicle) (string, error)
	UpdateVehicle(vehicle entities.Vehicle) (string, error)
	UpdateVehiclePartial(id uint, updates map[string]interface{}) (string, error)
	DeleteVehicle(id uint) error
}

type VehicleService struct {
	repo vehicles.VehicleRepositoryInterface
}

func NewVehicleService(repo vehicles.VehicleRepositoryInterface) VehicleServiceInterface {
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
	if !voPlate.IsValidFormat() {
		return nil, errors.New("invalid plate")
	}
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
func (s *VehicleService) CreateVehicle(vehicle entities.Vehicle) (string, error) {
	if !vehicle.Plate.IsValidFormat() {
		return "invalid plate format", errors.New("invalid plate format")
	}
	err := s.repo.Create(vehicle)
	if err != nil {
		return "error creating a new vehicle", err
	}
	return "Vehicle created successfully", nil
}

func (s *VehicleService) UpdateVehicle(vehicle entities.Vehicle) (string, error) {
	if !vehicle.Plate.IsValidFormat() {
		return "invalid plate format", errors.New("invalid plate format")
	}
	err := s.repo.Update(vehicle)
	if err != nil {
		return "error updating a new vehicle", err
	}
	return "Vehicle updated successfully", nil
}

func (s *VehicleService) UpdateVehiclePartial(id uint, updates map[string]interface{}) (string, error) {
	// First get the existing vehicle
	existingVehicle, err := s.repo.FindByID(id)
	if err != nil {
		return "", err
	}

	if !valueobject.ParsePlate(existingVehicle.Plate).IsValidFormat() {
		return "invalid plate format", errors.New("invalid plate format")
	}
	// Update only the fields that were provided
	if plate, ok := updates["plate"].(string); ok {
		existingVehicle.Plate = plate
	}
	if model, ok := updates["model"].(string); ok {
		existingVehicle.Model = model
	}
	if year, ok := updates["year"].(string); ok {
		existingVehicle.Year = year
	}
	if brand, ok := updates["brand"].(string); ok {
		existingVehicle.Brand = brand
	}
	if customerId, ok := updates["customer_id"].(float64); ok {
		existingVehicle.CustomerID = uint(customerId)
	}

	// Convert DTO to domain entity and update
	if err := s.repo.Update(*existingVehicle.ToDomain()); err != nil {
		return "", err
	}

	return "Vehicle updated successfully", nil
}

func (s *VehicleService) DeleteVehicle(id uint) error {
	return s.repo.Delete(id)
}
