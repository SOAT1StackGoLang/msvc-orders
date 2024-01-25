package persistence

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	kitlog "github.com/go-kit/log"
)

const categoriesTable = "lanchonete_categories"

type persistence struct {
	db  *gorm.DB
	log kitlog.Logger
}

func (p *persistence) InsertCategory(ctx context.Context, in *models.Category) (*models.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (p *persistence) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	cat := Category{}

	if err := p.db.WithContext(ctx).Table(categoriesTable).
		Select("*").Where("id = ?", id).First(&cat).Error; err != nil {
		_ = p.log.Log(
			"db failed getting category",
			zap.String("category_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	out := &models.Category{
		ID:        cat.ID,
		CreatedAt: cat.CreatedAt,
		Name:      cat.Name,
	}

	if cat.UpdatedAt.Valid {
		out.CreatedAt = cat.UpdatedAt.Time
	}

	return out, nil
}

func (p *persistence) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p *persistence) ListCategories(ctx context.Context, limit int, offset int) (*models.CategoryList, error) {
	//TODO implement me
	panic("implement me")
}

func NewCategoriesPersistence(ctx context.Context, db *gorm.DB, log kitlog.Logger) CategoriesRepository {
	return &persistence{db: db}
}
