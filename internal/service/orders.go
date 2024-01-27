package service

import (
	"context"
	"encoding/json"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	"github.com/SOAT1StackGoLang/msvc-orders/pkg/helpers"
	payments "github.com/SOAT1StackGoLang/msvc-payments/pkg/api"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore"
	production "github.com/SOAT1StackGoLang/msvc-production/pkg/api"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

const ordersPendingPayment = "orders_pending_payment"
const ordersPendingProduction = "orders_pending_production"

type pendingOrders struct {
	ID        string `json:"id"`
	PaymentID string `json:"payment_id"`
}

func MarshalPendingOrders(order models.Order) ([]byte, error) {
	return json.Marshal(&pendingOrders{
		ID:        order.ID.String(),
		PaymentID: order.PaymentID.String(),
	})
}

func UnmarshalPendingOrders(marshalled string) (*pendingOrders, error) {
	var out *pendingOrders
	err := json.Unmarshal([]byte(marshalled), out)
	return out, err
}

type ordersSvc struct {
	paymentsAPI   payments.PaymentAPI
	productionAPI production.ProductionAPI
	cache         datastore.RedisStore
	ordersRepo    persistence.OrdersRepository
	productsSvc   ProductsService
	paymentsSvc   PaymentsService
	log           kitlog.Logger
}

func NewOrdersService(
	repo persistence.OrdersRepository,
	prodSvc ProductsService,
	paySvc PaymentsService,
	log kitlog.Logger,
	cache datastore.RedisStore,
) OrdersService {
	return &ordersSvc{
		cache:       cache,
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

	orders, err := MarshalPendingOrders(*order)
	if err != nil {
		return nil, err
	}

	err = o.cache.RPush(ctx, ordersPendingPayment, orders)
	if err != nil {
		o.log.Log(
			"failed pushing into "+ordersPendingPayment,
			zap.Error(err),
		)
		return nil, err
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

func (o *ordersSvc) PoolOrderPaymentStatus() {
	var (
		ctx = context.Background()

		start, end    int64
		ids           []string
		err           error
		paymentStatus models.PaymentStatus
	)

	for {
		time.Sleep(1 * time.Second)
		ids, err = o.cache.LRange(ctx, ordersPendingPayment, start, end)
		if err != nil {
			o.log.Log(
				"error at cache",
				zap.Int64("start", start),
				zap.Int64("end", end),
				zap.Error(err),
			)
			break
		}
		if len(ids) == 0 {
			continue
		}

		paymentID, err := uuid.Parse(ids[0])
		if err != nil {
			o.log.Log(
				"error parsing uuid",
				zap.Int64("start", start),
				zap.Int64("end", end),
				zap.String("id", ids[0]),
				zap.Error(err),
			)
			break
		}

		payment, err := o.paymentsAPI.GetPayment(payments.GetPaymentRequest{PaymentID: paymentID})
		if err != nil {
			o.log.Log(
				"error getting payment from 3rd party",
				zap.Int64("start", start),
				zap.Int64("end", end),
				zap.String("id", paymentID.String()),
				zap.Error(err),
			)
			break
		}

		paymentStatus = models.PaymentStatusFromClearingService(payment.Status)
		switch payment.Status {
		case payments.PaymentStatusPaid:
			if _, err = o.paymentsSvc.UpdatePayment(ctx, payment.Payment.ID, paymentStatus); err != nil {
				o.log.Log("UpdatePayment INCONSISTENCY", zap.String("payment_id", payment.Payment.ID.String()), zap.Error(err))
				continue
			}

			if _, err = o.UpdateOrderStatus(ctx, payment.Payment.OrderID, models.ORDER_STATUS_RECEIVED); err != nil {
				o.log.Log("UpdateOrderStatus INCONSISTENCY", zap.String("order_id", payment.Payment.OrderID.String()), zap.Error(err))
				break
			}

			if err = o.cache.RPush(ctx, ordersPendingProduction, payment.Payment.OrderID.String()); err != nil {
				o.log.Log(ordersPendingProduction+" INCONSISTENCY", zap.String("order_id", payment.Payment.ID.String()), zap.Error(err))
			}

			err := o.cache.LIndex(ctx, ordersPendingPayment, start)
			if err != nil {
				o.log.Log(ordersPendingPayment+" INCONSISTENCY", zap.String("order_id", payment.Payment.ID.String()), zap.Error(err))
				break
			}
		}
	}
}

func (o *ordersSvc) NotifyProduction() {
	var (
		ctx = context.Background()

		start, end int64
		ids        []string
		err        error
	)

	for {
		time.Sleep(1 * time.Second)
		ids, err = o.cache.LRange(ctx, ordersPendingPayment, start, end)
		if err != nil {
			o.log.Log(
				"error at cache",
				zap.Int64("start", start),
				zap.Int64("end", end),
				zap.Error(err),
			)
			break
		}
		if len(ids) == 0 {
			continue
		}
		if err != nil {
			o.log.Log(
				"UpdateOrder failed",
				zap.String("order_id", orderID.String()),
				zap.Error(err),
			)
			return err
		}

		_, err = o.UpdateOrderStatus(ctx, orderID, models.ORDER_STATUS_PREPARING)
		if err != nil {
			o.log.Log(
				"UpdateOrderStatus failed",
				zap.String("order_id", orderID.String()),
				zap.Error(err),
			)
		}

		err := o.cache.LIndex(ctx, ordersPendingPayment, start)

		return err
	}
}
