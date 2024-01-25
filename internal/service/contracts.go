package service

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/google/uuid"
)

type CategoriesService interface {
	GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error)
	InsertCategory(ctx context.Context, userID uuid.UUID, in *models.Category) (*models.Category, error)
	DeleteCategory(ctx context.Context, userID, id uuid.UUID) error
	ListCategories(ctx context.Context, userID uuid.UUID, limit, offset int) (*models.CategoryList, error)
}
