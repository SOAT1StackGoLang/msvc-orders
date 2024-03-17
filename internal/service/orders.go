package service

import (
	"context"
	"encoding/json"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	"github.com/SOAT1StackGoLang/msvc-orders/pkg/helpers"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/messages"
	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
	productionmsgs "github.com/SOAT1StackGoLang/msvc-production/pkg/messages"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type ordersSvc struct {
	cache       datastore.RedisStore
	ordersRepo  persistence.OrdersRepository
	productsSvc ProductsService
	paymentsSvc PaymentsService
	log         kitlog.Logger
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

	go svc.SubscribeToPaymentUpdates()
	go svc.SubscribeToProductionUpdates()

	return svc
}

func (o *ordersSvc) SubscribeToProductionUpdates() {
	ctx := context.Background()
	sub, err := o.cache.Subscribe(ctx, productionmsgs.ProductionStatusChannel)
	if err != nil {
		logger.Info("error subscribing to order status updates")
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
	var in productionmsgs.ProductionStatusChangedMessage
	err := json.Unmarshal([]byte(message), &in)
	if err != nil {
		o.log.Log(
			"error unmarshalling order status update",
			zap.Error(err),
		)
		return
	}

	orderID, err := uuid.Parse(in.OrderID)
	if err != nil {
		o.log.Log(
			"error parsing order id",
			zap.Error(err),
		)
		return
	}

	status := models.OrderStatus(in.Status)
	out, err := o.UpdateOrderStatus(context.Background(), uuid.MustParse(in.OrderID), status)
	if err != nil {
		return
	}

	o.log.Log("Order Status updated",
		zap.String("status", string(in.Status)),
		zap.Any("order_id", orderID),
		zap.Time("updated_at", out.UpdatedAt),
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
	order.UserID = userID

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

func (o *ordersSvc) SubscribeToPaymentUpdates() {
	ctx := context.Background()
	sub, err := o.cache.Subscribe(ctx, messages.PaymentStatusResponseChannel)
	if err != nil {
		logger.Info("error subscribing to payment status updates")
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-sub:
			o.handlePaymentStatusChangedMessage(msg.Payload)
		}
	}
}

func (o *ordersSvc) handlePaymentStatusChangedMessage(msg string) {
	var (
		in  messages.PaymentStatusChangedMessage
		err error
		out *models.Order
	)

	err = json.Unmarshal([]byte(msg), &in)
	if err != nil {
		o.log.Log(
			"error unmarshalling order status update",
			err,
		)
		return
	}

	status := models.PaymentStatusFromClearingService(in.Status)

	switch status {
	case models.PAYMENT_STATUS_APPROVED:
		orderID, err := uuid.Parse(in.OrderID)
		if err != nil {
			o.log.Log(
				"error parsing order id",
				err,
			)
			return
		}
		_, err = o.paymentsSvc.UpdatePayment(context.Background(), uuid.MustParse(in.ID), models.PAYMENT_STATUS_APPROVED)
		if err != nil {
			return
		}

		if out, err = o.UpdateOrderStatus(context.Background(), orderID, models.ORDER_STATUS_RECEIVED); err != nil {
			return
		}
		err = o.publishMessage(context.Background(), productionmsgs.OrderSentMessage{
			OrderID: out.ID.String(),
			Status:  productionmsgs.OrderStatus(string(out.Status)),
		}, productionmsgs.ProductionChannel)

		if out, err = o.UpdateOrderStatus(context.Background(), uuid.MustParse(in.OrderID), models.ORDER_STATUS_PREPARING); err != nil {
			return
		}

	case models.PAYMENT_SATUS_REFUSED:
		_, err = o.paymentsSvc.UpdatePayment(context.Background(), uuid.MustParse(in.ID), models.PAYMENT_SATUS_REFUSED)
		if err != nil {
			return
		}
		if out, err = o.UpdateOrderStatus(context.Background(), uuid.MustParse(in.OrderID), models.ORDER_STATUS_CANCELED); err != nil {
			return
		}
		err = o.publishMessage(context.Background(), productionmsgs.OrderSentMessage{
			OrderID: out.ID.String(),
			Status:  productionmsgs.OrderStatus(string(out.Status)),
		}, productionmsgs.ProductionChannel)
	}
}

func (o *ordersSvc) publishMessage(ctx context.Context, msg any, channel string) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		o.log.Log(
			"error marshalling message",
			err,
		)
		return err

	}
	err = o.cache.Publish(ctx, channel, bytes)
	if err != nil {
		o.log.Log(
			"error publishing message",
			err,
		)
	}

	return err
}
