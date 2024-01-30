package service

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	kitlog "github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"time"
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
	in.ID = uuid.New()
	in.CreatedAt = time.Now()

	return c.persistence.InsertCategory(ctx, in)
}

func (c *categoriesSvc) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return c.persistence.DeleteCategory(ctx, id)
}

func (c *categoriesSvc) ListCategories(ctx context.Context, limit, offset int) (*models.CategoryList, error) {
	return c.persistence.ListCategories(ctx, limit, offset)
}
