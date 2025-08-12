package usecase

import (
	"context"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/repository/additional_repair"
	"mecanica_xpto/internal/domain/repository/parts_supply"

	"github.com/rs/zerolog/log"

	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/repository/service"
	serviceorder "mecanica_xpto/internal/domain/repository/service_order"
)

type IAdditionalRepairUseCase interface {
	CreateAdditionalRepair(ctx context.Context, adr entities.AdditionalRepair) error
	AddPartSupplyAndService(ctx context.Context, adrId uint, adr entities.AdditionalRepair) error
	RemovePartSupplyAndService(ctx context.Context, adrId uint, adr entities.AdditionalRepair) error
	GetAdditionalRepair(ctx context.Context, additionalRepairId uint) (entities.AdditionalRepair, error)
}

type AdditionalRepairUseCase struct {
	repo            additional_repair.IAdditionalRepairRepository
	repoOS          serviceorder.IServiceOrderRepository
	serviceRepo     service.IServiceRepo
	partsSupplyRepo parts_supply.IPartsSupplyRepo
}

var _ IAdditionalRepairUseCase = (*AdditionalRepairUseCase)(nil)

func NewSOAdditionalRepairUseCase(repo additional_repair.IAdditionalRepairRepository, repoOS serviceorder.IServiceOrderRepository, serviceRepo service.IServiceRepo, partsSupplyRepo parts_supply.IPartsSupplyRepo) *AdditionalRepairUseCase {
	return &AdditionalRepairUseCase{
		repo:            repo,
		repoOS:          repoOS,
		serviceRepo:     serviceRepo,
		partsSupplyRepo: partsSupplyRepo,
	}
}

func (u *AdditionalRepairUseCase) CreateAdditionalRepair(ctx context.Context, adr entities.AdditionalRepair) error {
	_, err := u.repoOS.GetByID(adr.ServiceOrderID)
	if err != nil {
		log.Error().Msgf("error finding service order with id %d: %v", adr.ServiceOrderID, err)
		return err
	}
	var estimatedPartsSupply float64
	var estimatedService float64

	arStatus := dto.AdditionalRepairStatusDTO{
		Description: "IN_ANALYSIS",
	}
	var listServices []dto.ServiceDTO
	var listPartsSupply []dto.PartsSupplyDTO
	listServices, estimatedService, err = u.addServiceToAdditionalRepair(ctx, adr.Services)
	if err != nil {
		return err
	}
	listPartsSupply, estimatedPartsSupply, err = u.addPartsSupplyToAdditionalRepair(ctx, adr.PartsSupplies)
	if err != nil {
		return err
	}

	additionalRepair := dto.AdditionalRepairDTO{
		ServiceOrderID: adr.ServiceOrderID,
		Description:    adr.Description,
		ARStatus:       arStatus,
		Estimate:       estimatedService + estimatedPartsSupply,
		Services:       listServices,
		PartsSupplies:  listPartsSupply,
	}

	err = u.repo.Create(&additionalRepair)
	if err != nil {
		log.Error().Msgf("Error creating additional repair: %v", err)
		return err
	}
	return nil
}

func (u *AdditionalRepairUseCase) AddPartSupplyAndService(ctx context.Context, adrId uint, adr entities.AdditionalRepair) error {
	var estimatedPartsSupply float64
	var estimatedService float64
	var err error

	var listServices []dto.ServiceDTO
	var listPartsSupply []dto.PartsSupplyDTO
	listServices, estimatedService, err = u.addServiceToAdditionalRepair(ctx, adr.Services)
	if err != nil {
		return err
	}
	listPartsSupply, estimatedPartsSupply, err = u.addPartsSupplyToAdditionalRepair(ctx, adr.PartsSupplies)
	if err != nil {
		return err
	}

	additionalRepair := dto.AdditionalRepairDTO{
		ServiceOrderID: adr.ServiceOrderID,
		Description:    adr.Description,
		Estimate:       estimatedService + estimatedPartsSupply,
		Services:       listServices,
		PartsSupplies:  listPartsSupply,
	}
	err = u.repo.AddPartSupplyAndService(adrId, &additionalRepair)
	if err != nil {
		log.Error().Msgf("Error adding part suplly and services for additional repair: %v", err)
		return err
	}
	return nil
}

