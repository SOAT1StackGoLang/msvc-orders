package endpoint

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	"github.com/SOAT1StackGoLang/msvc-orders/pkg/helpers"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	CategoryEndpoints struct {
		GetCategoryEndpoint endpoint.Endpoint
	}
)

func MakeCategoryEndpoints(svc service.CategoriesService) CategoryEndpoints {
	return CategoryEndpoints{
		GetCategoryEndpoint: makeGetCategoryEndpoint(svc),
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
