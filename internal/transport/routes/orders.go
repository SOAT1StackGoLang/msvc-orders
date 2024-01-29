package routes

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/endpoint"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"net/http"
)

func NewOrdersRouter(svc service.OrdersService, r *mux.Router, logger kitlog.Logger) *mux.Router {
	ordersEnpoints := endpoint.MakeOrdersEndpoint(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodGet).Path("/order/all").Queries("limit", "{limit:[0-9]+", "offset", "{offset:[0-9]+").Handler(httptransport.NewServer(
		ordersEnpoints.ListOrdersEndpoint,
		mockDecoder,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodGet).Path("/order/{id}").Handler(httptransport.NewServer(
		ordersEnpoints.GetOrderEndpoint,
		mockDecoder,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodPost).Path("/order").Handler(httptransport.NewServer(
		ordersEnpoints.CreateOrderEndpoint,
		mockDecoder,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodPut).Path("/order/items").Handler(httptransport.NewServer(
		ordersEnpoints.UpdateOrderItemsEndpoint,
		mockDecoder,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodDelete).Path("/order").Handler(httptransport.NewServer(
		ordersEnpoints.DeleteOrderEndpoint,
		mockDecoder,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodPost).Path("/order/checkout").Handler(httptransport.NewServer(
		ordersEnpoints.OrderCheckoutEndpoint,
		mockDecoder,
		encodeResponse,
		options...,
	))

	// TODO decoding

	return r
}

func mockDecoder(_ context.Context, r *http.Request) (request any, err error) {
	return nil, nil
}
