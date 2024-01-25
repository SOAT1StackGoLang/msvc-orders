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
)

type Payment struct {
	ID        uuid.UUID `gorm:"id,primaryKey"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Value     decimal.Decimal `json:"value"`
	OrderID   uuid.UUID
	Status    PaymentStatus
}

type PaymentStatus string

const (
	PAYMENT_SATUS_OPEN      PaymentStatus = "Aberto"
	PAYMENT_STATUS_APPROVED               = "Aprovado"
	PAYMENT_STATUS_REFUSED                = "Recusado"
)
