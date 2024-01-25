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
