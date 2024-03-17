package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SOAT1StackGoLang/msvc-orders/internal/endpoint"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"

	_ "github.com/SOAT1StackGoLang/msvc-orders/docs" // docs is generated by Swag CLI, you have to import it.
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewOrdersRouter(svc service.OrdersService, r *mux.Router, logger kitlog.Logger) *mux.Router {
	ordersEnpoints := endpoint.MakeOrdersEndpoint(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodGet).Path("/order/all").Queries("limit", "{limit:[0-9]+}", "offset", "{offset:[0-9]+}").Handler(httptransport.NewServer(
		ordersEnpoints.ListOrdersEndpoint,
		decodeListOrdersRequest,
		encodeResponse,
		options...,
	))

	r.Methods(http.MethodGet).Path("/order/{id}").Handler(httptransport.NewServer(
		ordersEnpoints.GetOrderEndpoint,
		decodeGetOrderRequest,
		encodeResponse,
		options...,
	))

	r.Methods(http.MethodPost).Path("/order").Handler(httptransport.NewServer(
		ordersEnpoints.CreateOrderEndpoint,
		decodeCreateOrderRequest,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodPut).Path("/order/items").Handler(httptransport.NewServer(
		ordersEnpoints.UpdateOrderItemsEndpoint,
		decodeAlterOrderItems,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodDelete).Path("/order/{id}").Handler(httptransport.NewServer(
		ordersEnpoints.DeleteOrderEndpoint,
		decodeDeleteOrder,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodGet).Path("/order/checkout/{id}").Handler(httptransport.NewServer(
		ordersEnpoints.OrderCheckoutEndpoint,
		decodeOrderCheckout,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodGet).PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // URL pointing to the API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))
	// redirect / to /swagger/index.html
	r.Methods(http.MethodGet).Path("/").Handler(http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))

	return r
}

// OrderCheckout godoc
//
//	@Summary	Checkout an order
//	@Tags		Orders
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{string}	string	"ok"
//	@Failure	400	{string}	string	"error"
//	@Failure	404	{string}	string	"error"
//	@Failure	500	{string}	string	"error"
//	@Router		/order/checkout/{id} [get]
func decodeOrderCheckout(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.CheckoutOrderRequest{ID: id}, nil
}

// DeleteOrder godoc
//
//	@Summary	Delete an order
//	@Tags		Orders
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{string}	string	"ok"
//	@Failure	400	{string}	string	"error"
//	@Failure	404	{string}	string	"error"
//	@Failure	500	{string}	string	"error"
//	@Router		/order/{id} [delete]
func decodeDeleteOrder(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.DeleteOrderRequest{ID: id}, nil
}

// UpdateOrderItems godoc
//
//	@Summary	Update order items
//	@Tags		Orders
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Success	200	{string}	string	"ok"
//	@Failure	400	{string}	string	"error"
//	@Failure	404	{string}	string	"error"
//	@Failure	500	{string}	string	"error"
//	@Router		/order/items [put]
func decodeAlterOrderItems(_ context.Context, r *http.Request) (request any, err error) {
	var req endpoint.UpdateOrderRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, ErrBadRequest
	}

	return endpoint.UpdateOrderRequest{
		ID:          req.ID,
		ProductsIDs: req.ProductsIDs,
	}, nil
}

// CreateOrder godoc
//
//	@Summary	Create an order
//	@Tags		Orders
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		user_id	header		string	false	"User ID"				default(123e4567-e89b-12d3-a456-426614174000)
//	@Param		request	body		string	true	"Order request data"	SchemaExample({\r\n "products_ids": ["b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12", "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"]\r\n})
//	@Success	200		{string}	string	"ok"
//	@Failure	400		{string}	string	"error"
//	@Failure	500		{string}	string	"error"
//	@Router		/order [post]
func decodeCreateOrderRequest(_ context.Context, r *http.Request) (request any, err error) {
	var (
		req endpoint.CreateOrderRequest
	)

	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, ErrBadRequest
	}
	if req.UserID != "" {
		return nil, ErrBadRequest
	}
	uID := r.Header.Get("user_id")

	return endpoint.CreateOrderRequest{
		UserID:      uID,
		ProductsIDs: req.ProductsIDs,
	}, nil
}

// GetOrder godoc
//
//	@Summary	Get an order
//	@Tags		Orders
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{string}	string	"ok"
//	@Failure	400	{string}	string	"error"
//	@Failure	404	{string}	string	"error"
//	@Failure	500	{string}	string	"error"
//	@Router		/order/{id} [get]
func decodeGetOrderRequest(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.GetOrderRequest{ID: id}, nil
}

// ListOrders godoc
//
//	@Summary	List all orders
//	@Tags		Orders
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		limit	query		int		true	"Limit"		default(10)
//	@Param		offset	query		int		true	"Offset"	default(0)
//	@Success	200		{string}	string	"ok"
//	@Failure	400		{string}	string	"error"
//	@Failure	500		{string}	string	"error"
//	@Router		/order/all [get]
func decodeListOrdersRequest(_ context.Context, r *http.Request) (request any, err error) {
	query := r.URL.Query()
	limit := query.Get("limit")

	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, err
	}
	if limitInt == 0 {
		return nil, ErrBadRequest
	}

	offset := query.Get("offset")

	offsetInt, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		return nil, err
	}

	return endpoint.ListOrderRequest{
		Limit:  int(limitInt),
		Offset: int(offsetInt),
	}, nil
}
