package service

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	"github.com/SOAT1StackGoLang/msvc-orders/pkg/helpers"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	return o.ordersRepo.GetOrder(ctx, id)
}

func (o *ordersSvc) GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*models.Order, error) {
	return o.ordersRepo.GetOrderByPaymentID(ctx, paymentID)
}

func (o *ordersSvc) CreateOrder(ctx context.Context, products []models.Product) (*models.Order, error) {
	var order *models.Order

	if len(products) == 0 {
		o.log.Log(
			"error at CreateOrder, must have at least one product in it",
			zap.Any("products", products),
			zap.Error(helpers.ErrInvalidInput),
		)
		return nil, helpers.ErrInvalidInput
	}

	for k, p := range products {
		fullProduct, err := o.productsSvc.GetProduct(ctx, p.ID)
		if err != nil {
			o.log.Log("CreateOrder failed due to invalid product",
				zap.String("product_id", p.ID.String()),
				zap.Any("requested_products", products),
				zap.Error(err),
			)
			return nil, err
		}
		products[k] = *fullProduct
	}

	order = &models.Order{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Status:    models.ORDER_STATUS_OPEN,
		Products:  products,
	}

	for _, v := range products {
		order.Price = order.Price.Add(v.Price)
	}

	return o.ordersRepo.CreateOrder(ctx, order)
}

func (o *ordersSvc) UpdateOrderItems(ctx context.Context, orderID uuid.UUID, products []models.Product) (*models.Order, error) {
	order, err := o.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		o.log.Log(
			"error at UpdateOrderItems, must have at least one product in it",
			zap.Any("inProducts", products),
			zap.Error(helpers.ErrBadRequest),
		)
		return nil, helpers.ErrBadRequest
	}

	for _, v := range products {
		order.Products = append(order.Products, v)
		order.Price = order.Price.Add(v.Price)
	}

	return o.ordersRepo.UpdateOrder(ctx, order)
}

func (o *ordersSvc) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	return o.ordersRepo.DeleteOrder(ctx, orderID)
}

func (o *ordersSvc) ListOrders(ctx context.Context, limit, offset int) (*models.OrderList, error) {
	return o.ordersRepo.ListOrders(ctx, limit, offset)
}

func (o *ordersSvc) Checkout(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order *models.Order

	order, err := o.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	order.Status = models.ORDER_STATUS_WAITING_PAYMENT

	payment, err := o.paymentsSvc.CreatePayment(ctx, order)
	if err != nil {
		return nil, err
	}
	order.PaymentID = payment.ID

	order.UpdatedAt = time.Now()
	order, err = o.ordersRepo.UpdateOrder(ctx, order)
	if err != nil {
		o.log.Log(
			"failed updating order status after checkout",
			zap.Error(err),
		)
	}

	return order, err
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
