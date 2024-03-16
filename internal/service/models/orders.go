package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PaymentID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Price     decimal.Decimal
	Status    OrderStatus
	Products  []Product
}

type OrderProductionNotification struct {
	ID        uuid.UUID   `json:"id"`
	UpdatedAt time.Time   `json:"updated_at"`
	Status    OrderStatus `json:"status"`
}

type OrderStatus string

const (
	ORDER_STATUS_UNSET           OrderStatus = ""
	ORDER_STATUS_OPEN                        = "Aberto"
	ORDER_STATUS_WAITING_PAYMENT             = "Aguardando Pagamento"
	ORDER_STATUS_RECEIVED                    = "Recebido"
	ORDER_STATUS_PREPARING                   = "Em Preparação"
	ORDER_STATUS_DONE                        = "Pronto"
	ORDER_STATUS_FINISHED                    = "Finalizado"
	ORDER_STATUS_CANCELED                    = "Cancelado"
	ORDER_STATUS_FAILED_PAYMENT              = "Falha no Pagamento"
)

type OrderList struct {
	Orders        []*Order
	Limit, Offset int
	Total         int64
}
