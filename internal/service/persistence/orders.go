package persistence

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ordersPersistence struct {
	db  *gorm.DB
	log kitlog.Logger
}

func (o *ordersPersistence) GetOrder(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersPersistence) GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersPersistence) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersPersistence) UpdateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersPersistence) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (o *ordersPersistence) SetOrderAsPaid(ctx context.Context, payment *models.Payment) error {
	//TODO implement me
	panic("implement me")
}

func (o *ordersPersistence) ListOrdersByUser(ctx context.Context, limit, offset int, userID uuid.UUID) (*models.OrderList, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersPersistence) ListOrders(ctx context.Context, limit, offset int) (*models.OrderList, error) {
	//TODO implement me
	panic("implement me")
}

func NewOrdersPersistence(db *gorm.DB, log kitlog.Logger) OrdersRepository {
	return &ordersPersistence{
		db:  db,
		log: log,
	}
}
