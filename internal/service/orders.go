package service

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/google/uuid"
)

type ordersSvc struct {
	repo any
	log  any
}

func (o *ordersSvc) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersSvc) GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersSvc) CreateOrder(ctx context.Context, products []models.Product) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersSvc) UpdateOrderItems(ctx context.Context, orderID uuid.UUID, products []models.Product) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersSvc) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (o *ordersSvc) ListOrders(ctx context.Context, limit, offset int) (*models.OrderList, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersSvc) Checkout(ctx context.Context, paymentID uuid.UUID) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersSvc) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func NewOrdersService(repo, log any) OrdersService {
	return &ordersSvc{
		repo: repo,
		log:  log,
	}
}
