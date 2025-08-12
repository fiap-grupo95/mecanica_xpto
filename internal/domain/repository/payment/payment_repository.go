package payment

import (
	"context"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"time"

	"gorm.io/gorm"
)

type IPaymentRepo interface {
	Create(ctx context.Context, so *entities.Payment) (*dto.PaymentDTO, error)
	GetByID(ctx context.Context, id uint) (*dto.PaymentDTO, error)
	GetByServiceOrderID(ctx context.Context, serviceOrderID uint) (*dto.PaymentDTO, error)
	List(ctx context.Context) ([]dto.PaymentDTO, error)
}

type PaymentRepository struct {
	db *gorm.DB
}

var _ IPaymentRepo = (*PaymentRepository)(nil)

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (p *PaymentRepository) Create(ctx context.Context, payment *entities.Payment) (*dto.PaymentDTO, error) {
	dto := dto.PaymentDTO{
		ServiceOrderID: payment.ServiceOrderID,
		PaymentDate:    time.Now(),
		Amount:         payment.Amount,
	}
	if err := p.db.Create(&dto).Error; err != nil {
		return nil, err
	}
	return &dto, nil
}

func (p *PaymentRepository) GetByID(ctx context.Context, id uint) (*dto.PaymentDTO, error) {
	var dto dto.PaymentDTO
	if err := p.db.Preload("ServiceOrder").First(&dto, id).Error; err != nil {
		return nil, err
	}
	return &dto, nil
}

func (p *PaymentRepository) GetByServiceOrderID(ctx context.Context, serviceOrderID uint) (*dto.PaymentDTO, error) {
	var dto dto.PaymentDTO
	if err := p.db.Preload("ServiceOrder").Where("service_order_id = ?", serviceOrderID).First(&dto).Error; err != nil {
		return nil, err
	}
	return &dto, nil
}

func (p *PaymentRepository) List(ctx context.Context) ([]dto.PaymentDTO, error) {
	var dtos []dto.PaymentDTO
	if err := p.db.Preload("ServiceOrder").Find(&dtos).Error; err != nil {
		return nil, err
	}
	return dtos, nil
}
