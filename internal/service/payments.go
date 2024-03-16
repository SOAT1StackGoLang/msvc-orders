package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/messages"
	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"time"
)

type paymentsSvc struct {
	repo  persistence.PaymentRepository
	log   kitlog.Logger
	redis datastore.RedisStore
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

	outPayment := messages.PaymentCreationRequestMessage{
		ID:        receipt.ID.String(),
		CreatedAt: receipt.CreatedAt.Format(time.RFC3339),
		UpdatedAt: receipt.UpdatedAt.Format(time.RFC3339),
		Price:     receipt.Price.InexactFloat64(),
		OrderID:   receipt.OrderID.String(),
		Status:    string(receipt.Status),
	}

	bytes, err := json.Marshal(outPayment)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed marshalling payment creation request", err.Error()))
		return nil, err
	}
	err = p.redis.Publish(ctx, messages.OrderPaymentCreationRequestChannel, bytes)
	if err != nil {
		return nil, err
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

func NewPaymentsService(
	repo persistence.PaymentRepository,
	log kitlog.Logger,
	cache datastore.RedisStore,
) PaymentsService {
	return &paymentsSvc{
		repo:  repo,
		log:   log,
		redis: cache,
	}
}
