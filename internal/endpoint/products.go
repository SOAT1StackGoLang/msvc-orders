package endpoint

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/helpers"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	ProductsEndpoints struct {
		GetProductEndpoint    endpoint.Endpoint
		InsertProductEndpoint endpoint.Endpoint
		UpdateProductEndpoint endpoint.Endpoint
		DeleteProductEndpoint endpoint.Endpoint
	}
)

func MakeProductsEndpoint(svc service.ProductsService) ProductsEndpoints {
	return ProductsEndpoints{
		GetProductEndpoint:    makeGetProductsEndpoint(svc),
		InsertProductEndpoint: makeInsertProductEndpoint(svc),
		UpdateProductEndpoint: makeUpdateProductEndpoint(svc),
		DeleteProductEndpoint: makeDeleteProductEndpoint(svc),
	}
}

func makeGetProductsEndpoint(svc service.ProductsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetProductRequest)

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, err
		}

		prod, err := svc.GetProduct(ctx, id)

		return ProductResponseFromModel(prod), nil
	}
}

func makeInsertProductEndpoint(svc service.ProductsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(InsertProductRequest)

		price, err := helpers.ParseDecimalFromString(req.Price)
		if err != nil {
			return nil, err
		}

		catID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return nil, err
		}

		product, err := svc.InsertProduct(ctx, &models.Product{
			CategoryID:  catID,
			Name:        req.Name,
			Description: req.Description,
			Price:       price,
		})
		if err != nil {
			return nil, err
		}

		return ProductResponseFromModel(product), nil
	}
}

func makeUpdateProductEndpoint(svc service.ProductsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateProductRequest)

		price, err := helpers.ParseDecimalFromString(req.Price)
		if err != nil {
			return nil, err
		}

		catID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return nil, err
		}

		ID, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, err
		}

		product, err := svc.UpdateProduct(ctx, &models.Product{
			ID:          ID,
			CategoryID:  catID,
			Name:        req.Name,
			Description: req.Description,
			Price:       price,
		})

		return ProductResponseFromModel(product), err
	}
}

func makeDeleteProductEndpoint(svc service.ProductsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteProductRequest)

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, err
		}

		err = svc.DeleteProduct(ctx, id)
		if err != nil {
			return nil, err
		}

		return DelectProductResponse{Deleted: true}, err
	}
}
