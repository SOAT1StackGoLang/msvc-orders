package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SOAT1StackGoLang/msvc-orders/internal/endpoint"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func NewProductsRouter(svc service.ProductsService, r *mux.Router, logger kitlog.Logger) *mux.Router {
	prodEndpoints := endpoint.MakeProductsEndpoint(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodGet).Path("/product/{id}").Handler(httptransport.NewServer(prodEndpoints.GetProductEndpoint,
		decodeGetProductsRequest,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodPost).Path("/product").Handler(httptransport.NewServer(prodEndpoints.InsertProductEndpoint,
		decodeInsertProductsRequest,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodPut).Path("/product").Handler(httptransport.NewServer(prodEndpoints.UpdateProductEndpoint,
		decodeUpdateProductsRequest,
		encodeResponse,
		options...))
	r.Methods(http.MethodDelete).Path("/product/{id}").Handler(httptransport.NewServer(prodEndpoints.DeleteProductEndpoint,
		decodeDeleteCategoriesRequest,
		encodeResponse,
		options...,
	))

	return r
}

// GetProduct
//
//	@Summary		Get a product by ID
//	@Tags			Products
//	@Description	Get a product by ID
//	@ID				get-product
//	@Produce		json
//	@Param			id	path		string	true	"Product ID"
//	@Success		200	{string}	string	"ok"
//	@Failure		400	{string}	string	"error"
//	@Failure		404	{string}	string	"error"
//	@Failure		500	{string}	string	"error"
//	@Router			/product/{id} [get]
func decodeGetProductsRequest(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.GetProductRequest{ID: id}, nil
}

// InsertProduct
//
//	@Summary		Insert a product
//	@Tags			Products
//	@Description	Insert a product
//	@ID				insert-product
//	@Accept			json
//	@Produce		json
//	@Param			request	body		string	true	"Product data"	SchemaExample({\r\n  "name": "Coca-Cola 2L",\r\n  "description": "Refrigerante Coca-Cola 2L",\r\n  "category_id": "a557b0c0-3bcf-11ee-be56-0242ac120002",\r\n  "price": "10.00"\r\n})
//	@Success		200		{string}	string	"ok"
//	@Failure		400		{string}	string	"error"
//	@Failure		404		{string}	string	"error"
//	@Failure		500		{string}	string	"error"
//	@Router			/product [post]
func decodeInsertProductsRequest(_ context.Context, r *http.Request) (request any, err error) {
	var req endpoint.InsertProductRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, ErrBadRequest
	}

	return endpoint.InsertProductRequest{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Price:       req.Price,
	}, nil
}
func decodeUpdateProductsRequest(_ context.Context, r *http.Request) (request any, err error) {
	var req endpoint.UpdateProductRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, ErrBadRequest
	}

	return endpoint.UpdateProductRequest{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Price:       req.Price,
	}, nil
}
func decodeDeleteProductsRequest(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.DeleteProductRequest{ID: id}, nil
}
