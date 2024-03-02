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

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /accounts/{id} [get]

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

	return r
}

func decodeOrderCheckout(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.CheckoutOrderRequest{ID: id}, nil
}

func decodeDeleteOrder(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.DeleteOrderRequest{ID: id}, nil
}

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

func decodeCreateOrderRequest(_ context.Context, r *http.Request) (request any, err error) {
	var req endpoint.CreateOrderRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, ErrBadRequest
	}

	return endpoint.CreateOrderRequest{
		UserID:      req.UserID,
		ProductsIDs: req.ProductsIDs,
	}, nil
}

func decodeGetOrderRequest(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.GetOrderRequest{ID: id}, nil
}

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
