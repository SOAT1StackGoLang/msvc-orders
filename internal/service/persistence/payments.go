package persistence

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentsPersistence struct {
	db  *gorm.DB
	log kitlog.Logger
}

func (p *paymentsPersistence) CreatePayment(ctx context.Context, in *models.Payment) (*models.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (p *paymentsPersistence) GetPayment(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (p *paymentsPersistence) UpdatePayment(ctx context.Context, in *models.Payment) (*models.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func NewPaymentsPersistence(db *gorm.DB, log kitlog.Logger) PaymentRepository {
	return &paymentsPersistence{
		db:  db,
		log: log,
	}
}
