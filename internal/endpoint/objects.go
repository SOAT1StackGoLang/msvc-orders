package endpoint

import (
	"github.com/SOAT1StackGoLang/msvc-orders/internal/helpers"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
)

type (
	// CATEGORIES
	CoreCategory struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at,omitempty"`
		Name      string `json:"name"`
	}
	GetCategoryRequest struct {
		ID string `json:"id"`
	}

	GetCategoryResponse struct {
		CoreCategory
	}

	InsertCategoryRequest struct {
		Name string `json:"name" description:"Nome da categoria de produto"`
	}

	InsertCategoryResponse struct {
		CoreCategory
	}

	DeleteCategoryRequest struct {
		ID string `json:"id"`
	}

	DeleteCategoryResponse struct {
		ID string `json:"id"`
	}

	ListCategoriesRequest struct {
		Limit  int `json:"limit" default:"10" description:"Quantidade de registros"`
		Offset int `json:"offset"`
	}

	ListCategoriesResponse struct {
		Categories []CoreCategory `json:"categories"`
		Limit      int            `json:"limit" default:"10"`
		Offset     int            `json:"offset"`
		Total      int64          `json:"total"`
	}
)

type (
	// PRODUCTS

	GetProductRequest struct {
		ID string `json:"id"`
	}

	DeleteProductRequest struct {
		ID string `json:"id"`
	}

	DelectProductResponse struct {
		Deleted bool `json:"Deleted"`
	}

	InsertProductRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
		Price       string `json:"price"`
	}

	UpdateProductRequest struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
		Price       string `json:"price"`
	}

	ProductResponse struct {
		ID          string `json:"id,omitempty"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
		Price       string `json:"price"`
		CreatedAt   string `json:"created_at,omitempty" readOnly:"true"`
		UpdatedAt   string `json:"updated_at,omitempty" readOnly:"true"`
		DeletedAt   string `json:"deleted_at,omitempty" readOnly:"true"`
	}

	ProductList struct {
		Products []ProductResponse `json:"products"`
		Limit    int               `json:"limit" default:"10"`
		Offset   int               `json:"offset"`
		Total    int64             `json:"total"`
	}
)

func ProductResponseFromModel(in *models.Product) ProductResponse {
	out := ProductResponse{
		ID:          in.ID.String(),
		Name:        in.Name,
		Description: in.Description,
		CategoryID:  in.CategoryID.String(),
		Price:       helpers.ParseDecimalToString(in.Price),
		CreatedAt:   in.CreatedAt.String(),
	}
	if !in.UpdatedAt.IsZero() {
		out.UpdatedAt = in.UpdatedAt.String()
	}

	return out
}
