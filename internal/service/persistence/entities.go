package persistence

import (
	"database/sql"
	"encoding/json"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"id,primaryKey"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
	Document  string
	Name      string
	Email     string
	IsAdmin   bool
}

type Product struct {
	ID          uuid.UUID       `gorm:"id,primaryKey" json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   sql.NullTime    `json:"updated_at"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CategoryID  uuid.UUID       `json:"category_id"`
	Price       decimal.Decimal `json:"price"`
}

func (p *Product) toModel() models.Product {
	out := models.Product{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}

	if p.UpdatedAt.Valid {
		out.UpdatedAt = p.UpdatedAt.Time
	}
	return out
}

func productsToModel(in json.RawMessage) []models.Product {
	var products []Product
	err := json.Unmarshal(in, &products)
	if err != nil {
		// TODO handle properly
		panic("failed to unmarshal products")
	}

	var outProducts []models.Product
	for _, v := range products {
		outProducts = append(outProducts, v.toModel())
	}

	return outProducts
}

func productFromModel(in []models.Product) json.RawMessage {
	var products []OrderProduct
	for _, p := range in {
		oP := OrderProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			CategoryID:  p.CategoryID,
			Price:       p.Price,
		}
		products = append(products, oP)
	}

	productsJSON, err := json.Marshal(products)
	if err != nil {
		//TODO handle properly
		panic("failed to marshal products")
	}

	return productsJSON
}

type OrderProduct struct {
	ID          uuid.UUID       `gorm:"id,primaryKey" json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CategoryID  uuid.UUID       `json:"category_id"`
	Price       decimal.Decimal `json:"price"`
}

type ProductList struct {
	products      []*models.Product
	limit, offset int
	total         int64
}

type Category struct {
	ID        uuid.UUID `gorm:"id,primaryKey"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Name      string
}

type Order struct {
	ID        uuid.UUID `gorm:"id,primaryKey"`
	UserID    uuid.UUID
	PaymentID uuid.UUID
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
	Price     decimal.Decimal
	Status    OrderStatus
	Products  json.RawMessage `json:"products" gorm:"type:jsonb"`
}

func (o *Order) toModels() *models.Order {
	out := &models.Order{
		ID:        o.ID,
		UserID:    o.UserID,
		PaymentID: o.PaymentID,
		CreatedAt: o.CreatedAt,
		Price:     o.Price,
	}
	if o.UpdatedAt.Valid {
		out.UpdatedAt = o.UpdatedAt.Time
	}
	if o.DeletedAt.Valid {
		out.DeletedAt = o.DeletedAt.Time
	}

	out.Status = orderStatusToModelStatus(o.Status)
	out.Products = productsToModel(o.Products)

	return out
}

func orderFromModels(in *models.Order) *Order {
	return &Order{
		ID:        in.ID,
		UserID:    in.UserID,
		PaymentID: in.PaymentID,
		CreatedAt: in.CreatedAt,
		Price:     in.Price,
		Status:    orderStatusFromModel(in.Status),
		Products:  productFromModel(in.Products),
	}
}

type OrderStatus int

const (
	ORDER_STATUS_UNSET OrderStatus = iota
	ORDER_STATUS_OPEN
	ORDER_STATUS_WAITING_PAYMENT
	ORDER_STATUS_RECEIVED
	ORDER_STATUS_PREPARING
	ORDER_STATUS_DONE
	ORDER_STATUS_FINISHED
	ORDER_STATUS_CANCELED
	ORDER_STATUS_FAILED_PAYMENT
)

func orderStatusToModelStatus(in OrderStatus) models.OrderStatus {
	switch in {
	case ORDER_STATUS_OPEN:
		return models.ORDER_STATUS_OPEN
	case ORDER_STATUS_WAITING_PAYMENT:
		return models.ORDER_STATUS_WAITING_PAYMENT
	case ORDER_STATUS_RECEIVED:
		return models.ORDER_STATUS_RECEIVED
	case ORDER_STATUS_PREPARING:
		return models.ORDER_STATUS_PREPARING
	case ORDER_STATUS_DONE:
		return models.ORDER_STATUS_DONE
	case ORDER_STATUS_FINISHED:
		return models.ORDER_STATUS_FINISHED
	case ORDER_STATUS_CANCELED:
		return models.ORDER_STATUS_CANCELED
	case ORDER_STATUS_FAILED_PAYMENT:
		return models.ORDER_STATUS_FAILED_PAYMENT
	default:
		return models.ORDER_STATUS_UNSET
	}
}

func orderStatusFromModel(in models.OrderStatus) OrderStatus {
	switch in {
	case models.ORDER_STATUS_OPEN:
		return ORDER_STATUS_OPEN
	case models.ORDER_STATUS_WAITING_PAYMENT:
		return ORDER_STATUS_WAITING_PAYMENT
	case models.ORDER_STATUS_RECEIVED:
		return ORDER_STATUS_RECEIVED
	case models.ORDER_STATUS_PREPARING:
		return ORDER_STATUS_PREPARING
	case models.ORDER_STATUS_DONE:
		return ORDER_STATUS_DONE
	case models.ORDER_STATUS_FINISHED:
		return ORDER_STATUS_FINISHED
	case models.ORDER_STATUS_CANCELED:
		return ORDER_STATUS_CANCELED
	case models.ORDER_STATUS_FAILED_PAYMENT:
		return ORDER_STATUS_FAILED_PAYMENT
	default:
		return ORDER_STATUS_UNSET
	}
}

type Payment struct {
	ID        uuid.UUID `gorm:"id,primaryKey"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Value     decimal.Decimal `json:"value"`
	OrderID   uuid.UUID
	Status    PaymentStatus
}

func paymentFromModels(in *models.Payment) *Payment {
	out := &Payment{
		ID:        in.ID,
		CreatedAt: in.CreatedAt,
		Value:     in.Price,
		OrderID:   in.OrderID,
		Status:    paymentStatusFromModel(in.Status),
	}
	if !in.UpdatedAt.IsZero() {
		out.UpdatedAt.Valid = true
		out.UpdatedAt.Time = in.UpdatedAt
	}

	return out
}

func (p *Payment) toModels() *models.Payment {
	out := &models.Payment{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Price:     p.Value,
		OrderID:   p.OrderID,
		Status:    paymentStatusToModel(p.Status),
	}
	if p.UpdatedAt.Valid {
		out.UpdatedAt = p.UpdatedAt.Time
	}

	return out
}

type PaymentStatus string

const (
	PAYMENT_SATUS_OPEN      PaymentStatus = "Aberto"
	PAYMENT_STATUS_APPROVED               = "Aprovado"
	PAYMENT_STATUS_REFUSED                = "Recusado"
)

func paymentStatusFromModel(in models.PaymentStatus) PaymentStatus {
	switch in {
	case models.PAYMENT_STATUS_OPEN:
		return PAYMENT_SATUS_OPEN
	case models.PAYMENT_STATUS_APPROVED:
		return PAYMENT_STATUS_APPROVED
	default:
		return PAYMENT_STATUS_REFUSED
	}
}

func paymentStatusToModel(in PaymentStatus) models.PaymentStatus {
	switch in {
	case PAYMENT_SATUS_OPEN:
		return models.PAYMENT_STATUS_OPEN
	case PAYMENT_STATUS_APPROVED:
		return models.PAYMENT_STATUS_APPROVED
	default:
		return models.PAYMENT_SATUS_REFUSED
	}
}
