package persistence

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const categoriesTable = "lanchonete_categories"

type catPersistence struct {
	db  *gorm.DB
	log kitlog.Logger
}

func (p *catPersistence) InsertCategory(ctx context.Context, in *models.Category) (*models.Category, error) {
	var out *models.Category
	cat := Category{
		ID:        in.ID,
		CreatedAt: in.CreatedAt,
		Name:      in.Name,
	}

	if err := p.db.WithContext(ctx).Table(categoriesTable).
		Create(&cat).Error; err != nil {
		p.log.Log(
			"db failed inserting category",
			zap.Any("in_category", in),
			zap.Error(err),
		)
		return nil, err
	}

	out = &models.Category{
		ID:        cat.ID,
		CreatedAt: cat.CreatedAt,
		Name:      cat.Name,
	}

	if cat.UpdatedAt.Valid {
		out.CreatedAt = cat.UpdatedAt.Time
	}

	return out, nil
}

func (p *catPersistence) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
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

func (p *catPersistence) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	cat := Category{ID: id}
	if err := p.db.WithContext(ctx).Table(categoriesTable).Delete(&cat).Error; err != nil {
		p.log.Log(
			"db failed deleting category",
			zap.Any("category_id", id.String()),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (p *catPersistence) ListCategories(ctx context.Context, limit int, offset int) (*models.CategoryList, error) {
	var total int64
	var savedCats []Category
	var err error

	// First, perform the count operation
	if err = p.db.WithContext(ctx).Table(categoriesTable).
		Count(&total).
		Error; err != nil {
		p.log.Log(
			"failed counting categories",
			zap.Error(err),
		)
		return nil, err
	}

	if err = p.db.WithContext(ctx).Table(categoriesTable).
		Limit(limit).
		Offset(offset).
		Find(&savedCats).
		Error; err != nil {
		p.log.Log(
			"failed listing categories",
			zap.Error(err),
		)
		return nil, err
	}

	out := &models.CategoryList{}
	outList := make([]*models.Category, 0)

	for _, c := range savedCats {
		out := &models.Category{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			Name:      c.Name,
		}

		if c.UpdatedAt.Valid {
			out.CreatedAt = c.UpdatedAt.Time
		}
		outList = append(outList, out)
	}

	out.Categories = outList
	out.Total = total
	out.Limit = limit
	out.Offset = offset

	return out, err
}

func NewCategoriesPersistence(db *gorm.DB, log kitlog.Logger) CategoriesRepository {
	return &catPersistence{db: db, log: log}
}
