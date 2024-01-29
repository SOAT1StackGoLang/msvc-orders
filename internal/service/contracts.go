package service

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=contracts.go -package=mocks -destination=../mocks/contracts_mock.go

type CategoriesService interface {
	GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error)
	InsertCategory(ctx context.Context, in *models.Category) (*models.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context, limit, offset int) (*models.CategoryList, error)
}

type ProductsService interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	InsertProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, uuid uuid.UUID) error
	ListProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) (*models.ProductList, error)
	GetProductsPriceSumByID(ctx context.Context, products []uuid.UUID) (*models.ProductsSum, error)
}

type OrdersService interface {
	GetOrder(ctx context.Context, orderID uuid.UUID) (*models.Order, error)
	GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*models.Order, error)
	CreateOrder(ctx context.Context, products []models.Product) (*models.Order, error)
	UpdateOrderItems(ctx context.Context, orderID uuid.UUID, products []models.Product) (*models.Order, error)
	DeleteOrder(ctx context.Context, orderID uuid.UUID) error
	ListOrders(ctx context.Context, limit, offset int) (*models.OrderList, error)
	Checkout(ctx context.Context, paymentID uuid.UUID) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) (*models.Order, error)
	ProcessPayment()
}

type PaymentsService interface {
	GetPayment(ctx context.Context, paymentID uuid.UUID) (*models.Payment, error)
	CreatePayment(ctx context.Context, orderID *models.Order) (*models.Payment, error)
	UpdatePayment(ctx context.Context, paymentID uuid.UUID, status models.PaymentStatus) (*models.Payment, error)
}
