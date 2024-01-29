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

	ListProductsByCategoryRequest struct {
		ID     string `json:"id"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
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
		Total    int               `json:"total"`
	}
)

type (
	// PAYMENT

	GetPaymentRequest struct {
		ID string `json:"id"`
	}

	GetPaymentResponse struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Price     string `json:"price"`
		OrderID   string `json:"order_id"`
		Status    string `json:"status"`
	}
)

type (
	ListOrderRequest struct {
		Limit  int `json:"limt"`
		Offset int `json:"offset"`
	}

	GetOrderRequest struct {
		ID string `json:"id"`
	}

	DeleteOrderRequest struct {
		ID string `json:"id"`
	}

	CheckoutOrderRequest struct {
		ID string `json:"id"`
	}

	DeleteOrderResponse struct {
		Deleted string `json:"deleted"`
	}

	GetOrderByPaymentIDRequest struct {
		PaymentID string `json:"payment_id"`
	}

	// ORDERS
	OrderResponse struct {
		ID        string            `json:"id" description:"ID do Pedido"`
		PaymentID string            `json:"payment_id,omitempty" description:"ID do pagamento"`
		CreatedAt string            `json:"created_at" description:"Data de criação"`
		UpdatedAt string            `json:"updated_at,omitempty" description:"Data de atualização"`
		DeletedAt string            `json:"deleted_at,omitempty" description:"Data de deleção"`
		Price     string            `json:"price" description:"Preço do pedido"`
		Status    string            `json:"status" description:"Status do pedido"`
		Products  []ProductResponse `json:"products" description:"Lista de Pedidos"`
	}

	CreateOrderRequest struct {
		UserID      string   `json:"user_id" description:"ID do dono do pedido"`
		ProductsIDs []string `json:"products_ids" description:"ID dos produtos"`
	}

	UpdateOrderRequest struct {
		ID          string   `json:"id"`
		ProductsIDs []string `json:"products_ids" description:"ID dos produtos"`
	}

	OrderList struct {
		Orders []OrderResponse `json:"orders"`
		Limit  int             `json:"limit" default:"10"`
		Offset int             `json:"offset"`
		Total  int             `json:"total"`
	}
)

func OrderResponseFromModel(in *models.Order) OrderResponse {
	out := OrderResponse{
		ID:        in.ID.String(),
		PaymentID: in.PaymentID.String(),
		CreatedAt: in.CreatedAt.String(),
		UpdatedAt: "",
		DeletedAt: "",
		Price:     helpers.ParseDecimalToString(in.Price),
		Status:    string(in.Status),
		Products:  nil,
	}

	if !in.UpdatedAt.IsZero() {
		out.UpdatedAt = in.UpdatedAt.String()
	}

	if !in.DeletedAt.IsZero() {
		out.DeletedAt = in.DeletedAt.String()
	}

	var prods []ProductResponse
	for _, p := range in.Products {
		prods = append(prods, ProductResponseFromModel(&p))
	}
	out.Products = prods
	return out
}

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
