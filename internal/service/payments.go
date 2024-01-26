package service

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/google/uuid"
)

type paymentsSvc struct {
	repo, log any
}

func (p *paymentsSvc) GetPayment(ctx context.Context, paymentID uuid.UUID) (*models.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (p *paymentsSvc) CreatePayment(ctx context.Context, orderID *models.Order) (*models.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (p *paymentsSvc) UpdatePayment(ctx context.Context, paymentID uuid.UUID, status models.PaymentStatus) (*models.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func NewPaymentsService(repo, log any) PaymentsService {
	return &paymentsSvc{
		repo: repo,
		log:  log,
	}
}
