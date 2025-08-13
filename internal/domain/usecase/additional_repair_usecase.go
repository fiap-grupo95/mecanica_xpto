package usecase

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/repository/additional_repair"
	"mecanica_xpto/internal/domain/repository/parts_supply"

	"github.com/rs/zerolog/log"

	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/repository/service"
	serviceorder "mecanica_xpto/internal/domain/repository/service_order"
)

var (
	ErrAdditionalRepairNotFound = errors.New("additional repair not found")
	ErrStatusNotPermitted       = errors.New("additional repair status not permitted")
)

type IAdditionalRepairUseCase interface {
	CreateAdditionalRepair(ctx context.Context, adr entities.AdditionalRepair) error
	AddPartSupplyAndService(ctx context.Context, adrId uint, adr entities.AdditionalRepair) error
	RemovePartSupplyAndService(ctx context.Context, adrId uint, adr entities.AdditionalRepair) error
	GetAdditionalRepair(ctx context.Context, additionalRepairId uint) (entities.AdditionalRepair, error)
	CustomerApprovalStatus(ctx context.Context, additionalRepairId uint, status entities.AdditionalRepairStatusDTO) error
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

func (u *AdditionalRepairUseCase) AddPartSupplyAndService(ctx context.Context, additionalRepairId uint, adr entities.AdditionalRepair) error {
	additionalRepairDto, err := u.repo.GetByID(additionalRepairId)
	if err != nil {
		log.Error().Msgf("error finding additional repair with id %d: %v", additionalRepairId, err)
		return err
	}
	if err := u.ValidateAdditionalRepairStatus(additionalRepairDto.ARStatus.Description); err != nil {
		return err
	}
	var estimatedPartsSupply float64
	var estimatedService float64

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

	updated := dto.AdditionalRepairDTO{
		ServiceOrderID: adr.ServiceOrderID,
		Description:    adr.Description,
		Estimate:       estimatedService + estimatedPartsSupply,
		Services:       listServices,
		PartsSupplies:  listPartsSupply,
	}
	err = u.repo.AddPartSupplyAndService(additionalRepairDto, &updated)
	if err != nil {
		log.Error().Msgf("Error adding part suplly and services for additional repair: %v", err)
		return err
	}
	return nil
}

func (u *AdditionalRepairUseCase) RemovePartSupplyAndService(ctx context.Context, additionalRepairId uint, adr entities.AdditionalRepair) error {
	additionalRepairDto, err := u.repo.GetByID(additionalRepairId)
	if err != nil {
		log.Error().Msgf("error finding additional repair with id %d: %v", additionalRepairId, err)
		return err
	}
	if err := u.ValidateAdditionalRepairStatus(additionalRepairDto.ARStatus.Description); err != nil {
		return err
	}
	var estimatedPartsSupply float64
	var estimatedService float64

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

	updated := dto.AdditionalRepairDTO{
		ServiceOrderID: adr.ServiceOrderID,
		Description:    adr.Description,
		Estimate:       estimatedService + estimatedPartsSupply,
		Services:       listServices,
		PartsSupplies:  listPartsSupply,
	}
	err = u.repo.AddPartSupplyAndService(additionalRepairDto, &updated)
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

func (u *AdditionalRepairUseCase) CustomerApprovalStatus(ctx context.Context, additionalRepairId uint, status entities.AdditionalRepairStatusDTO) error {
	additionalRepairDto, err := u.repo.GetByID(additionalRepairId)
	if err != nil {
		log.Error().Msgf("error finding additional repair with id %d: %v", additionalRepairId, err)
		return err
	}
	if err := u.ValidateAdditionalRepairStatus(additionalRepairDto.ARStatus.Description); err != nil {
		return err
	}

	if status.ApprovalStatus == "DENIED" {
		for _, ps := range additionalRepairDto.PartsSupplies {
			err := unreservePartsSupply(ctx, ps.ToDomain(), u.partsSupplyRepo)
			if err != nil {
				log.Error().Msgf("Error unreserving parts supply: %v", err)
			}
		}
		log.Info().Msgf("Customer rejected additional repair with id %d", additionalRepairId)
		return nil
	} else {
		for _, ps := range additionalRepairDto.PartsSupplies {
			relation, err := u.repo.GetPartsSupplyAdditionalRepair(ps.ID, additionalRepairDto.ID)
			if err != nil {
				log.Error().Msgf("Error getting parts supply additional repair relation: %v", err)
				return err
			}

			entity := entities.PartsSupply{
				ID:              ps.ID,
				QuantityReserve: relation.Quantity, // Use the quantity from the relationship
				QuantityTotal:   relation.Quantity,
			}
			err = releaseReservedPartsSupply(ctx, entity, u.partsSupplyRepo)
			if err != nil {
				log.Error().Msgf("Error releasing reserved parts supply: %v", err)
				return err
			}
		}
	}

	err = u.repo.CustomerApprovalStatus(additionalRepairId, status)
	if err != nil {
		log.Error().Msgf("error updating customer approval with id %d: %v", additionalRepairId, err)
		return err
	}

	err = u.repoOS.UpdateEstimate(additionalRepairDto.ServiceOrderID, additionalRepairDto.Estimate)
	if err != nil {
		log.Error().Msgf("error updating service order estimate with id %d: %v", additionalRepairDto.ServiceOrderID, err)
		return err
	}

	return nil
}

func (u *AdditionalRepairUseCase) addPartsSupplyToAdditionalRepair(ctx context.Context, partsSupplyId []entities.PartsSupply) ([]dto.PartsSupplyDTO, float64, error) {
	var listPartsSupply []dto.PartsSupplyDTO
	var estimatedPrice float64

	for _, ps := range partsSupplyId {
		err := reservePartsSupply(ctx, ps, u.partsSupplyRepo)
		if err != nil {
			log.Error().Msgf("Error reserving parts supply: %v", err)
			return nil, 0, err
		}
		psDto, err := u.partsSupplyRepo.GetByID(ctx, ps.ID)
		if err != nil {
			log.Error().Msgf("error finding parts supply with id %d: %v", ps.ID, err)
		}
		if psDto.ID == 0 {
			log.Error().Msgf("parts supply with id %d not found", ps.ID)
			return listPartsSupply, estimatedPrice, ErrServiceNotFound
		}
		psDto.QuantityReserve = ps.QuantityReserve
		estimatedPrice += psDto.Price
		listPartsSupply = append(listPartsSupply, dto.PartsSupplyDTO{ID: ps.ID, QuantityReserve: ps.QuantityReserve, QuantityTotal: psDto.QuantityTotal})
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

func (u *AdditionalRepairUseCase) ValidateAdditionalRepairStatus(status string) error {
	if status != "IN_ANALYSIS" {
		log.Error().Msgf("invalid additional repair status: %s", status)
		return ErrStatusNotPermitted
	}
	return nil
}

func (u *AdditionalRepairUseCase) releaseReservedPartsSupply(ctx context.Context, additionalRepairDto *dto.AdditionalRepairDTO) error {
	for _, ps := range additionalRepairDto.PartsSupplies {

		entity := entities.PartsSupply{
			ID:              ps.ID,
			QuantityReserve: ps.QuantityReserve,
			QuantityTotal:   ps.QuantityTotal,
		}
		err := releaseReservedPartsSupply(ctx, entity, u.partsSupplyRepo)
		if err != nil {
			log.Error().Msgf("Error releasing reserved parts supply: %v", err)
			return err
		}
	}
	return nil
}
