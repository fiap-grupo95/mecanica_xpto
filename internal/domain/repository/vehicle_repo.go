package repository

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"time"

	"gorm.io/gorm"
)

type VehicleDTO struct {
	ID            uint           `gorm:"primaryKey"`
	Plate         string         `gorm:"size:10;not null"`
	CustomerID    uint           `gorm:"not null"`
	Customer      CustomerDTO    `gorm:"foreignKey:CustomerID"`
	Model         string         `gorm:"size:50;not null"`
	Year          string         `gorm:"size:4"`
	Brand         string         `gorm:"size:50;not null"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	ServiceOrders []ServiceOrderDTO
}

func (v *VehicleDTO) ToDomain() *entities.Vehicle {
	return &entities.Vehicle{
		ID:        v.ID,
		Plate:     valueobject.ParsePlate(v.Plate),
		Customer:  v.Customer.ToDomain(),
		Model:     v.Model,
		Year:      v.Year,
		Brand:     v.Brand,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
		DeletedAt: func() *time.Time {
			if v.DeletedAt.Valid {
				return &v.DeletedAt.Time
			}
			return nil
		}(),
		ServiceOrders: nil, // This will be populated by the repository layer
	}
}

// TableName specifies the table name for VehicleDTO
func (v *VehicleDTO) TableName() string {
	return "db_mecanica_xpto.tb_vehicle"
}

type VehicleRepositoryInterface interface {
	FindAll() ([]VehicleDTO, error)
	FindByID(id uint) (*VehicleDTO, error)
	FindByPlate(plate valueobject.Plate) (*VehicleDTO, error)
	FindByCustomerID(customerID uint) ([]VehicleDTO, error)
	Create(vehicle entities.Vehicle) (*VehicleDTO, error)
	Update(vehicle entities.Vehicle) (*VehicleDTO, error)
	Delete(id uint) error
}

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) VehicleRepositoryInterface {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) FindAll() ([]VehicleDTO, error) {
	var vehicleDTOs []VehicleDTO
	if err := r.db.Preload("Customer").Find(&vehicleDTOs).Error; err != nil {
		return nil, err
	}

	vehicles := make([]VehicleDTO, len(vehicleDTOs))
	for i, dto := range vehicleDTOs {
		vehicles[i] = dto
	}
	return vehicles, nil
}

func (r *VehicleRepository) FindByID(id uint) (*VehicleDTO, error) {
	var vehicleDTO VehicleDTO
	if err := r.db.Preload("Customer").First(&vehicleDTO, id).Error; err != nil {
		return nil, err
	}
	return &vehicleDTO, nil
}

func (r *VehicleRepository) FindByPlate(plate valueobject.Plate) (*VehicleDTO, error) {
	var vehicleDTO VehicleDTO
	if err := r.db.Preload("Customer").Where("plate = ?", plate.String()).First(&vehicleDTO).Error; err != nil {
		return nil, err
	}
	return &vehicleDTO, nil
}

func (r *VehicleRepository) FindByCustomerID(customerID uint) ([]VehicleDTO, error) {
	var vehicleDTOs []VehicleDTO
	if err := r.db.Preload("Customer").Where("customer_id = ?", customerID).Find(&vehicleDTOs).Error; err != nil {
		return nil, err
	}

	vehicles := make([]VehicleDTO, len(vehicleDTOs))
	for i, dto := range vehicleDTOs {
		vehicles[i] = dto
	}
	return vehicles, nil
}

func (r *VehicleRepository) Create(vehicle entities.Vehicle) (*VehicleDTO, error) {
	vehicleDTO := VehicleDTO{
		Plate:      string(vehicle.Plate),
		CustomerID: vehicle.Customer.ID,
		Model:      vehicle.Model,
		Year:       vehicle.Year,
		Brand:      vehicle.Brand,
	}

	if err := r.db.Create(&vehicleDTO).Error; err != nil {
		return nil, err
	}

	return &vehicleDTO, nil
}

func (r *VehicleRepository) Update(vehicle entities.Vehicle) (*VehicleDTO, error) {
	vehicleDTO := VehicleDTO{
		ID:         vehicle.ID,
		Plate:      string(vehicle.Plate),
		CustomerID: vehicle.Customer.ID,
		Model:      vehicle.Model,
		Year:       vehicle.Year,
		Brand:      vehicle.Brand,
	}

	if err := r.db.Save(&vehicleDTO).Error; err != nil {
		return nil, err
	}

	return &vehicleDTO, nil
}

func (r *VehicleRepository) Delete(id uint) error {
	if err := r.db.Delete(&VehicleDTO{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Ensure VehicleRepository implements VehicleRepositoryInterface
var _ VehicleRepositoryInterface = (*VehicleRepository)(nil)
