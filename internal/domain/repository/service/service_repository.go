package service

import (
	"context"
	"fmt"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"

	"gorm.io/gorm"
)

type IServiceRepo interface {
	Create(ctx context.Context, so *entities.Service) (entities.Service, error)
	GetByID(ctx context.Context, id uint) (entities.Service, error)
	GetByName(ctx context.Context, name string) (entities.Service, error)
	Update(ctx context.Context, so *entities.Service) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]entities.Service, error)
}

type ServiceRepository struct {
	db *gorm.DB
}

var _ IServiceRepo = (*ServiceRepository)(nil)

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (s *ServiceRepository) Create(ctx context.Context, service *entities.Service) (entities.Service, error) {
	dto := dto.ServiceDTO{
		Name:        service.Name,
		Description: service.Description,
		Price:       service.Price,
	}

	fmt.Printf("Creating service: %+v\n", dto)
	if err := s.db.Create(&dto).Error; err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return entities.Service{}, err
	}
	return dto.ToDomain(), nil
}

func (s *ServiceRepository) GetByID(ctx context.Context, id uint) (entities.Service, error) {
	var dto dto.ServiceDTO
	if err := s.db.WithContext(ctx).First(&dto, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entities.Service{}, nil
		}
		return entities.Service{}, err
	}
	return dto.ToDomain(), nil
}

func (s *ServiceRepository) GetByName(ctx context.Context, name string) (entities.Service, error) {
	var dto dto.ServiceDTO
	if err := s.db.WithContext(ctx).Where("name = ?", name).First(&dto).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entities.Service{}, nil
		}
		return entities.Service{}, err
	}
	return dto.ToDomain(), nil
}

func (s *ServiceRepository) Update(ctx context.Context, service *entities.Service) error {
	var dtoDB dto.ServiceDTO
	if err := s.db.WithContext(ctx).First(&dtoDB, service.ID).Error; err != nil {
		return err
	}
	updates := make(map[string]interface{})
	if service.Name != "" {
		updates["name"] = service.Name
	}
	if service.Description != "" {
		updates["description"] = service.Description
	}
	if service.Price != 0 {
		updates["price"] = service.Price
	}

	return s.db.WithContext(ctx).
		Model(&dto.ServiceDTO{}).
		Where("id = ?", service.ID).
		Updates(updates).Error

}

func (s *ServiceRepository) Delete(ctx context.Context, id uint) error {
	if err := s.db.WithContext(ctx).Delete(&dto.ServiceDTO{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *ServiceRepository) List(ctx context.Context) ([]entities.Service, error) {
	var dtos []dto.ServiceDTO
	if err := s.db.WithContext(ctx).Find(&dtos).Error; err != nil {
		return nil, err
	}
	result := make([]entities.Service, len(dtos))
	for i, dto := range dtos {
		result[i] = dto.ToDomain()
	}
	return result, nil
}
