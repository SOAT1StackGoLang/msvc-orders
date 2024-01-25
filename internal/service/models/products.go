package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Product struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Name        string
	Description string
	Price       decimal.Decimal
}

type ProductList struct {
	Products      []*Product
	Limit, Offset int
	Total         int64
}

type ProductsSum struct {
	Products    []uuid.UUID
	RequestedAt time.Time
	Sum         decimal.Decimal
}
