package parts_supply

//go:generate mockgen -destination=../mocks/parts_supply_repository_mock.go -package=mocks mecanica_xpto/internal/domain PartsSupplyRepo

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"

	"gorm.io/gorm"
)

type IPartsSupplyRepo interface {
	Create(ctx context.Context, ps *entities.PartsSupply) (entities.PartsSupply, error)
	GetByID(ctx context.Context, id uint) (entities.PartsSupply, error)
	GetByName(ctx context.Context, name string) (entities.PartsSupply, error)
	Update(ctx context.Context, ps *entities.PartsSupply) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]entities.PartsSupply, error)
	GetByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error)
}

type PartsSupplyRepository struct {
	db *gorm.DB
}

var _ IPartsSupplyRepo = (*PartsSupplyRepository)(nil)

func NewPartsSupplyRepository(db *gorm.DB) *PartsSupplyRepository {
	return &PartsSupplyRepository{db: db}
}

func (s *PartsSupplyRepository) Create(ctx context.Context, ps *entities.PartsSupply) (entities.PartsSupply, error) {
	dto := dto.PartsSupplyDTO{
		Name:            ps.Name,
		Description:     ps.Description,
		Price:           ps.Price,
		QuantityTotal:   ps.QuantityTotal,
		QuantityReserve: ps.QuantityReserve,
	}
	if err := s.db.WithContext(ctx).Create(&dto).Error; err != nil {
		return entities.PartsSupply{}, err
	}
	return dto.ToDomain(), nil
}

func (s *PartsSupplyRepository) GetByID(ctx context.Context, id uint) (entities.PartsSupply, error) {
	var dto dto.PartsSupplyDTO
	if err := s.db.WithContext(ctx).First(&dto, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.PartsSupply{}, nil
		}
		return entities.PartsSupply{}, err
	}
	return dto.ToDomain(), nil
}

func (s *PartsSupplyRepository) GetByName(ctx context.Context, name string) (entities.PartsSupply, error) {
	var dto dto.PartsSupplyDTO
	if err := s.db.WithContext(ctx).Where("name = ?", name).Find(&dto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.PartsSupply{}, nil
		}
		return entities.PartsSupply{}, err
	}
	return dto.ToDomain(), nil
}

func (s *PartsSupplyRepository) Update(ctx context.Context, ps *entities.PartsSupply) error {
	var dtoDB dto.PartsSupplyDTO
	if err := s.db.First(&dtoDB, ps.ID).Error; err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if ps.Name != "" {
		updates["name"] = ps.Name
	}
	if ps.Description != "" {
		updates["description"] = ps.Description
	}
	if ps.Price != 0 {
		updates["price"] = ps.Price
	}
	if ps.QuantityTotal != 0 {
		updates["quantity_total"] = ps.QuantityTotal
	}
	if ps.QuantityReserve != 0 {
		updates["quantity_reserve"] = ps.QuantityReserve
	}

	if len(updates) == 0 {
		return nil
	}

	return s.db.WithContext(ctx).
		Model(&dto.PartsSupplyDTO{}).
		Where("id = ?", ps.ID).
		Updates(updates).Error
}

func (s *PartsSupplyRepository) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&dto.PartsSupplyDTO{}, id).Error
}

func (s *PartsSupplyRepository) List(ctx context.Context) ([]entities.PartsSupply, error) {
	var dtos []dto.PartsSupplyDTO
	if err := s.db.WithContext(ctx).Find(&dtos).Error; err != nil {
		return nil, err
	}
	result := make([]entities.PartsSupply, len(dtos))
	for i, dto := range dtos {
		result[i] = dto.ToDomain()
	}
	return result, nil
}

func (s *PartsSupplyRepository) GetByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error) {
	var dtos []dto.PartsSupplyDTO
	if err := s.db.WithContext(ctx).
		Joins("JOIN parts_supply_service_order_dtos ON parts_supply_service_order_dtos.parts_supply_id = parts_supply_dtos.id").
		Where("parts_supply_service_order_dtos.service_order_id = ?", serviceOrderID).
		Find(&dtos).Error; err != nil {
		return nil, err
	}

	var partsSupplies []entities.PartsSupply
	for _, dto := range dtos {
		partsSupplies = append(partsSupplies, dto.ToDomain())
	}
	return partsSupplies, nil
}
