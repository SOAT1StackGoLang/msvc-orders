package service

import (
	"context"
	"errors"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	paymentapi "github.com/SOAT1StackGoLang/msvc-payments/pkg/api"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"time"
)

type paymentsSvc struct {
	repo   persistence.PaymentRepository
	client paymentapi.PaymentAPI
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

	_, err = p.client.CreatePayment(paymentapi.CreatePaymentRequest{Payment: paymentapi.Payment{
		ID:        payment.ID,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
		Price:     payment.Price,
		OrderID:   payment.OrderID,
		Status:    paymentapi.PaymentStatusPending,
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

func NewPaymentsService(repo persistence.PaymentRepository, api paymentapi.PaymentAPI, log kitlog.Logger) PaymentsService {
	return &paymentsSvc{
		client: api,
		repo:   repo,
		log:    log,
	}
}
