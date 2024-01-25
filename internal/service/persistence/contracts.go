package persistence

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/google/uuid"
)

type ProductsRepository interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	InsertProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, uuid uuid.UUID) error
	ListProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) (*models.ProductList, error)
	GetProductsPriceSumByID(ctx context.Context, ids []uuid.UUID) (*models.ProductsSum, error)
}

type CategoriesRepository interface {
	InsertCategory(ctx context.Context, in *models.Category) (*models.Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context, limit int, offset int) (*models.CategoryList, error)
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error)
	GetPayment(ctx context.Context, paymentID uuid.UUID) (*models.Payment, error)
	UpdatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error)
}

type OrdersRepository interface {
	GetOrder(ctx context.Context, orderID uuid.UUID) (*models.Order, error)
	GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*models.Order, error)
	CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	DeleteOrder(ctx context.Context, orderID uuid.UUID) error
	SetOrderAsPaid(ctx context.Context, payment *models.Payment) error
	ListOrdersByUser(ctx context.Context, limit, offset int, userID uuid.UUID) (*models.OrderList, error)
	ListOrders(ctx context.Context, limit, offset int) (*models.OrderList, error)
}
