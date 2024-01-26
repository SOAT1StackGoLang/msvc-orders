package persistence

import (
	"context"
	"database/sql"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

const productsTable = "lanchonete_products"

type productsPersistence struct {
	db  *gorm.DB
	log kitlog.Logger
}

func (p *productsPersistence) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var out *models.Product
	product := Product{}

	if err := p.db.WithContext(ctx).Table(productsTable).
		Select("*").Where("id = ?", id).First(&product).Error; err != nil {
		p.log.Log(
			"db failed getting product",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	out.ID = product.ID
	out.Description = product.Description
	out.CategoryID = product.CategoryID
	out.Price = product.Price
	out.CreatedAt = product.CreatedAt

	if product.UpdatedAt.Valid {
		out.UpdatedAt = product.UpdatedAt.Time
	}

	return out, nil
}

func (p *productsPersistence) InsertProduct(ctx context.Context, in *models.Product) (*models.Product, error) {
	product := Product{
		ID:          in.ID,
		CreatedAt:   in.CreatedAt,
		Name:        in.Name,
		Description: in.Description,
		CategoryID:  in.CategoryID,
		Price:       in.Price,
	}

	if err := p.db.WithContext(ctx).Table(productsTable).Create(&product).Error; err != nil {
		p.log.Log(
			"db failed inserting product",
			zap.Any("in_product", in),
			zap.Error(err),
		)
		return nil, err
	}

	out := &models.Product{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	return out, nil
}

func (p *productsPersistence) UpdateProduct(ctx context.Context, in *models.Product) (*models.Product, error) {
	product := Product{
		ID:        in.ID,
		CreatedAt: in.CreatedAt,
		UpdatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
		Name:        in.Name,
		Description: in.Description,
		CategoryID:  in.CategoryID,
		Price:       in.Price,
	}
	product.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err := p.db.WithContext(ctx).Table(productsTable).Updates(&product).Where("id = ?", in.ID).Error; err != nil {
		p.log.Log(
			"db failed updating product",
			zap.Any("in_product", in),
			zap.Error(err),
		)
		return nil, err
	}

	out := &models.Product{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt.Time,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	return out, nil
}

func (p *productsPersistence) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	product := Product{ID: id}
	if err := p.db.WithContext(ctx).Table(productsTable).Delete(&product).Error; err != nil {
		p.log.Log(
			"failed deleting product",
			zap.String("product_id", id.String()),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (p *productsPersistence) ListProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) (*models.ProductList, error) {
	var products []Product
	var total int64

	err := p.db.WithContext(ctx).Table(productsTable).
		Where("category_id = ?", categoryID).Limit(limit).Offset(offset).Order("name ASC").Find(&products).Error
	if err != nil {
		p.log.Log(
			"failed listing products",
			zap.String("category", categoryID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	if err = p.db.WithContext(ctx).Table(productsTable).
		Where("category_id = ?", categoryID).
		Count(&total).Error; err != nil {
		p.log.Log(
			"failed counting products by category id",
			zap.String("category", categoryID.String()),
			zap.Error(err),
		)
	}

	pList := &models.ProductList{}
	out := make([]*models.Product, 0, len(products))
	for _, v := range products {
		product := &models.Product{
			ID:          v.ID,
			CategoryID:  v.CategoryID,
			CreatedAt:   v.CreatedAt,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
		}
		if v.UpdatedAt.Valid {
			product.UpdatedAt = v.UpdatedAt.Time
		}
		out = append(out, product)
	}

	pList.Products = out
	pList.Total = total
	pList.Limit = limit
	pList.Offset = offset

	return pList, err
}

func (p *productsPersistence) GetProductsPriceSumByID(ctx context.Context, ids []uuid.UUID) (*models.ProductsSum, error) {
	type IDAndPrice struct {
		ID    uuid.UUID
		Price decimal.Decimal
	}
	var itemsAndPrices []IDAndPrice

	if err := p.db.WithContext(ctx).Table(productsTable).
		Select("id, price").
		Where("id IN (?)", ids).
		Scan(&itemsAndPrices).
		Error; err != nil {
		p.log.Log(
			"db failed list products and price",
			zap.Any("ids", ids),
			zap.Error(err),
		)
		return nil, err
	}

	prodsSum := &models.ProductsSum{
		Products:    ids,
		RequestedAt: time.Now(),
	}

	var calc decimal.Decimal
	for _, v := range itemsAndPrices {
		calc = calc.Add(v.Price.Abs())
	}

	prodsSum.Sum = calc.Abs()

	return prodsSum, nil
}

func NewProductsPersistence(db *gorm.DB, log kitlog.Logger) ProductsRepository {
	return &productsPersistence{
		db:  db,
		log: log,
	}
}
