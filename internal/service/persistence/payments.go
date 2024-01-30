package persistence

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type paymentsPersistence struct {
	db  *gorm.DB
	log kitlog.Logger
}

const paymentTable = "lanchonete_payments"

func (p *paymentsPersistence) CreatePayment(ctx context.Context, in *models.Payment) (*models.Payment, error) {
	payment := paymentFromModels(in)

	if err := p.db.WithContext(ctx).Table(paymentTable).Create(&payment).Error; err != nil {
		p.log.Log(
			"db failed at CreatePayment",
			zap.Any("payment_input", in),
			zap.Error(err),
		)
		return nil, err
	}

	return payment.toModels(), nil
}

func (p *paymentsPersistence) GetPayment(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	payment := new(Payment)

	var err error
	if err = p.db.WithContext(ctx).Table(paymentTable).
		Select("*").
		Where("id = ?", id).
		First(payment).Error; err != nil {
		p.log.Log(
			"db failed getting payment",
			zap.Error(err),
		)
		return nil, err
	}

	out := payment.toModels()
	return out, err
}

func (p *paymentsPersistence) UpdatePayment(ctx context.Context, in *models.Payment) (*models.Payment, error) {
	payment := paymentFromModels(in)

	if err := p.db.WithContext(ctx).Table(paymentTable).
		Updates(&payment).
		Where("id = ?", in.ID).
		Error; err != nil {
		p.log.Log(
			"db failed updating payment",
			zap.Any("in_payment", in),
			zap.Any("repo_payment", payment),
			zap.Error(err),
		)
		return nil, err
	}

	return payment.toModels(), nil
}

func NewPaymentsPersistence(db *gorm.DB, log kitlog.Logger) PaymentRepository {
	return &paymentsPersistence{
		db:  db,
		log: log,
	}
}
