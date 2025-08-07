package vehicles

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"

	"gorm.io/gorm"
)

type VehicleRepositoryInterface interface {
	FindAll() ([]dto.VehicleDTO, error)
	FindByID(id uint) (*dto.VehicleDTO, error)
	FindByPlate(plate valueobject.Plate) (*dto.VehicleDTO, error)
	FindByCustomerID(customerID uint) ([]dto.VehicleDTO, error)
	Create(vehicle entities.Vehicle) error
	Update(vehicle entities.Vehicle) error
	Delete(id uint) error
}

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) VehicleRepositoryInterface {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) FindAll() ([]dto.VehicleDTO, error) {
	var vehicleDTOs []dto.VehicleDTO
	if err := r.db.Preload("Customer").Find(&vehicleDTOs).Error; err != nil {
		return nil, err
	}

	vehicles := make([]dto.VehicleDTO, len(vehicleDTOs))
	for i, dto := range vehicleDTOs {
		vehicles[i] = dto
	}
	return vehicles, nil
}

func (r *VehicleRepository) FindByID(id uint) (*dto.VehicleDTO, error) {
	var vehicleDTO dto.VehicleDTO
	if err := r.db.Preload("Customer").First(&vehicleDTO, id).Error; err != nil {
		return nil, err
	}
	return &vehicleDTO, nil
}

func (r *VehicleRepository) FindByPlate(plate valueobject.Plate) (*dto.VehicleDTO, error) {
	var vehicleDTO dto.VehicleDTO
	if err := r.db.Preload("Customer").Where("plate = ?", plate.String()).First(&vehicleDTO).Error; err != nil {
		return nil, err
	}
	return &vehicleDTO, nil
}

func (r *VehicleRepository) FindByCustomerID(customerID uint) ([]dto.VehicleDTO, error) {
	var vehicleDTOs []dto.VehicleDTO
	if err := r.db.Preload("Customer").Where("customer_id = ?", customerID).Find(&vehicleDTOs).Error; err != nil {
		return nil, err
	}

	vehicles := make([]dto.VehicleDTO, len(vehicleDTOs))
	for i, dto := range vehicleDTOs {
		vehicles[i] = dto
	}
	return vehicles, nil
}

func (r *VehicleRepository) Create(vehicle entities.Vehicle) error {
	vehicleDTO := dto.VehicleDTO{
		Plate:      string(vehicle.Plate),
		CustomerID: vehicle.Customer.ID,
		Model:      vehicle.Model,
		Year:       vehicle.Year,
		Brand:      vehicle.Brand,
	}

	if err := r.db.Create(&vehicleDTO).Error; err != nil {
		return err
	}

	return nil
}

func (r *VehicleRepository) Update(vehicle entities.Vehicle) error {
	vehicleDTO := dto.VehicleDTO{
		ID:         vehicle.ID,
		Plate:      string(vehicle.Plate),
		CustomerID: vehicle.Customer.ID,
		Model:      vehicle.Model,
		Year:       vehicle.Year,
		Brand:      vehicle.Brand,
	}

	if err := r.db.Save(&vehicleDTO).Error; err != nil {
		return err
	}

	return nil
}

func (r *VehicleRepository) Delete(id uint) error {
	// Using GORM's soft delete functionality which will automatically set the deleted_at timestamp
	if err := r.db.Model(&dto.VehicleDTO{}).Where("id = ?", id).Delete("").Error; err != nil {
		return err
	}
	return nil
}

// Ensure VehicleRepository implements VehicleRepositoryInterface
var _ VehicleRepositoryInterface = (*VehicleRepository)(nil)
