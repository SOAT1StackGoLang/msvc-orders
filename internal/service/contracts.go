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
