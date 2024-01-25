package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Payment struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Price     decimal.Decimal
	OrderID   uuid.UUID
	Status    PaymentStatus
}

type PaymentStatus string

const (
	PAYMENT_STATUS_OPEN     PaymentStatus = "Aberto"
	PAYMENT_STATUS_APPROVED               = "Aprovado"
	PAYMENT_SATUS_REFUSED                 = "Recusado"
)

type PaymentStatusNotification struct {
	PaymentID uuid.UUID
	OrderID   uuid.UUID
	Status    PaymentStatus // Can be "approved" or "denied"
}
