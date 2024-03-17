package persistence

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

const ordersTable = "lanchonete_orders"

type ordersPersistence struct {
	db  *gorm.DB
	log kitlog.Logger
}

func (o *ordersPersistence) GetOrder(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	var out *models.Order
	order := &Order{}

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).
		Select("*").
		Where("id = ?", orderID).
		First(order).Error; err != nil {
		o.log.Log(
			"db failed getting order",
			zap.Error(err),
		)
		return nil, err
	}

	out = &models.Order{
		ID:        order.ID,
		UserID:    order.UserID,
		PaymentID: order.PaymentID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: time.Time{},
		DeletedAt: time.Time{},
		Price:     order.Price,
		Status:    orderStatusToModelStatus(order.Status),
		Products:  nil,
	}
	out = order.toModels()
	return out, err
}

func (o *ordersPersistence) GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*models.Order, error) {
	var out *models.Order
	order := &Order{}

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).
		Select("*").
		Where("payment_id = ?", paymentID).
		First(order).Error; err != nil {
		o.log.Log(
			"db failed getting order",
			zap.Error(err),
		)
		return nil, err
	}

	out = order.toModels()
	return out, err
}

func (o *ordersPersistence) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	in := orderFromModels(order)
	in.Status = ORDER_STATUS_OPEN

	columns := []string{"updated_at"}
	if in.UserID == uuid.Nil {
		columns = append(columns, "user_id")
	}

	if err := o.db.WithContext(ctx).Table(ordersTable).Omit(columns...).Create(&in).Error; err != nil {
		o.log.Log(
			"db failed at CreateOrder",
			zap.Any("order_input", order),
			zap.Error(err),
		)
		return nil, err
	}

	return in.toModels(), nil
}

func (o *ordersPersistence) UpdateOrder(ctx context.Context, in *models.Order) (*models.Order, error) {
	order := orderFromModels(in)

	order.UpdatedAt = sql.NullTime{
		Time:  in.UpdatedAt,
		Valid: true,
	}

	order.Status = orderStatusFromModel(in.Status)

	if err := o.db.WithContext(ctx).Table(ordersTable).
		Updates(&order).
		Where("id = ?", in.ID).
		Error; err != nil {
		o.log.Log(
			"db failed updating order",
			zap.Any("in_order", in),
			zap.Any("repo_order", order),
			zap.Error(err),
		)
		return nil, err
	}

	return order.toModels(), nil
}

func (o *ordersPersistence) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	deletedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if err := o.db.WithContext(ctx).Table(ordersTable).
		UpdateColumn("deleted_at", deletedAt).
		Where("order_id", orderID).
		Error; err != nil {
		o.log.Log(
			"db failed deleting order",
			zap.String("order_id", orderID.String()),
			zap.Error(err),
		)
	}
	return nil
}

func (o *ordersPersistence) SetOrderAsPaid(ctx context.Context, payment *models.Payment) error {
	return errors.New("Unimplemented")
}

func (o *ordersPersistence) ListOrdersByUser(ctx context.Context, limit, offset int, userID uuid.UUID) (*models.OrderList, error) {
	var orders []Order
	var total int64

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at ASC").
		Find(&orders).Error; err != nil {
		o.log.Log(
			"failed listing orders",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	if err = o.db.WithContext(ctx).Table(ordersTable).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		o.log.Log(
			"failed counting orders by user_id",
			zap.String("category", userID.String()),
			zap.Error(err),
		)
	}

	oList := &models.OrderList{}
	out := make([]*models.Order, 0, len(orders))

	for _, v := range orders {
		out = append(out, v.toModels())
	}

	oList.Orders = out
	oList.Total = total
	oList.Limit = limit
	oList.Offset = offset

	return oList, err
}

func (o *ordersPersistence) ListOrders(ctx context.Context, limit, offset int) (*models.OrderList, error) {
	var total int64

	var saveOrders []Order

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).
		Limit(limit).
		Offset(offset).
		Order("status DESC").
		Find(&saveOrders).Error; err != nil {
		o.log.Log(
			"failed listing orders",
			zap.Error(err),
		)
		return nil, err
	}

	if err = o.db.WithContext(ctx).Table(ordersTable).
		Where("status > ? AND status < ? ", ORDER_STATUS_WAITING_PAYMENT, ORDER_STATUS_FINISHED).
		Count(&total).Error; err != nil {
		o.log.Log(
			"failed counting orders",
			zap.Error(err),
		)
	}

	oList := &models.OrderList{}
	out := make([]*models.Order, 0, len(saveOrders))

	for _, v := range saveOrders {
		out = append(out, v.toModels())
	}

	oList.Orders = out
	oList.Total = total
	oList.Limit = limit
	oList.Offset = offset

	return oList, err
}

func NewOrdersPersistence(db *gorm.DB, log kitlog.Logger) OrdersRepository {
	return &ordersPersistence{
		db:  db,
		log: log,
	}
}
