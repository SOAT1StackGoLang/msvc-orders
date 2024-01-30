package models

import (
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/api"
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

func PaymentStatusFromClearingService(status api.PaymentStatus) PaymentStatus {
	switch status {
	case api.PaymentStatusPaid:
		return PAYMENT_STATUS_APPROVED
	case api.PaymentStatusPending:
		return PAYMENT_STATUS_OPEN
	default:
		return PAYMENT_SATUS_REFUSED
	}
}

type PaymentStatusNotification struct {
	PaymentID uuid.UUID
	OrderID   uuid.UUID
	Status    PaymentStatus // Can be "approved" or "denied"
}
