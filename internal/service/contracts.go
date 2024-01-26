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
