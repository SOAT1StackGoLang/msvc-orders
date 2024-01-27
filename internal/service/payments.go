package service

import (
	"context"
	"errors"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/api"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"time"
)

type paymentsSvc struct {
	repo   persistence.PaymentRepository
	client *api.Client
	log    kitlog.Logger
}

func (p *paymentsSvc) GetPayment(ctx context.Context, paymentID uuid.UUID) (*models.Payment, error) {
	return p.repo.GetPayment(ctx, paymentID)

}

func (p *paymentsSvc) CreatePayment(ctx context.Context, order *models.Order) (*models.Payment, error) {
	payment := &models.Payment{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Price:     order.Price,
		OrderID:   order.ID,
		Status:    models.PAYMENT_STATUS_OPEN,
	}

	receipt, err := p.repo.CreatePayment(ctx, payment)

	_, err = p.client.CreatePayment(api.CreatePaymentRequest{Payment: api.Payment{
		ID:        payment.ID,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
		Price:     payment.Price,
		OrderID:   payment.OrderID,
		Status:    api.PaymentStatusPending,
	}})

	if err != nil {
		return nil, errors.New("failed creating payment")
	}

	return receipt, nil
}

func (p *paymentsSvc) UpdatePayment(ctx context.Context, paymentID uuid.UUID, status models.PaymentStatus) (*models.Payment, error) {
	payment, err := p.GetPayment(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	payment.UpdatedAt = time.Now()
	payment.Status = status

	updated, err := p.repo.UpdatePayment(ctx, payment)

	return updated, err
}

func NewPaymentsService(repo persistence.PaymentRepository, log kitlog.Logger) PaymentsService {
	return &paymentsSvc{
		repo: repo,
		log:  log,
	}
}
