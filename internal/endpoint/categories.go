package endpoint

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/SOAT1StackGoLang/msvc-orders/pkg/helpers"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	CategoryEndpoints struct {
		GetCategoryEndpoint    endpoint.Endpoint
		InsertCategoryEndpoint endpoint.Endpoint
		ListCategoriesEndpoint endpoint.Endpoint
		DeleteCategoryEndpoint endpoint.Endpoint
	}
)

func MakeCategoryEndpoints(svc service.CategoriesService) CategoryEndpoints {
	return CategoryEndpoints{
		GetCategoryEndpoint:    makeGetCategoryEndpoint(svc),
		InsertCategoryEndpoint: makeInsertCategoryEndpoint(svc),
		DeleteCategoryEndpoint: makeDeleteCategoryEndpoint(svc),
		ListCategoriesEndpoint: makeListCategoriesEndpoint(svc),
	}
}

func makeListCategoriesEndpoint(svc service.CategoriesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListCategoriesRequest)

		out := ListCategoriesResponse{
			Limit:  int(req.Limit),
			Offset: int(req.Offset),
		}

		cats, err := svc.ListCategories(ctx, int(req.Limit), int(req.Offset))
		if err != nil {
			return nil, err
		}

		for _, c := range cats.Categories {
			cat := CoreCategory{
				ID:        c.ID.String(),
				CreatedAt: c.CreatedAt.String(),
				Name:      c.Name,
			}
			if !c.UpdatedAt.IsZero() {
				cat.UpdatedAt = c.UpdatedAt.String()
			}
			out.Categories = append(out.Categories, cat)
		}

		return out, nil
	}
}

func makeDeleteCategoryEndpoint(svc service.CategoriesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var out DeleteCategoryResponse
		req := request.(InsertCategoryResponse)

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, err
		}

		err = svc.DeleteCategory(ctx, id)
		if err != nil {
			return nil, err
		}

		out.ID = id.String()
		return out, nil
	}
}
func makeInsertCategoryEndpoint(svc service.CategoriesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var out InsertCategoryResponse
		req := request.(InsertCategoryRequest)

		in := &models.Category{Name: req.Name}

		cat, err := svc.InsertCategory(ctx, in)
		if err != nil {
			return nil, helpers.ErrInvalidInput
		}

		out.ID = cat.ID.String()
		out.Name = cat.Name
		out.CreatedAt = cat.CreatedAt.String()

		return out, nil
	}
}

func makeGetCategoryEndpoint(svc service.CategoriesService) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req := request.(GetCategoryRequest)

		categoryID, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, helpers.ErrInvalidInput
		}

		cat, err := svc.GetCategory(ctx, categoryID)
		if err != nil {
			return nil, err
		}

		var out GetCategoryResponse
		out.ID = cat.ID.String()
		out.CreatedAt = cat.CreatedAt.String()
		out.Name = cat.Name

		if !cat.UpdatedAt.IsZero() {
			out.UpdatedAt = cat.UpdatedAt.String()
		}
		return out, nil
	}
}
