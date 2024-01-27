package service

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"time"
)

type ordersSvc struct {
	ordersRepo  persistence.OrdersRepository
	productsSvc ProductsService
	paymentsSvc PaymentsService
	log         kitlog.Logger
}

func NewOrdersService(repo persistence.OrdersRepository, prodSvc ProductsService, paySvc PaymentsService, log kitlog.Logger) OrdersService {
	return &ordersSvc{
		ordersRepo:  repo,
		productsSvc: prodSvc,
		paymentsSvc: paySvc,
		log:         log,
	}
}

func (o *ordersSvc) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersSvc) GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*models.Order, error) {
	return o.ordersRepo.GetOrderByPaymentID(ctx, paymentID)
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
	order, err := o.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	order.Status = status
	order.UpdatedAt = time.Now()

	return o.ordersRepo.UpdateOrder(ctx, order)
}
