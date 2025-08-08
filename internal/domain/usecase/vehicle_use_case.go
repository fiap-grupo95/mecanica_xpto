package usecase

import (
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/domain/repository/vehicles"
)

var (
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrInvalidPlateFormat   = errors.New("invalid plate format")
	ErrVehicleAlreadyExists = errors.New("vehicle already exists")
	ErrInvalidID            = errors.New("invalid vehicle ID")
)

var (
	MessageVehicleCreatedSuccessfully = "Vehicle created successfully"
	MessageVehicleUpdatedSuccessfully = "Vehicle updated successfully"
	MessageVehicleNotFound            = "Vehicle not found"
	MessageInvalidPlateFormat         = "Invalid plate format"
	MessageErrorCreatingVehicle       = "Error creating a new vehicle"
	MessageErrorUpdatingVehicle       = "Error updating the vehicle"
	MessageVehicleAlreadyExists       = "vehicle already exists with the same plate number"
	MessageErrorSearch                = "Error searching existing vehicles"
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

	vehicleList, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, v := range vehicleList {
		vehiclesList = append(vehiclesList, *v.ToDomain())
	}
	return vehiclesList, nil
}
func (s *VehicleService) GetVehicleByID(id uint) (*entities.Vehicle, error) {
	vehicle, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, ErrVehicleNotFound
	}
	return vehicle.ToDomain(), nil
}
func (s *VehicleService) GetVehicleByPlate(plate string) (*entities.Vehicle, error) {
	voPlate := valueobject.ParsePlate(plate)
	if !voPlate.IsValidFormat() {
		return nil, ErrInvalidPlateFormat
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
	result, err := s.repo.FindByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	var vehiclesList []entities.Vehicle
	for _, v := range result {
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
		return MessageInvalidPlateFormat, ErrInvalidPlateFormat
	}

	existingVehicle, err := s.repo.FindByPlate(vehicle.Plate)
	if err != nil {
		return MessageErrorSearch, err
	}
	if existingVehicle != nil {
		return MessageVehicleAlreadyExists, ErrVehicleAlreadyExists
	}
	err = s.repo.Create(vehicle)
	if err != nil {
		return MessageErrorCreatingVehicle, err
	}
	return MessageVehicleCreatedSuccessfully, nil
}

func (s *VehicleService) UpdateVehicle(vehicle entities.Vehicle) (string, error) {
	if !vehicle.Plate.IsValidFormat() {
		return MessageInvalidPlateFormat, ErrInvalidPlateFormat
	}
	existingVehicle, err := s.repo.FindByID(vehicle.ID)
	if err != nil {
		return MessageErrorSearch, err
	}
	if existingVehicle == nil {
		return MessageVehicleNotFound, ErrVehicleNotFound
	}
	err = s.repo.Update(vehicle)
	if err != nil {
		return MessageErrorUpdatingVehicle, err
	}
	return MessageVehicleUpdatedSuccessfully, nil
}

func (s *VehicleService) UpdateVehiclePartial(id uint, updates map[string]interface{}) (string, error) {
	// First get the existing vehicle
	existingVehicle, err := s.repo.FindByID(id)
	if err != nil {
		return "", err
	}
	if existingVehicle == nil {
		return MessageVehicleNotFound, ErrVehicleNotFound
	}

	// Update only the fields that were provided
	if plate, ok := updates["plate"].(string); ok {
		if !valueobject.ParsePlate(plate).IsValidFormat() {
			return MessageInvalidPlateFormat, ErrInvalidPlateFormat
		}
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

	return MessageVehicleUpdatedSuccessfully, nil
}

func (s *VehicleService) DeleteVehicle(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	vehicle, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if vehicle == nil {
		return ErrVehicleNotFound
	}

	err = s.repo.Delete(id)
	return err
}