func (u *AdditionalRepairUseCase) RemovePartSupplyAndService(ctx context.Context, adrId uint, adr entities.AdditionalRepair) error {
	var estimatedPartsSupply float64
	var estimatedService float64
	var err error

	var listServices []dto.ServiceDTO
	var listPartsSupply []dto.PartsSupplyDTO
	listServices, estimatedService, err = u.addServiceToAdditionalRepair(ctx, adr.Services)
	if err != nil {
		return err
	}
	listPartsSupply, estimatedPartsSupply, err = u.addPartsSupplyToAdditionalRepair(ctx, adr.PartsSupplies)
	if err != nil {
		return err
	}

	additionalRepair := dto.AdditionalRepairDTO{
		ServiceOrderID: adr.ServiceOrderID,
		Description:    adr.Description,
		Estimate:       estimatedService + estimatedPartsSupply,
		Services:       listServices,
		PartsSupplies:  listPartsSupply,
	}
	err = u.repo.AddPartSupplyAndService(adrId, &additionalRepair)
	if err != nil {
		log.Error().Msgf("Error adding part suplly and services for additional repair: %v", err)
		return err
	}
	return nil
}

func (u *AdditionalRepairUseCase) GetAdditionalRepair(ctx context.Context, additionalRepairId uint) (entities.AdditionalRepair, error) {
	additionalRepairDto, err := u.repo.GetByID(additionalRepairId)
	if err != nil {
		log.Error().Msgf("error finding additional repair with id %d: %v", additionalRepairId, err)
		return entities.AdditionalRepair{}, err
	}
	if additionalRepairDto == nil {
		log.Error().Msgf("additional repair with id %d not found", additionalRepairId)
		return entities.AdditionalRepair{}, ErrServiceOrderNotFound
	}
	return additionalRepairDto.ToDomain(), nil
}

func (u *AdditionalRepairUseCase) addPartsSupplyToAdditionalRepair(ctx context.Context, partsSupplyId []entities.PartsSupply) ([]dto.PartsSupplyDTO, float64, error) {
	var listPartsSupply []dto.PartsSupplyDTO
	var estimatedPrice float64

	for _, ps := range partsSupplyId {
		psDto, err := u.partsSupplyRepo.GetByID(ctx, ps.ID)
		if err != nil {
			log.Error().Msgf("error finding parts supply with id %d: %v", ps.ID, err)
		}
		if psDto.ID == 0 {
			log.Error().Msgf("parts supply with id %d not found", ps.ID)
			return listPartsSupply, estimatedPrice, ErrServiceNotFound
		}
		estimatedPrice += psDto.Price
		listPartsSupply = append(listPartsSupply, dto.PartsSupplyDTO{ID: ps.ID})
	}
	return listPartsSupply, estimatedPrice, nil
}

func (u *AdditionalRepairUseCase) addServiceToAdditionalRepair(ctx context.Context, services []entities.Service) ([]dto.ServiceDTO, float64, error) {
	var listServices []dto.ServiceDTO
	var estimatedPrice float64

	for _, s := range services {
		serviceDto, err := u.serviceRepo.GetByID(ctx, s.ID)
		if err != nil {
			log.Error().Msgf("error finding service with id %d: %v", s.ID, err)
		}
		if serviceDto.ID == 0 {
			log.Error().Msgf("service with id %d not found", s.ID)
			return listServices, estimatedPrice, ErrServiceNotFound
		}
		estimatedPrice += serviceDto.Price
		listServices = append(listServices, dto.ServiceDTO{ID: s.ID})
	}
	return listServices, estimatedPrice, nil
}
