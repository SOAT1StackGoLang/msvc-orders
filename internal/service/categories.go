package service

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	kitlog "github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

type categoriesSvc struct {
	log         kitlog.Logger
	persistence persistence.CategoriesRepository
}

func NewCategoriesService(persistence persistence.CategoriesRepository, log kitlog.Logger) CategoriesService {
	return &categoriesSvc{
		log:         log,
		persistence: persistence,
	}
}

func (c *categoriesSvc) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	return c.persistence.GetCategoryByID(ctx, id)
}

func (c *categoriesSvc) InsertCategory(ctx context.Context, in *models.Category) (*models.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (c *categoriesSvc) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (c *categoriesSvc) ListCategories(ctx context.Context, limit, offset int) (*models.CategoryList, error) {
	//TODO implement me
	panic("implement me")
}
