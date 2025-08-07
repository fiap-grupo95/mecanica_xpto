package use_case

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/repository/parts_supply"
)

type IPartsSupplyUseCase interface {
	GetPartsSupplyByID(ctx context.Context, id uint) (entities.PartsSupply, error)
	CreatePartsSupply(ctx context.Context, partsSupply *entities.PartsSupply) (entities.PartsSupply, error)
	UpdatePartsSupply(ctx context.Context, partsSupply *entities.PartsSupply) error
	DeletePartsSupply(ctx context.Context, id uint) error
	ListPartsSupplies(ctx context.Context) ([]entities.PartsSupply, error)
}
type PartsSupplyUseCase struct {
	repo parts_supply.IPartsSupplyRepo
}

var _ IPartsSupplyUseCase = (*PartsSupplyUseCase)(nil)

func NewPartsSupplyUseCase(repo parts_supply.IPartsSupplyRepo) *PartsSupplyUseCase {
	return &PartsSupplyUseCase{repo: repo}
}

var ErrPartsSupplyNotFound = errors.New("parts supply not found")
var ErrPartsSupplyAlreadyExists = errors.New("parts supply already exists")

func (h *PartsSupplyUseCase) GetPartsSupplyByID(ctx context.Context, id uint) (entities.PartsSupply, error) {
	foundPartsSupply, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return entities.PartsSupply{}, err
	}

	if foundPartsSupply.ID == 0 {
		return entities.PartsSupply{}, ErrPartsSupplyNotFound
	}

	return foundPartsSupply, nil
}

func (h *PartsSupplyUseCase) CreatePartsSupply(ctx context.Context, partsSupply *entities.PartsSupply) (entities.PartsSupply, error) {
	//valide se já existe uma peça com o mesmo nome
	existingPartsSupply, err := h.repo.GetByName(ctx, partsSupply.Name)
	if err == nil && existingPartsSupply.ID != 0 {
		return entities.PartsSupply{}, ErrPartsSupplyAlreadyExists
	}

	return h.repo.Create(ctx, partsSupply)

}

func (h *PartsSupplyUseCase) UpdatePartsSupply(ctx context.Context, partsSupply *entities.PartsSupply) error {

	existingPartsSupply, err := h.repo.GetByID(ctx, partsSupply.ID)
	if err != nil {
		return errors.New("failed to retrieve parts supply")
	}

	if existingPartsSupply.ID == 0 {
		return ErrPartsSupplyNotFound
	}

	return h.repo.Update(ctx, partsSupply)

}

func (h *PartsSupplyUseCase) DeletePartsSupply(ctx context.Context, id uint) error {
	return h.repo.Delete(ctx, id)
}

func (h *PartsSupplyUseCase) ListPartsSupplies(ctx context.Context) ([]entities.PartsSupply, error) {
	return h.repo.List(ctx)
}
