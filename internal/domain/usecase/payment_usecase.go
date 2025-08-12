package usecase

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/repository/payment"
	serviceorder "mecanica_xpto/internal/domain/repository/service_order"

	"gorm.io/gorm"
)

var (
	ErrorPaymentNotFound         = errors.New("payment not found")
	ErrPaymentAlreadyExists      = errors.New("payment already exists")
	ErrPaymentAmountDoesNotMatch = errors.New("payment amount does not match service order estimate")
)

type IPaymentUseCase interface {
	CreatePayment(ctx context.Context, payment *entities.Payment) (*entities.Payment, error)
	GetPaymentByID(ctx context.Context, id uint) (*entities.Payment, error)
	ListPayments(ctx context.Context) ([]entities.Payment, error)
}

type PaymentUseCase struct {
	repo             payment.IPaymentRepo
	serviceOrderRepo serviceorder.IServiceOrderRepository
}

var _ IPaymentUseCase = (*PaymentUseCase)(nil)

func NewPaymentUseCase(repo payment.IPaymentRepo, serviceOrderRepo serviceorder.IServiceOrderRepository) *PaymentUseCase {
	return &PaymentUseCase{
		repo:             repo,
		serviceOrderRepo: serviceOrderRepo,
	}
}

func (p *PaymentUseCase) CreatePayment(ctx context.Context, payment *entities.Payment) (*entities.Payment, error) {
	serviceOrder, err := p.serviceOrderRepo.GetByID(payment.ServiceOrderID)
	if err != nil {
		return nil, err
	}
	if serviceOrder.Estimate != payment.Amount {
		return nil, ErrPaymentAmountDoesNotMatch
	}
	existingPayment, err := p.repo.GetByServiceOrderID(ctx, payment.ServiceOrderID)
	if err == nil && existingPayment.ID != 0 {
		return nil, ErrPaymentAlreadyExists
	}

	dto, err := p.repo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}
	return dto.ToDomain(), nil
}

func (p *PaymentUseCase) GetPaymentByID(ctx context.Context, id uint) (*entities.Payment, error) {
	paymentDTO, err := p.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entities.Payment{}, ErrorPaymentNotFound
		}
		return &entities.Payment{}, err
	}
	return paymentDTO.ToDomain(), nil
}

func (p *PaymentUseCase) ListPayments(ctx context.Context) ([]entities.Payment, error) {
	dtos, err := p.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	payments := make([]entities.Payment, 0, len(dtos))
	for _, dto := range dtos {
		payments = append(payments, *dto.ToDomain())
	}
	return payments, nil
}
