package service

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	"github.com/SOAT1StackGoLang/msvc-orders/pkg/helpers"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type productsSvc struct {
	productRepo persistence.ProductsRepository
	log         kitlog.Logger
}

func (p *productsSvc) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	return p.productRepo.GetProduct(ctx, id)
}

func (p *productsSvc) InsertProduct(ctx context.Context, in *models.Product) (*models.Product, error) {
	if in.Price == decimal.Zero {
		return nil, helpers.ErrBadRequest
	}
	in.ID = uuid.New()
	out, err := p.productRepo.InsertProduct(ctx, in)
	return out, err
}

func (p *productsSvc) UpdateProduct(ctx context.Context, in *models.Product) (*models.Product, error) {
	if in.Price == decimal.Zero {
		return nil, helpers.ErrBadRequest
	}
	return p.productRepo.UpdateProduct(ctx, in)
}

func (p *productsSvc) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return p.productRepo.DeleteProduct(ctx, id)
}

func (p *productsSvc) ListProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) (*models.ProductList, error) {
	return p.productRepo.ListProductsByCategory(ctx, categoryID, limit, offset)
}

func (p *productsSvc) GetProductsPriceSumByID(ctx context.Context, products []uuid.UUID) (*models.ProductsSum, error) {
	return p.productRepo.GetProductsPriceSumByID(ctx, products)
}

func NewProductsService(repo persistence.ProductsRepository, log kitlog.Logger) ProductsService {
	return &productsSvc{
		productRepo: repo,
		log:         log,
	}
}
