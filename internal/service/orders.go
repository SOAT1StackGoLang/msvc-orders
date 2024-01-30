package service

import (
	"context"
	"encoding/json"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	"github.com/SOAT1StackGoLang/msvc-orders/pkg/helpers"
	payments "github.com/SOAT1StackGoLang/msvc-payments/pkg/api"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore"
	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
	productionpkg "github.com/SOAT1StackGoLang/msvc-production/pkg"
	production "github.com/SOAT1StackGoLang/msvc-production/pkg/api"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ordersPendingPayment = "orders_pending_payment"
const ordersProcessingPayment = "orders_processing_payment"
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
	svc := &ordersSvc{
		cache:       cache,
		ordersRepo:  repo,
		productsSvc: prodSvc,
		paymentsSvc: paySvc,
		log:         log,
	}

	go svc.subscribeToProductionSvc()

	return svc
}

func (o *ordersSvc) subscribeToProductionSvc() {
	ctx := context.Background()
	sub, err := o.cache.Subscribe(ctx, productionpkg.OrderStatusChannel)
	if err != nil {
		o.log.Log(
			"error subscribing to order status updates",
			zap.Error(err),
		)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-sub:
			o.handleProductionMessage(msg.Payload)
		}
	}
}

func (o *ordersSvc) handleProductionMessage(message string) {
	var in models.OrderProductionNotification
	err := json.Unmarshal([]byte(message), &in)
	if err != nil {
		o.log.Log(
			"error unmarshalling order status update",
			zap.Error(err),
		)
		return
	}

	_, err = o.UpdateOrderStatus(context.Background(), in.ID, in.Status)
	if err != nil {
		return
	}

	o.log.Log("Order Status updated",
		zap.String("status", string(in.Status)),
		zap.Any("order_id", in.ID),
		zap.Time("updated_at", in.UpdatedAt),
	)

	return
}

func (o *ordersSvc) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return o.ordersRepo.GetOrder(ctx, id)
}

func (o *ordersSvc) GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*models.Order, error) {
	return o.ordersRepo.GetOrderByPaymentID(ctx, paymentID)
}

func (o *ordersSvc) CreateOrder(ctx context.Context, products []models.Product, userID uuid.UUID) (*models.Order, error) {
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

	if userID != uuid.Nil {
		order.UserID = userID
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

	err = o.cache.LPush(ctx, ordersPendingPayment, orders)
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

func (o *ordersSvc) ProcessPayment() {
	var paidChannel = make(chan *pendingOrders, 1)
	logger.Info("Starting payment pooling...")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go o.processPayments(paidChannel)

	for {
		select {
		case <-shutdown:
			logger.Info("Shutting down ProcessPayment...")
			return
		case paid := <-paidChannel:
			logger.Info(zap.String("payment_id", paid.PaymentID).String)
			_, err := o.paymentsSvc.UpdatePayment(context.Background(), uuid.MustParse(paid.PaymentID), models.PAYMENT_STATUS_APPROVED)
			if err != nil {
				logger.Error("failed updating payment status")
				continue
			}

			_, err = o.productionAPI.UpdateOrder(production.UpdateOrderRequest{
				OrderID: paid.ID,
				Status:  production.ORDER_STATUS_PREPARING,
			})
			if err != nil {
				logger.Error("failed sending to production")
				continue
			}

			_, err = o.UpdateOrderStatus(context.Background(), uuid.MustParse(paid.ID), models.ORDER_STATUS_PREPARING)
			if err != nil {
				logger.Error("failed updating oredr status")
				continue
			}
		}
	}

}

func (o *ordersSvc) processPayments(paidChannel chan *pendingOrders) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-shutdown:
			logger.Info("Shutting down processPayments...")
			return
		default:
			time.Sleep(1 * time.Second)
			registers, err := o.cache.LRange(context.Background(), ordersPendingPayment, 0, -1)
			if err != nil {
				logger.Error("error retrieving from cache")
				return
			}

			for _, v := range registers {
				var pO *pendingOrders
				err = json.Unmarshal([]byte(v), pO)
				if err != nil {
					logger.Error("error unmarshalling from cache")
					continue
				}

				pUID, err := uuid.Parse(pO.PaymentID)
				if err != nil {
					logger.Error("error parsing uuid from cache")
					continue
				}

				payment, err := o.paymentsAPI.GetPayment(payments.GetPaymentRequest{PaymentID: pUID})
				if err != nil {
					logger.Error("error getting payment status")
					continue
				}

				if payment.Payment.Status == payments.PaymentStatusPaid {
					paidChannel <- &pendingOrders{
						ID:        payment.Payment.OrderID.String(),
						PaymentID: payment.Payment.ID.String(),
					}
					err = o.cache.LREM(context.Background(), ordersPendingPayment, 1, v)
					if err != nil {
						logger.Error("error deleting from cache")
						return
					}
				}
			}
		}
	}
}

// TODO HANDLE CONSUMED PRODUCITON MSG AND UPDATE ORDER STATUS
